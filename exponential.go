package tinybackoff

import (
	"math"
	"sync"
	"time"
)

// Creates Instance of Exponential Back-off
// delay will be calculated as:
// `maxDelay` / exp(`maxAttempts` - n), where n is lesser from either attempt number or `maxAttempts`
func NewExponentialBackOff(maxDelay time.Duration, maxAttempts uint64) ContinuableResettableBackOff {
	return &exponentialBackOff{
		maxAttempts:   maxAttempts,
		attemptsCount: 0,
		maxDelay:      maxDelay}
}

type exponentialBackOff struct {
	rwM sync.RWMutex

	maxAttempts, attemptsCount uint64
	maxDelay                   time.Duration
}

func (e *exponentialBackOff) NextDelay() time.Duration {
	e.rwM.RLock()
	defer e.rwM.RUnlock()

	e.attemptsCount++

	return time.Duration(float64(e.maxDelay) / math.Exp(float64(e.maxAttempts-e.getCount())))
}

func (e *exponentialBackOff) Continue() bool {
	e.rwM.RLock()
	defer e.rwM.RUnlock()

	return e.maxAttempts > e.attemptsCount
}

func (e *exponentialBackOff) Reset() {
	e.rwM.Lock()
	e.attemptsCount = 0
	e.rwM.Unlock()
}

func (e *exponentialBackOff) getCount() uint64 {
	if e.maxAttempts <= e.attemptsCount {
		return e.maxAttempts
	}

	return e.attemptsCount
}
