package future

import "context"

// Option applies a single retry option to the strategy config
type Option func(o *Options)

// Callback is a future callback (can be used to notify the status of the operation)
type Callback func(err error)

// Factory is a future factory
type Factory func(fn func() error) func() error

// Options is a future config
type Options struct {
	context  context.Context
	callback Callback
}

// generate default config
func defaultOptions() *Options {
	return &Options{
		callback: func(error) {},
		context:  context.Background(),
	}
}

// WithContext defines future context
func WithContext(ctx context.Context) Option {
	return func(options *Options) {
		options.context = ctx
	}
}

// NewFactory creates new future factory with provided options
func NewFactory(opts ...Option) Factory {
	options := defaultOptions()
	for _, opt := range opts {
		opt(options)
	}

	return func(fn func() error) func() error {
		var err error
		done := make(chan struct{}, 1)
		go func() {
			defer close(done)
			err = fn()
		}()

		return func() error {
			select {
			case <-options.context.Done():
				return options.context.Err()
			// wait until callable returns the result
			case <-done:
				go options.callback(err)
				return err
			}
		}
	}
}

// New is a wrapper around Factory with no options provided returning a new future
func New(fn func() error) func() error {
	return NewFactory()(fn)
}

// NewWithContext creates a future func with provided callable and context
func NewWithContext(ctx context.Context, fn func() error) func() error {
	return NewFactory(WithContext(ctx))(fn)
}
