package tinybackoff

import (
	"sync"
	"time"
)

func NewLinearBackOff(delay time.Duration, maxRetries, multiplier uint64) BackOff {
	return &linearBackOff{
		delay:        delay,
		maxRetries:   maxRetries,
		retriesCount: 0,
		multiplier:   multiplier}
}

type linearBackOff struct {
	mu                       sync.Mutex
	maxRetries, retriesCount uint64
	multiplier               uint64
	delay                    time.Duration
}

func (l *linearBackOff) NextDelay() time.Duration {
	l.mu.Lock()
	l.retriesCount++
	l.mu.Unlock()
	return l.delay * time.Duration(l.multiplier) * time.Duration(l.getCount())
}

func (l *linearBackOff) Continue() bool {
	return l.maxRetries > l.retriesCount
}

func (l *linearBackOff) Reset() {
	l.mu.Lock()
	l.retriesCount = 0
	l.mu.Unlock()
}

func (l *linearBackOff) getCount() uint64 {
	if l.maxRetries <= l.retriesCount {
		return l.maxRetries
	}

	return l.retriesCount
}
