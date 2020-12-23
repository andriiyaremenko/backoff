package tinybackoff

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRepeat(t *testing.T) {
	t.Run("Repeat should return first encountered error", testRepeatFail)
	t.Run("Repeat should run until all attempts where taken", testRepeatSuccess)
	t.Run("Repeat should return first encountered error event if succeeded at first",
		testRepeatSuccessThenFail)
	t.Run("RepeatUtilCancelled should run until first encountered error",
		testRepeatUntilCancelledErrorReceived)
	t.Run("RepeatUtilCancelled should run until context is cancelled if no error was encountered",
		testRepeatUntilCancelledContextCancelled)
}

func testRepeatFail(t *testing.T) {
	assert := assert.New(t)
	backOff := WithMaxAttempts(Randomize(NewConstantBackOff(delay), time.Millisecond*100), uint64(attempts))
	failF := func() error { return fmt.Errorf("failed") }
	err := Repeat(backOff, failF)

	assert.NotNil(err)
	assert.True(backOff.Continue())
}

func testRepeatSuccess(t *testing.T) {
	assert := assert.New(t)
	backOff := WithMaxAttempts(Randomize(NewConstantBackOff(delay), time.Millisecond*100), uint64(attempts))
	failF := func() error { return nil }
	err := Repeat(backOff, failF)

	assert.Nil(err)
	assert.False(backOff.Continue())
}

func testRepeatSuccessThenFail(t *testing.T) {
	assert := assert.New(t)
	backOff := WithMaxAttempts(Randomize(NewConstantBackOff(delay), time.Millisecond*100), uint64(attempts))
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

func testRepeatUntilCancelledErrorReceived(t *testing.T) {
	assert := assert.New(t)
	ctx := context.TODO()
	backOff := WithMaxAttempts(Randomize(NewConstantBackOff(delay), time.Millisecond*100), uint64(attempts))
	err := errors.New("test error")
	failF := func() error { return err }
	done := RepeatUntilCancelled(ctx, backOff, failF)

	assert.Eventually(func() bool { return <-done == err }, time.Millisecond*100*2, time.Millisecond)
}

func testRepeatUntilCancelledContextCancelled(t *testing.T) {
	assert := assert.New(t)
	ctx := context.TODO()
	backOff := WithMaxAttempts(Randomize(NewConstantBackOff(delay), time.Millisecond*100), uint64(attempts))
	failF := func() error { return nil }
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*100*2)

	defer cancel()

	done := RepeatUntilCancelled(ctx, backOff, failF)

	assert.Eventually(func() bool { return <-done == nil }, time.Millisecond*100*4, time.Millisecond)
}
