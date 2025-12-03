package accumulator

import "github.com/Ephram84/stream/stream/common"

// Accumulator is a function type that takes a slice of values and reduces them to a single value.
// It returns the accumulated result and an error if the operation fails.
type Accumulator[T any] func(values []T) (T, error)

// Sum calculates the sum of all numeric values in the slice.
// It returns the total sum and nil error.
func Sum[T common.Numbers](values []T) (T, error) {
	var sum T
	for _, v := range values {
		sum = sum + v
	}
	return sum, nil
}

// Avg calculates the average (mean) of all float64 values in the slice.
// It returns 0 if the slice is empty, otherwise returns the average and nil error.
func Avg(values []float64) (float64, error) {
	var sum float64
	if len(values) == 0 {
		return sum, nil
	}
	for _, v := range values {
		sum = sum + v
	}

	return sum / float64(len(values)), nil
}

// Min finds the minimum value in the slice of numeric values.
// It returns the zero value if the slice is empty, otherwise returns the minimum and nil error.
func Min[T common.Numbers](values []T) (T, error) {
	if len(values) == 0 {
		var zero T
		return zero, nil
	}

	min := values[0]
	for _, v := range values[1:] {
		min = common.Min(min, v)
	}
	return min, nil
}

// Max finds the maximum value in the slice of numeric values.
// It returns the zero value if the slice is empty, otherwise returns the maximum and nil error.
func Max[T common.Numbers](values []T) (T, error) {
	if len(values) == 0 {
		var zero T
		return zero, nil
	}
	max := values[0]
	for _, v := range values[1:] {
		max = common.Max(max, v)
	}
	return max, nil
}

// AccumulatorSeq is a function type for sequential accumulation operations.
// It takes the current index, the previous accumulated value, and the current value,
// returning the new accumulated value and an error if the operation fails.
type AccumulatorSeq[T any] func(index int, prev T, current T) (T, error)

// SumSeq is a sequential accumulator that adds the current value to the previous sum.
func SumSeq[T common.Numbers](index int, prev T, current T) (T, error) {
	return prev + current, nil
}

// AvgSeq is a sequential accumulator that calculates the running average.
// It uses the index to weight the previous average and incorporates the current value.
func AvgSeq(index int, prev float64, current float64) (float64, error) {
	return (prev*float64(index-1) + current) / float64(index), nil
}

// MinSeq is a sequential accumulator that returns the minimum of the previous and current values.
func MinSeq[T common.Numbers](index int, prev T, current T) (T, error) {
	return common.Min(prev, current), nil
}

// MaxSeq is a sequential accumulator that returns the maximum of the previous and current values.
func MaxSeq[T common.Numbers](index int, prev T, current T) (T, error) {
	return common.Max(prev, current), nil
}
