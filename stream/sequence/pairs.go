package sequence

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

func (p *pairs[K, V]) Flatten() *seq[V] {
	arr := make([]V, 0)
	if p.err != nil {
		return From[V](nil, p.err)
	}

	for _, values := range p.m {
		arr = append(arr, values)
	}

	return From(arr, nil)
}

func (p *pairs[K, V]) Keys() *seq[K] {
	s := make([]K, 0, len(p.m))
	if p.err != nil {
		return From[K](nil, p.err)
	}

	for key := range p.m {
		s = append(s, key)
	}

	return From(s, nil)
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
