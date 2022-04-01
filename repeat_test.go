package backoff_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/andriiyaremenko/backoff"
	"github.com/stretchr/testify/assert"
)

func TestRepeat(t *testing.T) {
	t.Run("ShouldReturnFirstEncounteredError", testRepeatFail)
	t.Run("ShouldRunUntilAllAttemptsWhereTaken", testRepeatSuccess)
	t.Run("ShouldReturnFirstEncounteredErrorEventIfSucceededAtFirst", testRepeatSuccessThenFail)
	t.Run("ShouldAcceptSeveralBackOffs", testRepeatSeveralBackOffs)
}

func testRepeatFail(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	failF := func() error { return fmt.Errorf("failed") }
	err := backoff.Repeat(failF, attempts, backoff.Constant(delay).Randomize(time.Millisecond*100))

	assert.Error(err)
}

func testRepeatSuccess(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	failF := func() error { return nil }
	err := backoff.Repeat(failF, attempts, backoff.Constant(delay).Randomize(time.Millisecond*100))

	assert.NoError(err)
}

func testRepeatSuccessThenFail(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	failF := func() func() error {
		i := attempts
		return func() error {
			if i--; i == 0 {
				return fmt.Errorf("failed")
			}

			return nil
		}
	}
	err := backoff.Repeat(failF(), attempts, backoff.Constant(delay).Randomize(time.Millisecond*100))

	assert.Error(err)
}

func testRepeatSeveralBackOffs(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	err := backoff.Repeat(
		func() error { return nil },
		[]int{1, 3, attempts, 2},
		backoff.Constant(delay).Randomize(time.Millisecond*100),
		backoff.Linear(time.Millisecond*100, time.Millisecond*10),
		backoff.Exponential(time.Millisecond*300),
		backoff.Power(time.Millisecond*100, 2),
		backoff.Constant(delay),
	)

	assert.Nil(err)
}
