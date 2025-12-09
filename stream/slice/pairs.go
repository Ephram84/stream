package slice

import (
	"maps"

	"github.com/Ephram84/stream/stream/common"
)

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

func FromMap[K common.Key, V any](m map[K]V) *pairs[K, V] {
	return &pairs[K, V]{
		m:   m,
		err: nil,
	}
}

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

func (p *pairs[K, V]) ToMap() (map[K]V, error) {
	if p.err != nil {
		return nil, p.err
	}

	m := make(map[K]V)

	maps.Copy(m, p.m)

	return m, nil
}

func (p *pairs[K, V]) Count() (int, error) {
	if p.err != nil {
		return 0, p.err
	}

	return len(p.m), nil
}

func (p *pairs[K, V]) ForEach(action func(key K, value V)) error {
	if p.err != nil {
		return p.err
	}

	for key, value := range p.m {
		action(key, value)
	}
	return nil
}
