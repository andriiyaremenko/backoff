package tinybackoff

import (
	"sync"
	"time"
)

// Creates Instance of Combine Back-off.
// Delay will be calculated by the rules of current `BackOff`.
// Each time `current.Continue()` returns `false` `next` `BackOff` becomes `current`
// until no more `next` left.
// `Reset()` will reset current to first `next` `BackOff`
// and call `Reset()` on any of `next` (or anything it wraps) if it can be converted to `ResettableBackOff`
func Combine(delay time.Duration, next ...ContinuableBackOff) ContinuableResettableBackOff {
	return &combine{current: &once{delay: delay}, backOffs: next, i: -1}
}

type once struct {
	rwM sync.RWMutex

	delay  time.Duration
	isUsed bool
}

func (o *once) NextDelay() time.Duration {
	o.rwM.Lock()
	defer o.rwM.Unlock()

	o.isUsed = true

	return o.delay
}

func (o *once) Continue() bool {
	o.rwM.RLock()
	defer o.rwM.RUnlock()

	return !o.isUsed
}

type combine struct {
	rwM sync.RWMutex

	i         int
	current   ContinuableBackOff
	backOffs  []ContinuableBackOff
	lastDelay time.Duration
}

func (c *combine) NextDelay() time.Duration {
	if !c.Continue() {
		return c.lastDelay
	}

	c.rwM.Lock()
	defer c.rwM.Unlock()

	c.lastDelay = c.nextDelay()

	return c.lastDelay
}

func (c *combine) Continue() bool {
	c.rwM.RLock()
	defer c.rwM.RUnlock()

	return c.current.Continue() || len(c.backOffs)-1 > c.i
}

func (c *combine) Reset() {
	c.rwM.Lock()
	c.i = 0

	if len(c.backOffs) > 0 {
		c.current = c.backOffs[c.i]
	}

	for _, b := range c.backOffs {
		if resettable := AsResettable(b); resettable != nil {
			resettable.Reset()
		}
	}

	c.rwM.Unlock()
}

func (c *combine) UnwrapBackOff() BackOff {
	return c.current
}

func (c *combine) nextDelay() time.Duration {
	if c.current.Continue() ||
		len(c.backOffs)-1 <= c.i {
		return c.current.NextDelay()
	}

	c.i++
	c.current = c.backOffs[c.i]

	return c.nextDelay()
}
