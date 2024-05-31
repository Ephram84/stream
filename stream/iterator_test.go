package stream

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIteratorInt(t *testing.T) {
	numbers, err := Iterator(0, incrementInt).WithSkip(3).WithLimit(5).Generate().ToSlice()
	assert.NoError(t, err)
	assert.Equal(t, []int{3, 4, 5, 6, 7}, numbers)

	numbers, err = Iterator(0, incrementInt).WithLimit(5).Generate().ToSlice()
	assert.NoError(t, err)
	assert.Equal(t, []int{0, 1, 2, 3, 4}, numbers)
}

func incrementInt(prev int) int {
	return prev + 1
}

func TestIteratorFloat(t *testing.T) {
	numbers, err := Iterator(0.0, incrementFloat).WithSkip(3).WithLimit(5).Generate().ToSlice()
	assert.NoError(t, err)
	assert.Equal(t, []float64{3.0, 4.0, 5.0, 6.0, 7.0}, numbers)
}

func incrementFloat(prev float64) float64 {
	return prev + 1.0
}
