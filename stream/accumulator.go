package stream

type Accumulator[T any] interface {
	Apply(partialResult, elem T) T
}

type sum[N Numbers] struct{}

func Sum[N Numbers]() *sum[N] {
	return &sum[N]{}
}

func (s *sum[N]) Apply(partialResult, elem N) N {
	return partialResult + elem
}

type avg[N Numbers] struct {
	count N
}

func Avg() *avg[float64] {
	return &avg[float64]{
		count: 0,
	}
}

func (a *avg[N]) Apply(partialResult, elem N) N {
	number := a.increment()
	return (partialResult*(number-1) + elem) / number
}

func (a *avg[N]) increment() N {
	a.count = a.count + 1
	return a.count
}

type max[T any] struct {
	comparator func(max, elem T) bool
}

func Max[T any](comparator func(max, elem T) bool) *max[T] {
	return &max[T]{
		comparator: comparator,
	}
}

func (m *max[T]) Apply(partialResult, elem T) T {
	if m.comparator(partialResult, elem) {
		return elem
	}

	return partialResult
}

type min[T any] struct {
	comparator func(max, elem T) bool
}

func Min[T any](comparator func(max, elem T) bool) *min[T] {
	return &min[T]{
		comparator: comparator,
	}
}

func (m *min[T]) Apply(partialResult, elem T) T {
	if m.comparator(partialResult, elem) {
		return elem
	}

	return partialResult
}
