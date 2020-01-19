package sync

import "testing"

const total = 10

func Test_WaitGroup(t *testing.T) {
	t.Run("given a WaitGroup", func(t *testing.T) {
		var wg WaitGroup

		t.Run("should wait until all functions exit", func(t *testing.T) {
			resultChan := make(chan struct{})
			for i := 0; i < total; i++ {
				wg.Run(func() {
					resultChan <- struct{}{}
				})
			}
			go func() {
				for i := 0; i < total; i++ {
					<-resultChan
				}
			}()
			wg.Wait()
		})
	})
}
