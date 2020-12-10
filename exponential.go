package tinybackoff

import (
	"math"
	"sync"
	"time"
)

// Creates Instance of Exponential Back-off
// delay will be calculated as:
// `maxDelay` / exp(`maxAttempts` - n), where n is lesser from either attempt number or `maxAttempts`
func NewExponentialBackOff(maxDelay time.Duration, maxAttempts uint64) BackOff {
	return &exponentialBackOff{
		maxAttempts:   maxAttempts,
		attemptsCount: 0,
		maxDelay:      maxDelay}
}

type exponentialBackOff struct {
	mu                         sync.Mutex
	maxAttempts, attemptsCount uint64
	maxDelay                   time.Duration
}

func (e *exponentialBackOff) NextDelay() time.Duration {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.attemptsCount++
	return time.Duration(float64(e.maxDelay) / math.Exp(float64(e.maxAttempts-e.getCount())))
}

func (e *exponentialBackOff) Continue() bool {
	return e.maxAttempts > e.attemptsCount
}

func (e *exponentialBackOff) Reset() {
	e.mu.Lock()
	e.attemptsCount = 0
	e.mu.Unlock()
}

func (e *exponentialBackOff) getCount() uint64 {
	if e.maxAttempts <= e.attemptsCount {
		return e.maxAttempts
	}

	return e.attemptsCount
}
