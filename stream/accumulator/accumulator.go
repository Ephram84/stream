package accumulator

import "github.com/Ephram84/stream/stream/common"

type Accumulator[T any] func(values []T) (T, error)

func Sum[T common.Numbers](values []T) (T, error) {
	var sum T
	for _, v := range values {
		sum = sum + v
	}
	return sum, nil
}

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

type AccumulatorSeq[T any] func(index int, prev T, current T) (T, error)

func SumSeq[T common.Numbers](index int, prev T, current T) (T, error) {
	return prev + current, nil
}

func AvgSeq(index int, prev float64, current float64) (float64, error) {
	return (prev*float64(index-1) + current) / float64(index), nil
}

func MinSeq[T common.Numbers](index int, prev T, current T) (T, error) {
	return common.Min(prev, current), nil
}

func MaxSeq[T common.Numbers](index int, prev T, current T) (T, error) {
	return common.Max(prev, current), nil
}
