package tinybackoff

import (
	"sync"
	"time"
)

// Creates Instance of WithMaxAttempts Back-off.
// Delay will be calculated by the rules of `base` `BackOff`.
// `Continue()` will return `false` all attempts were spent or
// if `base` (or anything it wraps) can be converted to `ContinuableBackOff`
// and `base.Continue()` will return false.
// `Reset()` will reset attempts counter to 0
// and call `base.Reset()` if it can be converted to `ResettableBackOff`
func WithMaxAttempts(backOff BackOff, maxAttempts uint64) ContinuableResettableBackOff {
	return &withMaxAttempts{base: backOff, maxAttempts: maxAttempts, attemptsCount: 0}
}

type withMaxAttempts struct {
	rwM sync.RWMutex

	base          BackOff
	maxAttempts   uint64
	attemptsCount uint64
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

	baseContinue := true
	if continuable := AsContinuable(wma.base); continuable != nil {
		baseContinue = continuable.Continue()
	}

	return baseContinue && wma.maxAttempts > wma.attemptsCount
}

func (wma *withMaxAttempts) Reset() {
	wma.rwM.Lock()

	wma.attemptsCount = 0

	if resettable := AsResettable(wma.base); resettable != nil {
		resettable.Reset()
	}

	wma.rwM.Unlock()
}

func (wma *withMaxAttempts) UnwrapBackOff() BackOff {
	return wma.base
}
