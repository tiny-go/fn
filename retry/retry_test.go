package retry

import (
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
}
