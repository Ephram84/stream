package slice

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilterPairs(t *testing.T) {
	pairsMap := map[string]int{
		"key1":  1,
		"key2":  2,
		"key3":  3,
		"key11": 11,
	}

	result, err := FromMap(pairsMap).Filter(func(key string, value int) (bool, error) {
		return strings.HasSuffix(key, "1") && value > 10, nil
	}).ToMap()
	assert.NoError(t, err)
	expected := map[string]int{
		"key11": 11,
	}
	assert.Equal(t, expected, result)
}

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
