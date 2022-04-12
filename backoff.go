package backoff

import (
	"math"
	"math/rand"
	"time"
)

type Backoff func(attempt, attempts int) time.Duration

// Adds random delay within deviation limit to the base delay.
func (backOff Backoff) Randomize(deviation time.Duration) Backoff {
	return func(attempt, attempts int) time.Duration {
		r := rand.New(rand.NewSource(time.Now().Unix()))

		return backOff(attempt, attempts) + time.Duration(r.Int63n(int64(deviation)))
	}
}

// Back-off with constant delay.
func Constant(delay time.Duration) Backoff {
	return func(_, _ int) time.Duration {
		return delay
	}
}

// Back-off with delay exponentially growing from the smallest till specified maximum.
// Uses natural exponent.
func NaturalExp(delay time.Duration) Backoff {
	return func(attempt, attempts int) time.Duration {
		return time.Duration(float64(delay) / math.Exp(float64(attempts-attempt)))
	}
}

// Back-off with delay growing linearly growing by delta.
func Linear(delay, delta time.Duration) Backoff {
	return func(attempt, _ int) time.Duration {
		return delay + (delta * time.Duration(attempt))
	}
}

// Back-off with delay growing by power if base.
func Exponential(delay time.Duration, base float64) Backoff {
	return func(attempt, _ int) time.Duration {
		return delay * time.Duration(math.Pow(base, float64(attempt)))
	}
}
