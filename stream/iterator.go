package stream

type iterator[N Numbers] struct {
	start N
	next  func(current N) N
	skip  int
	limit int
}

func Iterator[N Numbers](start N, next func(current N) N) *iterator[N] {
	return &iterator[N]{
		start: start,
		next:  next,
		limit: 100,
	}
}

func (i *iterator[N]) WithSkip(skip int) *iterator[N] {
	if skip > 0 {
		i.skip = skip
	}

	return i
}

func (i *iterator[N]) WithLimit(limit int) *iterator[N] {
	if limit > 0 {
		i.limit = limit
	}

	return i
}

func (i *iterator[N]) Generate() *slice[N] {
	slice := &slice[N]{
		slice: make([]N, 0),
	}
	current := i.start
	for i.skip > 0 {
		current = i.next(current)
		i.skip--
	}

	for elem := 0; elem < i.limit; elem++ {
		slice.slice = append(slice.slice, current)
		current = i.next(current)
	}

	return slice
}
