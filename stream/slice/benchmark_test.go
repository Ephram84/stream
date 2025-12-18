package slice

import (
	"testing"

	"github.com/Ephram84/stream/stream/accumulator"
	"github.com/Ephram84/stream/stream/common"
)

func BenchmarkFilter_Small(b *testing.B) {
	data := make([]int, 100)
	for i := range data {
		data[i] = i
	}
	b.ResetTimer()
	for b.Loop() {
		_, _ = From(data).Filter(isEven).ToSlice()
	}
}

func BenchmarkFilter_Medium(b *testing.B) {
	data := make([]int, 10000)
	for i := range data {
		data[i] = i
	}
	b.ResetTimer()
	for b.Loop() {
		_, _ = From(data).Filter(isEven).ToSlice()
	}
}

func BenchmarkFilter_Large(b *testing.B) {
	data := make([]int, 1000000)
	for i := range data {
		data[i] = i
	}
	b.ResetTimer()
	for b.Loop() {
		_, _ = From(data).Filter(isEven).ToSlice()
	}
}

func BenchmarkMap_Small(b *testing.B) {
	data := make([]int, 100)
	for i := range data {
		data[i] = i
	}
	b.ResetTimer()
	for b.Loop() {
		_, _ = From(data).Map(func(x int) (int, error) { return x * 2, nil }).ToSlice()
	}
}

func BenchmarkMap_Medium(b *testing.B) {
	data := make([]int, 10000)
	for i := range data {
		data[i] = i
	}
	b.ResetTimer()
	for b.Loop() {
		_, _ = From(data).Map(func(x int) (int, error) { return x * 2, nil }).ToSlice()
	}
}

func BenchmarkMap_Large(b *testing.B) {
	data := make([]int, 1000000)
	for i := range data {
		data[i] = i
	}
	b.ResetTimer()
	for b.Loop() {
		_, _ = From(data).Map(func(x int) (int, error) { return x * 2, nil }).ToSlice()
	}
}

func BenchmarkReduce_Small(b *testing.B) {
	data := make([]int, 100)
	for i := range data {
		data[i] = i
	}
	b.ResetTimer()
	for b.Loop() {
		_, _ = From(data).Reduce(accumulator.Sum[int])
	}
}

func BenchmarkReduce_Medium(b *testing.B) {
	data := make([]int, 10000)
	for i := range data {
		data[i] = i
	}
	b.ResetTimer()
	for b.Loop() {
		_, _ = From(data).Reduce(accumulator.Sum[int])
	}
}

func BenchmarkReduce_Large(b *testing.B) {
	data := make([]int, 1000000)
	for i := range data {
		data[i] = i
	}
	b.ResetTimer()
	for b.Loop() {
		_, _ = From(data).Reduce(accumulator.Sum[int])
	}
}

func BenchmarkFlatMap_Small(b *testing.B) {
	data := make([]int, 100)
	for i := range data {
		data[i] = i
	}
	b.ResetTimer()
	for b.Loop() {
		_, _ = FlatMap(From(data), func(x int) ([]int, error) {
			return []int{x, x * 2}, nil
		}).ToSlice()
	}
}

func BenchmarkFlatMap_Medium(b *testing.B) {
	data := make([]int, 10000)
	for i := range data {
		data[i] = i
	}
	b.ResetTimer()
	for b.Loop() {
		_, _ = FlatMap(From(data), func(x int) ([]int, error) {
			return []int{x, x * 2}, nil
		}).ToSlice()
	}
}

func BenchmarkChaining_FilterMapReduce_Small(b *testing.B) {
	data := make([]int, 100)
	for i := range data {
		data[i] = i
	}
	b.ResetTimer()
	for b.Loop() {
		_, _ = From(data).
			Filter(isEven).
			Map(func(x int) (int, error) { return x * 2, nil }).
			Reduce(accumulator.Sum[int])
	}
}

func BenchmarkChaining_FilterMapReduce_Medium(b *testing.B) {
	data := make([]int, 10000)
	for i := range data {
		data[i] = i
	}
	b.ResetTimer()
	for b.Loop() {
		_, _ = From(data).
			Filter(isEven).
			Map(func(x int) (int, error) { return x * 2, nil }).
			Reduce(accumulator.Sum[int])
	}
}

func BenchmarkChaining_FilterMapReduce_Large(b *testing.B) {
	data := make([]int, 1000000)
	for i := range data {
		data[i] = i
	}
	b.ResetTimer()
	for b.Loop() {
		_, _ = From(data).
			Filter(isEven).
			Map(func(x int) (int, error) { return x * 2, nil }).
			Reduce(accumulator.Sum[int])
	}
}

func BenchmarkDistinct_Small(b *testing.B) {
	data := make([]int, 100)
	for i := range data {
		data[i] = i % 10
	}
	b.ResetTimer()
	for b.Loop() {
		_, _ = From(data).Distinct(func(a, b int) bool { return a == b }).ToSlice()
	}
}

func BenchmarkDistinct_Medium(b *testing.B) {
	data := make([]int, 10000)
	for i := range data {
		data[i] = i % 100
	}
	b.ResetTimer()
	for b.Loop() {
		_, _ = From(data).Distinct(func(a, b int) bool { return a == b }).ToSlice()
	}
}

func BenchmarkSort_Small(b *testing.B) {
	data := make([]int, 100)
	for i := range data {
		data[i] = 100 - i
	}
	b.ResetTimer()
	for b.Loop() {
		_, _ = From(data).Sort(common.Sort[int]).ToSlice()
	}
}

func BenchmarkSort_Medium(b *testing.B) {
	data := make([]int, 10000)
	for i := range data {
		data[i] = 10000 - i
	}
	b.ResetTimer()
	for b.Loop() {
		_, _ = From(data).Sort(common.Sort[int]).ToSlice()
	}
}

func BenchmarkReverse_Small(b *testing.B) {
	data := make([]int, 100)
	for i := range data {
		data[i] = i
	}
	b.ResetTimer()
	for b.Loop() {
		_, _ = From(data).Reverse().ToSlice()
	}
}

func BenchmarkReverse_Medium(b *testing.B) {
	data := make([]int, 10000)
	for i := range data {
		data[i] = i
	}
	b.ResetTimer()
	for b.Loop() {
		_, _ = From(data).Reverse().ToSlice()
	}
}
