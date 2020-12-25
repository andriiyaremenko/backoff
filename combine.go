package tinybackoff

import (
	"sync"
	"time"
)

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

	i        int
	current  ContinuableBackOff
	backOffs []ContinuableBackOff
}

func (c *combine) NextDelay() time.Duration {
	c.rwM.Lock()
	defer c.rwM.Unlock()

	return c.nextDelay()
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

func (c *combine) nextDelay() time.Duration {
	if c.current.Continue() ||
		len(c.backOffs)-1 <= c.i {
		return c.current.NextDelay()
	}

	c.i++
	c.current = c.backOffs[c.i]

	return c.nextDelay()
}
