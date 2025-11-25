package sequence

import (
	"github.com/Ephram84/stream/stream/accumulator"
	"github.com/Ephram84/stream/stream/common"
)

type pairsSlice[K common.Key, V any] struct {
	m   map[K][]V
	err error
}

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

func (p *pairsSlice[K, V]) Flatten() *seq[V] {
	arr := make([]V, 0)
	if p.err != nil {
		return From[V](nil, p.err)
	}

	for _, values := range p.m {
		arr = append(arr, values...)
	}

	return From(arr, nil)
}

func (p *pairsSlice[K, V]) Keys() *seq[K] {
	s := make([]K, len(p.m))
	if p.err != nil {
		return From[K](nil, p.err)
	}

	idx := 0
	for key := range p.m {
		s[idx] = key
		idx++
	}

	return From(s, nil)
}

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

func (p *pairsSlice[K, V]) Reduce(identity V, acc accumulator.Accumulator[V]) *pairs[K, V] {
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
			pairs.m[key] = acc.Apply(identity, values[0])
		default:
			result := identity
			for idx := range values {
				result = acc.Apply(result, values[idx])
			}
			pairs.m[key] = result
		}
	}

	return pairs
}

func (p *pairsSlice[K, V]) Aggregate(acc accumulator.Accumulator[V]) *pairs[K, V] {
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
			result := values[0]
			for idx := 1; idx < len(values); idx++ {
				result = acc.Apply(result, values[idx])
			}
			pairs.m[key] = result
		}
	}

	return pairs
}

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

func (p *pairsSlice[K, V]) Count() (int, error) {
	if p.err != nil {
		return 0, p.err
	}

	return len(p.m), nil
}

func (p *pairsSlice[K, V]) ForEach(action func(key K, value []V)) *pairsSlice[K, V] {
	if p.err != nil {
		return p
	}

	for key, value := range p.m {
		action(key, value)
	}

	return p
}
