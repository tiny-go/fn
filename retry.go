package fn

import (
	"context"
	"time"
)

// RetryStrategy is a configurable strategy that calls retryable function
type RetryStrategy func(Callable) error

// RetryOption applies a single retry option to the strategy config
type RetryOption func(o *RetryOptions)

// RetryCallback is called after each attempt failure (you can use it to log the
// intermediary error)
type RetryCallback func(n uint, err error)

// RetryOptions is a retry strategy config
type RetryOptions struct {
	attempts uint
	callback RetryCallback
	context  context.Context
	interval time.Duration
}

// generate default config
// NOTE: there is no strategy by default, it wraps provided retryable, calls it
// once and propagates the error as if there is no retry strategy
func defaultOptions() *RetryOptions {
	return &RetryOptions{
		attempts: 1,
		callback: func(uint, error) {},
		context:  context.Background(),
		interval: 1,
	}
}

// WithAttempts sets the number of retries
func WithAttempts(attempts uint) RetryOption {
	return func(options *RetryOptions) {
		options.attempts = attempts
	}
}

// WithInterval allows to configure retry interval
func WithInterval(interval time.Duration) RetryOption {
	return func(options *RetryOptions) {
		options.interval = interval
	}
}

// WithCallback is a function which is going to be called after every failure
func WithCallback(callback RetryCallback) RetryOption {
	return func(options *RetryOptions) {
		options.callback = callback
	}
}

// WithContext defines the context of retry strategy
func WithContext(ctx context.Context) RetryOption {
	return func(options *RetryOptions) {
		options.context = ctx
	}
}

// NewRetryStrategy creates retry strategy with provided options, note that the
// strategy can be defined globally and reused multiple times
func NewRetryStrategy(opts ...RetryOption) RetryStrategy {
	return func(fn Callable) (err error) {
		options := defaultOptions()
		for _, opt := range opts {
			opt(options)
		}

		tick := time.NewTicker(options.interval)
		defer tick.Stop()

		for try := uint(1); try <= options.attempts; try++ {
			if err = fn(); err == nil {
				break
			}
			// no need to wait if last attempt
			if try == options.attempts {
				break
			}
			options.callback(try, err)
			select {
			case <-options.context.Done():
				return err
			case <-tick.C:
			}
		}
		return err
	}
}
