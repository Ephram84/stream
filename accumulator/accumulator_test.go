package accumulator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSum(t *testing.T) {
	ints := []int{1, 2, 3, 4, 5}

	result, err := Sum(ints)
	assert.NoError(t, err)
	assert.Equal(t, 15, result)

	floats := []float64{1.5, 2.5, 3.0}

	resultF, err := Sum(floats)
	assert.NoError(t, err)
	assert.Equal(t, 7.0, resultF)
}

func TestAvg(t *testing.T) {
	floats := []float64{2.0, 4.3, 6.4, 8.1}
	resultF, err := Avg(floats)
	assert.NoError(t, err)
	assert.InDelta(t, 5.2, resultF, 0.1)
}

func TestMin(t *testing.T) {
	ints := []int{5, 3, 8, 1, 4}
	result, err := Min(ints)
	assert.NoError(t, err)
	assert.Equal(t, 1, result)

	floats := []float64{5.5, 3.3, 8.8, 1.1, 4.4}
	resultF, err := Min(floats)
	assert.NoError(t, err)
	assert.Equal(t, 1.1, resultF)
}

func TestMax(t *testing.T) {
	ints := []int{5, 3, 8, 1, 4}
	result, err := Max(ints)
	assert.NoError(t, err)
	assert.Equal(t, 8, result)

	floats := []float64{5.5, 3.3, 8.8, 1.1, 4.4}
	resultF, err := Max(floats)
	assert.NoError(t, err)
	assert.Equal(t, 8.8, resultF)
}
