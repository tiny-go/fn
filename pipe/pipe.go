package pipe

import "sync"

// Extract elements from a slice to a channel.
func Extract[T any](elements ...T) <-chan T {
	output := make(chan T)
	go func() {
		defer close(output)

		for _, element := range elements {
			output <- element
		}
	}()

	return output
}

// Collect elements from the channel to a slice.
func Collect[T any](input <-chan T) []T {
	output := []T{}
	for element := range input {
		output = append(output, element)
	}

	return output
}

// Pipe reads elements from input channel and sends to all provided outputs.
// It discards the input channel when no outputs provided.
func Pipe[T any](in <-chan T, outs ...chan<- T) {
	for v := range in {
		for _, out := range outs {
			go func(ch chan<- T) { ch <- v }(out)
		}
	}
}

// Map higher-order function implementation with channels.
func Map[I, O any](input <-chan I, mapFunc func(element I) O) <-chan O {
	output := make(chan O)
	go func() {
		defer close(output)

		for element := range input {
			output <- mapFunc(element)
		}
	}()

	return output
}

// Filter higher-order function implementation with channels.
func Filter[I any](input <-chan I, filterFunc func(element I) bool) <-chan I {
	output := make(chan I)
	go func() {
		defer close(output)

		for element := range input {
			if filterFunc(element) {
				output <- element
			}
		}
	}()

	return output
}

// Reduce higher-order function implementation with channels.
func Reduce[I, A any](input <-chan I, reduceFunc func(accum A, element I) A) A {
	var accum A
	for element := range input {
		accum = reduceFunc(accum, element)
	}

	return accum
}

// Merge multiple channels to a single output channnel.
func Merge[T any](outputs ...<-chan T) chan T {
	var (
		wg     sync.WaitGroup
		merged = make(chan T)
	)

	wg.Add(len(outputs))

	for _, output := range outputs {
		go func(output <-chan T) {
			for v := range output {
				merged <- v
			}

			wg.Done()
		}(output)
	}

	go func() { wg.Wait(); close(merged) }()

	return merged
}
