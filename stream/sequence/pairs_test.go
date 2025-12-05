package sequence

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFlattenPairs(t *testing.T) {
	pairsMap := map[string]int{
		"a": 1,
		"b": 2,
		"c": 3,
	}

	result, err := FromMap(pairsMap).Flatten().ToSlice()
	assert.NoError(t, err)
	assert.ElementsMatch(t, []int{1, 2, 3}, result)
}

func TestKeysPairs(t *testing.T) {
	pairsMap := map[string]int{
		"a": 1,
		"b": 2,
		"c": 3,
	}

	result, err := FromMap(pairsMap).Keys().ToSlice()
	assert.NoError(t, err)
	assert.ElementsMatch(t, []string{"a", "b", "c"}, result)
}

func TestCountPairs(t *testing.T) {
	pairsMap := map[string]int{
		"a": 1,
		"b": 2,
		"c": 3,
	}

	result, err := FromMap(pairsMap).Count()
	assert.NoError(t, err)
	assert.Equal(t, 3, result)
}
