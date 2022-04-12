package backoff_test

import (
	"math"
	"time"

	"github.com/andriiyaremenko/backoff"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Backoff", func() {
	const (
		attempts = 4
		delay    = time.Millisecond * 100
	)

	It("constant should return same delay", func() {
		backOff := backoff.Constant(delay)
		for i := 1; i <= attempts; i++ {
			Expect(backOff(i, attempts)).To(Equal(delay))
		}
	})

	It("linear should returns delay with linear growth", func() {
		multiplier := time.Millisecond * 50
		backOff := backoff.Linear(delay, multiplier)
		for i := 1; i <= attempts; i++ {
			expected := delay + multiplier*time.Duration(i)
			Expect(backOff(i, attempts)).To(Equal(expected))
		}
	})

	It("exponential should return delay multiplied by base power of attempt number", func() {
		var base float64 = 2
		backOff := backoff.Exponential(delay, base)
		for i := 1; i <= attempts; i++ {
			expected := delay * time.Duration(math.Pow(base, float64(i)))
			Expect(backOff(i, attempts)).To(Equal(expected))
		}
	})

	It("natural exponent should return delay with exponential growth", func() {
		maxDelay := time.Hour * 24
		attempts := 7
		backOff := backoff.NaturalExp(maxDelay)
		for i := 1; i <= attempts; i++ {
			expected := time.Duration(float64(maxDelay) / math.Exp(float64(attempts-i)))
			Expect(backOff(i, attempts)).To(Equal(expected))
		}
	})
})
