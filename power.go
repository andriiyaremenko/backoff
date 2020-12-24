package tinybackoff

import (
	"math"
	"sync"
	"time"
)

// Creates Instance of Power Back-off with base of `base`
// delay will be calculated as:
// `delay` * (`base` ^ n), where n is lesser from either attempt number or `stopGrowthAfter`
func NewPowerBackOff(delay time.Duration, base, stopGrowthAfter uint64) ResettableBackOff {
	return &powerBackOff{
		delay:           delay,
		stopGrowthAfter: stopGrowthAfter,
		attemptsCount:   0,
		base:            base}
}

type powerBackOff struct {
	rwM sync.RWMutex

	stopGrowthAfter uint64
	attemptsCount   uint64
	base            uint64
	delay           time.Duration
}

func (p *powerBackOff) NextDelay() time.Duration {
	p.rwM.RLock()
	defer p.rwM.RUnlock()

	p.attemptsCount++

	return p.delay * time.Duration(math.Pow(float64(p.base), float64(p.getCount())))
}

func (p *powerBackOff) Reset() {
	p.rwM.Lock()
	p.attemptsCount = 0
	p.rwM.Unlock()
}

func (p *powerBackOff) getCount() uint64 {
	if p.stopGrowthAfter <= p.attemptsCount {
		return p.stopGrowthAfter
	}

	return p.attemptsCount
}
