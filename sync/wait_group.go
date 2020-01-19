package sync

import "sync"

// WaitGroup is a wrapper over sync.WaitGroup providing additional functionality
type WaitGroup struct{ sync.WaitGroup }

// Run adds provided functions to the group (incrementing the group counter)
// running each callable in a separate goroutine
func (w *WaitGroup) Run(fns ...func()) {
	for _, fn := range fns {
		w.Add(1)
		go func(callable func()) { callable(); w.Done() }(fn)
	}
}
