# backoff

[![GoDoc](https://img.shields.io/badge/pkg.go.dev-doc-blue)](http://pkg.go.dev/github.com/andriiyaremenko/backoff)

This package provides a simple backoff retry (and repeat) mechanism.

### To install pipelines:
```go
go get -u github.com/andriiyaremenko/backoff
```

### How to use:
```go
import (
	"fmt"
	"testing"
	"time"

	"github.com/andriiyaremenko/backoff"
	"github.com/stretchr/testify/assert"
)

func main() {
	v, err := backoff.Retry(
		// it can return any type
		func() (any, error) {
			// your function to process and get the result
		},
		1,
		backoff.Constant(time.Second).Randomize(time.Millisecond*100),
		backoff.Linear(time.Millisecond*100, time.Millisecond*10).With(2),
	)
	// check if err is nil and process response v

	err := backoff.Repeat(
		func() error {
			// your function to process
		},
		2,
		backoff.Constant(time.Second).Randomize(time.Millisecond*100),
		backoff.Linear(time.Millisecond*100, time.Millisecond*10).With(4),
	)
	// check if err is nil
}
```
"
