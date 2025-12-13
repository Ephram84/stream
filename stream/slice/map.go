package slice

import (
	"github.com/Ephram84/stream/stream/accumulator"
	"github.com/Ephram84/stream/stream/common"
)

// pairsSlice represents a map with slice values (map[K][]V).
// It supports eager evaluation with operations like Flatten, Reduce, and MapToFloat64.
type pairsSlice[K common.Key, V any] struct {
	m   map[K][]V
	err error
}

// FromMapWithSlices creates a new pairsSlice from an existing map with slice values.
// The map is stored internally and can be processed using various operations.
func FromMapWithSlices[K common.Key, V any](m map[K][]V) *pairsSlice[K, V] {
	return &pairsSlice[K, V]{
		m:   m,
		err: nil,
	}
}

func emptyMapWithSlices[K common.Key, V any]() *pairsSlice[K, V] {
	return &pairsSlice[K, V]{
		m:   make(map[K][]V),
		err: nil,
	}
}

// CountValues returns a simple pairs map containing the count of elements for each key.
// Converts map[K][]V to map[K]int where the int is the length of each slice.
func (p pairsSlice[K, V]) CountValues() *pairs[K, int] {
	newPairs := emptyPairs[K, int]()

	if p.err != nil {
		newPairs.err = p.err
		return newPairs
	}

	for key, values := range p.m {
		newPairs.m[key] = len(values)
	}

	return newPairs
}

// Flatten returns a slice containing all values from all slices in the map.
// The order of elements is not guaranteed due to map iteration.
func (p *pairsSlice[K, V]) Flatten() *slice[V] {
	arr := emptySlice[V](0)
	if p.err != nil {
		arr.err = p.err
		return arr
	}

	for _, values := range p.m {
		arr.slice = append(arr.slice, values...)
	}

	return arr
}

// Keys returns a slice containing all keys from the map.
// The order of keys is not guaranteed due to map iteration.
func (p *pairsSlice[K, V]) Keys() *slice[K] {
	s := emptySlice[K](len(p.m))
	if p.err != nil {
		s.err = p.err
		return s
	}

	idx := 0
	for key := range p.m {
		s.slice[idx] = key
		idx++
	}

	return s
}

// MapToFloat64 transforms all values in all slices to float64 using the provided mapper function.
// Returns a new pairsSlice with float64 values while preserving the map structure.
func (p *pairsSlice[K, V]) MapToFloat64(mapper func(elem V) (float64, error)) *pairsSlice[K, float64] {
	pairsSlice := emptyMapWithSlices[K, float64]()
	if pairsSlice.err != nil {
		pairsSlice.err = p.err
		return pairsSlice
	}

	for key, values := range p.m {
		floats := make([]float64, len(values))
		for idx := range values {
			result, err := mapper(values[idx])
			if err != nil {
				pairsSlice.err = err
				return pairsSlice
			}

			floats[idx] = result
		}
		pairsSlice.m[key] = floats
	}

	return pairsSlice
}

// Reduce applies a batch accumulator to each slice in the map, converting it to a single value.
// Returns a simple pairs map (map[K]V) where each key maps to the reduced value of its slice.
// Slices with zero elements are skipped, slices with one element use that element directly.
func (p *pairsSlice[K, V]) Reduce(acc accumulator.Accumulator[V]) *pairs[K, V] {
	pairs := emptyPairs[K, V]()
	if pairs.err != nil {
		pairs.err = p.err
		return pairs
	}

	for key, values := range p.m {
		switch len(values) {
		case 0:
			continue
		case 1:
			pairs.m[key] = values[0]
		default:
			result, err := acc(values)
			if err != nil {
				pairs.err = err
				return pairs
			}
			pairs.m[key] = result
		}
	}

	return pairs
}

// ToMap converts the pairsSlice back to a regular Go map with slice values.
func (p *pairsSlice[K, V]) ToMap() (map[K][]V, error) {
	if p.err != nil {
		return nil, p.err
	}

	m := make(map[K][]V)

	for key, values := range p.m {
		m[key] = values
	}

	return m, nil
}

// Count returns the number of keys in the map.
func (p *pairsSlice[K, V]) Count() (int, error) {
	if p.err != nil {
		return 0, p.err
	}

	return len(p.m), nil
}

// ForEach executes the provided action function for each key-slice pair in the map.
// Returns the pairsSlice for method chaining.
func (p *pairsSlice[K, V]) ForEach(action func(key K, value []V)) *pairsSlice[K, V] {
	if p.err != nil {
		return p
	}

	for key, value := range p.m {
		action(key, value)
	}

	return p
}
