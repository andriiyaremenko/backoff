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
	Restart()
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

type StoppableResettableBackOff interface {
	BackOff
	Stoppable
	Resettable
}

type Operation func() error

func UnwrapBackOff(backOff BackOff) BackOff {
	b, ok := backOff.(interface{ UnwrapBackOff() BackOff })
	if !ok {
		return nil
	}

	return b.UnwrapBackOff()
}

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

func AsResettable(backOff BackOff) ResettableBackOff {
	continuable, ok := backOff.(ResettableBackOff)
	if ok {
		return continuable
	}

	backOff = UnwrapBackOff(backOff)
	if backOff != nil {
		return AsResettable(backOff)
	}

	return nil
}

func AsStoppable(backOff BackOff) StoppableBackOff {
	continuable, ok := backOff.(StoppableBackOff)
	if ok {
		return continuable
	}

	backOff = UnwrapBackOff(backOff)
	if backOff != nil {
		return AsStoppable(backOff)
	}

	return nil
}
