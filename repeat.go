package backoff

import (
	"time"
)

// Repeats fn as long as it is successful.
func Repeat[A Attempts](
	fn func() error,
	attempts A,
	backOff Backoff,
	backOffs ...Backoff,
) error {
	backOffs = append([]Backoff{backOff}, backOffs...)

	for n, backOff := range backOffs {
		backOffAttempts := attempts.Next(n)

		for j := 1; j <= backOffAttempts; j++ {
			err := fn()

			if err != nil {
				return err
			}

			time.Sleep(backOff(j, backOffAttempts))
		}
	}

	return nil
}
