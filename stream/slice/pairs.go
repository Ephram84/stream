package slice

import (
	"maps"

	"github.com/Ephram84/stream/stream/common"
)

// pairs represents a simple map with single values (map[K]V).
// It supports eager evaluation with operations like Filter, Flatten, and Keys.
type pairs[K common.Key, V any] struct {
	m   map[K]V
	err error
}

func emptyPairs[K common.Key, V any]() *pairs[K, V] {
	return &pairs[K, V]{
		m:   make(map[K]V),
		err: nil,
	}
}

// FromMap creates a new pairs from an existing map with single values.
// The map is stored internally and can be processed using various operations.
func FromMap[K common.Key, V any](m map[K]V) *pairs[K, V] {
	return &pairs[K, V]{
		m:   m,
		err: nil,
	}
}

// Filter returns a new pairs containing only the key-value pairs that match the predicate.
// The predicate receives both the key and value for evaluation.
func (p *pairs[K, V]) Filter(predicate func(key K, value V) (bool, error)) *pairs[K, V] {
	pairs := emptyPairs[K, V]()
	if p.err != nil {
		pairs.err = p.err
		return pairs
	}

	for key, value := range p.m {
		ok, err := predicate(key, value)
		if err != nil {
			pairs.err = err
			return pairs
		}
		if ok {
			pairs.m[key] = value
		}
	}

	return pairs
}

// Flatten returns a slice containing all values from the map.
// The order of values is not guaranteed due to map iteration.
func (p *pairs[K, V]) Flatten() *slice[V] {
	s := emptySlice[V](0)
	if p.err != nil {
		s.err = p.err
		return s
	}

	for _, value := range p.m {
		s.slice = append(s.slice, value)
	}

	return s
}

// Keys returns a slice containing all keys from the map.
// The order of keys is not guaranteed due to map iteration.
func (p *pairs[K, V]) Keys() *slice[K] {
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

// ToMap converts the pairs back to a regular Go map.
func (p *pairs[K, V]) ToMap() (map[K]V, error) {
	if p.err != nil {
		return nil, p.err
	}

	m := make(map[K]V)

	maps.Copy(m, p.m)

	return m, nil
}

// Count returns the number of key-value pairs in the map.
func (p *pairs[K, V]) Count() (int, error) {
	if p.err != nil {
		return 0, p.err
	}

	return len(p.m), nil
}

// ForEach executes the provided action function for each key-value pair in the map.
func (p *pairs[K, V]) ForEach(action func(key K, value V)) error {
	if p.err != nil {
		return p.err
	}

	for key, value := range p.m {
		action(key, value)
	}
	return nil
}
