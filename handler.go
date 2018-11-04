package main

import (
	"errors"
	"log"
	"strings"
	"time"
)

type CheckDomainHandlerType func(domain string, expireThresholdDays int64, whoisBackend WhoisBackend) (freeDate *time.Time, err error)

func CheckDomainHandler(domain string, expireThresholdDays int64, whoisBackend WhoisBackend) (freeDate *time.Time, err error) {
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
