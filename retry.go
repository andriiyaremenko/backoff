package backoff

import (
	"time"
)

// Retries fn until it is successful.
// Retry will call fn at least once.
// First fn call is done without delay.
// Total amount of fn calls will equal attempts + 1.
func Retry[T any, Fn func() (T, error)](
	fn Fn,
	attempts int,
	backOff Backoff,
	backOffs ...BackoffOption,
) (T, error) {
	backOffAttempts := attempts

	v, err := fn()
	for i := -1; i < len(backOffs); i++ {
		if i != -1 {
			backOff, backOffAttempts = backOffs[i]()
		}

		if backOffAttempts == -1 {
			backOffAttempts = attempts
		}

		for i := 0; i < backOffAttempts; i++ {
			if err == nil {
				return v, nil
			}

			time.Sleep(backOff(i, backOffAttempts-1))

			v, err = fn()
		}
	}

	return v, err
}

// Lifts function with single error return to one acceptable by Retry
func Lift(fn func() error) func() (any, error) {
	return func() (any, error) { return nil, fn() }
}
