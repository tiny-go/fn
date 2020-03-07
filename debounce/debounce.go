package debounce

import (
	"time"

	"github.com/tiny-go/fn"
)

// New wraps a callable that needs to be called only once for a group of events
func New(callable fn.Callable, duration time.Duration) (func(), <-chan error) {
	var (
		called  = make(chan struct{})
		errChan = make(chan error)
	)

	go func() {
		t := time.NewTimer(duration)
		t.Stop()

		for {
			select {
			case <-called:
				t.Reset(duration)
			case <-t.C:
				go func() {
					if err := callable(); err != nil {
						errChan <- err
					}
				}()
			}
		}
	}()

	return func() { go func() { called <- struct{}{} }() }, errChan
}
