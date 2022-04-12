package backoff_test

import (
	"fmt"
	"time"

	"github.com/andriiyaremenko/backoff"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Retry", func() {
	const (
		delay = time.Millisecond * 100
	)

	var counter *int
	var defaultBackoff backoff.Backoff
	successCounter := func(counter *int) func() (int, error) {
		return func() (int, error) {
			*counter += 1
			return 42, nil
		}
	}
	retryCounter := func(counter *int) func() error {
		return func() error {
			*counter += 1
			return fmt.Errorf("failed")
		}
	}

	BeforeEach(func() {
		counter = func() *int { i := 0; return &i }()
		defaultBackoff = backoff.Constant(delay).Randomize(time.Millisecond * 100)
	})

	It("should not return until all attempts were taken", func() {
		_, err := backoff.Retry(backoff.Lift(retryCounter(counter)), 4, defaultBackoff)

		Expect(err).Should(HaveOccurred())
		Expect(*counter).To(Equal(5))
	})

	It("should return first successful result", func() {
		v, err := backoff.Retry(successCounter(counter), 4, defaultBackoff)

		Expect(err).ShouldNot(HaveOccurred())
		Expect(v).To(Equal(42))
		Expect(*counter).To(Equal(1))
	})

	It("should retry on error until first success", func() {
		failF := func() func() (int, error) {
			i := 4
			return func() (int, error) {
				if i--; i == 0 {
					return 42, nil
				}

				return 0, retryCounter(counter)()
			}
		}
		v, err := backoff.Retry(failF(), 4, defaultBackoff)

		Expect(err).ShouldNot(HaveOccurred())
		Expect(v).To(Equal(42))
		Expect(*counter).To(Equal(4 - 1))
	})

	It("should accept several backoffs", func() {
		failF := func() func() (int, error) {
			i := 1 + 1 + 1 + 3 + 1 + 4 + 1 + 2 + 1 + 2
			return func() (int, error) {
				if i--; i == 0 {
					return 42, nil
				}

				return 0, retryCounter(counter)()
			}
		}

		check := func(expected int) backoff.Backoff {
			return func(_, _ int) time.Duration {
				*counter -= 1
				Expect(*counter).To(Equal(expected))
				return time.Duration(0)
			}
		}

		v, err := backoff.Retry(
			failF(),
			1,
			backoff.Constant(delay).Randomize(time.Millisecond*100),
			check(1).AsIs(),
			backoff.Linear(time.Millisecond*100, time.Millisecond*10).With(3),
			check(1+3).AsIs(),
			backoff.NaturalExp(time.Millisecond*300).With(4),
			check(1+3+4).AsIs(),
			backoff.Exponential(time.Millisecond*100, 2).With(2),
			check(1+3+4+2).AsIs(),
			backoff.Constant(delay).With(2),
		)

		Expect(*counter).To(Equal(1 + 3 + 4 + 2 + 2))
		Expect(v).To(Equal(42))
		Expect(err).ShouldNot(HaveOccurred())
	})
})
