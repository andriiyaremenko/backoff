package tinybackoff

import (
	"math"
	"sync"
	"time"
)

func NewPowerBackOff(delay time.Duration, maxRetries, base uint64) BackOff {
	return &powerBackOff{
		delay:        delay,
		maxRetries:   maxRetries,
		retriesCount: 0,
		base:         base}
}

type powerBackOff struct {
	mu                       sync.Mutex
	maxRetries, retriesCount uint64
	base                     uint64
	delay                    time.Duration
}

func (p *powerBackOff) NextDelay() time.Duration {
	p.mu.Lock()
	p.retriesCount++
	p.mu.Unlock()
	return p.delay * time.Duration(math.Pow(float64(p.base), float64(p.getCount())))
}

func (p *powerBackOff) Continue() bool {
	return p.maxRetries > p.retriesCount
}

func (p *powerBackOff) Reset() {
	p.mu.Lock()
	p.retriesCount = 0
	p.mu.Unlock()
}

func (p *powerBackOff) getCount() uint64 {
	if p.maxRetries <= p.retriesCount {
		return p.maxRetries
	}

	return p.retriesCount
}
