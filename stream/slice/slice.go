package slice

import (
	"encoding/json"
	"io"
	"os"
	"slices"
	"sort"
	"strings"

	"github.com/Ephram84/stream/stream/accumulator"
	"github.com/Ephram84/stream/stream/common"
	"github.com/Ephram84/stream/stream/tupel"
)

type slice[T any] struct {
	slice []T
	err   error
}

func emptySlice[T any](len int) *slice[T] {
	return &slice[T]{
		slice: make([]T, len),
		err:   nil,
	}
}

func From[T any](tokens []T) *slice[T] {
	arr := slice[T]{
		slice: make([]T, len(tokens)),
		err:   nil,
	}

	copy(arr.slice, tokens)

	return &arr
}

func FromFile(path string) *slice[string] {
	f, err := os.Open(path)
	if err != nil {
		return &slice[string]{
			err: err,
		}
	}
	input, err := io.ReadAll(f)
	if err != nil {
		return &slice[string]{
			err: err,
		}
	}

	fields := strings.Fields(string(input))
	return From(fields)
}

func (s *slice[T]) Filter(filter func(elem T) (bool, error)) *slice[T] {
	if s.err != nil {
		return s
	}

	newSlice := emptySlice[T](0)
	for idx := range s.slice {
		result, err := filter(s.slice[idx])
		if err != nil {
			newSlice.err = err
			return newSlice
		}
		if result {
			newSlice.slice = append(newSlice.slice, s.slice[idx])
		}
	}

	return newSlice
}

func (s *slice[T]) Map(mapper func(elem T) (T, error)) *slice[T] {
	if s.err != nil {
		return s
	}

	newSlice := emptySlice[T](len(s.slice))

	for idx := range s.slice {
		newValue, err := mapper(s.slice[idx])
		if err != nil {
			newSlice.err = err
			return newSlice
		}
		newSlice.slice[idx] = newValue
	}

	return newSlice
}

func (s *slice[T]) MapToInt(mapper func(elem T) (int, error)) *slice[int] {
	return Map(s, mapper)
}

func (s *slice[T]) MapToInt64(mapper func(elem T) (int64, error)) *slice[int64] {
	return Map(s, mapper)
}

func (s *slice[T]) MapToFloat(mapper func(elem T) (float64, error)) *slice[float64] {
	return Map(s, mapper)
}

func (s *slice[T]) MapToString(mapper func(elem T) (string, error)) *slice[string] {
	return Map(s, mapper)
}

func (s *slice[T]) PartitioningBy(predicate func(elem T) (bool, error)) *pairsSlice[bool, T] {
	mapped := emptyMapWithSlices[bool, T]()
	if s.err != nil {
		mapped.err = s.err
		return mapped
	}

	for idx := range s.slice {
		key, err := predicate(s.slice[idx])
		if err != nil {
			mapped.err = err
			return mapped
		}

		mapped.m[key] = append(mapped.m[key], s.slice[idx])
	}

	return mapped
}

func (s *slice[T]) GroupByString(grouper func(elem T) (string, error)) *pairsSlice[string, T] {
	mapped := emptyMapWithSlices[string, T]()
	if s.err != nil {
		mapped.err = s.err
		return mapped
	}

	for idx := range s.slice {
		key, err := grouper(s.slice[idx])
		if err != nil {
			mapped.err = err
			return mapped
		}
		mapped.m[key] = append(mapped.m[key], s.slice[idx])
	}

	return mapped
}

func (s *slice[T]) AssociateByString(mapper func(elem T) (string, error)) *pairs[string, T] {
	p := emptyPairs[string, T]()
	if s.err != nil {
		p.err = s.err
		return p
	}

	for idx := range s.slice {
		key, err := mapper(s.slice[idx])
		if err != nil {
			p.err = err
			return p
		}
		p.m[key] = s.slice[idx]
	}

	return p
}

func (s *slice[T]) Sort(sortFunc common.SortFunc[T]) *slice[T] {
	if s.err != nil {
		return s
	}

	sorted := emptySlice[T](len(s.slice))
	copy(sorted.slice, s.slice)
	sort.Slice(sorted.slice, sortFunc(sorted.slice))

	return sorted
}

func (s *slice[T]) Take(n int) *slice[T] {
	if s.err != nil {
		return s
	}

	newSlice := emptySlice[T](0)
	if n <= 0 {
		return newSlice
	}

	if n >= len(s.slice) {
		newSlice.slice = append(newSlice.slice, s.slice...)
	} else {
		newSlice.slice = append(newSlice.slice, s.slice[:n]...)
	}

	return newSlice
}

func (s *slice[T]) Skip(n int) *slice[T] {
	if s.err != nil {
		return s
	}

	newSlice := emptySlice[T](0)
	if n < len(s.slice) {
		newSlice.slice = append(newSlice.slice, s.slice[n:]...)
	}
	return newSlice
}

func (s *slice[T]) Reverse() *slice[T] {
	if s.err != nil {
		return s
	}

	reversed := emptySlice[T](0)
	for i := len(s.slice) - 1; i >= 0; i-- {
		reversed.slice = append(reversed.slice, s.slice[i])
	}
	return reversed
}

func (s *slice[T]) Distinct(eq func(a, b T) bool) *slice[T] {
	if s.err != nil {
		return s
	}

	distinct := &slice[T]{
		slice: slices.CompactFunc(s.slice, eq),
	}

	return distinct
}

func (s *slice[T]) Concat(slices ...*slice[T]) *slice[T] {
	if s.err != nil {
		return s
	}

	for _, slice := range slices {
		if slice.err != nil {
			s.err = slice.err
			return s
		}
		s.slice = append(s.slice, slice.slice...)
	}

	return s
}

// terminal functions

func (s *slice[T]) Write(writer io.Writer) (int, error) {
	if s.err != nil {
		return 0, s.err
	}

	bytes, err := json.Marshal(s.slice)
	if err != nil {
		return 0, err
	}

	return writer.Write(bytes)
}

func (s *slice[T]) ToSlice() ([]T, error) {
	if s.err != nil {
		return nil, s.err
	}
	slice := make([]T, len(s.slice))

	copy(slice, s.slice)

	return slice, nil
}

func (s *slice[T]) Count() (int, error) {
	if s.err != nil {
		return 0, s.err
	}

	return len(s.slice), nil
}

func (s *slice[T]) First(orElse ...T) (*T, error) {
	if s.err != nil {
		return nil, s.err
	}

	count, _ := s.Count()
	if count == 0 {
		if len(orElse) > 0 {
			return &orElse[0], nil
		}

		return nil, nil
	}

	return &s.slice[0], nil
}

func (s *slice[T]) FirstOrNil(predicate func(elem T) (bool, error)) (*T, error) {
	if s.err != nil {
		return nil, s.err
	}

	for idx := range s.slice {
		result, err := predicate(s.slice[idx])
		if err != nil {
			return nil, err
		}
		if result {
			return &s.slice[idx], nil
		}
	}

	return nil, nil
}

func (s *slice[T]) Last(orElse ...T) (*T, error) {
	if s.err != nil {
		return nil, s.err
	}

	count, _ := s.Count()
	if count == 0 {
		if len(orElse) > 0 {
			return &orElse[0], nil
		}

		return nil, nil
	}

	return &s.slice[len(s.slice)-1], nil
}

func (s *slice[T]) LastOrNil(predicate func(elem T) (bool, error)) (*T, error) {
	if s.err != nil {
		return nil, s.err
	}
	for i := len(s.slice) - 1; i >= 0; i-- {
		result, err := predicate(s.slice[i])
		if err != nil {
			return nil, err
		}
		if result {
			return &s.slice[i], nil
		}
	}

	return nil, nil
}

func (s *slice[T]) AnyMatch(predicate func(elem T) (bool, error)) (bool, error) {
	if s.err != nil {
		return false, s.err
	}

	for idx := range s.slice {
		result, err := predicate(s.slice[idx])
		if err != nil {
			return false, err
		}

		if result {
			return true, nil
		}
	}

	return false, nil
}

func (s *slice[T]) AllMatch(predicate func(elem T) (bool, error)) (bool, error) {
	if s.err != nil {
		return false, s.err
	}

	for idx := range s.slice {
		result, err := predicate(s.slice[idx])
		if err != nil {
			return false, err
		}

		if !result {
			return false, nil
		}
	}

	return true, nil
}

func (s *slice[T]) NoneMatch(predicate func(elem T) (bool, error)) (bool, error) {
	if s.err != nil {
		return false, s.err
	}

	for idx := range s.slice {
		result, err := predicate(s.slice[idx])
		if err != nil {
			return false, err
		}

		if result {
			return false, nil
		}
	}

	return true, nil
}

func (s *slice[T]) Reduce(accumulator accumulator.Accumulator[T]) (T, error) {
	if s.err != nil {
		var zero T
		return zero, s.err
	}

	return accumulator(s.slice)
}

func (s *slice[T]) ForEach(consumer func(elem T) error) error {
	if s.err != nil {
		return s.err
	}

	for idx := range s.slice {
		if err := consumer(s.slice[idx]); err != nil {
			return err
		}
	}

	return nil
}

// free functions

func Map[T, R any](s *slice[T], mapper func(elem T) (R, error)) *slice[R] {
	if s == nil {
		return emptySlice[R](0)
	}
	newSlice := emptySlice[R](len(s.slice))

	if s.err != nil {
		newSlice.err = s.err
		return newSlice
	}

	for idx := range s.slice {
		newValue, err := mapper(s.slice[idx])
		if err != nil {
			newSlice.err = err
			return newSlice
		}
		newSlice.slice[idx] = newValue
	}

	return newSlice
}

func FlatMap[T1, T2 any](s *slice[T1], mapper func(elem T1) ([]T2, error)) *slice[T2] {
	newSlice := emptySlice[T2](0)
	if s == nil {
		return newSlice
	}

	if s.err != nil {
		newSlice.err = s.err
		return newSlice
	}

	for idx := range s.slice {
		newValues, err := mapper(s.slice[idx])
		if err != nil {
			newSlice.err = err
			return newSlice
		}

		newSlice.slice = append(newSlice.slice, newValues...)
	}

	return newSlice
}

func Zip[T1, T2 any](s1 *slice[T1], s2 *slice[T2]) *slice[tupel.Tupel[T1, T2]] {
	newSlice := &slice[tupel.Tupel[T1, T2]]{}
	if s1 == nil || s2 == nil {
		return newSlice
	}

	if s1.err != nil {
		newSlice.err = s1.err
		return newSlice
	}

	if s2.err != nil {
		newSlice.err = s2.err
		return newSlice
	}

	minLen := common.Min(len(s1.slice), len(s2.slice))
	for i := range minLen {
		newSlice.slice = append(newSlice.slice, *tupel.New(s1.slice[i], s2.slice[i]))
	}

	return newSlice
}

func GroupBy[K common.Key, T any](s *slice[T], keyMapper func(elem T) (K, error)) *pairsSlice[K, T] {
	pairs := emptyMapWithSlices[K, T]()
	if s == nil {
		return pairs
	}
	if s.err != nil {
		pairs.err = s.err
		return pairs
	}

	for idx := range s.slice {
		key, err := keyMapper(s.slice[idx])
		if err != nil {
			pairs.err = err
			return pairs
		}
		pairs.m[key] = append(pairs.m[key], s.slice[idx])
	}

	return pairs
}
