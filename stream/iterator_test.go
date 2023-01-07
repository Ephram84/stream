package stream

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIterator(t *testing.T) {
	numbers, err := Iterator(0, increment).WithSkip(3).WithLimit(5).ToStream().ToArray()
	assert.NoError(t, err)
	assert.Equal(t, []int{3, 4, 5, 6, 7}, numbers)

	numbers, err = Iterator(0, increment).WithLimit(5).ToStream().ToArray()
	assert.NoError(t, err)
	assert.Equal(t, []int{0, 1, 2, 3, 4}, numbers)
}

func increment(prev int) int {
	return prev + 1
}
