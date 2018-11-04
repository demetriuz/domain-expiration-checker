package main

import (
	"time"
)

type WorkRequest struct {
	Domain              string
	ExpireThresholdDays int64
}

type WorkResult struct {
	Domain   string
	freeDate *time.Time
	err      error
}

type WorkQueue chan WorkRequest

type ResultQueue chan WorkResult

type Worker struct {
	id           int
	workQueue    *WorkQueue
	resultQueue  *ResultQueue
	handler      CheckDomainHandlerType
	whoisBackend WhoisBackend
}

func (w *Worker) Start() {
	go func() {
		for {
			select {
			case work := <-*w.workQueue:
				//log.Printf("Worker%d received work %s\n", w.id, work)

				freeDate, err := w.handler(work.Domain, work.ExpireThresholdDays, w.whoisBackend)
				if err == nil {
					*w.resultQueue <- WorkResult{work.Domain, freeDate, nil}
				} else {
					*w.resultQueue <- WorkResult{work.Domain, nil, err}
				}
			}
		}
	}()
}

func NewWorker(
	id int,
	workQueue *WorkQueue,
	resultQueue *ResultQueue,
	handler CheckDomainHandlerType,
	whoisBackend WhoisBackend) Worker {
	return Worker{id, workQueue, resultQueue, handler, whoisBackend}
}

func AddWork(domain string, expireThresholdDays int64, queue WorkQueue) {
	queue <- WorkRequest{domain, expireThresholdDays}
}
