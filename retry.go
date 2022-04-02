package backoff

import (
	"time"
)

// Retries fn until it is successful.
func Retry[T any, A int | []int](
	fn func() (T, error),
	attempts A,
	backOff Backoff,
	backOffs ...Backoff,
) (v T, err error) {
	backOffs = append([]Backoff{backOff}, backOffs...)
	attemptsSlice := func(v any) []int {
		attemptsSlice, ok := v.([]int)
		if ok {
			return attemptsSlice
		}

		return []int{v.(int)}
	}(attempts)

	for i, backOff := range backOffs {
		next := i
		if next >= len(attemptsSlice) {
			next = len(attemptsSlice) - 1
		}

		for j := 1; j <= attemptsSlice[next]; j++ {
			v, err = fn()

			if err == nil {
				return
			}

			time.Sleep(backOff(j, attemptsSlice[next]))
		}
	}

	return
}

// Lifts function with single error return to one acceptable by Retry
func Lift(fn func() error) func() (any, error) {
	return func() (any, error) { return nil, fn() }
}
