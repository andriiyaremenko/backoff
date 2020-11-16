package tinybackoff

import (
	"time"
)

// Creates Instance of Power Back-off with base of 1
// delay will be calculated as: `delay`
func NewConstantBackOff(delay time.Duration, maxAttempts uint64) BackOff {
	return &powerBackOff{
		delay:         delay,
		maxAttempts:   maxAttempts,
		attemptsCount: 0,
		base:          1}
}
