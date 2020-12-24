package tinybackoff

import (
	"context"
	"time"
)

// Runs `operation` function based on `backOff` configuration
// returns last `error` returned from `operation` if all attempts have failed
func Retry(backOff ContinuableBackOff, operation Operation) (err error) {
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
// Runs until either `ctx` is cancelled or `operation` returns `nil`
// returns `chan error` (`nil` if operation was successful or `ctx.Err()` if `ctx` was cancelled)
func RetryUntilSucceeded(ctx context.Context, backOff BackOff, operation Operation) <-chan error {
	done := make(chan error)
	go func() {
		defer close(done)
		for {
			select {
			case <-ctx.Done():
				done <- ctx.Err()

				return
			default:
				if continuable := AsContinuable(backOff); continuable != nil && !continuable.Continue() {
					if resettable := AsResettable(continuable); resettable != nil {
						resettable.Reset()

						continue
					}

					done <- CannotResetContinuable

					return
				}

				if err := operation(); err != nil {
					time.Sleep(backOff.NextDelay())
					continue
				}

				done <- nil

				return
			}
		}
	}()

	return done
}
