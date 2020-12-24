package tinybackoff

import (
	"sync"
	"time"
)

// Creates Instance of Linear Back-off
// delay will be calculated as:
// `delay` + `delta` * n, where n is lesser from either attempt number or `stopGrowthAfter`
func NewLinearBackOff(delay, delta time.Duration, stopGrowthAfter uint64) ResettableBackOff {
	return &linearBackOff{
		delay:           delay,
		delta:           delta,
		stopGrowthAfter: stopGrowthAfter,
		attemptsCount:   0}
}

type linearBackOff struct {
	rwM sync.RWMutex

	stopGrowthAfter uint64
	attemptsCount   uint64
	delta           time.Duration
	delay           time.Duration
}

func (l *linearBackOff) NextDelay() time.Duration {
	l.rwM.RLock()
	defer l.rwM.RUnlock()

	l.attemptsCount++

	return l.delay + (l.delta * time.Duration(l.getCount()))
}

func (l *linearBackOff) Reset() {
	l.rwM.Lock()
	l.attemptsCount = 0
	l.rwM.Unlock()
}

func (l *linearBackOff) getCount() uint64 {
	if l.stopGrowthAfter <= l.attemptsCount {
		return l.stopGrowthAfter
	}

	return l.attemptsCount
}
