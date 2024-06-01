package stream

import (
	"encoding/json"
	"io"
	"os"
	"sort"
	"strings"
)

type Numbers interface {
	~int | ~int64 | ~float64
}

type slice[T any] struct {
	slice []T
	err   error
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

func FlatMap[T1, T2 any](s *slice[T1], mapper func(elem T1) ([]T2, error)) *slice[T2] {
	newSlice := &slice[T2]{}
	if s == nil {
		return newSlice
	}

	if s.err != nil {
		newSlice.setError(s.err)
		return newSlice
	}

	for idx := range s.slice {
		newValues, err := mapper(s.slice[idx])
		if err != nil {
			newSlice.setError(err)
			return newSlice
		}

		newSlice.slice = append(newSlice.slice, newValues...)
	}

	return newSlice
}

func (a *slice[T]) setError(err error) {
	if err != nil && a.err == nil {
		a.err = err
	}
}

func (s *slice[T]) Filter(filter func(elem T) (bool, error)) *slice[T] {
	if s.err != nil {
		return s
	}

	newSlice := slice[T]{}
	for idx := range s.slice {
		result, err := filter(s.slice[idx])
		newSlice.setError(err)
		if result {
			newSlice.slice = append(newSlice.slice, s.slice[idx])
		}
	}

	return &newSlice
}

func (s *slice[T]) ForEach(action func(elem T)) *slice[T] {
	if s.err != nil {
		return s
	}

	for idx := range s.slice {
		action(s.slice[idx])
	}

	return s
}

func (s *slice[T]) Map(mapper func(elem T) (T, error)) *slice[T] {
	if s.err != nil {
		return s
	}

	newSlice := &slice[T]{
		slice: make([]T, len(s.slice)),
	}

	for idx := range s.slice {
		newValue, err := mapper(s.slice[idx])
		newSlice.setError(err)
		newSlice.slice[idx] = newValue
	}

	return newSlice
}

func (s *slice[T]) PartitioningBy(predicate func(elem T) (bool, error)) *pairsSlice[bool, T] {
	mapped := newPairsSlice[bool, T](s.err)
	if mapped.err != nil {
		return mapped
	}

	for idx := range s.slice {
		key, err := predicate(s.slice[idx])
		if err != nil {
			mapped.setError(err)
			return mapped
		}

		mapped.m[key] = append(mapped.m[key], s.slice[idx])
	}

	return mapped
}

func (s *slice[T]) GroupByString(grouper func(elem T) string) *pairsSlice[string, T] {
	mapped := newPairsSlice[string, T](s.err)
	if mapped.err != nil {
		return mapped
	}

	for idx := range s.slice {
		key := grouper(s.slice[idx])
		mapped.m[key] = append(mapped.m[key], s.slice[idx])
	}

	return mapped
}

func (s *slice[T]) MapToInt(mapper func(elem T) (int, error)) *slice[int] {
	newSlice := &slice[int]{
		err: s.err,
	}
	if newSlice.err != nil {
		return newSlice
	}

	mapped, err := mapTo(s.slice, mapper)
	newSlice.setError(err)
	newSlice.slice = mapped

	return newSlice
}

func (s *slice[T]) MapToInt64(mapper func(elem T) (int64, error)) *slice[int64] {
	newSlice := &slice[int64]{
		err: s.err,
	}
	if newSlice.err != nil {
		return newSlice
	}

	mapped, err := mapTo(s.slice, mapper)
	newSlice.setError(err)
	newSlice.slice = mapped

	return newSlice
}

func (s *slice[T]) MapToFloat64(mapper func(elem T) (float64, error)) *slice[float64] {
	newSlice := &slice[float64]{
		err: s.err,
	}
	if newSlice.err != nil {
		return newSlice
	}

	mapped, err := mapTo(s.slice, mapper)
	newSlice.setError(err)
	newSlice.slice = mapped

	return newSlice
}

func mapTo[T1, T2 any](elems []T1, mapper func(elem T1) (T2, error)) ([]T2, error) {
	slice := make([]T2, len(elems))

	for idx := range elems {
		value, err := mapper(elems[idx])
		if err != nil {
			return nil, err
		}
		slice[idx] = value
	}

	return slice, nil
}

func (s *slice[T]) Sort(sortFunc func(slice []T) func(i, j int) bool) *slice[T] {
	if s.err != nil {
		return s
	}

	sorted := &slice[T]{
		slice: make([]T, len(s.slice)),
	}
	copy(sorted.slice, s.slice)

	sort.Slice(sorted.slice, sortFunc(sorted.slice))

	return sorted
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

func (s *slice[T]) Reduce(identity T, accumulator Accumulator[T]) (T, error) {
	if s.err != nil {
		return identity, s.err
	}

	result := identity
	for idx := range s.slice {
		result = accumulator.Apply(result, s.slice[idx])
	}

	return result, nil
}
