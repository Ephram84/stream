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

## Documentation

The main functions and methods are documented in GoDoc comments in the code. More examples and tests can be found in the respective `*_test.go` files.

## Contributing

Pull requests, bug reports, and feature requests are welcome! Please follow the Go code style and add tests for new features whenever possible.

## License

MIT License
