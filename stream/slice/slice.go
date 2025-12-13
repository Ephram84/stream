// Package slice provides eager, in-memory stream processing operations for Go slices.
// All operations are evaluated immediately when called, making it suitable for batch processing
// and when all data can fit comfortably in memory.
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

// slice represents an eagerly evaluated stream of elements.
// All operations are performed immediately on the internal slice.
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

// From creates a new slice stream from an existing Go slice.
// The input slice is copied to avoid external modifications.
func From[T any](tokens []T) *slice[T] {
	arr := slice[T]{
		slice: make([]T, len(tokens)),
		err:   nil,
	}

	copy(arr.slice, tokens)

	return &arr
}

// FromFile creates a slice stream of strings from a file.
// The file content is split by whitespace into individual string elements.
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

// Filter returns a new slice containing only elements that match the predicate.
// Elements are evaluated eagerly and filtered immediately.
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

// Map transforms each element using the mapper function.
// Returns a new slice with transformed elements of the same type.
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

// MapToInt maps each element to an int using the provided mapper function.
func (s *slice[T]) MapToInt(mapper func(elem T) (int, error)) *slice[int] {
	return Map(s, mapper)
}

// MapToInt64 maps each element to an int64 using the provided mapper function.
func (s *slice[T]) MapToInt64(mapper func(elem T) (int64, error)) *slice[int64] {
	return Map(s, mapper)
}

// MapToFloat maps each element to a float64 using the provided mapper function.
func (s *slice[T]) MapToFloat(mapper func(elem T) (float64, error)) *slice[float64] {
	return Map(s, mapper)
}

// MapToString maps each element to a string using the provided mapper function.
func (s *slice[T]) MapToString(mapper func(elem T) (string, error)) *slice[string] {
	return Map(s, mapper)
}

// PartitioningBy partitions elements into two groups based on the predicate.
// Returns a pairsSlice with boolean keys (true/false) mapping to slices of elements.
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

// GroupByString groups elements by string keys using the grouper function.
// Returns a pairsSlice mapping each string key to a slice of elements.
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

// AssociateByString creates a map associating string keys to single values.
// If multiple elements map to the same key, the last one wins.
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

// Sort returns a new sorted slice using the provided comparison function.
// The original slice remains unchanged.
func (s *slice[T]) Sort(sortFunc common.SortFunc[T]) *slice[T] {
	if s.err != nil {
		return s
	}

	sorted := emptySlice[T](len(s.slice))
	copy(sorted.slice, s.slice)
	sort.Slice(sorted.slice, sortFunc(sorted.slice))

	return sorted
}

// Take returns a new slice containing at most the first n elements.
// If n is greater than the slice length, returns all elements.
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

// Skip returns a new slice with the first n elements removed.
// If n is greater than the slice length, returns an empty slice.
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

// Reverse returns a new slice with elements in reverse order.
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

// ToSlice returns the internal slice as a regular Go slice.
// The returned slice is a copy to prevent external modifications.
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

// First returns a pointer to the first element in the slice.
// If the slice is empty, returns the optional default value (orElse) if provided.
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

// FirstOrNil returns a pointer to the first element that matches the predicate.
// Returns nil if no element matches.
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

// Last returns a pointer to the last element in the slice.
// If the slice is empty, returns the optional default value (orElse) if provided.
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

// LastOrNil returns a pointer to the last element that matches the predicate.
// Returns nil if no element matches.
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

// AnyMatch returns true if at least one element matches the predicate.
// Short-circuits on the first match.
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

// AllMatch returns true if all elements match the predicate.
// Short-circuits on the first non-match.
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

// NoneMatch returns true if no elements match the predicate.
// Short-circuits on the first match.
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

// Reduce applies a batch accumulator to reduce all elements to a single value.
func (s *slice[T]) Reduce(accumulator accumulator.Accumulator[T]) (T, error) {
	if s.err != nil {
		var zero T
		return zero, s.err
	}

	return accumulator(s.slice)
}

// ForEach executes the consumer function for each element in the slice.
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

// Map transforms a slice of type T to a slice of type R using the mapper function.
// This free function allows type transformation across different types.
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

// FlatMap transforms each element to a slice and flattens the results.
// Useful for expanding nested structures into a single flat slice.
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

// Zip combines two slices into a slice of tuples.
// The resulting slice length is the minimum of the two input slices.
func Zip[T1, T2 any](s1 *slice[T1], s2 *slice[T2]) *slice[tupel.Tupel[T1, T2]] {
	newSlice := emptySlice[tupel.Tupel[T1, T2]](0)
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

// GroupBy groups elements by keys extracted using the keyMapper function.
// Returns a pairsSlice mapping each key to a slice of elements with that key.
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
