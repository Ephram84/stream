package stream

import (
	"encoding/json"
	"io"
	"os"
	"sort"
	"strings"
	"sync"
)

type stream[T any] struct {
	stream chan T
	err    error
	size   int
}

func (s *stream[T]) setError(err error) {
	if err != nil && s.err == nil {
		s.err = err
	}
}

func FromFile(path string, sizes ...int) *stream[string] {
	f, err := os.Open(path)
	if err != nil {
		return &stream[string]{
			err:  err,
			size: getSize(sizes),
		}
	}
	input, err := io.ReadAll(f)
	if err != nil {
		return &stream[string]{
			err:  err,
			size: getSize(sizes),
		}
	}

	fields := strings.Fields(string(input))
	return From(fields)
}

func From[T any](tokens []T, sizes ...int) *stream[T] {
	size := getSize(sizes)

	out := make(chan T, size)

	go func() {
		for idx := range tokens {
			out <- tokens[idx]
		}
		close(out)
	}()

	return &stream[T]{
		stream: out,
		size:   size,
	}
}

func (s *stream[T]) WithSize(size int) *stream[T] {
	s.size = size
	return s
}

func (s *stream[T]) Filter(filter func(elem T) (bool, error)) *stream[T] {
	if s.err != nil {
		return s
	}

	out := make(chan T, s.size)

	go func() {
		for elem := range s.stream {
			if result, err := filter(elem); err != nil {
				s.setError(err)
			} else {
				if result {
					out <- elem
				}
			}
		}
		close(out)
	}()

	return &stream[T]{
		stream: out,
		err:    s.err,
		size:   s.size,
	}
}

func (s *stream[T]) ForEach(f func(elem T) T) *stream[T] {
	if s.err != nil {
		return s
	}

	out := make(chan T, s.size)

	go func() {
		for elem := range s.stream {
			out <- f(elem)
		}
		close(out)
	}()

	return &stream[T]{
		stream: out,
		err:    s.err,
		size:   s.size,
	}
}

func (s *stream[T]) PartitioningBy(predicate func(elem T) (bool, error)) *streamM[bool, T] {
	if s.err != nil {
		return &streamM[bool, T]{
			err: s.err,
		}
	}

	out := make(chan pair[bool, T], s.size)
	newStream := &streamM[bool, T]{
		stream: out,
		size:   s.size,
	}
	go func() {
		for elem := range s.stream {
			if b, err := predicate(elem); err != nil {
				newStream.setError(err)
			} else {
				out <- pair[bool, T]{
					key:   b,
					value: elem,
				}
			}
		}
		close(out)
	}()

	return newStream
}

func (s *stream[T]) GroupByString(grouper func(elem T) string) *streamM[string, T] {
	if s.err != nil {
		return &streamM[string, T]{
			err: s.err,
		}
	}

	out := make(chan pair[string, T], s.size)
	go func() {
		for elem := range s.stream {
			out <- pair[string, T]{
				key:   grouper(elem),
				value: elem,
			}
		}
		close(out)
	}()

	return &streamM[string, T]{
		stream: out,
		size:   s.size,
	}
}

func (s *stream[T]) GroupByInt(grouper func(elem T) int) *streamM[int, T] {
	if s.err != nil {
		return &streamM[int, T]{
			err: s.err,
		}
	}

	out := make(chan pair[int, T], s.size)
	go func() {
		for elem := range s.stream {
			out <- pair[int, T]{
				key:   grouper(elem),
				value: elem,
			}
		}
		close(out)
	}()

	return &streamM[int, T]{
		stream: out,
		size:   s.size,
	}
}

func (s *stream[T]) GroupByInt64(grouper func(elem T) int64) *streamM[int64, T] {
	if s.err != nil {
		return &streamM[int64, T]{
			err: s.err,
		}
	}

	out := make(chan pair[int64, T], s.size)
	go func() {
		for elem := range s.stream {
			out <- pair[int64, T]{
				key:   grouper(elem),
				value: elem,
			}
		}
		close(out)
	}()

	return &streamM[int64, T]{
		stream: out,
		size:   s.size,
	}
}

func (s *stream[T]) GroupByFloat(grouper func(elem T) float64) *streamM[float64, T] {
	if s.err != nil {
		return &streamM[float64, T]{
			err: s.err,
		}
	}

	out := make(chan pair[float64, T], s.size)
	go func() {
		for elem := range s.stream {
			out <- pair[float64, T]{
				key:   grouper(elem),
				value: elem,
			}
		}
		close(out)
	}()

	return &streamM[float64, T]{
		stream: out,
		size:   s.size,
	}
}

func (s *stream[T]) MapToFloat(m func(elem T) (float64, error)) *stream[float64] {
	newStream := &stream[float64]{
		err:  s.err,
		size: s.size,
	}

	if s.err != nil {
		return newStream
	}

	out := make(chan float64, s.size)
	newStream.stream = out
	go func() {
		for elem := range s.stream {
			if f, err := m(elem); err != nil {
				newStream.setError(err)
			} else {
				out <- f
			}
		}
		close(out)
	}()

	return newStream
}

func (s *stream[T]) MapToInt(m func(elem T) (int, error)) *stream[int] {
	newStream := &stream[int]{
		err:  s.err,
		size: s.size,
	}

	if s.err != nil {
		return newStream
	}

	out := make(chan int, s.size)
	newStream.stream = out
	go func() {
		for elem := range s.stream {
			if f, err := m(elem); err != nil {
				newStream.setError(err)
			} else {
				out <- f
			}
		}
		close(out)
	}()

	return newStream
}

func Map[T1, T2 any](s *stream[T1], mapper func(elem T1) (T2, error)) *stream[T2] {
	if s.err != nil {
		return &stream[T2]{
			err: s.err,
		}
	}

	out := make(chan T2, s.size)
	newStream := &stream[T2]{
		stream: out,
		size:   s.size,
	}

	go func() {
		for elem := range s.stream {
			if newElem, err := mapper(elem); err != nil {
				newStream.setError(err)
			} else {
				out <- newElem
			}
		}
		close(out)
	}()

	return newStream
}

func MapSlice[T1, T2 any](s *stream[T1], mapper func(elem T1) ([]T2, error)) *stream[T2] {
	if s.err != nil {
		return &stream[T2]{
			err: s.err,
		}
	}

	out := make(chan T2, s.size)
	newStream := &stream[T2]{
		stream: out,
		size:   s.size,
	}

	go func() {
		wg := &sync.WaitGroup{}
		for elem := range s.stream {
			newElems, err := mapper(elem)
			if err != nil {
				newStream.setError(err)
			} else {
				wg.Add(1)
				go func(wg *sync.WaitGroup) {
					defer wg.Done()
					for idx := range newElems {
						out <- newElems[idx]
					}
				}(wg)
			}
		}
		wg.Wait()
		close(out)
	}()

	return newStream
}

// terminal functions

func (s *stream[T]) FindFirst(orElse ...T) (*T, error) {
	if s.err != nil {
		return nil, s.err
	}

	firstElem, isOpen := <-s.stream
	if !isOpen {
		if len(orElse) > 0 {
			return &orElse[0], nil
		}
		return nil, nil
	}
	return &firstElem, nil
}

func (s *stream[T]) Count() (int, error) {
	slice, err := s.ToArray()
	if err != nil {
		return 0, s.err
	}

	return len(slice), nil
}

func (s *stream[T]) AnyMatch(predicate func(elem T) (bool, error)) (bool, error) {
	if s.err != nil {
		return false, s.err
	}

	for elem := range s.stream {
		result, err := predicate(elem)
		if err != nil {
			return false, err
		}
		if result {
			return true, nil
		}
	}

	return false, nil
}

func (s *stream[T]) AllMatch(predicate func(elem T) (bool, error)) (bool, error) {
	if s.err != nil {
		return false, s.err
	}

	for elem := range s.stream {
		result, err := predicate(elem)
		if err != nil {
			return false, err
		}
		if !result {
			return false, nil
		}
	}

	return true, nil
}

func (s *stream[T]) NoneMatch(predicate func(elem T) (bool, error)) (bool, error) {
	if s.err != nil {
		return false, s.err
	}

	for elem := range s.stream {
		result, err := predicate(elem)
		if err != nil {
			return false, err
		}
		if result {
			return false, nil
		}
	}

	return true, nil
}

func (s *stream[T]) ToArray() ([]T, error) {
	if s.err != nil {
		return nil, s.err
	}

	array := []T{}

	for elem := range s.stream {
		array = append(array, elem)
	}

	return array, nil
}

func (s *stream[T]) ToSortedArray(sortFunc func(array []T) func(i, j int) bool) ([]T, error) {
	array, err := s.ToArray()
	if err != nil {
		return nil, s.err
	}

	sort.Slice(array, sortFunc(array))

	return array, nil
}

func (s *stream[T]) Write(writer io.Writer) (int, error) {
	slice, err := s.ToArray()
	if err != nil {
		return 0, s.err
	}

	bytes, err := json.Marshal(slice)
	if err != nil {
		return 0, err
	}

	return writer.Write(bytes)
}

func (s *stream[T]) Max(comparator func(max, elem T) bool) (*T, error) {
	slice, err := s.ToArray()
	if err != nil {
		return nil, err
	}

	switch len(slice) {
	case 0:
		return nil, nil
	case 1:
		return &slice[0], nil
	case 2:
		if comparator(slice[0], slice[1]) {
			return &slice[0], nil
		} else {
			return &slice[1], nil
		}
	default:
		max := slice[0]
		for _, elem := range slice[1:] {
			if comparator(max, elem) {
				max = elem
			}
		}
		return &max, nil
	}
}

func (s *stream[T]) Min(comparator func(min, elem T) bool) (*T, error) {
	slice, err := s.ToArray()
	if err != nil {
		return nil, err
	}

	switch len(slice) {
	case 0:
		return nil, nil
	case 1:
		return &slice[0], nil
	case 2:
		if comparator(slice[0], slice[1]) {
			return &slice[0], nil
		} else {
			return &slice[1], nil
		}
	default:
		min := slice[0]
		for _, elem := range slice[1:] {
			if comparator(min, elem) {
				min = elem
			}
		}
		return &min, nil
	}
}
