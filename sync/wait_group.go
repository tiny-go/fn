package sync

import "sync"

// WaitGroup is a wrapper over sync.WaitGroup providing additional functionality
type WaitGroup struct{ *sync.WaitGroup }

// NewWaitGroup creates new WaitGroup
func NewWaitGroup() *WaitGroup { return &WaitGroup{&sync.WaitGroup{}} }

// Run adds function to the group (incrementing the group) running in a goroutine
func (w WaitGroup) Run(fn func()) {
	w.Add(1)
	go func() { fn(); w.Done() }()
}
