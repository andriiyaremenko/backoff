# TinyBackOff

Simple back-off library for Golang

This library allows you to run retries and repeats based on configured back-off.

## Key Features

* Configure back-off:
    * Constant back-off
    * Linear back-off
    * Power back-off
    * Exponential back-off
	* Random back-off based on any other back-off
    * Define your own back-off
* Retry operation based on back-off configuration
* Repeat operation based on back-off configuration

## Available functions
#### Back-off
###### Back-off with constant delay
````go
tinybackoff.NewConstantBackOff(delay) tinybackoff.BackOff
````
###### Back-off with delay linear growth
````go
tinybackoff.NewLinearBackOff(delay, delta, attmempts) tinybackoff.ResettableBackOff
````
###### Back-off with delay power growth
````go
tinybackoff.NewPowerBackOff(delay, base, attempts) tinybackoff.ResettableBackOff
````
###### Back-off with delay expotential growth
````go
tinybackoff.NewExponentialBackOff(maxDelay, attempts) tinybackoff.ResettableBackOff
````
###### Back-off with randomized delay
````go
tinybackoff.Randomize(backOff, maxDeviation) tinybackoff.BackOff
````
###### Back-off with max attempts
````go
tinybackoff.WithMaxAttempts(backOff, maxAttempts) tinybackoff.ContinuableResettableBackOff
````
#### Retry
###### Standard retry
````go
tinybackoff.Retry(backOff, operation) error
````
###### Retry until context is cancelled
````go
tinybackoff.RetryUntilSucceeded(ctx, backOff, operation) <-chan error
````
#### Repeat
###### Standard repeat
````go
tinybackoff.Repeat(backOff, operation) error
````
###### Repeat until context is cancelled
````go
tinybackoff.RepeatUntilCancelled(ctx, backOff, operation) <- chan error
````

## Installing
````bash
go get github.com/andriiyaremenko/tinybackoff
````

## Configure BackOff
````go
package main

import (
	"fmt"
	"time"

	"github.com/andriiyaremenko/tinybackoff"
)

var (
	attempts   uint64 = 3
	multiplier uint64 = 2
	base       uint64 = 3
	delay             = time.Second
	delta             = time.Second
	maxDelay          = time.Minute * 30
)

func main() {
	constantBackOff := tinybackoff.NewConstantBackOff(delay)
	linearBackOff := tinybackoff.NewLinearBackOff(delay, delta, attempts)
	powerBackOff := tinybackoff.NewPowerBackOff(delay, attempts, base)
	exponentialBackOff := tinybackoff.NewExponentialBackOff(maxDelay, attempts)

	fmt.Println("First attempt")
	fmt.Println(constantBackOff.NextDelay())
	// Output: 1s
	fmt.Println(linearBackOff.NextDelay())
	// Output: 2s
	fmt.Println(powerBackOff.NextDelay())
	// Output: 3s
	fmt.Println(exponentialBackOff.NextDelay())
	// Output: 4m3.603509825s

	fmt.Println("Second attempt")
	fmt.Println(constantBackOff.NextDelay())
	// Output: 1s
	fmt.Println(linearBackOff.NextDelay())
	// Output: 3s
	fmt.Println(powerBackOff.NextDelay())
	// Output: 9s
	fmt.Println(exponentialBackOff.NextDelay())
	// Output: 11m2.182994108s

	fmt.Println("Third attempt")
	fmt.Println(constantBackOff.NextDelay())
	// Output: 1s
	fmt.Println(linearBackOff.NextDelay())
	// Output: 4s
	fmt.Println(powerBackOff.NextDelay())
	// Output: 27s
	fmt.Println(exponentialBackOff.NextDelay())
	// Output: 30m0s

	fmt.Println(`All attempts after last attempt (> than "attempts")`)
	fmt.Println(constantBackOff.NextDelay())
	// Output: 1s
	fmt.Println(linearBackOff.NextDelay())
	// Output: 4s
	fmt.Println(powerBackOff.NextDelay())
	// Output: 27s
	fmt.Println(exponentialBackOff.NextDelay())
	// Output: 30m0s
}
````

## Run Retry
````go
package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/andriiyaremenko/tinybackoff"
)

var (
	attempts   uint64 = 3
	multiplier uint64 = 2
	delay             = time.Second
	delta             = time.Second
)

func main() {
	linearBackOff := tinybackoff.WithMaxAttempts(tinybackoff.NewLinearBackOff(delay, delta, attempts), attempts)
	failedOperation := func() error { return errors.New("my error") }
	succeededOperation := func() error { return nil }

	// will retunrn nil on first try
	fmt.Println(tinybackoff.Retry(linearBackOff, succeededOperation))
	// Output: <nil>
	// will retunrn error after all attempts were spent
	fmt.Println(tinybackoff.Retry(linearBackOff, failedOperation))
	// Output: "my error"
}
````

## Run RetryUntilSucceeded
````go
package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/andriiyaremenko/tinybackoff"
)

var (
	attempts   uint64 = 3
	multiplier uint64 = 2
	delay             = time.Second
	delta             = time.Second
)

func main() {
	linearBackOff := tinybackoff.NewLinearBackOff(delay, delta, attempts)
	failedOperation := func() error { return errors.New("my error") }
	succeededOperation := func() error { return nil }
	ctx := context.TODO()
	ctx, cancel := context.WithTimeout(ctx, time.Second*2)

	defer cancel()

	// will retunrn nil on first try
	fmt.Println(<-tinybackoff.RetryUntilSucceeded(ctx, linearBackOff, succeededOperation))
	// Output: {}
	// will retunrn context cancellation error after context was cancelled
	fmt.Println(<-tinybackoff.RetryUntilSucceeded(ctx, linearBackOff, failedOperation))
	// Output: "context deadline exceeded"
}
````

## Run Repeat
````go
package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/andriiyaremenko/tinybackoff"
)

var (
	attempts   uint64 = 3
	multiplier uint64 = 2
	delay             = time.Second
	delta             = time.Second
)

func main() {
	linearBackOff := tinybackoff.WithMaxAttempts(tinybackoff.NewLinearBackOff(delay, delta, attempts), attempts)
	failedOperation := func() error { return errors.New("my error") }
	succeededOperation := func() error { return nil }

	// will retunrn error on first try
	fmt.Println(tinybackoff.Repeat(linearBackOff, failedOperation))
	// Output: "my error"
	// will retunrn nil after all attempts were spent
	fmt.Println(tinybackoff.Repeat(linearBackOff, succeededOperation))
	// Output: <nil>
}
````

## Run RepeatUntilCancelled
````go
package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/andriiyaremenko/tinybackoff"
)

var (
	attempts   uint64 = 3
	multiplier uint64 = 2
	delay             = time.Second
	delta             = time.Second
)

func main() {
	linearBackOff := tinybackoff.NewLinearBackOff(delay, delta, attempts)
	failedOperation := func() error { return errors.New("my error") }
	succeededOperation := func() error { return nil }
	ctx := context.TODO()
	ctx, cancel := context.WithTimeout(ctx, time.Second*2)

	defer cancel()

	// will retunrn error on first try
	fmt.Println(<-tinybackoff.RepeatUntilCancelled(ctx, linearBackOff, failedOperation))
	// Output: "my error"
	// will retunrn nil after context was cancelled
	fmt.Println(<-tinybackoff.RepeatUntilCancelled(ctx, linearBackOff, succeededOperation))
	// Output: {}
}
`````
