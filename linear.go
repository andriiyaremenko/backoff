package tinybackoff

import (
	"sync"
	"time"
)

// Creates Instance of Linear Back-off
// delay will be calculated as:
// `delay` * `multiplier` * n, where n is lesser from either attempt number or `maxAttempts`
func NewLinearBackOff(delay time.Duration, maxAttempts, multiplier uint64) ContinuableResettableBackOff {
	return &linearBackOff{
		delay:         delay,
		maxAttempts:   maxAttempts,
		attemptsCount: 0,
		multiplier:    multiplier}
}

type linearBackOff struct {
	rwM sync.RWMutex

	maxAttempts, attemptsCount uint64
	multiplier                 uint64
	delay                      time.Duration
}

func (l *linearBackOff) NextDelay() time.Duration {
	l.rwM.RLock()
	defer l.rwM.RUnlock()

	l.attemptsCount++

	return l.delay * time.Duration(l.multiplier) * time.Duration(l.getCount())
}

func (l *linearBackOff) Continue() bool {
	l.rwM.RLock()
	defer l.rwM.RUnlock()

	return l.maxAttempts > l.attemptsCount
}

func (l *linearBackOff) Reset() {
	l.rwM.Lock()
	l.attemptsCount = 0
	l.rwM.Unlock()
}

func (l *linearBackOff) getCount() uint64 {
	if l.maxAttempts <= l.attemptsCount {
		return l.maxAttempts
	}

	return l.attemptsCount
}
