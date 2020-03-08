package debounce

import (
	"errors"
	"testing"
	"time"
)

const interval = 1000

func Test_Debounce(t *testing.T) {
	t.Run("given debounced function", func(t *testing.T) {
		var (
			eventChan   = make(chan struct{})
			errExpected = errors.New("failed")
		)

		debounced, errChan := New(
			func() error { eventChan <- struct{}{}; return errExpected },
			time.Millisecond*interval,
		)

		for i := 0; i < 5; i++ {
			debounced()
			time.Sleep(time.Millisecond * interval / 10)
		}

		t.Run("test if callable was called only once", func(t *testing.T) {
			<-eventChan

			err := <-errChan
			if err == nil {
				t.Error("debounced function should have failed")
			}
			if err != errExpected {
				t.Errorf("unexpected error: %s", err)
			}
			// test there are no more events
			select {
			case <-eventChan:
				t.Errorf("channel should be empty")
			default:
			}
		})
	})
}
