# Iterator - Lazy Sequence Generation

The `stream/iterator` package provides lazy sequence generation for numeric types. Unlike eager evaluation in the `slice` package, iterators generate values on-demand only when `Generate()` is called.

**Key Features:**
- Lazy evaluation: values are generated only when needed
- Memory efficient for large sequences
- Supports numeric types (int, int64, float64)
- Configurable skip and limit operations
- Custom sequence generation functions

**Recommended Use Cases:**
- Generating large numeric sequences without loading all values into memory upfront
- Creating infinite or very long sequences with limits
- Memory-constrained environments
- When you need only a subset of a sequence (skip/take patterns)

## Creating Iterators

### `Iterator[N common.Numbers](start N, next NextFunc[N])`
Creates a new iterator starting at the given value with a custom next function.

```go
// Create an iterator starting at 1 with increment by 1
iter := iterator.Iterator(1, iterator.IncrementInt())

// Create an iterator with custom logic (multiply by 2)
iter := iterator.Iterator(1, func(current int) int {
    return current * 2
})
```

## Configuration Operations

### `WithSkip(skip int)`
Skips the first n elements in the sequence.

```go
// Generate sequence starting from element 5: [5, 6, 7, 8, 9]
result := iterator.Iterator(0, iterator.IncrementInt()).
    WithSkip(5).
    WithLimit(5).
    Generate()
// result = [5, 6, 7, 8, 9]
```

### `WithLimit(limit int)`
Limits the number of elements to generate. Default limit is 100.

```go
// Generate first 10 elements
result := iterator.Iterator(1, iterator.IncrementInt()).
    WithLimit(10).
    Generate()
// result = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]
```

## Terminal Operation

### `Generate()`
Produces a slice containing the generated sequence. This is when the actual generation happens.

```go
iter := iterator.Iterator(1, iterator.IncrementInt()).
    WithLimit(5)
    
result := iter.Generate()
// result = [1, 2, 3, 4, 5]
```

## Built-in Next Functions

### `IncrementInt()`
Returns a NextFunc that increments an integer by 1.

```go
iter := iterator.Iterator(0, iterator.IncrementInt()).
    WithLimit(5)
    
result := iter.Generate()
// result = [0, 1, 2, 3, 4]
```

### `IncrementFloat()`
Returns a NextFunc that increments a float64 by 1.0.

```go
iter := iterator.Iterator(0.0, iterator.IncrementFloat()).
    WithLimit(5)
    
result := iter.Generate()
// result = [0.0, 1.0, 2.0, 3.0, 4.0]
```

## Custom Next Functions

You can define custom sequence generators for various patterns:

### Arithmetic Progression
```go
// Generate even numbers: 2, 4, 6, 8, 10
iter := iterator.Iterator(2, func(current int) int {
    return current + 2
}).WithLimit(5)

result := iter.Generate()
// result = [2, 4, 6, 8, 10]
```

### Geometric Progression
```go
// Generate powers of 2: 1, 2, 4, 8, 16
iter := iterator.Iterator(1, func(current int) int {
    return current * 2
}).WithLimit(5)

result := iter.Generate()
// result = [1, 2, 4, 8, 16]
```

### Fibonacci Sequence
```go
// Note: For Fibonacci, you'd need to track state externally
// or use a closure with mutable state
prev := 0
iter := iterator.Iterator(1, func(current int) int {
    next := prev + current
    prev = current
    return next
}).WithLimit(10)

result := iter.Generate()
// result = [1, 1, 2, 3, 5, 8, 13, 21, 34, 55]
```

### Countdown
```go
// Generate countdown: 10, 9, 8, 7, 6
iter := iterator.Iterator(10, func(current int) int {
    return current - 1
}).WithLimit(5)

result := iter.Generate()
// result = [10, 9, 8, 7, 6]
```

### Float Sequences
```go
// Generate sequence with decimal increment: 0.0, 0.5, 1.0, 1.5, 2.0
iter := iterator.Iterator(0.0, func(current float64) float64 {
    return current + 0.5
}).WithLimit(5)

result := iter.Generate()
// result = [0.0, 0.5, 1.0, 1.5, 2.0]
```

## Complete Examples

### Generate Range with Skip
```go
// Generate numbers from 100 to 109, skipping first 100
iter := iterator.Iterator(0, iterator.IncrementInt()).
    WithSkip(100).
    WithLimit(10)
    
result := iter.Generate()
// result = [100, 101, 102, 103, 104, 105, 106, 107, 108, 109]
```

### Large Sequence with Memory Efficiency
```go
// Generate first 1000 even numbers
iter := iterator.Iterator(0, func(current int) int {
    return current + 2
}).WithLimit(1000)

// Values are only generated when Generate() is called
result := iter.Generate()
// result = [0, 2, 4, 6, ..., 1998]
```

### Combining with Slice Operations
```go
// Generate sequence and then apply slice operations
import "github.com/Ephram84/stream/slice"

numbers := iterator.Iterator(1, iterator.IncrementInt()).
    WithLimit(100).
    Generate()

// Now use slice operations
result, _ := slice.From(numbers).
    Filter(func(x int) (bool, error) {
        return x%2 == 0, nil
    }).
    Take(10).
    ToSlice()
// result = [2, 4, 6, 8, 10, 12, 14, 16, 18, 20]
```

## Iterator vs Slice Package

**Use Iterator when:**
- You need to generate large numeric sequences
- Memory efficiency is important
- You want lazy evaluation
- You're working with numeric types only

**Use Slice when:**
- You already have data in memory
- You need complex transformations (map, filter, etc.)
- You're working with any data type
- You need eager evaluation and immediate results

## Performance Considerations

- Iterators are memory-efficient for generation but still produce a slice at the end
- The actual generation happens only when `Generate()` is called
- Configuration methods (`WithSkip`, `WithLimit`) are zero-cost until generation
- For very large sequences, consider processing in chunks rather than generating all at once
