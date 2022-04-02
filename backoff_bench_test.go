package backoff_test

import (
	"testing"
	"time"

	"github.com/andriiyaremenko/backoff"
)

func BenchmarkConstant(b *testing.B) {
	backOff := backoff.Constant(time.Millisecond * 10)
	for i := 0; i < b.N; i++ {
		_ = backOff(i, b.N)
	}
}

func BenchmarkLinear(b *testing.B) {
	backOff := backoff.Linear(time.Millisecond*10, time.Millisecond)
	for i := 0; i < b.N; i++ {
		_ = backOff(i, b.N)
	}
}

func BenchmarkPower(b *testing.B) {
	backOff := backoff.Power(time.Millisecond*10, 2)
	for i := 0; i < b.N; i++ {
		_ = backOff(i, b.N)
	}
}

func BenchmarkExpotential(b *testing.B) {
	backOff := backoff.Exponential(time.Millisecond * 10)
	for i := 0; i < b.N; i++ {
		_ = backOff(i, b.N)
	}
}

func BenchmarkRandom(b *testing.B) {
	backOff := backoff.Constant(time.Millisecond * 10).Randomize(time.Millisecond * 100)
	for i := 0; i < b.N; i++ {
		_ = backOff(i, b.N)
	}
}
