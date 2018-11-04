package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/demetriuz/domain-expiration-checker/whois_backends"
)

const FREEDATE_FIELD_PREFIX = "free-date:"
const WORKERS = 5

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

	var workQueue = make(WorkQueue, 100)
	var resultQueue = make(ResultQueue, 100)

	for i := 0; i < WORKERS; i++ {
		w := NewWorker(i, &workQueue, &resultQueue, CheckDomainHandler, whoisBackend)
		w.Start()
	}

	for _, domain := range domains {
		AddWork(domain, *expireThresholdDays, workQueue)
	}

	domainsCount := len(domains)
	var resultsCount = 0
	for res := range resultQueue {
		resultsCount += 1

		if res.err != nil {
			fmt.Printf("%s: %s\n", res.Domain, res.err)
		} else {
			fmt.Printf("%s: %s\n", res.Domain, *res.freeDate)
		}

		if resultsCount == domainsCount {
			break
		}
	}
}
