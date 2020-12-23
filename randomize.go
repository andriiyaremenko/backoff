package tinybackoff

import (
	"math/rand"
	"time"
)

// Creates Instance of Rand Back-off
// delay will be calculated by the rules of `base` `BackOff` + random delay (less than `maxDeviation`)
// `delay` * `multiplier` * n, where n is lesser from either attempt number or `maxAttempts`
func Randomize(base BackOff, maxDeviation time.Duration) BackOff {
	return &randomize{
		maxDeviation: maxDeviation,
		base:         base}
}

type randomize struct {
	base         BackOff
	maxDeviation time.Duration
}

func (r *randomize) NextDelay() time.Duration {
	return r.base.NextDelay() + time.Duration(rand.Int63n(int64(r.maxDeviation)))
}

func (r *randomize) BackOff() BackOff {
	return r.base
}
