package pipe_test

import (
	"reflect"
	"testing"

	"github.com/tiny-go/fn/pipe"
)

func TestPipe(t *testing.T) {
	t.Run("map", func(t *testing.T) {
		var (
			input  = pipe.Extract(1, 2, 3, 4, 5, 6, 7, 8, 9)
			mapped = pipe.Map(input, func(i int) float64 {
				return float64(i * 2)
			})
			actual   = pipe.Collect(mapped)
			expected = []float64{2, 4, 6, 8, 10, 12, 14, 16, 18}
		)

		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("actual output %#v does not match expected %#v", actual, expected)
		}
	})

	t.Run("filter", func(t *testing.T) {
		var (
			input    = pipe.Extract(1, 2, 3, 4, 5, 6, 7, 8, 9)
			filtered = pipe.Filter(input, func(i int) bool {
				return i > 3 && i < 9
			})
			actual   = pipe.Collect(filtered)
			expected = []int{4, 5, 6, 7, 8}
		)

		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("actual output %#v does not match expected %#v", actual, expected)
		}
	})

	t.Run("reduce", func(t *testing.T) {
		var (
			input  = pipe.Extract(1, 2, 3, 4, 5, 6, 7, 8, 9)
			actual = pipe.Reduce(input, func(accum int, element int) int {
				return accum + element
			})
			expected = 45
		)

		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("actual output %#v does not match expected %#v", actual, expected)
		}
	})

	t.Run("pipe", func(t *testing.T) {
		var (
			in   = make(chan struct{}, 1)
			out1 = make(chan struct{})
			out2 = make(chan struct{})
		)

		go pipe.Pipe(in, out1, out2)

		in <- struct{}{}
		<-out1
		<-out2
	})
}
