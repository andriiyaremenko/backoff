package tinybackoff

import (
	"math"
	"sync"
	"time"
)

// Creates Instance of Power Back-off with base of `base`
// delay will be calculated as:
// `delay` * (`base` ^ n), where n is lesser from either attempt number or `maxAttempts`
func NewPowerBackOff(delay time.Duration, maxAttempts, base uint64) BackOff {
	return &powerBackOff{
		delay:         delay,
		maxAttempts:   maxAttempts,
		attemptsCount: 0,
		base:          base}
}

type powerBackOff struct {
	mu                         sync.Mutex
	maxAttempts, attemptsCount uint64
	base                       uint64
	delay                      time.Duration
}

func (p *powerBackOff) NextDelay() time.Duration {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.attemptsCount++
	return p.delay * time.Duration(math.Pow(float64(p.base), float64(p.getCount())))
}

func (p *powerBackOff) Continue() bool {
	return p.maxAttempts > p.attemptsCount
}

func (p *powerBackOff) Reset() {
	p.mu.Lock()
	p.attemptsCount = 0
	p.mu.Unlock()
}

func (p *powerBackOff) getCount() uint64 {
	if p.maxAttempts <= p.attemptsCount {
		return p.maxAttempts
	}

	return p.attemptsCount
}
