package stream

type Key interface {
	int | int64 | float64 | string | bool
}

type pair[K Key, V any] struct {
	key   K
	value V
}

type streamM[K Key, V any] struct {
	stream chan pair[K, V]
	err    error
}

func (s *streamM[K, V]) setError(err error) {
	if err != nil && s.err == nil {
		s.err = err
	}
}

func (m *streamM[K, V]) StreamKeys() *stream[K] {
	out := make(chan K, 1)

	go func() {
		for pair := range m.stream {
			out <- pair.key
		}
		close(out)
	}()

	return &stream[K]{
		stream: out,
	}
}

func (m *streamM[K, V]) FlatMapValues() *stream[V] {
	out := make(chan V, 1)

	go func() {
		for pair := range m.stream {
			out <- pair.value
		}
		close(out)
	}()

	return &stream[V]{
		stream: out,
	}
}

func (m *streamM[K, V]) ToMap() (map[K][]V, error) {
	if m.err != nil {
		return nil, m.err
	}

	result := map[K][]V{}

	for pair := range m.stream {
		values, contains := result[pair.key]
		if !contains {
			values = make([]V, 0)
		}

		result[pair.key] = append(values, pair.value)
	}

	return result, nil
}

func (m *streamM[K, V]) Reducing(reducer func(key K, values []V) (K, V)) (map[K]V, error) {
	aggregatedMap, err := m.ToMap()
	if err != nil {
		return nil, m.err
	}

	result := map[K]V{}
	for key, values := range aggregatedMap {
		nKey, nValue := reducer(key, values)
		result[nKey] = nValue
	}

	return result, nil
}
