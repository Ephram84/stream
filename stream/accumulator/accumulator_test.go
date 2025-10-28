package accumulator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAvgFloats(t *testing.T) {
	avg := Avg()

	result := avg.Apply(0.0, 1.0)
	assert.Equal(t, 1.0, result)
	result = avg.Apply(result, 2.0)
	assert.Equal(t, 1.5, result)
}
