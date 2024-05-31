package stream

type Key interface {
	int | int64 | float64 | string | bool
}

type pairsSlice[K Key, V any] struct {
	m   map[K][]V
	err error
}

func newPairsSlice[K Key, V any](err error) *pairsSlice[K, V] {
	return &pairsSlice[K, V]{
		m:   make(map[K][]V),
		err: err,
	}
}

func (p *pairsSlice[K, V]) setError(err error) {
	if err != nil && p.err == nil {
		p.err = err
	}
}

func (m pairsSlice[K, V]) CountValues() *pairs[K, int] {
	newPairs := newPairs[K, int](m.err)

	if newPairs.err != nil {
		return newPairs
	}

	for key, values := range m.m {
		newPairs.m[key] = len(values)
	}

	return newPairs
}

func (p *pairsSlice[K, V]) Flatten() *slice[V] {
	arr := &slice[V]{
		slice: make([]V, 0),
		err:   p.err,
	}

	if arr.err != nil {
		return arr
	}

	for _, values := range p.m {
		arr.slice = append(arr.slice, values...)
	}

	return arr
}

func (p *pairsSlice[K, V]) Keys() *slice[K] {
	s := &slice[K]{
		slice: make([]K, len(p.m)),
	}
	s.setError(p.err)

	idx := 0
	for key := range p.m {
		s.slice[idx] = key
		idx++
	}

	return s
}

func (p *pairsSlice[K, V]) MapToFloat64(mapper func(elem V) (float64, error)) *pairsSlice[K, float64] {
	pairsSlice := newPairsSlice[K, float64](p.err)
	if pairsSlice.err != nil {
		return pairsSlice
	}

	for key, values := range p.m {
		floats := make([]float64, len(values))
		for idx := range values {
			result, err := mapper(values[idx])
			if err != nil {
				pairsSlice.setError(err)
				return pairsSlice
			}

			floats[idx] = result
		}
		pairsSlice.m[key] = floats
	}

	return pairsSlice
}

func (p *pairsSlice[K, V]) Reduce(identity V, accumulator Accumulator[V]) *pairs[K, V] {
	pairs := newPairs[K, V](p.err)
	if pairs.err != nil {
		return pairs
	}

	for key, values := range p.m {
		switch len(values) {
		case 0:
			continue
		case 1:
			pairs.m[key] = accumulator.Apply(identity, values[0])
		default:
			result := identity
			for idx := range values {
				result = accumulator.Apply(result, values[idx])
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
