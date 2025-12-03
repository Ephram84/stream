# Lazy Processing

The `stream/sequence` package provides lazy, iterator-based processing for Go data.

**Key Features:**
- Lazy evaluation: elements are produced and processed only as needed
- Iterator-based API, using a `next()` function
- Supports chaining of operations like `Filter`, `Map`, `Take`, etc.
- Efficient for large or infinite data streams

**Recommended Use Cases:**
- You want to process data on-demand, not all at once
- You are working with large or infinite data sources
- You want to build a pipeline of operations that are only executed when iterating
- Memory efficiency is important

## Creating Sequences

### `From[T any](slice []T, errs ...error)`
Creates a new lazy sequence from an existing slice.

```go
numbers := sequence.From([]int{1, 2, 3, 4, 5}, nil)
words := sequence.From([]string{"hello", "world", "go"}, nil)
```

## Transformation Operations

### `Filter(filter func(elem T) (bool, error))`
Filters elements based on a predicate (lazy evaluation).

```go
numbers := sequence.From([]int{1, 2, 3, 4, 5, 6}, nil)
evens := numbers.Filter(func(x int) (bool, error) {
    return x%2 == 0, nil
})
result, _ := evens.ToSlice() // [2, 4, 6]
```

### `Map(mapper func(elem T) (T, error))`
Transforms each element using a mapper function (lazy evaluation).

```go
numbers := sequence.From([]int{1, 2, 3, 4}, nil)
squared := numbers.Map(func(x int) (int, error) {
    return x * x, nil
})
result, _ := squared.ToSlice() // [1, 4, 9, 16]
```

### `MapToInt(mapper func(elem T) (int, error))`
Maps elements to integers (lazy evaluation).

```go
words := sequence.From([]string{"a", "bb", "ccc"}, nil)
lengths := words.MapToInt(func(s string) (int, error) {
    return len(s), nil
})
result, _ := lengths.ToSlice() // [1, 2, 3]
```

### `MapToInt64(mapper func(elem T) (int64, error))`
Maps elements to int64 (lazy evaluation).

```go
numbers := sequence.From([]int{1, 2, 3}, nil)
int64s := numbers.MapToInt64(func(x int) (int64, error) {
    return int64(x * 1000), nil
})
result, _ := int64s.ToSlice() // [1000, 2000, 3000]
```

### `MapToFloat64(mapper func(elem T) (float64, error))`
Maps elements to float64 (lazy evaluation).

```go
numbers := sequence.From([]int{1, 2, 3}, nil)
floats := numbers.MapToFloat64(func(x int) (float64, error) {
    return float64(x) / 2.0, nil
})
result, _ := floats.ToSlice() // [0.5, 1.0, 1.5]
```

### `MapToString(mapper func(elem T) (string, error))`
Maps elements to strings (lazy evaluation).

```go
numbers := sequence.From([]int{1, 2, 3}, nil)
strings := numbers.MapToString(func(x int) (string, error) {
    return fmt.Sprintf("num_%d", x), nil
})
result, _ := strings.ToSlice() // ["num_1", "num_2", "num_3"]
```

### `FlatMapSeq[T, R any](s *seq[T], mapper func(elem T) (*seq[R], error))`
Flattens nested sequences (lazy evaluation).

```go
words := sequence.From([]string{"hello", "world"}, nil)
chars := sequence.FlatMapSeq(words, func(s string) (*sequence.seq[string], error) {
    return sequence.From(strings.Split(s, ""), nil), nil
})
result, _ := chars.ToSlice() // ["h", "e", "l", "l", "o", "w", "o", "r", "l", "d"]
```

## Slicing Operations

### `Take(n int)`
Takes the first n elements (lazy evaluation).

```go
numbers := sequence.From([]int{1, 2, 3, 4, 5}, nil)
first3 := numbers.Take(3)
result, _ := first3.ToSlice() // [1, 2, 3]
```

### `Skip(n int)`
Skips the first n elements (lazy evaluation).

```go
numbers := sequence.From([]int{1, 2, 3, 4, 5}, nil)
afterFirst2 := numbers.Skip(2)
result, _ := afterFirst2.ToSlice() // [3, 4, 5]
```

## Set Operations

### `Distinct()`
Removes duplicate elements (lazy evaluation).

```go
numbers := sequence.From([]int{1, 2, 2, 3, 3, 3, 4}, nil)
unique := numbers.Distinct()
result, _ := unique.ToSlice() // [1, 2, 3, 4]
```

### `Concat(sequences ...*seq[T])`
Concatenates multiple sequences (lazy evaluation).

```go
seq1 := sequence.From([]int{1, 2, 3}, nil)
seq2 := sequence.From([]int{4, 5, 6}, nil)
seq3 := sequence.From([]int{7, 8, 9}, nil)
combined := seq1.Concat(seq2, seq3)
result, _ := combined.ToSlice() // [1, 2, 3, 4, 5, 6, 7, 8, 9]
```

### `Reverse()`
Reverses the order of elements (requires full evaluation).

```go
numbers := sequence.From([]int{1, 2, 3, 4, 5}, nil)
reversed := numbers.Reverse()
result, _ := reversed.ToSlice() // [5, 4, 3, 2, 1]
```

## Grouping Operations

### `GroupByString(grouper func(elem T) (string, error))`
Groups elements by string keys.

```go
words := sequence.From([]string{"apple", "banana", "apricot", "blueberry"}, nil)
grouped := words.GroupByString(func(s string) (string, error) {
    return string(s[0]), nil // Group by first letter
})
result, _ := grouped.ToMap()
// Result: map["a": ["apple", "apricot"], "b": ["banana", "blueberry"]]
```

### `GroupByInt(grouper func(elem T) (int, error))`
Groups elements by int keys.

```go
numbers := sequence.From([]int{1, 2, 3, 4, 5, 6}, nil)
grouped := numbers.GroupByInt(func(x int) (int, error) {
    return x % 3, nil // Group by modulo 3
})
result, _ := grouped.ToMap()
// Result: map[0: [3, 6], 1: [1, 4], 2: [2, 5]]
```

### `GroupByInt64(grouper func(elem T) (int64, error))`
Groups elements by int64 keys.

```go
numbers := sequence.From([]int{10, 20, 30, 15, 25}, nil)
grouped := numbers.GroupByInt64(func(x int) (int64, error) {
    return int64(x / 10), nil // Group by tens
})
result, _ := grouped.ToMap()
// Result: map[1: [10, 15], 2: [20, 25], 3: [30]]
```

### `GroupByFloat(grouper func(elem T) (float64, error))`
Groups elements by float64 keys.

```go
numbers := sequence.From([]float64{1.1, 1.9, 2.1, 2.9}, nil)
grouped := numbers.GroupByFloat(func(x float64) (float64, error) {
    return math.Floor(x), nil // Group by integer part
})
result, _ := grouped.ToMap()
// Result: map[1.0: [1.1, 1.9], 2.0: [2.1, 2.9]]
```

### `PartitioningBy(predicate func(elem T) (bool, error))`
Partitions elements into two groups based on a predicate.

```go
numbers := sequence.From([]int{1, 2, 3, 4, 5, 6}, nil)
partitioned := numbers.PartitioningBy(func(x int) (bool, error) {
    return x%2 == 0, nil
})
result, _ := partitioned.ToMap()
// Result: map[true: [2, 4, 6], false: [1, 3, 5]]
```

### `AssociateByString(mapper func(elem T) (string, error))`
Creates a map with string keys.

```go
words := sequence.From([]string{"cat", "dog", "elephant"}, nil)
associated := words.AssociateByString(func(s string) (string, error) {
    return string(s[0]), nil
})
result, _ := associated.ToMap()
// Result: map["c": "cat", "d": "dog", "e": "elephant"]
```

### `AssociateByInt(mapper func(elem T) (int, error))`
Creates a map with int keys.

```go
words := sequence.From([]string{"a", "bb", "ccc"}, nil)
associated := words.AssociateByInt(func(s string) (int, error) {
    return len(s), nil
})
result, _ := associated.ToMap()
// Result: map[1: "a", 2: "bb", 3: "ccc"]
```

### `AssociateByInt64(mapper func(elem T) (int64, error))`
Creates a map with int64 keys.

```go
numbers := sequence.From([]int{10, 20, 30}, nil)
associated := numbers.AssociateByInt64(func(x int) (int64, error) {
    return int64(x), nil
})
result, _ := associated.ToMap()
// Result: map[10: 10, 20: 20, 30: 30]
```

### `AssociateByFloat(mapper func(elem T) (float64, error))`
Creates a map with float64 keys.

```go
numbers := sequence.From([]float64{1.5, 2.5, 3.5}, nil)
associated := numbers.AssociateByFloat(func(x float64) (float64, error) {
    return x * 2, nil
})
result, _ := associated.ToMap()
// Result: map[3.0: 1.5, 5.0: 2.5, 7.0: 3.5]
```

## Terminal Operations

### `ToSlice() ([]T, error)`
Converts the sequence to a regular Go slice (evaluates the entire sequence).

```go
seq := sequence.From([]int{1, 2, 3}, nil).Filter(func(x int) (bool, error) {
    return x > 1, nil
})
result, err := seq.ToSlice() // [2, 3], nil
```

### `Count() (int, error)`
Returns the number of elements (evaluates the entire sequence).

```go
numbers := sequence.From([]int{1, 2, 3, 4, 5}, nil)
count, _ := numbers.Count() // 5
```

### `First(orElse ...T) (*T, error)`
Returns the first element or a default value.

```go
numbers := sequence.From([]int{1, 2, 3}, nil)
first, _ := numbers.First() // *1

empty := sequence.From([]int{}, nil)
firstOrDefault, _ := empty.First(99) // *99
```

### `FirstOrNil(predicate func(elem T) (bool, error)) (*T, error)`
Returns the first element matching a predicate.

```go
numbers := sequence.From([]int{1, 2, 3, 4, 5}, nil)
firstEven, _ := numbers.FirstOrNil(func(x int) (bool, error) {
    return x%2 == 0, nil
}) // *2
```

## Matching Operations

### `AnyMatch(predicate func(elem T) (bool, error)) (bool, error)`
Checks if any element matches the predicate (short-circuits on first match).

```go
numbers := sequence.From([]int{1, 3, 5, 7}, nil)
hasEven, _ := numbers.AnyMatch(func(x int) (bool, error) {
    return x%2 == 0, nil
}) // false
```

### `AllMatch(predicate func(elem T) (bool, error)) (bool, error)`
Checks if all elements match the predicate (short-circuits on first non-match).

```go
numbers := sequence.From([]int{2, 4, 6, 8}, nil)
allEven, _ := numbers.AllMatch(func(x int) (bool, error) {
    return x%2 == 0, nil
}) // true
```

### `NoneMatch(predicate func(elem T) (bool, error)) (bool, error)`
Checks if no elements match the predicate (short-circuits on first match).

```go
numbers := sequence.From([]int{1, 3, 5, 7}, nil)
noneEven, _ := numbers.NoneMatch(func(x int) (bool, error) {
    return x%2 == 0, nil
}) // true
```

## Reduction Operations

### `Reduce(accumulator accumulator.AccumulatorSeq[T]) (T, error)`
Reduces the sequence to a single value. See [accumulator](stream/accumulator/README.md#2-sequential-accumulator-accumulatorseqt-any)

```go
numbers := sequence.From([]int{1, 2, 3, 4, 5}, nil)
sum, _ := numbers.Reduce(accumulator.SumSeq)
// sum = 15
```

### `ForEach(consumer func(elem T) error) error`
Executes a function for each element.

```go
numbers := sequence.From([]int{1, 2, 3}, nil)
err := numbers.ForEach(func(x int) error {
    fmt.Printf("Number: %d\n", x)
    return nil
})
// Prints: Number: 1, Number: 2, Number: 3
```

## Free Functions
Unfortunately, Go does not allow something like `func (s *seq[T]) Map[R any](mapper func(elem T) (R, error))`. Therefore, there are a few helper functions, where T becomes R

### `MapSeq[T, R any](s *seq[T], mapper func(elem T) (R, error))`
Maps a sequence to another type T -> R (lazy evaluation).

```go
numbers := sequence.From([]int{1, 2, 3, 4, 5}, nil)
result, err := sequence.MapSeq(numbers, func(elem int) (string, error) {
    return "Number: " + strconv.Itoa(elem), nil
}).ToSlice()
// Result: ["Number: 1", "Number: 2", "Number: 3", "Number: 4", "Number: 5"]
```

### `ZipSeqs[T1, T2 any](seq1 *seq[T1], seq2 *seq[T2])`
Combines two sequences into tuples (lazy evaluation).

```go
numbers := sequence.From([]int{1, 2, 3}, nil)
letters := sequence.From([]string{"a", "b", "c"}, nil)
zipped := sequence.ZipSeqs(numbers, letters)
// Result: [(1,"a"), (2,"b"), (3,"c")]
```

### `GroupBySeq[K common.Key, T any](seq *seq[T], keyMapper func(elem T) (K, error))`
Groups sequence elements by a key (evaluates the entire sequence).

```go
words := sequence.From([]string{"apple", "banana", "apricot"}, nil)
grouped := sequence.GroupBySeq(words, func(s string) (string, error) {
    return string(s[0]), nil
})
result, _ := grouped.ToMap()
// Result: map["a": ["apple", "apricot"], "b": ["banana"]]
```