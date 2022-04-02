package backoff

import (
	"time"
)

// Repeats fn as long as it is successful.
func Repeat[A int | []int](
	fn func() error,
	attempts A,
	backOff Backoff,
	backOffs ...Backoff,
) (err error) {
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
			err = fn()

			if err != nil {
				return
			}

			time.Sleep(backOff(j, attemptsSlice[next]))
		}
	}

	return
}
