package tinybackoff

import "time"

type BackOff interface {
	NextDelay() time.Duration
	Continue() bool
	Reset()
}

type Operation func() error
