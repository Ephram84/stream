package sequence

import (
	"testing"

	"github.com/Ephram84/stream/stream/accumulator"
	"github.com/Ephram84/stream/stream/common"
	"github.com/Ephram84/stream/stream/tupel"
	"github.com/stretchr/testify/assert"
)

// help functions, structures and variables for tests
func isEven(elem int) (bool, error) {
	return elem%2 == 0, nil
}

func mapper(s string) (string, error) {
	return s, nil
}

type Account struct {
	ID           string
	Transactions []Transaction
}

type Transaction struct {
	ID          string
	Amount      float64
	BookingDate int64
	Tags        []string
}

func TestError(t *testing.T) {
	numbers := []string{"1", "2", "a", "4"}
	result, err := From(numbers).MapToInt(common.StringToInt).Filter(isEven).ToSlice()
	assert.Error(t, err)
	assert.Empty(t, result)
}

func TestFilterEven(t *testing.T) {
	result, err := From([]int{1, 2, 3, 4, 5}).
		Filter(isEven).ToSlice()
	assert.NoError(t, err)
	assert.Equal(t, []int{2, 4}, result)
}

func TestTake(t *testing.T) {
	result, err := From([]int{1, 2, 3, 4, 5}).
		Take(3).
		ToSlice()
	assert.NoError(t, err)
	assert.Equal(t, []int{1, 2, 3}, result)
}

func TestSkip(t *testing.T) {
	result, err := From([]int{1, 2, 3, 4, 5}).
		Skip(2).
		ToSlice()
	assert.NoError(t, err)
	assert.Equal(t, []int{3, 4, 5}, result)
}

func TestReverse(t *testing.T) {
	result, err := From([]int{1, 2, 3, 4, 5}).
		Reverse().
		ToSlice()
	assert.NoError(t, err)
	assert.Equal(t, []int{5, 4, 3, 2, 1}, result)
}

func TestDistinct(t *testing.T) {
	result, err := From([]int{1, 2, 2, 3, 3, 3, 4, 4, 4, 4}).
		Distinct().
		ToSlice()
	assert.NoError(t, err)
	assert.Equal(t, []int{1, 2, 3, 4}, result)
}

func TestCount(t *testing.T) {
	count, err := From([]int{1, 2, 3, 4, 5}).Count()
	assert.NoError(t, err)
	assert.Equal(t, 5, count)
}

func TestFirst(t *testing.T) {
	result, err := From([]int{1, 2, 3, 4, 5}).First()
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, *result)
}

func TestFirstReturnsNil(t *testing.T) {
	result, err := From([]int{1, 2, 3, 4, 5}).Filter(func(elem int) (bool, error) {
		return elem > 5, nil
	}).First()
	assert.NoError(t, err)
	assert.Nil(t, result)
}

func TestFirstReturnsElse(t *testing.T) {
	result, err := From([]int{}).First(42)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 42, *result)
}

func TestAnyMatch(t *testing.T) {
	result, err := From([]int{1, 2, 3, 4, 5}).AnyMatch(func(elem int) (bool, error) {
		return elem == 3, nil
	})
	assert.NoError(t, err)
	assert.True(t, result)
}

func TestAnyMatchReturnsFalse(t *testing.T) {
	result, err := From([]int{1, 2, 3, 4, 5}).AnyMatch(func(elem int) (bool, error) {
		return elem == 42, nil
	})
	assert.NoError(t, err)
	assert.False(t, result)
}

func TestAllMatch(t *testing.T) {
	result, err := From([]int{2, 4, 6, 8}).AllMatch(func(elem int) (bool, error) {
		return elem%2 == 0, nil
	})
	assert.NoError(t, err)
	assert.True(t, result)
}

func TestAllMatchReturnsFalse(t *testing.T) {
	result, err := From([]int{2, 4, 6, 7, 8}).AllMatch(func(elem int) (bool, error) {
		return elem%2 == 0, nil
	})
	assert.NoError(t, err)
	assert.False(t, result)
}

func TestNoneMatch(t *testing.T) {
	result, err := From([]int{1, 3, 5, 7}).NoneMatch(func(elem int) (bool, error) {
		return elem%2 == 0, nil
	})
	assert.NoError(t, err)
	assert.True(t, result)
}

func TestNoneMatchReturnsFalse(t *testing.T) {
	result, err := From([]int{1, 2, 3, 5, 7}).NoneMatch(func(elem int) (bool, error) {
		return elem%2 == 0, nil
	})
	assert.NoError(t, err)
	assert.False(t, result)
}

func TestReduceSum(t *testing.T) {
	result, err := From([]int{1, 2, 3, 4, 5}).Reduce(accumulator.SumSeq)
	assert.NoError(t, err)
	assert.Equal(t, 15, result)
}

func TestReduceAvg(t *testing.T) {
	result, err := From([]float64{2.1, 4.2, 6.3, 8.4}).Reduce(accumulator.AvgSeq)
	assert.NoError(t, err)
	assert.Equal(t, 5.25, result)
}

func TestReduceMin(t *testing.T) {
	result, err := From([]int{5, 3, 8, 1, 4}).Reduce(accumulator.MinSeq)
	assert.NoError(t, err)
	assert.Equal(t, 1, result)
}

func TestReduceMax(t *testing.T) {
	result, err := From([]int{5, 3, 8, 1, 4}).Reduce(accumulator.MaxSeq)
	assert.NoError(t, err)
	assert.Equal(t, 8, result)
}

func TestConcat(t *testing.T) {
	result, err := From([]int{4, 5}).Concat(From([]int{6, 7, 8})).ToSlice()
	assert.NoError(t, err)
	assert.Equal(t, []int{4, 5, 6, 7, 8}, result)
}

func TestMap(t *testing.T) {
	numbers := From([]int{1, 2, 3, 4, 5})
	result, err := Map(numbers, func(elem int) (string, error) {
		return "Number: " + string(rune('0'+elem)), nil
	}).ToSlice()
	assert.NoError(t, err)
	assert.Equal(t, []string{"Number: 1", "Number: 2", "Number: 3", "Number: 4", "Number: 5"}, result)
}

func TestMap2(t *testing.T) {
	numbers := From([]int{1, 2, 3, 4, 5})
	result, err := numbers.Map(func(elem int) (int, error) {
		return elem * 2, nil
	}).ToSlice()
	assert.NoError(t, err)
	assert.Equal(t, []int{2, 4, 6, 8, 10}, result)
}

func TestGroupBy(t *testing.T) {
	numbers := From([]int{1, 2, 3, 4, 5, 6})
	result := GroupBy(numbers, func(elem int) (string, error) {
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

func TestFlatMapWithInts(t *testing.T) {
	intSlice := [][]int{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	}

	result, err := FlatMap(From(intSlice), func(elem []int) ([]int, error) {
		return elem, nil
	}).ToSlice()
	assert.NoError(t, err)
	assert.Equal(t, []int{1, 2, 3, 4, 5, 6, 7, 8, 9}, result)
}

func TestFlatMap(t *testing.T) {
	accounts := []Account{
		{
			Transactions: []Transaction{
				{
					ID:     "1",
					Amount: 120.0,
				},
				{
					ID:     "2",
					Amount: -120.0,
				},
			},
		},
		{
			Transactions: []Transaction{
				{
					ID:     "3",
					Amount: 20.0,
				},
				{
					ID:     "4",
					Amount: 10.0,
				},
			},
		},
	}

	transactions, err := FlatMap(From(accounts), func(elem Account) ([]Transaction, error) {
		return elem.Transactions, nil
	}).Filter(func(elem Transaction) (bool, error) {
		return elem.Amount > 0.0, nil
	}).ToSlice()
	assert.NoError(t, err)
	assert.Equal(t, []Transaction{
		{ID: "1", Amount: 120.0},
		{ID: "3", Amount: 20.0},
		{ID: "4", Amount: 10.0},
	}, transactions)
}

func TestZip(t *testing.T) {
	result, err := Zip(From([]int{1, 2, 3}), From([]string{"a", "b", "c", "d"})).ToSlice()
	assert.NoError(t, err)
	assert.Equal(t, []tupel.Tupel[int, string]{
		{First: 1, Second: "a"},
		{First: 2, Second: "b"},
		{First: 3, Second: "c"},
	}, result)
}
