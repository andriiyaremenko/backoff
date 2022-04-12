package backoff_test

import (
	"errors"
	"testing"

	"github.com/andriiyaremenko/backoff"
)

func BenchmarkRetryTwoAttempts(b *testing.B) {
	backOff := backoff.Constant(0)
	err := errors.New("failed")
	failF := func() (any, error) { return nil, err }
	for i := 0; i < b.N; i++ {
		_, _ = backoff.Retry(failF, 2, backOff)
	}
}

func BenchmarkRetryTwoAttemptsWithTwoBackoffs(b *testing.B) {
	backOff := backoff.Constant(0)
	err := errors.New("failed")
	failF := func() (any, error) { return nil, err }
	for i := 0; i < b.N; i++ {
		_, _ = backoff.Retry(failF, 2, backOff, backOff.AsIs())
	}
}

func BenchmarkRetryTwoAttemptsWithThreeBackoffs(b *testing.B) {
	backOff := backoff.Constant(0)
	err := errors.New("failed")
	failF := func() (any, error) { return nil, err }
	for i := 0; i < b.N; i++ {
		_, _ = backoff.Retry(failF, 2, backOff, backOff.AsIs(), backOff.AsIs())
	}
}

func BenchmarkRetryTwoAttemptsWithLift(b *testing.B) {
	backOff := backoff.Constant(0)
	err := errors.New("failed")
	failF := func() error { return err }
	for i := 0; i < b.N; i++ {
		_, _ = backoff.Retry(backoff.Lift(failF), 2, backOff)
	}
}
