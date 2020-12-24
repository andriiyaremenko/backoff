package tinybackoff

import (
	"time"
)

// Creates Instance of Constant Back-off
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
