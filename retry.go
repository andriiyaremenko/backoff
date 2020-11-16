package tinybackoff

import (
	"context"
	"time"
)

// Runs `operation` function based on `backOff` configuration
// returns last error returned from `operation` if all attempts have failed
func Retry(backOff BackOff, operation Operation) (err error) {
	for backOff.Continue() {
		if err = operation(); err != nil {
			time.Sleep(backOff.NextDelay())
			continue
		}

		return nil
	}

	return
}

// Runs `operation` function based on `backOff` configuration
// Runs until either `ctx` is cancelled or `operation` returns nil
// returns channel to inform when operation has been completed
func RetryUntilSucceeded(ctx context.Context, backOff BackOff, operation Operation) <-chan struct{} {
	done := make(chan struct{})
	go func() {
		for {
			select {
			case <-ctx.Done():
				done <- struct{}{}
				return
			default:
				if err := operation(); err != nil {
					time.Sleep(backOff.NextDelay())
					continue
				}

				done <- struct{}{}
				return
			}
		}
	}()

	return done
}
