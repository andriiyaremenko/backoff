package tinybackoff

import (
	"sync"
	"time"
)

func WithMaxAttempts(backOff BackOff, maxAttempts uint64) ContinuableResettableBackOff {
	return &withMaxAttempts{base: backOff, maxAttempts: maxAttempts, attemptsCount: 0}
}

type withMaxAttempts struct {
	rwM sync.RWMutex

	base                       BackOff
	maxAttempts, attemptsCount uint64
}

func (wma *withMaxAttempts) NextDelay() time.Duration {
	wma.rwM.Lock()
	defer wma.rwM.Unlock()

	wma.attemptsCount++

	return wma.base.NextDelay()
}

func (wma *withMaxAttempts) Continue() bool {
	wma.rwM.RLock()
	defer wma.rwM.RUnlock()

	return wma.maxAttempts > wma.attemptsCount
}

func (wma *withMaxAttempts) Reset() {
	wma.rwM.Lock()
	wma.attemptsCount = 0
	wma.rwM.Unlock()
}
