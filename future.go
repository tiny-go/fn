package fn

// NewFuture creates a future func with provided callable (closure)
func NewFuture(fn Callable) func() error {
	var err error
	done := make(chan struct{}, 1)
	go func() {
		defer close(done)
		err = fn()
	}()
	return func() error {
		<-done // wait until callable returns the result
		return err
	}
}
