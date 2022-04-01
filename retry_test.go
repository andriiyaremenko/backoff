package backoff_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/andriiyaremenko/backoff"
	"github.com/stretchr/testify/assert"
)

func TestRetry(t *testing.T) {
	t.Run("ShouldRunUntilAllAttemptsWhereTaken", testRetryFail)
	t.Run("ShouldReturnFirstEncounteredSuccess", testRetrySuccess)
	t.Run("ShouldReturnFirstEncounteredSuccessEventIfFailedAtFirs", testRetryFailThenSuccess)
	t.Run("ShouldAcceptSeveralBackOffs", testRetrySeveralBackOffs)
}

func testRetryFail(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	failF := func() (any, error) { return nil, fmt.Errorf("failed") }
	_, err := backoff.Retry(failF, attempts, backoff.Constant(delay).Randomize(time.Millisecond*100))

	assert.NotNil(err)
}

func testRetrySuccess(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	failF := func() (any, error) { return nil, nil }
	_, err := backoff.Retry(failF, attempts, backoff.Constant(delay).Randomize(time.Millisecond*100))

	assert.NoError(err)
}

func testRetryFailThenSuccess(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	failF := func() func() (any, error) {
		i := attempts
		return func() (any, error) {
			if i--; i == 0 {
				return nil, nil
			}

			return nil, fmt.Errorf("failed")
		}
	}
	_, err := backoff.Retry(failF(), attempts, backoff.Constant(delay).Randomize(time.Millisecond*100))

	assert.Nil(err)
}

func testRetrySeveralBackOffs(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	failF := func() func() (any, error) {
		i := 1 + 3 + attempts + 2
		return func() (any, error) {
			if i--; i == 0 {
				return nil, nil
			}

			return nil, fmt.Errorf("failed")
		}
	}
	_, err := backoff.Retry(
		failF(),
		[]int{1, 3, attempts, 2},
		backoff.Constant(delay).Randomize(time.Millisecond*100),
		backoff.Linear(time.Millisecond*100, time.Millisecond*10),
		backoff.Exponential(time.Millisecond*300),
		backoff.Power(time.Millisecond*100, 2),
		backoff.Constant(delay),
	)

	assert.Nil(err)
}
