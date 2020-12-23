package tinybackoff

import (
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	attempts = 4
	delay    = time.Millisecond * 100
)

func TestBackOff(t *testing.T) {
	t.Run("ConstantBackOff returns same delay", testConstantBackOff)
	t.Run("LinearBackOff returns delay with linear growth", testLinearBackOff)
	t.Run("PowerBackOff returns delay with growth equals power of `base`", testPowerBackOff)
	t.Run("ExponentialBackOff returns delay with exponential growth until it reaches `maxDelay`",
		testExponentialBackOff)
}

func testConstantBackOff(t *testing.T) {
	assert := assert.New(t)
	delay := time.Millisecond * 100 * 10
	backOff := WithMaxAttempts(NewConstantBackOff(delay), uint64(attempts))

	for i := 0; i < attempts; i++ {
		assert.Equal(true, backOff.Continue())
		assert.Equal(delay, backOff.NextDelay())
	}

	assert.Equal(false, backOff.Continue())
}

func testLinearBackOff(t *testing.T) {
	assert := assert.New(t)
	multiplier := 3
	delay := time.Millisecond * 100 * 10
	backOff := NewLinearBackOff(delay, uint64(attempts), uint64(multiplier))

	for i := 0; i < attempts; i++ {
		attempt := i + 1
		expected := delay * time.Duration(multiplier) * time.Duration(attempt)

		assert.Equal(true, backOff.Continue())
		assert.Equal(expected, backOff.NextDelay())
	}

	assert.Equal(false, backOff.Continue())
}

func testPowerBackOff(t *testing.T) {
	assert := assert.New(t)
	base := 2
	backOff := NewPowerBackOff(delay, uint64(attempts), uint64(base))

	for i := 0; i < attempts; i++ {
		attempt := i + 1
		expected := delay * time.Duration(math.Pow(float64(base), float64(attempt)))

		assert.Equal(true, backOff.Continue())
		assert.Equal(expected, backOff.NextDelay())
	}

	assert.Equal(false, backOff.Continue())
}

func testExponentialBackOff(t *testing.T) {
	assert := assert.New(t)
	maxDelay := time.Hour * 24
	attempts := 7
	backOff := NewExponentialBackOff(maxDelay, uint64(attempts))

	for i := 0; i < attempts; i++ {
		attempt := i + 1
		expected := time.Duration(float64(maxDelay) / math.Exp(float64(attempts-attempt)))

		assert.Equal(true, backOff.Continue())
		assert.Equal(expected, backOff.NextDelay())
	}

	assert.Equal(false, backOff.Continue())
}
