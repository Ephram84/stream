package stream

type iterator struct {
	start int
	next  func(current int) int
	skip  int
	limit int
}

func Iterator(start int, next func(current int) int) *iterator {
	return &iterator{
		start: start,
		next:  next,
		limit: 100,
	}
}

func (i *iterator) WithSkip(skip int) *iterator {
	if skip > 0 {
		i.skip = skip
	}

	return i
}

func (i *iterator) WithLimit(limit int) *iterator {
	if limit > 0 {
		i.limit = limit
	}

	return i
}

func (i *iterator) ToStream(sizes ...int) *stream[int] {
	size := getSize(sizes)
	out := make(chan int, size)

	go func() {
		current := i.start
		for i.skip > 0 {
			current = i.next(current)
			i.skip--
		}

		for elem := 0; elem < i.limit; elem++ {
			out <- current
			current = i.next(current)
		}
		close(out)
	}()

	return &stream[int]{
		stream: out,
		size:   size,
	}
}
