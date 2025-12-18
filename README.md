# stream

A Go package that brings the Kotlin collection and stream feeling to Go – as much as possible within Go's language constraints.

## Motivation

This repository was created out of the desire to use the expressive and functional collection and stream operations known from Kotlin in Go as well. The goal is to make working with data streams and collections in Go more elegant, declarative, and productive.

## Installation

Add the module to your project:

```bash
go get github.com/Ephram84/stream
```

## Features

- Functional stream API for Go
- Methods like `map`, `filter`, `reduce`, `sort`, `groupBy`, and more
- Chainable operations with both eager and lazy evaluation
- Easy integration into existing Go projects
- Two main processing approaches: eager (slice-based) and lazy (sequence-based)

## Architecture

This package provides two complementary approaches for data processing:

- [Eager Processing](stream/slice/README.md) - Immediate evaluation with in-memory slices
- [Lazy Processing](stream/sequence/README.md) - On-demand evaluation with iterators

### Supporting Packages

- [Accumulator](stream/accumulator/README.md) - Reduction operations for aggregating values
- [Iterator](stream/Iterator/README.md) - Lazy sequence generation for numeric types

## Complete Example

```go
package main

import (
    "fmt"
    "github.com/Ephram84/stream"
)

func main() {
    // Process a list of numbers
    result, err := stream.EagerFrom([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}).
        Filter(func(x int) (bool, error) { return x%2 == 0, nil }). // Keep even numbers
        Map(func(x int) (int, error) { return x * x, nil }).        // Square them
        Take(3).                                                    // Take first 3
        ToSlice()                                                   // Convert to slice
    
    if err != nil {
        panic(err)
    }
    
    fmt.Println(result) // [4, 16, 36]
}
```

## Performance Benchmarks

Performance comparison between eager (`slice`) and lazy (`sequence`) evaluation. All benchmarks run on Intel Pentium Gold 8505.

### 🚀 Performance Comparison (Lower is Better)

#### Filter Operation - Time (ns/op)
```
Small (100 elements)
  Slice    █████████████████████████ 1,239 ns
  Sequence ████████████████████████  1,141 ns  🏆 8% faster

Medium (10K elements)
  Slice    ████████████████████████████████████████████████████ 101,786 ns
  Sequence ████████████████████████████████████████████████     96,048 ns  🏆 6% faster

Large (1M elements)
  Slice    ████████████████████████████████████████████████████ 11,862,396 ns
  Sequence ███████████████████████████████████████████          9,376,324 ns  🏆 21% faster
```

#### Map Operation - Time (ns/op)
```
Small (100 elements)
  Slice    ████████████████████ 1,104 ns  🏆
  Sequence ████████████████████████████ 1,547 ns

Medium (10K elements)
  Slice    ████████████████████████████ 84,930 ns  🏆
  Sequence ████████████████████████████████████████████████████ 148,608 ns

Large (1M elements)
  Slice    ███████████████████████████████ 8,965,579 ns  🏆
  Sequence ████████████████████████████████████████████████████ 15,122,503 ns
```

#### Chained Operations (Filter → Map → Reduce) - Time (ns/op)
```
Small (100 elements)
  Slice    ████████████████████████████ 1,401 ns
  Sequence ███████████████████          947 ns  🏆 32% faster

Medium (10K elements)
  Slice    ████████████████████████████████████████████████████ 119,342 ns
  Sequence ███████████████████████████                           64,412 ns  🏆 46% faster

Large (1M elements)
  Slice    ████████████████████████████████████████████████████ 12,580,187 ns
  Sequence ██████████████████████████                            6,532,707 ns  🏆 48% faster
```

### 💾 Memory Efficiency (Lower is Better)

#### Filter Operation - Memory (Bytes/op)
```
Small (100 elements)
  Slice    ████████████████████████████████████████████████████ 2,424 B
  Sequence █████████████████████████                             1,152 B  💚 52% less

Medium (10K elements)
  Slice    ████████████████████████████████████████████████████ 251,225 B
  Sequence █████████████████████████                            128,384 B  💚 49% less

Large (1M elements)
  Slice    ████████████████████████████████████████████████████ 33,092,984 B
  Sequence █████████████████████████████████                    21,083,593 B  💚 36% less
```

#### Chained Operations (Filter → Map → Reduce) - Memory (Bytes/op)
```
Small (100 elements)
  Slice    ████████████████████████████████████████████████████ 2,472 B
  Sequence ▏                                                        176 B  💚 93% less

Medium (10K elements)
  Slice    ████████████████████████████████████████████████████ 251,273 B
  Sequence ▏                                                        176 B  💚 99.9% less

Large (1M elements)
  Slice    ████████████████████████████████████████████████████ 33,093,027 B
  Sequence ▏                                                           176 B  💚 99.9% less
```

### 📊 Key Insights

- **🎯 Single Operations**: Eager evaluation (`slice`) is typically faster for isolated operations like Map and Reduce
- **⚡ Chained Operations**: Lazy evaluation (`sequence`) significantly outperforms eager evaluation when chaining multiple operations
- **💾 Memory Efficiency**: Lazy evaluation uses dramatically less memory, especially for chained operations (up to 99.9% less!)
  - **Important**: Chained operations ending in `Reduce()` use less memory than single `Filter().ToSlice()` because `Reduce()` doesn't materialize intermediate results - it only keeps the accumulated value!
  - `ToSlice()` must allocate memory for all result elements, while `Reduce()` only needs memory for the final aggregated value
- **📊 Best Practice**: Use `slice` for simple transformations, `sequence` for complex pipelines and large datasets

## Documentation

The main functions and methods are documented in GoDoc comments in the code. More examples and tests can be found in the respective `*_test.go` files.

## Contributing

Pull requests, bug reports, and feature requests are welcome! Please follow the Go code style and add tests for new features whenever possible.

## License

MIT License
