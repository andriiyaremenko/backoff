package backoff_test

import (
	"testing"

	"github.com/andriiyaremenko/backoff"
)

func BenchmarkRepeatTwoAttempts(b *testing.B) {
	backOff := backoff.Constant(0)
	failF := func() error { return nil }
	for i := 0; i < b.N; i++ {
		_ = backoff.Repeat(failF, 2, backOff)
	}
}
