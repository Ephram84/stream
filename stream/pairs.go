package stream

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

func (p *pairs[K, V]) setError(err error) {
	if err != nil && p.err == nil {
		p.err = err
	}
}

func (p *pairs[K, V]) Flatten() *slice[V] {
	arr := emptySlice[V]()
	if p.err != nil {
		arr.setError(p.err)
		return arr
	}

	for _, value := range p.m {
		arr.slice = append(arr.slice, value)
	}

	return arr
}

func (p *pairs[K, V]) Keys() *slice[K] {
	s := emptySlice[K]()
	if p.err != nil {
		s.setError(p.err)
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
