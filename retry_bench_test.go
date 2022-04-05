package backoff_test

import (
	"errors"
	"testing"

	"github.com/andriiyaremenko/backoff"
)

func BenchmarkRetryTwoAttempts(b *testing.B) {
	backOff := backoff.Constant(0)
	failF := func() (any, error) { return nil, errors.New("failed") }
	for i := 0; i < b.N; i++ {
		_, _ = backoff.Retry[backoff.SameAttempts](failF, 2, backOff)
	}
}

func BenchmarkRetryTwoAttemptsWithLift(b *testing.B) {
	backOff := backoff.Constant(0)
	failF := func() error { return errors.New("failed") }
	for i := 0; i < b.N; i++ {
		_, _ = backoff.Retry[backoff.SameAttempts](backoff.Lift(failF), 2, backOff)
	}
}
