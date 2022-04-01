package backoff_test

import (
	"math"
	"testing"
	"time"

	"github.com/andriiyaremenko/backoff"
	"github.com/stretchr/testify/assert"
)

const (
	attempts = 4
	delay    = time.Millisecond * 100
)

func TestBackOff(t *testing.T) {
	t.Run("ConstantReturnsSameDelay", testConstantBackOff)
	t.Run("LinearReturnsDelayWithLinearGrowth", testLinearBackOff)
	t.Run("PowerReturnsDelayWithGrowthEqualsPowerOfBase", testPowerBackOff)
	t.Run("ExponentialReturnsDelayWithExponentialGrowth", testExponentialBackOff)
}

func testConstantBackOff(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	delay := time.Millisecond * 100
	backOff := backoff.Constant(delay)

	for i := 1; i <= attempts; i++ {
		assert.Equal(delay, backOff(i, attempts))
	}
}

func testLinearBackOff(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	multiplier := time.Millisecond * 50
	delay := time.Millisecond * 100
	backOff := backoff.Linear(delay, multiplier)

	for i := 1; i <= attempts; i++ {
		expected := delay + multiplier*time.Duration(i)
		assert.Equal(expected, backOff(i, attempts))
	}
}

func testPowerBackOff(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	var base float64 = 2
	backOff := backoff.Power(delay, base)

	for i := 1; i <= attempts; i++ {
		expected := delay * time.Duration(math.Pow(base, float64(i)))
		assert.Equal(expected, backOff(i, attempts))
	}
}

func testExponentialBackOff(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	maxDelay := time.Hour * 24
	attempts := 7
	backOff := backoff.Exponential(maxDelay)

	for i := 1; i <= attempts; i++ {
		expected := time.Duration(float64(maxDelay) / math.Exp(float64(attempts-i)))

		assert.Equal(expected, backOff(i, attempts))
	}
}
