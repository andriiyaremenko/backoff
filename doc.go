//This package provides a simple backoff retry (and repeat) mechanism.

// To install pipelines:
// 	go get -u github.com/andriiyaremenko/backoff

/// How to use:
//
// import (
// 	"fmt"
// 	"testing"
// 	"time"
//
// 	"github.com/andriiyaremenko/backoff"
// 	"github.com/stretchr/testify/assert"
// )
//
// func main() {
// 	v, err := backoff.Retry(
// 		// it can return any type
// 		func() (any, error) {
// 			// your function to process and get the result
// 		},
// 		[]int{1, 3, attempts, 2}, // also possible to use just any single int value
// 		backoff.Constant(delay).Randomize(time.Millisecond*100),
// 		backoff.Linear(time.Millisecond*100, time.Millisecond*10),
// 		backoff.Exponential(time.Millisecond*300),
// 		backoff.Power(time.Millisecond*100, 2),
// 		backoff.Constant(delay),
// 	)
// 	// check if err is nil and process response v
//
// 	v, err := backoff.Repeat(
// 		// it can return any type
// 		func() error {
// 			// your function to process
// 		},
// 		[]int{1, 3, attempts, 2}, // also possible to use just any single int value
// 		backoff.Constant(delay).Randomize(time.Millisecond*100),
// 		backoff.Linear(time.Millisecond*100, time.Millisecond*10),
// 		backoff.Exponential(time.Millisecond*300),
// 		backoff.Power(time.Millisecond*100, 2),
// 		backoff.Constant(delay),
// 	)
// 	// check if err is nil
// }
package backoff
