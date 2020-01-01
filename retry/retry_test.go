package retry

import (
	"context"
	"errors"
	"testing"

	"github.com/tiny-go/fn"
)

const tries = 10

var errExpected = errors.New("error")

func createCallable(successOn uint) fn.Callable {
	var try uint
	return func() error {
		try++
		if try < successOn {
			return errExpected
		}
		return nil
	}
}

func Test_NewStrategy(t *testing.T) {
	t.Run("given a strategy with default options", func(t *testing.T) {
		retry := NewStrategy()

		t.Run("should return `nil` if callable succeeds", func(t *testing.T) {
			err := retry(createCallable(1))
			if err != nil {
				t.Error("error was not expected")
			}
		})

		t.Run("should propagate an error when callable fails", func(t *testing.T) {
			err := retry(createCallable(2))
			if err != errExpected {
				t.Error("error was expected")
			}
		})
	})

	t.Run("given a strategy with custom options", func(t *testing.T) {
		var called int

		retry := NewStrategy(
			WithAttempts(tries),
			WithInterval(1),
			WithCallback(func(uint, error) { called++ }),
		)

		t.Run("should return `nil` if callable succeeds", func(t *testing.T) {
			called = 0

			err := retry(createCallable(tries))
			if err != nil {
				t.Error("error was not expected")
			}

			if called != tries-1 {
				t.Errorf("callback was called %d times (was expected to be called exactly %d times)", called, tries-1)
			}
		})

		t.Run("should propagate an error when callable fails", func(t *testing.T) {
			called = 0

			err := retry(createCallable(tries + 1))
			if err != errExpected {
				t.Error("error was expected")
			}

			// on the last try callback should not be called
			if called != tries-1 {
				t.Errorf("callback was called %d times (was expected to be called exactly %d times)", called, tries-1)
			}
		})
	})

	t.Run("given a strategy with context", func(t *testing.T) {
		var called int

		ctx, cancel := context.WithCancel(context.Background())

		retry := NewStrategy(
			WithAttempts(tries),
			WithInterval(1),
			WithCallback(func(uint, error) { called++ }),
			WithContext(ctx),
		)

		cancel()

		t.Run("should return last error when context is done", func(t *testing.T) {
			called = 0

			err := retry(createCallable(tries))
			if err != errExpected {
				t.Error("error was expected")
			}

			if called != 1 {
				t.Error("callback was expected to be called exactly once")
			}
		})
	})
}
