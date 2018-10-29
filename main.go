package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/demetriuz/domain-expiration-checker/whois_backends"
)

const FREEDATE_FIELD_PREFIX = "free-date:"

type domainsType []string

func (i *domainsType) String() string {
	return strings.Join(*i, " ")
}

func (i *domainsType) Set(value string) error {
	*i = append(*i, value)
	return nil
}

type WhoisBackend interface {
	Fetch(domain string) (string, error)
}

func main() {
	var domains domainsType
	var expireThresholdDays *int64
	//var whoisBackend = whois_backends.SystemWhoisBackend{}
	var whoisBackend = whois_backends.RawWhoisBackend{}

	flag.Var(&domains, "d", "Domains")
	expireThresholdDays = flag.Int64("t", 30, "Expire Threshold Days")

	flag.Parse()

	for _, domain := range domains {
		freeDate, err := checkDomain(domain, *expireThresholdDays, whoisBackend)

		if err != nil {
			if freeDate != nil {
				fmt.Printf("%s: %s\n", domain, *freeDate)
			} else {
				fmt.Printf("%s: %s\n", domain, err)
			}
		}
	}
}

func checkDomain(domain string, expireThresholdDays int64, whoisBackend WhoisBackend) (freeDate *time.Time, err error) {
	out, err := whoisBackend.Fetch(domain)

	if err != nil {
		log.Fatal(err)
	}

	outStrings := strings.Split(out, "\n")
	for _, line := range outStrings {
		if strings.Contains(line, FREEDATE_FIELD_PREFIX) {
			line = strings.Replace(line, FREEDATE_FIELD_PREFIX, "", -1)
			line = strings.Trim(line, " ")

			// https://golang.org/src/time/format.go
			// layout: stdLongYear-stdZeroMonth-stdZeroDay
			expirationDate, err := time.Parse("2006-01-02", line)
			if err != nil {
				return nil, errors.New("can't parse date")
			}
			timeDelta := int64(expirationDate.Unix()) - int64(time.Now().Unix())
			if timeDelta < (expireThresholdDays * 24 * 60 * 60) {
				return &expirationDate, errors.New("expiration threshold is reached")
			} else {
				return &expirationDate, nil
			}
		}
	}
	return nil, errors.New("NotFound")
}
