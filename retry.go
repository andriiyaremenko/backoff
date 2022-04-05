package backoff

import (
	"time"
)

// Retries fn until it is successful.
func Retry[A Attempts, T any, Fn func() (T, error)](
	fn Fn,
	attempts A,
	backOff Backoff,
	backOffs ...Backoff,
) (T, error) {
	backOffs = append([]Backoff{backOff}, backOffs...)
	var (
		v   T
		err error
	)

	for n, backOff := range backOffs {
		backOffAttempts := attempts.Next(n)

		for i := 1; i <= backOffAttempts; i++ {
			v, err = fn()

			if err == nil {
				return v, nil
			}

			time.Sleep(backOff(i, backOffAttempts))
		}
	}

	return v, err
}

// Lifts function with single error return to one acceptable by Retry
func Lift(fn func() error) func() (any, error) {
	return func() (any, error) { return nil, fn() }
}
