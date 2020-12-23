package tinybackoff

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRetry(t *testing.T) {
	t.Run("Retry should run until all attempts where taken", testRetryFail)
	t.Run("Retry should return first encountered success", testRetrySuccess)
	t.Run("Retry should return first encountered success event if failed at firs",
		testRetryFailThenSuccess)
	t.Run("RetryUntilSucceeded should run until first encountered success", testRetryUntilSucceeded)
	t.Run("RetryUntilSucceeded should run until context is cancelled if no success was encountered",
		testRetryUntilSucceededContextCancelled)
}

func testRetryFail(t *testing.T) {
	assert := assert.New(t)
	backOff := WithMaxAttempts(Randomize(NewConstantBackOff(delay), time.Millisecond*100), uint64(attempts))
	failF := func() error { return fmt.Errorf("failed") }
	err := Retry(backOff, failF)

	assert.NotNil(err)
	assert.False(backOff.Continue())
}

func testRetrySuccess(t *testing.T) {
	assert := assert.New(t)
	backOff := WithMaxAttempts(Randomize(NewConstantBackOff(delay), time.Millisecond*100), uint64(attempts))
	failF := func() error { return nil }
	err := Retry(backOff, failF)

	assert.Nil(err)
	assert.True(backOff.Continue())
}

func testRetryFailThenSuccess(t *testing.T) {
	assert := assert.New(t)
	backOff := WithMaxAttempts(Randomize(NewConstantBackOff(delay), time.Millisecond*100), uint64(attempts))
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

func testRetryUntilSucceeded(t *testing.T) {
	assert := assert.New(t)
	ctx := context.TODO()
	backOff := WithMaxAttempts(Randomize(NewConstantBackOff(delay), time.Millisecond*100), uint64(attempts))
	failF := func() error { return nil }
	done := RetryUntilSucceeded(ctx, backOff, failF)

	assert.Eventually(func() bool { return <-done == nil }, time.Millisecond*100*2, time.Millisecond)
}

func testRetryUntilSucceededContextCancelled(t *testing.T) {
	assert := assert.New(t)
	ctx := context.TODO()
	backOff := WithMaxAttempts(Randomize(NewConstantBackOff(delay), time.Millisecond*100), uint64(attempts))
	failF := func() error { return fmt.Errorf("failed") }
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*100)

	defer cancel()

	done := RetryUntilSucceeded(ctx, backOff, failF)

	assert.Eventually(func() bool { return <-done == context.DeadlineExceeded }, time.Millisecond*100*4, time.Millisecond)
}
