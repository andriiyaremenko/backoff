package backoff

import (
	"time"
)

// Repeats fn as long as it is successful.
// Repeat will call fn at least once.
// First fn call is done without delay.
// Total amount of fn calls will equal attempts + 1.
func Repeat(fn func() error, attempts int, backOff Backoff, backOffs ...BackoffOption) error {
	backOffAttempts := attempts

	err := fn()
	for i := -1; i < len(backOffs); i++ {
		if i != -1 {
			backOff, backOffAttempts = backOffs[i]()
		}

		if backOffAttempts == -1 {
			backOffAttempts = attempts
		}

		for i := 0; i < backOffAttempts; i++ {
			if err != nil {
				return err
			}

			time.Sleep(backOff(i, backOffAttempts-1))

			err = fn()
		}
	}

	return err
}
