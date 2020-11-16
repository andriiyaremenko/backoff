package tinybackoff

import (
	"math"
	"sync"
	"time"
)

func NewExponentialBackOff(maxDelay time.Duration, maxRetries uint64) BackOff {
	return &exponentialBackOff{
		maxRetries:   maxRetries,
		retriesCount: 0,
		maxDelay:     maxDelay}
}

type exponentialBackOff struct {
	mu                       sync.Mutex
	maxRetries, retriesCount uint64
	maxDelay                 time.Duration
}

func (e *exponentialBackOff) NextDelay() time.Duration {
	e.mu.Lock()
	e.retriesCount++
	e.mu.Unlock()
	return time.Duration(float64(e.maxDelay) / math.Exp(float64(e.maxRetries-e.getCount())))
}

func (e *exponentialBackOff) Continue() bool {
	return e.maxRetries > e.retriesCount
}

func (e *exponentialBackOff) Reset() {
	e.mu.Lock()
	e.retriesCount = 0
	e.mu.Unlock()
}

func (e *exponentialBackOff) getCount() uint64 {
	if e.maxRetries <= e.retriesCount {
		return e.maxRetries
	}

	return e.retriesCount
}
