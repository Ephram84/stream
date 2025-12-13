package iterator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIteratorInt(t *testing.T) {
	numbers := Iterator(0, IncrementInt()).WithSkip(3).WithLimit(5).Generate()
	assert.Equal(t, []int{3, 4, 5, 6, 7}, numbers)

	numbers = Iterator(0, IncrementInt()).WithLimit(5).Generate()
	assert.Equal(t, []int{0, 1, 2, 3, 4}, numbers)
}

func TestIteratorFloat(t *testing.T) {
	numbers := Iterator(0.0, IncrementFloat()).WithSkip(3).WithLimit(5).Generate()
	assert.Equal(t, []float64{3.0, 4.0, 5.0, 6.0, 7.0}, numbers)

	numbers = Iterator(0.0, IncrementFloat()).WithLimit(5).Generate()
	assert.Equal(t, []float64{0.0, 1.0, 2.0, 3.0, 4.0}, numbers)
}

func TestIteratorWithCustomNexFunc(t *testing.T) {
	ints := Iterator(0, func(current int) int { return current + 2 }).WithLimit(5).Generate()
	assert.Equal(t, []int{0, 2, 4, 6, 8}, ints)

	floats := Iterator(1.0, func(current float64) float64 { return current * 0.5 }).WithLimit(5).Generate()
	assert.Equal(t, []float64{1.0, 0.5, 0.25, 0.125, 0.0625}, floats)
}

func TestFibonacciIterator(t *testing.T) {
	prev := 0
	fib := Iterator(1, func(current int) int {
		next := prev + current
		prev = current
		return next
	}).WithLimit(10).Generate()
	assert.Equal(t, []int{1, 1, 2, 3, 5, 8, 13, 21, 34, 55}, fib)
}

func TestCountdown(t *testing.T) {
	countdown := Iterator(10, func(current int) int {
		return current - 1
	}).WithLimit(11).Generate()
	assert.Equal(t, []int{10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0}, countdown)
}
