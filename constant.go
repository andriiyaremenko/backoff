package tinybackoff

import (
	"time"
)

func NewConstantBackOff(delay time.Duration, maxRetries uint64) BackOff {
	return &powerBackOff{
		delay:        delay,
		maxRetries:   maxRetries,
		retriesCount: 0,
		base:         1}
}
