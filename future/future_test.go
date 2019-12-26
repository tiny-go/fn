package future

import "testing"

func Test_NewFuture(t *testing.T) {
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
