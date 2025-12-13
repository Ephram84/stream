// Package sequence provides lazy, iterator-based stream processing operations.
// Operations are evaluated on-demand, making it memory-efficient for large or infinite data streams.
package sequence

import (
	"io"
	"os"
	"strings"

	"github.com/Ephram84/stream/stream/accumulator"
	"github.com/Ephram84/stream/stream/common"
	"github.com/Ephram84/stream/stream/tupel"
)

// seq represents a lazily evaluated sequence of elements.
// Elements are produced on-demand through the next function.
type seq[T any] struct {
	next func() (T, bool, error)
}

// From creates a new lazy sequence from an existing Go slice.
// Elements are produced on-demand when the sequence is iterated.
// Optional error parameter can be provided to create a sequence in error state.
func From[T any](slice []T, errs ...error) *seq[T] {
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

// FromFile creates a lazy sequence of strings from a file.
// The file content is split by whitespace into individual string elements.
func FromFile(path string) *seq[string] {
	f, err := os.Open(path)
	if err != nil {
		return &seq[string]{
			next: func() (string, bool, error) {
				return "", false, err
			},
		}
	}
	input, err := io.ReadAll(f)
	if err != nil {
		return &seq[string]{
			next: func() (string, bool, error) {
				return "", false, err
			},
		}
	}

	fields := strings.Fields(string(input))
	return From(fields)
}

// Filter returns a new lazy sequence containing only elements that match the predicate.
// Elements are evaluated lazily as the sequence is iterated.
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

// Map transforms each element using the mapper function lazily.
// Returns a new sequence with transformed elements of the same type.
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

// MapToInt maps each element to an int using the provided mapper function lazily.
func (s *seq[T]) MapToInt(mapper func(elem T) (int, error)) *seq[int] {
	return Map(s, mapper)
}

// MapToInt64 maps each element to an int64 using the provided mapper function lazily.
func (s *seq[T]) MapToInt64(mapper func(elem T) (int64, error)) *seq[int64] {
	return Map(s, mapper)
}

// MapToFloat64 maps each element to a float64 using the provided mapper function lazily.
func (s *seq[T]) MapToFloat64(mapper func(elem T) (float64, error)) *seq[float64] {
	return Map(s, mapper)
}

// MapToString maps each element to a string using the provided mapper function lazily.
func (s *seq[T]) MapToString(mapper func(elem T) (string, error)) *seq[string] {
	return Map(s, mapper)
}

// PartitioningBy partitions elements into two groups based on the predicate.
// Returns a pairsSlice with boolean keys (true/false) mapping to slices of elements.
// Note: This is a terminal operation that evaluates the entire sequence.
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

// GroupByString groups elements by string keys using the grouper function.
// Returns a pairsSlice mapping each string key to a slice of elements.
func (s *seq[T]) GroupByString(grouper func(elem T) (string, error)) *pairsSlice[string, T] {
	return GroupBy(s, grouper)
}

// GroupByInt groups elements by int keys using the grouper function.
// Returns a pairsSlice mapping each int key to a slice of elements.
func (s *seq[T]) GroupByInt(grouper func(elem T) (int, error)) *pairsSlice[int, T] {
	return GroupBy(s, grouper)
}

// GroupByInt64 groups elements by int64 keys using the grouper function.
// Returns a pairsSlice mapping each int64 key to a slice of elements.
func (s *seq[T]) GroupByInt64(grouper func(elem T) (int64, error)) *pairsSlice[int64, T] {
	return GroupBy(s, grouper)
}

// GroupByFloat groups elements by float64 keys using the grouper function.
// Returns a pairsSlice mapping each float64 key to a slice of elements.
func (s *seq[T]) GroupByFloat(grouper func(elem T) (float64, error)) *pairsSlice[float64, T] {
	return GroupBy(s, grouper)
}

// AssociateByString creates a map associating string keys to single values.
// If multiple elements map to the same key, the last one wins.
// Note: This is a terminal operation that evaluates the entire sequence.
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

// AssociateByInt creates a map associating int keys to single values.
// If multiple elements map to the same key, the last one wins.
// Note: This is a terminal operation that evaluates the entire sequence.
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

// AssociateByInt64 creates a map associating int64 keys to single values.
// If multiple elements map to the same key, the last one wins.
// Note: This is a terminal operation that evaluates the entire sequence.
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

// AssociateByFloat creates a map associating float64 keys to single values.
// If multiple elements map to the same key, the last one wins.
// Note: This is a terminal operation that evaluates the entire sequence.
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

// Take returns a new lazy sequence containing at most the first n elements.
// Elements beyond n are not evaluated.
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

// Skip returns a new lazy sequence with the first n elements removed.
// The first n elements are discarded when the sequence is iterated.
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

// Reverse returns a new sequence with elements in reverse order.
// Note: This requires full evaluation of the sequence and stores all elements in memory.
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

// Concat concatenates this sequence with additional sequences lazily.
// Elements from each sequence are produced in order.
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

// ToSlice converts the lazy sequence to a regular Go slice.
// Note: This is a terminal operation that evaluates the entire sequence.
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

// First returns a pointer to the first element in the sequence.
// If the sequence is empty, returns the optional default value (orElse) if provided.
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

// FirstOrNil returns a pointer to the first element that matches the predicate.
// Returns nil if no element matches.
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

// AllMatch returns true if all elements match the predicate.
// Short-circuits on the first non-match.
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

// NoneMatch returns true if no elements match the predicate.
// Short-circuits on the first match.
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

// Reduce applies a sequential accumulator to reduce all elements to a single value.
// Uses sequential accumulation which processes elements one by one with an index.
// Note: This is a terminal operation that evaluates the entire sequence.
func (s *seq[T]) Reduce(acc accumulator.AccumulatorSeq[T]) (T, error) {
	var identity T
	var count int
	var first = true
	for {
		val, ok, err := s.next()
		if err != nil {
			return identity, err
		}
		if !ok {
			break
		}
		count++
		if first {
			identity = val
			first = false
			continue
		}
		identity, err = acc(count, identity, val)
		if err != nil {
			return identity, err
		}
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

// Map transforms a lazy sequence of type T to a lazy sequence of type R using the mapper function.
// This free function allows type transformation across different types.
func Map[T, R any](s *seq[T], mapper func(elem T) (R, error)) *seq[R] {
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

// FlatMap transforms each element to a sequence and flattens the results lazily.
// Useful for expanding nested structures into a single flat sequence.
func FlatMap[T, R any](s *seq[T], mapper func(elem T) (*seq[R], error)) *seq[R] {
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

// Zip combines two lazy sequences into a sequence of tuples.
// The resulting sequence ends when either input sequence ends.
func Zip[T1, T2 any](seq1 *seq[T1], seq2 *seq[T2]) *seq[tupel.Tupel[T1, T2]] {
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

// GroupBy groups sequence elements by keys extracted using the keyMapper function.
// Returns a pairsSlice mapping each key to a slice of elements with that key.
// Note: This is a terminal operation that evaluates the entire sequence.
func GroupBy[K common.Key, T any](seq *seq[T], keyMapper func(elem T) (K, error)) *pairsSlice[K, T] {
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
