package tinybackoff

import (
	"context"
	"time"
)

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
