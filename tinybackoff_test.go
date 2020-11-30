package tinybackoff

import (
	"context"
	"errors"
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	attempts = 4
	delay    = time.Second
)

func TestConstantBackOff(t *testing.T) {
	assert := assert.New(t)
	delay := time.Second * 10
	backOff := NewConstantBackOff(delay, uint64(attempts))

	for i := 0; i < attempts; i++ {
		assert.Equal(true, backOff.Continue())
		assert.Equal(delay, backOff.NextDelay())
	}

	assert.Equal(false, backOff.Continue())
}

func TestLinearBackOff(t *testing.T) {
	assert := assert.New(t)
	multiplier := 3
	delay := time.Second * 10
	backOff := NewLinearBackOff(delay, uint64(attempts), uint64(multiplier))

	for i := 0; i < attempts; i++ {
		attempt := i + 1
		expected := delay * time.Duration(multiplier) * time.Duration(attempt)

		assert.Equal(true, backOff.Continue())
		assert.Equal(expected, backOff.NextDelay())
	}

	assert.Equal(false, backOff.Continue())
}

func TestPowerBackOff(t *testing.T) {
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

func TestExponentialBackOff(t *testing.T) {
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

func TestRetryFail(t *testing.T) {
	assert := assert.New(t)
	backOff := NewConstantBackOff(delay, uint64(attempts))
	failF := func() error { return fmt.Errorf("failed") }
	err := Retry(backOff, failF)

	assert.NotNil(err)
	assert.False(backOff.Continue())
}

func TestRetrySuccess(t *testing.T) {
	assert := assert.New(t)
	backOff := NewConstantBackOff(delay, uint64(attempts))
	failF := func() error { return nil }
	err := Retry(backOff, failF)

	assert.Nil(err)
	assert.True(backOff.Continue())
}

func TestRetryFailThenSuccess(t *testing.T) {
	assert := assert.New(t)
	backOff := NewConstantBackOff(delay, uint64(attempts))
	failF := func() func() error {
		i := attempts
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

func TestRetryUntilSucceeded(t *testing.T) {
	assert := assert.New(t)
	ctx := context.TODO()
	backOff := NewConstantBackOff(delay, uint64(attempts))
	failF := func() error { return nil }
	done := RetryUntilSucceeded(ctx, backOff, failF)

	assert.Eventually(func() bool { return <-done == nil }, time.Second*2, time.Millisecond)
}

func TestRetryUntilSucceededContextCancelled(t *testing.T) {
	assert := assert.New(t)
	ctx := context.TODO()
	backOff := NewConstantBackOff(delay, uint64(attempts))
	failF := func() error { return fmt.Errorf("failed") }
	ctx, cancel := context.WithTimeout(ctx, time.Second*2)

	defer cancel()

	done := RetryUntilSucceeded(ctx, backOff, failF)

	assert.Eventually(func() bool { return <-done == context.DeadlineExceeded }, time.Second*4, time.Millisecond)
}

func TestRepeatFail(t *testing.T) {
	assert := assert.New(t)
	backOff := NewConstantBackOff(delay, uint64(attempts))
	failF := func() error { return fmt.Errorf("failed") }
	err := Repeat(backOff, failF)

	assert.NotNil(err)
	assert.True(backOff.Continue())
}

func TestRepeatSuccess(t *testing.T) {
	assert := assert.New(t)
	backOff := NewConstantBackOff(delay, uint64(attempts))
	failF := func() error { return nil }
	err := Repeat(backOff, failF)

	assert.Nil(err)
	assert.False(backOff.Continue())
}

func TestRepeatSuccessThenFail(t *testing.T) {
	assert := assert.New(t)
	backOff := NewConstantBackOff(delay, uint64(attempts))
	failF := func() func() error {
		i := attempts
		return func() error {
			if i--; i == 0 {
				return fmt.Errorf("failed")
			}

			return nil
		}
	}
	err := Repeat(backOff, failF())

	assert.NotNil(err)
}

func TestRepeatUntilCancelledErrorReceived(t *testing.T) {
	assert := assert.New(t)
	ctx := context.TODO()
	backOff := NewConstantBackOff(delay, uint64(attempts))
	err := errors.New("test error")
	failF := func() error { return err }
	done := RepeatUntilCancelled(ctx, backOff, failF)

	assert.Eventually(func() bool { return <-done == err }, time.Second*2, time.Millisecond)
}

func TestRepeatUntilCancelledContextCancelled(t *testing.T) {
	assert := assert.New(t)
	ctx := context.TODO()
	backOff := NewConstantBackOff(delay, uint64(attempts))
	failF := func() error { return nil }
	ctx, cancel := context.WithTimeout(ctx, time.Second*2)

	defer cancel()

	done := RepeatUntilCancelled(ctx, backOff, failF)

	assert.Eventually(func() bool { return <-done == nil }, time.Second*4, time.Millisecond)
}
