package tinybackoff

import (
	"time"
)

type Resettable interface {
	// resets attempts count to `0`
	Reset()
}

type Continuable interface {
	// returns `false` if all attempts were spent
	Continue() bool
}

type Stoppable interface {
	Continuable
	// makes `Continue()` return `false`
	Stop()
	// restores `Continue()` calculation logic
	CarryOn()
}

type BackOff interface {
	// returns next delay
	NextDelay() time.Duration
}

type ContinuableBackOff interface {
	BackOff
	Continuable
}

type ResettableBackOff interface {
	BackOff
	Resettable
}

type ContinuableResettableBackOff interface {
	BackOff
	Continuable
	Resettable
}

type StoppableBackOff interface {
	BackOff
	Stoppable
}

type Operation func() error

// Unwraps `backOff` if it wraps another `BackOff`
func UnwrapBackOff(backOff BackOff) BackOff {
	b, ok := backOff.(interface{ UnwrapBackOff() BackOff })
	if !ok {
		return nil
	}

	return b.UnwrapBackOff()
}

// Returns `ContinuableBackOff` if `backOff` (or anything it wraps) can be converted to it
// or `nil` otherwise
func AsContinuable(backOff BackOff) ContinuableBackOff {
	continuable, ok := backOff.(ContinuableBackOff)
	if ok {
		return continuable
	}

	backOff = UnwrapBackOff(backOff)
	if backOff != nil {
		return AsContinuable(backOff)
	}

	return nil
}

// Returns `ResettableBackOff` if `backOff` (or anything it wraps) can be converted to it
// or `nil` otherwise
func AsResettable(backOff BackOff) ResettableBackOff {
	resettable, ok := backOff.(ResettableBackOff)
	if ok {
		return resettable
	}

	backOff = UnwrapBackOff(backOff)
	if backOff != nil {
		return AsResettable(backOff)
	}

	return nil
}

// Returns `StoppableBackOff` if `backOff` (or anything it wraps) can be converted to it
// or `nil` otherwise
func AsStoppable(backOff BackOff) StoppableBackOff {
	stoppable, ok := backOff.(StoppableBackOff)
	if ok {
		return stoppable
	}

	backOff = UnwrapBackOff(backOff)
	if backOff != nil {
		return AsStoppable(backOff)
	}

	return nil
}
