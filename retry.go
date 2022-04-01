package backoff

import (
	"time"
)

// Retries fn until it is successful.
func Retry[T any, A int | []int](
	fn func() (T, error),
	attempts A,
	backOff BackOff,
	backOffs ...BackOff,
) (v T, err error) {
	backOffs = append([]BackOff{backOff}, backOffs...)
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
