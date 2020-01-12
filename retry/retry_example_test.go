package retry

import (
	"log"
	"time"
)

func ExampleNewStrategy() {
	retry := NewStrategy(
		WithAttempts(10),
		WithInterval(time.Second),
		WithCallback(func(attempt uint, err error) {
			log.Printf("attempt %d failed with error: %s", attempt, err)
		}),
	)

	retryableFunc := func() (err error) {
		// do some retryable stuff here
		return
	}

	if err := retry(retryableFunc); err != nil {
		// handle error: all attempts have failed  :(
	}
}
