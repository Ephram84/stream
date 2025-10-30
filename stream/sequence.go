package stream

import (
	"github.com/Ephram84/stream/stream/accumulator"
	"github.com/Ephram84/stream/stream/common"
	"github.com/Ephram84/stream/stream/tupel"
)

type seq[T any] struct {
	next func() (T, bool, error)
}

func LazyFrom[T any](slice []T, errs ...error) *seq[T] {
	var i int
	return &seq[T]{
		next: func() (T, bool, error) {
			if len(errs) > 0 && errs[0] != nil {
				var zero T
				return zero, false, errs[0]
			}
			if i < len(slice) {
				val := slice[i]
				i++
				return val, true, nil
			}
			var zero T
			return zero, false, nil
		},
	}
}

func (s *seq[T]) Filter(filter func(elem T) (bool, error)) *seq[T] {
	return &seq[T]{
		next: func() (T, bool, error) {
			for {
				var zero T
				val, ok, err := s.next()
				if err != nil {
					return zero, false, err
				}
				if !ok {
					return zero, false, nil
				}
				keep, err := filter(val)
				if err != nil {
					return zero, false, err
				}
				if keep {
					return val, true, nil
				}
			}
		},
	}
}

func (s *seq[T]) Map(mapper func(elem T) (T, error)) *seq[T] {
	return &seq[T]{
		next: func() (T, bool, error) {
			var zero T
			val, ok, err := s.next()
			if err != nil {
				return zero, false, err
			}
			if !ok {
				return zero, false, nil
			}

			newVal, err := mapper(val)
			if err != nil {
				return zero, false, err
			}

			return newVal, true, nil
		},
	}
}

func (s *seq[T]) MapToInt(mapper func(elem T) (int, error)) *seq[int] {
	return MapSeq(s, mapper)
}

func (s *seq[T]) MapToInt64(mapper func(elem T) (int64, error)) *seq[int64] {
	return MapSeq(s, mapper)
}

func (s *seq[T]) MapToFloat64(mapper func(elem T) (float64, error)) *seq[float64] {
	return MapSeq(s, mapper)
}

func (s *seq[T]) MapToString(mapper func(elem T) (string, error)) *seq[string] {
	return MapSeq(s, mapper)
}

func (s *seq[T]) PartitioningBy(predicate func(elem T) (bool, error)) *pairsSlice[bool, T] {
	result := emptyMapWithSlices[bool, T]()
	for {
		val, ok, err := s.next()
		if err != nil {
			result.err = err
			return result
		}
		if !ok {
			break
		}
		key, err := predicate(val)
		if err != nil {
			result.err = err
			return result
		}
		result.m[key] = append(result.m[key], val)
	}
	return result
}

func (s *seq[T]) GroupByString(grouper func(elem T) (string, error)) *pairsSlice[string, T] {
	return GroupBySeq(s, grouper)
}

func (s *seq[T]) GroupByInt(grouper func(elem T) (int, error)) *pairsSlice[int, T] {
	return GroupBySeq(s, grouper)
}

func (s *seq[T]) GroupByInt64(grouper func(elem T) (int64, error)) *pairsSlice[int64, T] {
	return GroupBySeq(s, grouper)
}

func (s *seq[T]) GroupByFloat(grouper func(elem T) (float64, error)) *pairsSlice[float64, T] {
	return GroupBySeq(s, grouper)
}

func (s *seq[T]) AssociateByString(mapper func(elem T) (string, error)) *pairs[string, T] {
	result := emptyPairs[string, T]()
	for {
		val, ok, err := s.next()
		if err != nil {
			result.err = err
			return result
		}
		if !ok {
			break
		}
		key, err := mapper(val)
		if err != nil {
			result.err = err
			return result
		}
		result.m[key] = val
	}
	return result
}

func (s *seq[T]) AssociateByInt(mapper func(elem T) (int, error)) *pairs[int, T] {
	result := emptyPairs[int, T]()
	for {
		val, ok, err := s.next()
		if err != nil {
			result.err = err
			return result
		}
		if !ok {
			break
		}
		key, err := mapper(val)
		if err != nil {
			result.err = err
			return result
		}
		result.m[key] = val
	}
	return result
}

func (s *seq[T]) AssociateByInt64(mapper func(elem T) (int64, error)) *pairs[int64, T] {
	result := emptyPairs[int64, T]()
	for {
		val, ok, err := s.next()
		if err != nil {
			result.err = err
			return result
		}
		if !ok {
			break
		}
		key, err := mapper(val)
		if err != nil {
			result.err = err
			return result
		}
		result.m[key] = val
	}
	return result
}

func (s *seq[T]) AssociateByFloat(mapper func(elem T) (float64, error)) *pairs[float64, T] {
	result := emptyPairs[float64, T]()
	for {
		val, ok, err := s.next()
		if err != nil {
			result.err = err
			return result
		}
		if !ok {
			break
		}
		key, err := mapper(val)
		if err != nil {
			result.err = err
			return result
		}
		result.m[key] = val
	}
	return result
}

func (s *seq[T]) Take(n int) *seq[T] {
	if n <= 0 {
		return &seq[T]{
			next: func() (T, bool, error) {
				var zero T
				return zero, false, nil
			},
		}
	}

	count := 0
	return &seq[T]{
		next: func() (T, bool, error) {
			if count >= n {
				var zero T
				return zero, false, nil
			}
			val, ok, err := s.next()
			if err != nil {
				var zero T
				return zero, false, err
			}
			if !ok {
				var zero T
				return zero, false, nil
			}
			count++
			return val, true, nil
		},
	}
}

func (s *seq[T]) Skip(n int) *seq[T] {
	if n <= 0 {
		return s
	}

	skipped := 0
	return &seq[T]{
		next: func() (T, bool, error) {
			for {
				val, ok, err := s.next()
				if err != nil {
					var zero T
					return zero, false, err
				}
				if !ok {
					var zero T
					return zero, false, nil
				}
				if skipped < n {
					skipped++
					continue
				}
				return val, true, nil
			}
		},
	}
}

func (s *seq[T]) Reverse() *seq[T] {
	values, err := s.ToSlice()
	if err != nil {
		return &seq[T]{
			next: func() (T, bool, error) {
				var zero T
				return zero, false, err
			},
		}
	}

	index := len(values)

	return &seq[T]{
		next: func() (T, bool, error) {
			if index == 0 {
				var zero T
				return zero, false, nil
			}
			index--
			return values[index], true, nil
		},
	}
}

func (s *seq[T]) Distinct() *seq[T] {
	seen := make(map[any]struct{})
	return &seq[T]{
		next: func() (T, bool, error) {
			for {
				val, ok, err := s.next()
				if err != nil {
					var zero T
					return zero, false, err
				}
				if !ok {
					var zero T
					return zero, false, nil
				}
				if _, exists := seen[val]; !exists {
					seen[val] = struct{}{}
					return val, true, nil
				}
			}
		},
	}
}

func (s *seq[T]) Concat(sequences ...*seq[T]) *seq[T] {
	return &seq[T]{
		next: func() (T, bool, error) {
			var zero T

			val, ok, err := s.next()
			if err != nil {
				return zero, false, err
			}
			if ok {
				return val, true, nil
			}

			for _, seq := range sequences {
				val, ok, err := seq.next()
				if err != nil {
					return zero, false, err
				}

				if ok {
					return val, true, nil
				}

			}

			return zero, false, nil
		},
	}
}

// terminal functions

func (s *seq[T]) Count() (int, error) {
	count := 0
	for {
		_, ok, err := s.next()
		if err != nil {
			return 0, err
		}
		if !ok {
			break
		}
		count++
	}
	return count, nil
}

func (s *seq[T]) ToSlice() ([]T, error) {
	var slice []T
	for {
		val, ok, err := s.next()
		if err != nil {
			return nil, err
		}
		if !ok {
			break
		}
		slice = append(slice, val)
	}
	return slice, nil
}

func (s *seq[T]) First(orElse ...T) (*T, error) {
	val, ok, err := s.next()
	if err != nil {
		return nil, err
	}
	if !ok {
		if len(orElse) > 0 {
			return &orElse[0], nil
		}
		return nil, nil
	}
	return &val, nil
}

func (s *seq[T]) FirstOrNil(predicate func(elem T) (bool, error)) (*T, error) {
	for {
		val, ok, err := s.next()
		if err != nil {
			return nil, err
		}
		if !ok {
			return nil, nil
		}
		match, err := predicate(val)
		if err != nil {
			return nil, err
		}
		if match {
			return &val, nil
		}
	}
}

func (s *seq[T]) AnyMatch(predicate func(elem T) (bool, error)) (bool, error) {
	for {
		val, ok, err := s.next()
		if err != nil {
			return false, err
		}
		if !ok {
			break
		}
		match, err := predicate(val)
		if err != nil {
			return false, err
		}
		if match {
			return true, nil
		}
	}
	return false, nil
}

func (s *seq[T]) AllMatch(predicate func(elem T) (bool, error)) (bool, error) {
	for {
		val, ok, err := s.next()
		if err != nil {
			return false, err
		}
		if !ok {
			break
		}
		match, err := predicate(val)
		if err != nil {
			return false, err
		}
		if !match {
			return false, nil
		}
	}
	return true, nil
}

func (s *seq[T]) NoneMatch(predicate func(elem T) (bool, error)) (bool, error) {
	for {
		val, ok, err := s.next()
		if err != nil {
			return false, err
		}
		if !ok {
			break
		}
		match, err := predicate(val)
		if err != nil {
			return false, err
		}
		if match {
			return false, nil
		}
	}
	return true, nil
}

func (s *seq[T]) Reduce(identity T, acc accumulator.Accumulator[T]) (T, error) {
	for {
		val, ok, err := s.next()
		if err != nil {
			return identity, err
		}
		if !ok {
			break
		}
		newValue := acc.Apply(identity, val)
		identity = newValue
	}
	return identity, nil
}

func (s *seq[T]) ForEach(consumer func(elem T) error) error {
	for {
		val, ok, err := s.next()
		if err != nil {
			return err
		}
		if !ok {
			break
		}
		if err := consumer(val); err != nil {
			return err
		}
	}

	return nil
}

// free functions

func MapSeq[T, R any](s *seq[T], mapper func(elem T) (R, error)) *seq[R] {
	return &seq[R]{
		next: func() (R, bool, error) {
			var zero R
			if s == nil {
				return zero, false, nil
			}
			val, ok, err := s.next()
			if err != nil {
				return zero, false, err
			}
			if !ok {
				return zero, false, nil
			}
			newValue, err := mapper(val)
			if err != nil {
				return zero, false, err
			}
			return newValue, true, nil
		},
	}
}

func FlatMapSeq[T, R any](s *seq[T], mapper func(elem T) (*seq[R], error)) *seq[R] {
	var currentSeq *seq[R]
	var hasCurrent bool
	return &seq[R]{
		next: func() (R, bool, error) {
			if s == nil {
				var zero R
				return zero, false, nil
			}
			for {
				if hasCurrent {
					val, ok, _ := currentSeq.next()
					if ok {
						return val, true, nil
					}
					hasCurrent = false
				}
				val, ok, err := s.next()
				if err != nil {
					var zero R
					return zero, false, err
				}
				if !ok {
					var zero R
					return zero, false, nil
				}
				currentSeq, err = mapper(val)
				if err != nil {
					var zero R
					return zero, false, err
				}
				hasCurrent = true
			}
		},
	}
}

func ZipSeqs[T1, T2 any](seq1 *seq[T1], seq2 *seq[T2]) *seq[tupel.Tupel[T1, T2]] {
	return &seq[tupel.Tupel[T1, T2]]{
		next: func() (tupel.Tupel[T1, T2], bool, error) {
			if seq1 == nil || seq2 == nil {
				var zero tupel.Tupel[T1, T2]
				return zero, false, nil
			}
			val1, ok1, err1 := seq1.next()
			val2, ok2, err2 := seq2.next()
			if err1 != nil {
				return tupel.Tupel[T1, T2]{}, false, err1
			}
			if err2 != nil {
				return tupel.Tupel[T1, T2]{}, false, err2
			}
			if !ok1 || !ok2 {
				return tupel.Tupel[T1, T2]{}, false, nil
			}
			return *tupel.New(val1, val2), true, nil
		},
	}
}

func GroupBySeq[K common.Key, T any](seq *seq[T], keyMapper func(elem T) (K, error)) *pairsSlice[K, T] {
	result := emptyMapWithSlices[K, T]()
	if seq == nil {
		return result
	}
	for {
		val, ok, err := seq.next()
		if err != nil {
			result.err = err
			return result
		}
		if !ok {
			break
		}
		key, err := keyMapper(val)
		if err != nil {
			result.err = err
			return result
		}
		result.m[key] = append(result.m[key], val)
	}
	return result
}
