package tinybackoff

import (
	"sync"
	"time"
)

// Creates Instance of WithStopAndCarryOn Back-off.
// Adds ability to stop and carry on wrapped `backOff`
func WithStopAndCarryOn(backOff BackOff) StoppableBackOff {
	return &withStopRestart{base: backOff, isRunning: true}
}

type withStopRestart struct {
	rwM sync.RWMutex

	base      BackOff
	isRunning bool
	lastDelay time.Duration
}

func (wsr *withStopRestart) NextDelay() time.Duration {
	if !wsr.Continue() {
		return wsr.lastDelay
	}

	wsr.rwM.Lock()
	defer wsr.rwM.Unlock()

	wsr.lastDelay = wsr.base.NextDelay()

	return wsr.lastDelay
}

func (wsr *withStopRestart) Continue() bool {
	wsr.rwM.RLock()
	defer wsr.rwM.RUnlock()

	baseContinue := true
	if continuable := AsContinuable(wsr.base); continuable != nil {
		baseContinue = continuable.Continue()
	}

	return wsr.isRunning && baseContinue
}

func (wsr *withStopRestart) CarryOn(delay time.Duration) {
	time.Sleep(delay)
	wsr.rwM.Lock()

	wsr.isRunning = true

	wsr.rwM.Unlock()
}

func (wsr *withStopRestart) Stop() {
	wsr.rwM.Lock()

	wsr.isRunning = false

	wsr.rwM.Unlock()
}

func (wsr *withStopRestart) UnwrapBackOff() BackOff {
	return wsr.base
}
