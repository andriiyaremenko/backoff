package backoff

import (
	"math"
	"math/rand"
	"time"
)

type BackOff func(attempt, attempts int) time.Duration

// Adds random delay within deviation limit to the base delay.
func (backOff BackOff) Randomize(deviation time.Duration) BackOff {
	return func(attempt, attempts int) time.Duration {
		r := rand.New(rand.NewSource(time.Now().Unix()))

		return backOff(attempt, attempts) + time.Duration(r.Int63n(int64(deviation)))
	}
}

// Back-off with constant delay.
func Constant(delay time.Duration) BackOff {
	return func(_, _ int) time.Duration {
		return delay
	}
}

// Back-off with delay exponentially growing from the smallest till specified maximum.
func Exponential(delay time.Duration) BackOff {
	return func(attempt, attempts int) time.Duration {
		return time.Duration(float64(delay) / math.Exp(float64(attempts-attempt)))
	}
}

// Back-off with delay growing linearly growing by delta.
func Linear(delay, delta time.Duration) BackOff {
	return func(attempt, _ int) time.Duration {
		return delay + (delta * time.Duration(attempt))
	}
}

// Back-off with delay growing by power if base.
func Power(delay time.Duration, base float64) BackOff {
	return func(attempt, _ int) time.Duration {
		return delay * time.Duration(math.Pow(base, float64(attempt)))
	}
}
