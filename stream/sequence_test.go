package stream

import (
	"testing"

	"github.com/Ephram84/stream/stream/accumulator"
	"github.com/Ephram84/stream/stream/common"
	"github.com/Ephram84/stream/stream/tupel"
	"github.com/stretchr/testify/assert"
)

func TestErrorSequence(t *testing.T) {
	numbers := []string{"1", "2", "a", "4"}
	result, err := LazyFrom(numbers).MapToInt(common.StringToInt).Filter(isEven).ToSlice()
	assert.Error(t, err)
	assert.Empty(t, result)
}

func TestFilterEvenSequence(t *testing.T) {
	result, err := LazyFrom([]int{1, 2, 3, 4, 5}).
		Filter(isEven).ToSlice()
	assert.NoError(t, err)
	assert.Equal(t, []int{2, 4}, result)
}

func TestTakeSequence(t *testing.T) {
	result, err := LazyFrom([]int{1, 2, 3, 4, 5}).
		Take(3).
		ToSlice()
	assert.NoError(t, err)
	assert.Equal(t, []int{1, 2, 3}, result)
}

func TestSkipSequence(t *testing.T) {
	result, err := LazyFrom([]int{1, 2, 3, 4, 5}).
		Skip(2).
		ToSlice()
	assert.NoError(t, err)
	assert.Equal(t, []int{3, 4, 5}, result)
}

func TestReverseSequence(t *testing.T) {
	result, err := LazyFrom([]int{1, 2, 3, 4, 5}).
		Reverse().
		ToSlice()
	assert.NoError(t, err)
	assert.Equal(t, []int{5, 4, 3, 2, 1}, result)
}

func TestDistinctSequence(t *testing.T) {
	result, err := LazyFrom([]int{1, 2, 2, 3, 3, 3, 4, 4, 4, 4}).
		Distinct().
		ToSlice()
	assert.NoError(t, err)
	assert.Equal(t, []int{1, 2, 3, 4}, result)
}

func TestCountSequence(t *testing.T) {
	count, err := LazyFrom([]int{1, 2, 3, 4, 5}).Count()
	assert.NoError(t, err)
	assert.Equal(t, 5, count)
}

func TestFirstSequence(t *testing.T) {
	result, err := LazyFrom([]int{1, 2, 3, 4, 5}).First()
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, *result)
}

func TestFirstReturnsNilSequence(t *testing.T) {
	result, err := LazyFrom([]int{1, 2, 3, 4, 5}).Filter(func(elem int) (bool, error) {
		return elem > 5, nil
	}).First()
	assert.NoError(t, err)
	assert.Nil(t, result)
}

func TestFirstReturnsElseSequence(t *testing.T) {
	result, err := LazyFrom([]int{}).First(42)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 42, *result)
}

func TestAnyMatchSequence(t *testing.T) {
	result, err := LazyFrom([]int{1, 2, 3, 4, 5}).AnyMatch(func(elem int) (bool, error) {
		return elem == 3, nil
	})
	assert.NoError(t, err)
	assert.True(t, result)
}

func TestAnyMatchReturnsFalseSequence(t *testing.T) {
	result, err := LazyFrom([]int{1, 2, 3, 4, 5}).AnyMatch(func(elem int) (bool, error) {
		return elem == 42, nil
	})
	assert.NoError(t, err)
	assert.False(t, result)
}

func TestAllMatchSequence(t *testing.T) {
	result, err := LazyFrom([]int{2, 4, 6, 8}).AllMatch(func(elem int) (bool, error) {
		return elem%2 == 0, nil
	})
	assert.NoError(t, err)
	assert.True(t, result)
}

func TestAllMatchReturnsFalseSequence(t *testing.T) {
	result, err := LazyFrom([]int{2, 4, 6, 7, 8}).AllMatch(func(elem int) (bool, error) {
		return elem%2 == 0, nil
	})
	assert.NoError(t, err)
	assert.False(t, result)
}

func TestNoneMatchSequence(t *testing.T) {
	result, err := LazyFrom([]int{1, 3, 5, 7}).NoneMatch(func(elem int) (bool, error) {
		return elem%2 == 0, nil
	})
	assert.NoError(t, err)
	assert.True(t, result)
}

func TestNoneMatchReturnsFalseSequence(t *testing.T) {
	result, err := LazyFrom([]int{1, 2, 3, 5, 7}).NoneMatch(func(elem int) (bool, error) {
		return elem%2 == 0, nil
	})
	assert.NoError(t, err)
	assert.False(t, result)
}

func TestReduceSumSequence(t *testing.T) {
	result, err := LazyFrom([]int{1, 2, 3, 4, 5}).Reduce(0, accumulator.Sum[int]())
	assert.NoError(t, err)
	assert.Equal(t, 15, result)
}

func TestConcatSequences(t *testing.T) {
	result, err := LazyFrom([]int{4, 5}).Concat(LazyFrom([]int{6, 7, 8})).ToSlice()
	assert.NoError(t, err)
	assert.Equal(t, []int{4, 5, 6, 7, 8}, result)
}

func TestMapSequence(t *testing.T) {
	numbers := LazyFrom([]int{1, 2, 3, 4, 5})
	result, err := MapSeq(numbers, func(elem int) (string, error) {
		return "Number: " + string(rune('0'+elem)), nil
	}).ToSlice()
	assert.NoError(t, err)
	assert.Equal(t, []string{"Number: 1", "Number: 2", "Number: 3", "Number: 4", "Number: 5"}, result)
}

func TestMapSequence2(t *testing.T) {
	numbers := LazyFrom([]int{1, 2, 3, 4, 5})
	result, err := numbers.Map(func(elem int) (int, error) {
		return elem * 2, nil
	}).ToSlice()
	assert.NoError(t, err)
	assert.Equal(t, []int{2, 4, 6, 8, 10}, result)
}

func TestGroupBySequence(t *testing.T) {
	numbers := LazyFrom([]int{1, 2, 3, 4, 5, 6})
	result := GroupBySeq(numbers, func(elem int) (string, error) {
		if elem%2 == 0 {
			return "even", nil
		}
		return "odd", nil
	})
	assert.NotNil(t, result)
	m, err := result.ToMap()
	assert.NoError(t, err)
	assert.Equal(t, []int{1, 3, 5}, m["odd"])
	assert.Equal(t, []int{2, 4, 6}, m["even"])
}

func TestFlatMapSequence(t *testing.T) {
	words := LazyFrom([]string{"hi", "go"})

	result, err := FlatMapSeq(words, func(elem string) (*seq[rune], error) {
		return LazyFrom([]rune(elem)), nil
	}).ToSlice()
	assert.NoError(t, err)
	assert.Equal(t, []rune{'h', 'i', 'g', 'o'}, result)
}

func TestFlatMapSeqWithInts(t *testing.T) {
	intSlice := [][]int{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	}

	result, err := FlatMapSeq(LazyFrom(intSlice), func(elem []int) (*seq[int], error) {
		return LazyFrom(elem), nil
	}).ToSlice()
	assert.NoError(t, err)
	assert.Equal(t, []int{1, 2, 3, 4, 5, 6, 7, 8, 9}, result)
}

func TestZipSequences(t *testing.T) {
	result, err := ZipSeqs(LazyFrom([]int{1, 2, 3}), LazyFrom([]string{"a", "b", "c", "d"})).ToSlice()
	assert.NoError(t, err)
	assert.Equal(t, []tupel.Tupel[int, string]{
		{First: 1, Second: "a"},
		{First: 2, Second: "b"},
		{First: 3, Second: "c"},
	}, result)
}
