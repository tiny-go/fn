package fn

import "net/http"

func ExampleNewFuture() {
	// define the response globally to be set in the callable
	var rsp *http.Response

	future := NewFuture(func() (err error) {
		rsp, err = http.Get("http://example.com")
		return err
	})
	// do some stuff here

	if err := future(); err != nil {
		// handle error
	}
	defer rsp.Body.Close()
	// use the response
}
