package tinybackoff

import (
	"context"
	"time"
)

// Runs `operation` function based on `backOff` configuration until all attempts were spent
// returns first `error` returned from `operation` or `nil`
func Repeat(backOff BackOff, operation Operation) error {
	for backOff.Continue() {
		if err := operation(); err != nil {
			return err
		}

		time.Sleep(backOff.NextDelay())
	}

	return nil
}

// Runs `operation` function based on `backOff` configuration
// Runs until either `ctx` is cancelled or `operation` returns `error`
// returns `chan error` (`nil` if operation was successful or `ctx` was cancelled, `error` if `operation` has failed)
func RepeatUntilCancelled(ctx context.Context, backOff BackOff, operation Operation) <-chan error {
	done := make(chan error)
	go func() {
		defer close(done)
		for {
			select {
			case <-ctx.Done():
				done <- nil

				return
			default:
				if err := operation(); err != nil {
					done <- err

					return
				}

				time.Sleep(backOff.NextDelay())
			}
		}
	}()

	return done
}
