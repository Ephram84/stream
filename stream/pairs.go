package stream

type pairs[K Key, V any] struct {
	m   map[K]V
	err error
}

func newPairs[K Key, V any](err error) *pairs[K, V] {
	return &pairs[K, V]{
		m:   make(map[K]V),
		err: err,
	}
}

func (p *pairs[K, V]) setError(err error) {
	if err != nil && p.err == nil {
		p.err = err
	}
}

func (p *pairs[K, V]) Flatten() *slice[V] {
	arr := &slice[V]{
		slice: make([]V, 0),
		err:   p.err,
	}

	if arr.err != nil {
		return arr
	}

	for _, value := range p.m {
		arr.slice = append(arr.slice, value)
	}

	return arr
}

func (p *pairs[K, V]) Keys() *slice[K] {
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

func (p *pairs[K, V]) ToMap() (map[K]V, error) {
	if p.err != nil {
		return nil, p.err
	}

	m := make(map[K]V)

	for key, values := range p.m {
		m[key] = values
	}

	return m, nil
}

func (p *pairs[K, V]) Count() (int, error) {
	if p.err != nil {
		return 0, p.err
	}

	return len(p.m), nil
}

func (p *pairs[K, V]) ForEach(action func(key K, value V)) *pairs[K, V] {
	if p.err != nil {
		return p
	}

	for key, value := range p.m {
		action(key, value)
	}

	return p
}
