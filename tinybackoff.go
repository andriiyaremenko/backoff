package tinybackoff

import "time"

type BackOff interface {
	// returns next delay
	NextDelay() time.Duration
	// returns `false` if all attempts were spent
	Continue() bool
	// resets attempts count to `0`
	Reset()
}

type Operation func() error
