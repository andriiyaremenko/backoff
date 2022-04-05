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
	attemptsCounter := func(counter *int, backOff backoff.Backoff) backoff.Backoff {
		return func(attempt, attempts int) time.Duration {
			*counter += 1
			return backOff(attempt, attempts)
		}
	}

	BeforeEach(func() {
		counter = func() *int { i := 0; return &i }()
		defaultBackoff = attemptsCounter(counter, backoff.Constant(delay).Randomize(time.Millisecond*100))
	})

	It("should not return until all attempts were taken", func() {
		failF := func() error { return fmt.Errorf("failed") }
		_, err := backoff.Retry[backoff.SameAttempts](backoff.Lift(failF), 4, defaultBackoff)

		Expect(err).Should(HaveOccurred())
		Expect(*counter).To(Equal(4))
	})

	It("should return first successful result", func() {
		failF := func() (int, error) { return 42, nil }
		v, err := backoff.Retry[backoff.SameAttempts](failF, 4, defaultBackoff)

		Expect(err).ShouldNot(HaveOccurred())
		Expect(v).To(Equal(42))
		Expect(*counter).To(Equal(0))
	})

	It("should retry on error until first success", func() {
		failF := func() func() (int, error) {
			i := 4
			return func() (int, error) {
				if i--; i == 0 {
					return 42, nil
				}

				return 0, fmt.Errorf("failed")
			}
		}
		v, err := backoff.Retry[backoff.SameAttempts](failF(), 4, defaultBackoff)

		Expect(err).ShouldNot(HaveOccurred())
		Expect(v).To(Equal(42))
		Expect(*counter).To(Equal(4 - 1))
	})

	It("should accept several backoffs", func() {
		failF := func() func() (int, error) {
			i := 1 + 3 + 4 + 2 + 2
			return func() (int, error) {
				if i--; i == 0 {
					return 42, nil
				}

				return 0, fmt.Errorf("failed")
			}
		}
		counter0 := func() *int { i := 0; return &i }()
		counter1 := func() *int { i := 0; return &i }()
		counter2 := func() *int { i := 0; return &i }()
		counter3 := func() *int { i := 0; return &i }()
		counter4 := func() *int { i := 0; return &i }()
		v, err := backoff.Retry[backoff.DifferentAttempts](
			failF(),
			[]int{1, 3, 4, 2},
			attemptsCounter(counter0, backoff.Constant(delay).Randomize(time.Millisecond*100)),
			attemptsCounter(counter1, backoff.Linear(time.Millisecond*100, time.Millisecond*10)),
			attemptsCounter(counter2, backoff.Exponential(time.Millisecond*300)),
			attemptsCounter(counter3, backoff.Power(time.Millisecond*100, 2)),
			attemptsCounter(counter4, backoff.Constant(delay)),
		)

		Expect(*counter0).To(Equal(1))
		Expect(*counter1).To(Equal(3))
		Expect(*counter2).To(Equal(4))
		Expect(*counter3).To(Equal(2))
		Expect(*counter4).To(Equal(1))
		Expect(v).To(Equal(42))
		Expect(err).ShouldNot(HaveOccurred())
	})
})
