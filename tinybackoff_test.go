package tinybackoff

import (
	"context"
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConstantBackOff(t *testing.T) {
	assert := assert.New(t)
	delay := time.Second * 10
	retries := 4
	backOff := NewConstantBackOff(delay, uint64(retries))

	for i := 0; i < retries; i++ {
		assert.Equal(true, backOff.Continue())
		assert.Equal(delay, backOff.NextDelay())
	}

	assert.Equal(false, backOff.Continue())
}

func TestLinearBackOff(t *testing.T) {
	assert := assert.New(t)
	delay := time.Second * 10
	retries := 4
	multiplier := 3
	backOff := NewLinearBackOff(delay, uint64(retries), uint64(multiplier))

	for i := 0; i < retries; i++ {
		attempt := i + 1
		expected := delay * time.Duration(multiplier) * time.Duration(attempt)

		assert.Equal(true, backOff.Continue())
		assert.Equal(expected, backOff.NextDelay())
	}

	assert.Equal(false, backOff.Continue())
}

func TestPowerBackOff(t *testing.T) {
	assert := assert.New(t)
	delay := time.Second
	retries := 4
	base := 2
	backOff := NewPowerBackOff(delay, uint64(retries), uint64(base))

	for i := 0; i < retries; i++ {
		attempt := i + 1
		expected := delay * time.Duration(math.Pow(float64(base), float64(attempt)))

		assert.Equal(true, backOff.Continue())
		assert.Equal(expected, backOff.NextDelay())
	}

	assert.Equal(false, backOff.Continue())
}

func TestExponentialBackOff(t *testing.T) {
	assert := assert.New(t)
	maxDelay := time.Hour * 24
	retries := 7
	backOff := NewExponentialBackOff(maxDelay, uint64(retries))

	for i := 0; i < retries; i++ {
		attempt := i + 1
		expected := time.Duration(float64(maxDelay) / math.Exp(float64(retries-attempt)))

		assert.Equal(true, backOff.Continue())
		assert.Equal(expected, backOff.NextDelay())
	}

	assert.Equal(false, backOff.Continue())
}

func TestRetryFail(t *testing.T) {
	assert := assert.New(t)
	delay := time.Second
	retries := 4
	backOff := NewConstantBackOff(delay, uint64(retries))
	failF := func() error { return fmt.Errorf("failed") }
	err := Retry(backOff, failF)

	assert.NotNil(err)
}

func TestRetrySuccess(t *testing.T) {
	assert := assert.New(t)
	delay := time.Second
	retries := 4
	backOff := NewConstantBackOff(delay, uint64(retries))
	failF := func() error { return nil }
	err := Retry(backOff, failF)

	assert.Nil(err)
}

func TestRetryFailThanSuccess(t *testing.T) {
	assert := assert.New(t)
	delay := time.Second
	retries := 4
	backOff := NewConstantBackOff(delay, uint64(retries))
	failF := func() func() error {
		i := retries
		return func() error {
			if i--; i == 0 {
				return nil
			}

			return fmt.Errorf("failed")
		}
	}
	err := Retry(backOff, failF())

	assert.Nil(err)
}

func TestRetryUntilSucceed(t *testing.T) {
	assert := assert.New(t)
	ctx := context.TODO()
	delay := time.Second
	retries := 4
	backOff := NewConstantBackOff(delay, uint64(retries))
	failF := func() error { return nil }
	done := RetryUntilSucceeded(ctx, backOff, failF)

	assert.Eventually(func() bool { <-done; return true }, time.Second*2, time.Millisecond)
}

func TestRetryUntilSucceedContextCancelled(t *testing.T) {
	assert := assert.New(t)
	ctx := context.TODO()
	delay := time.Second
	retries := 4
	backOff := NewConstantBackOff(delay, uint64(retries))
	failF := func() error { return fmt.Errorf("failed") }
	ctx, cancel := context.WithTimeout(ctx, time.Second*2)

	defer cancel()

	done := RetryUntilSucceeded(ctx, backOff, failF)

	assert.Eventually(func() bool { <-done; return true }, time.Second*4, time.Millisecond)
}
