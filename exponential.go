package tinybackoff

import (
	"math"
	"sync"
	"time"
)

// Creates Instance of Exponential Back-off
// delay will be calculated as:
// `maxDelay` / exp(`attemptsToReachMax` - n), where n is lesser from either attempt number or `attemptsToReachMax`
func NewExponentialBackOff(maxDelay time.Duration, attemptsToReachMax uint64) ResettableBackOff {
	return &exponentialBackOff{
		attemptsToReachMax: attemptsToReachMax,
		attemptsCount:      0,
		maxDelay:           maxDelay}
}

type exponentialBackOff struct {
	rwM sync.RWMutex

	attemptsToReachMax uint64
	attemptsCount      uint64
	maxDelay           time.Duration
}

func (e *exponentialBackOff) NextDelay() time.Duration {
	e.rwM.RLock()
	defer e.rwM.RUnlock()

	e.attemptsCount++

	return time.Duration(float64(e.maxDelay) / math.Exp(float64(e.attemptsToReachMax-e.getCount())))
}

func (e *exponentialBackOff) Reset() {
	e.rwM.Lock()
	e.attemptsCount = 0
	e.rwM.Unlock()
}

func (e *exponentialBackOff) getCount() uint64 {
	if e.attemptsToReachMax <= e.attemptsCount {
		return e.attemptsToReachMax
	}

	return e.attemptsCount
}
