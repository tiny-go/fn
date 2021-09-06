package retry

import (
	"context"
	"time"
)

type Retryable func() (abort bool, err error)

// Strategy is a configurable strategy that calls retryable function
type Strategy func(Retryable) error

// Option applies a single retry option to the strategy config
type Option func(options *Options)

// Callback is called after each attempt failure (you can use it to log the
// intermediary error)
type Callback func(n uint, err error)

// Options is a retry strategy config
type Options struct {
	attempts uint
	callback Callback
	context  context.Context
	interval time.Duration
}

// generate default config
// NOTE: there is no strategy by default, it wraps provided retryable, calls it
// once and propagates the error as if there is no retry strategy
func defaultOptions() *Options {
	return &Options{
		attempts: 1,
		callback: func(uint, error) {},
		context:  context.Background(),
		interval: 1,
	}
}

// WithAttempts sets the number of retries
func WithAttempts(attempts uint) Option {
	return func(options *Options) {
		options.attempts = attempts
	}
}

// WithInterval allows to configure retry interval
func WithInterval(interval time.Duration) Option {
	return func(options *Options) {
		options.interval = interval
	}
}

// WithCallback is a function which is going to be called after every failure
func WithCallback(callback Callback) Option {
	return func(options *Options) {
		options.callback = callback
	}
}

// WithContext defines the context of retry strategy
func WithContext(ctx context.Context) Option {
	return func(options *Options) {
		options.context = ctx
	}
}

// NewStrategy creates retry strategy with provided options, note that the
// strategy can be defined globally and reused multiple times
func NewStrategy(opts ...Option) Strategy {
	options := defaultOptions()
	for _, opt := range opts {
		opt(options)
	}

	return func(fn Retryable) (err error) {
		var (
			tick  = time.NewTicker(options.interval)
			abort bool
		)

		defer tick.Stop()

		for try := uint(1); try <= options.attempts; try++ {
			if abort, err = fn(); err == nil || abort {
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
