package tinybackoff

import (
	"time"
)

// Creates Instance of Power Back-off with base of 1
// delay will be calculated as: `delay`
func NewConstantBackOff(delay time.Duration) BackOff {
	return &constantBackOff{delay: delay}
}

type constantBackOff struct {
	delay time.Duration
}

func (c *constantBackOff) NextDelay() time.Duration {
	return c.delay
}
