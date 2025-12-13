// Package iterator provides functionality for generating sequences of numbers lazily.
// Unlike eager evaluation, iterators generate values on-demand when Generate is called.
package iterator

import (
	"github.com/Ephram84/stream/stream/common"
)

// NextFunc is a function type that generates the next value in a sequence
// based on the current value.
type NextFunc[N common.Numbers] func(current N) N

// iterator represents a lazy sequence generator for numeric values.
// It supports configuration through WithSkip and WithLimit methods.
type iterator[N common.Numbers] struct {
	start N
	next  func(current N) N
	skip  int
	limit int
}

// Iterator creates a new iterator starting at the given value.
// The next function determines how to generate subsequent values.
// The default limit is 100 elements.
func Iterator[N common.Numbers](start N, next NextFunc[N]) *iterator[N] {
	return &iterator[N]{
		start: start,
		next:  next,
		limit: 100,
	}
}

// WithSkip sets the number of elements to skip at the beginning of the sequence.
// Only positive values are accepted; zero or negative values are ignored.
func (i *iterator[N]) WithSkip(skip int) *iterator[N] {
	if skip > 0 {
		i.skip = skip
	}

	return i
}

// WithLimit sets the maximum number of elements to generate.
// Only positive values are accepted; zero or negative values are ignored.
func (i *iterator[N]) WithLimit(limit int) *iterator[N] {
	if limit > 0 {
		i.limit = limit
	}

	return i
}

// Generate produces a slice containing the generated sequence.
// It first skips the configured number of elements, then generates
// up to the configured limit of elements.
func (i *iterator[N]) Generate() []N {
	s := make([]N, 0)
	current := i.start
	for i.skip > 0 {
		current = i.next(current)
		i.skip--
	}

	for elem := 0; elem < i.limit; elem++ {
		s = append(s, current)
		current = i.next(current)
	}

	return s
}

// IncrementInt returns a NextFunc that increments an integer by 1.
// Useful for generating sequences like 1, 2, 3, 4, ...
func IncrementInt() NextFunc[int] {
	return func(current int) int {
		return current + 1
	}
}

// IncrementFloat returns a NextFunc that increments a float64 by 1.0.
// Useful for generating sequences like 1.0, 2.0, 3.0, 4.0, ...
func IncrementFloat() NextFunc[float64] {
	return func(current float64) float64 {
		return current + 1.0
	}
}
