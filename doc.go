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
// 		1,
// 		backoff.Constant(delay).Randomize(time.Millisecond*100),
// 		backoff.Linear(time.Millisecond*100, time.Millisecond*10).WithAttempts(2),
// 	)
// 	// check if err is nil and process response v
//
// 	err := backoff.Repeat(
// 		func() error {
// 			// your function to process
// 		},
// 		2,
// 		backoff.Constant(time.Second).Randomize(time.Millisecond*100),
// 		backoff.Linear(time.Millisecond*100, time.Millisecond*10).WithAttempts(4),
// 	)
// 	// check if err is nil
// }
package backoff
