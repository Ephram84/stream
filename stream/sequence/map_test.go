package sequence

import (
	"testing"
	"time"

	"github.com/Ephram84/stream/stream/accumulator"
	"github.com/stretchr/testify/assert"
)

func TestWords(t *testing.T) {
	words, err := FromFile("../../assets/words.txt").GroupByString(mapper).CountValues().ToMap()
	assert.NoError(t, err)

	assert.Equal(t, 110, words["a"])
	assert.Equal(t, 123, words["ac"])
	assert.Equal(t, 55, words["luctus"])
}

func TestFlatten(t *testing.T) {
	ints, err := FromMapWithSlices(map[string][]int{
		"first":  {0, 1, 2},
		"second": {3, 4, 5},
		"third":  {6},
	}).Flatten().ToSlice()
	assert.NoError(t, err)
	assert.ElementsMatch(t, []int{0, 1, 2, 3, 4, 5, 6}, ints)
}

func TestKeys(t *testing.T) {
	keys, err := FromMapWithSlices(map[string][]int{
		"first":  {0, 1, 2},
		"second": {3, 4, 5},
		"third":  {6},
	}).Keys().ToSlice()
	assert.NoError(t, err)
	assert.ElementsMatch(t, []string{"first", "second", "third"}, keys)
}

func TestReducingToFloat(t *testing.T) {
	transactionsMap := map[string][]Transaction{
		"2023-9": {
			{
				Amount:      10.0,
				BookingDate: time.Date(2023, time.September, 9, 0, 0, 0, 0, time.UTC).Unix(),
			},
			{
				Amount:      15.0,
				BookingDate: time.Date(2023, time.September, 10, 0, 0, 0, 0, time.UTC).Unix(),
			},
		},
		"2023-8": {
			{
				Amount:      150.0,
				BookingDate: time.Date(2023, time.August, 20, 0, 0, 0, 0, time.UTC).Unix(),
			},
		},
	}

	result, err := FromMapWithSlices(transactionsMap).MapToFloat64(func(elem Transaction) (float64, error) { return elem.Amount, nil }).Reduce(accumulator.SumSeq).ToMap()
	assert.NoError(t, err)
	assert.Equal(t, map[string]float64{
		"2023-9": 25.0,
		"2023-8": 150.0,
	}, result)
}

func TestReduce(t *testing.T) {
	m := map[string][]float64{
		"positive": {10.0, 15.0},
		"negative": {-5.0, -10.0},
	}
	sum, err := FromMapWithSlices(m).Reduce(accumulator.SumSeq).ToMap()
	assert.NoError(t, err)
	assert.Equal(t, map[string]float64{
		"positive": 25.0,
		"negative": -15.0,
	}, sum)

	avg, err := FromMapWithSlices(m).Reduce(accumulator.AvgSeq).ToMap()
	assert.NoError(t, err)
	assert.Equal(t, map[string]float64{
		"positive": 12.5,
		"negative": -7.5,
	}, avg)

	min, err := FromMapWithSlices(m).Reduce(accumulator.MinSeq).ToMap()
	assert.NoError(t, err)
	assert.Equal(t, map[string]float64{
		"positive": 10.0,
		"negative": -10.0,
	}, min)

	max, err := FromMapWithSlices(m).Reduce(accumulator.MaxSeq).ToMap()
	assert.NoError(t, err)
	assert.Equal(t, map[string]float64{
		"positive": 15.0,
		"negative": -5.0,
	}, max)
}

func TestCountWithMaps(t *testing.T) {
	m := map[string][]int{
		"a": {1, 2, 3},
		"b": {4, 5},
		"c": {},
	}
	countValues, err := FromMapWithSlices(m).CountValues().ToMap()
	assert.NoError(t, err)
	assert.Equal(t, map[string]int{
		"a": 3,
		"b": 2,
		"c": 0,
	}, countValues)

	count, err := FromMapWithSlices(m).Count()
	assert.NoError(t, err)
	assert.Equal(t, 3, count)
}
