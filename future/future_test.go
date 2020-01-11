package future

import (
	"context"
	"testing"
)

func Test_New(t *testing.T) {
	t.Run("future should wait until callable returns the result", func(t *testing.T) {
		called := make(chan struct{}, 1)

		callable := func() error {
			called <- struct{}{}
			return nil
		}

		New(callable)()

		select {
		case <-called:
		default:
			t.Error("callable func was expected to be called")
		}
	})
}

func Test_NewWithContext(t *testing.T) {
	t.Run("future should wait until callable returns the result", func(t *testing.T) {
		called := make(chan struct{}, 1)

		callable := func() error {
			called <- struct{}{}
			return nil
		}

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		err := NewWithContext(ctx, callable)()
		if err.Error() != "context canceled" {
			t.Error(`"context canceled" error was expected`)
		}

		select {
		case <-called:
			t.Error("callable func was not expected to be called")
		default:
		}
	})
}
