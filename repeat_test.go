package backoff_test

import (
	"fmt"
	"time"

	"github.com/andriiyaremenko/backoff"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Repeat", func() {
	const (
		delay = time.Millisecond * 100
	)

	var counter *int
	var defaultBackoff backoff.Backoff
	errorCounter := func(counter *int) func() error {
		return func() error {
			*counter += 1
			return fmt.Errorf("failed")
		}
	}
	repeatCounter := func(counter *int) func() error {
		return func() error {
			*counter += 1
			return nil
		}
	}

	BeforeEach(func() {
		counter = func() *int { i := 0; return &i }()
		defaultBackoff = backoff.Constant(delay).Randomize(time.Millisecond * 100)
	})

	It("should return first encountered error", func() {
		err := backoff.Repeat(errorCounter(counter), 4, defaultBackoff)

		Expect(err).Should(HaveOccurred())
		Expect(*counter).To(Equal(1))
	})

	It("should not return until all attempts were taken", func() {
		err := backoff.Repeat(repeatCounter(counter), 4, defaultBackoff)

		Expect(err).ShouldNot(HaveOccurred())
		Expect(*counter).To(Equal(5))
	})

	It("should return first error after successful attempts", func() {
		failF := func() func() error {
			i := 4
			return func() error {
				if i--; i == 0 {
					return fmt.Errorf("failed")
				}

				return repeatCounter(counter)()
			}
		}
		err := backoff.Repeat(failF(), 4, defaultBackoff)

		Expect(err).Should(HaveOccurred())
		Expect(*counter).To(Equal(4 - 1))
	})

	It("should accept several backoffs", func() {
		check := func(expected int) backoff.Backoff {
			return func(_, _ int) time.Duration {
				*counter -= 1

				Expect(*counter).To(Equal(expected))
				return time.Duration(0)
			}
		}

		err := backoff.Repeat(
			repeatCounter(counter),
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

		Expect(*counter).To(Equal(1 + 1 + 3 + 4 + 2 + 2))
		Expect(err).ShouldNot(HaveOccurred())
	})
})
