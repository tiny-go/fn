package future

import (
	"context"
	"testing"
)

func Test_New(t *testing.T) {
	t.Run("future should wait until callable returns the result", func(t *testing.T) {
		var externalVar bool

		callable := func() error {
			externalVar = true
			return nil
		}

		New(callable)()

		if externalVar != true {
			t.Error("eternal variable was expected to be modified")
		}
	})
}

func Test_NewWithContext(t *testing.T) {
	t.Run("future should wait until callable returns the result", func(t *testing.T) {
		var externalVar bool

		callable := func() error {
			externalVar = true
			return nil
		}

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		err := NewWithContext(ctx, callable)()
		if err.Error() != "context canceled" {
			t.Error(`"context canceled" error was expected`)
		}

		if externalVar != false {
			t.Error("external variable should not be modified")
		}
	})
}
