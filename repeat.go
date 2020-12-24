package tinybackoff

import (
	"context"
	"time"
)

// Runs `operation` function based on `backOff` configuration until all attempts were taken
// returns first `error` returned from `operation` or `nil`
func Repeat(backOff ContinuableBackOff, operation Operation) error {
	for backOff.Continue() {
		if err := operation(); err != nil {
			return err
		}

		time.Sleep(backOff.NextDelay())
	}

	return nil
}

// Runs `operation` function based on `backOff` configuration.
// Runs until either `ctx` is cancelled or `operation` returns `error`.
// Returns `chan error` (`nil` if operation was successful or `ctx` was cancelled).
func RepeatUntilCancelled(ctx context.Context, backOff BackOff, operation Operation) <-chan error {
	done := make(chan error)

	go func() {
		defer close(done)

		var lastErr error
		for {
			select {
			case <-ctx.Done():
				done <- nil

				return
			default:
				if continuable := AsContinuable(backOff); continuable != nil && !continuable.Continue() {
					if resettable := AsResettable(continuable); resettable != nil {
						resettable.Reset()
					}

					if continuable.Continue() {
						continue
					}

					done <- lastErr

					return
				}

				if lastErr = operation(); lastErr != nil {
					done <- lastErr

					return
				}

				time.Sleep(backOff.NextDelay())
			}
		}
	}()

	return done
}
