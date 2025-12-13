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

### `Map[T, R any](s *seq[T], mapper func(elem T) (R, error))`
Maps a sequence to another type T -> R (lazy evaluation).

```go
numbers := sequence.From([]int{1, 2, 3, 4, 5}, nil)
result, err := sequence.MapSeq(numbers, func(elem int) (string, error) {
    return "Number: " + strconv.Itoa(elem), nil
}).ToSlice()
// Result: ["Number: 1", "Number: 2", "Number: 3", "Number: 4", "Number: 5"]
```

### `Zip[T1, T2 any](seq1 *seq[T1], seq2 *seq[T2])`
Combines two sequences into tuples (lazy evaluation).

```go
numbers := sequence.From([]int{1, 2, 3}, nil)
letters := sequence.From([]string{"a", "b", "c"}, nil)
zipped := sequence.ZipSeqs(numbers, letters)
// Result: [(1,"a"), (2,"b"), (3,"c")]
```

### `GroupBy[K common.Key, T any](seq *seq[T], keyMapper func(elem T) (K, error))`
Groups sequence elements by a key (evaluates the entire sequence).

```go
words := sequence.From([]string{"apple", "banana", "apricot"}, nil)
grouped := sequence.GroupBySeq(words, func(s string) (string, error) {
    return string(s[0]), nil
})
result, _ := grouped.ToMap()
// Result: map["a": ["apple", "apricot"], "b": ["banana"]]
```

## Creating Maps

There are two types of map operations in the sequence package:

1. **Simple Maps** (`pairs[K, V]`) - Maps with single values: `map[K]V`
2. **Maps with Slices** (`pairsSlice[K, V]`) - Maps with slice values: `map[K][]V`

### `FromMap[K common.Key, V any](m map[K]V)`
Creates a new sequence from an existing map with single values.

```go
// Simple map: map[string]int
simpleMap := map[string]int{
    "a": 1,
    "bb": 2,
    "ccc": 3,
}
pairs := sequence.FromMap(simpleMap)
```

### `FromMapWithSlices[K common.Key, V any](m map[K][]V)`
Creates a new sequence from an existing map with slice values.

```go
// Map with slices: map[string][]string
mapWithSlices := map[string][]string{
    "a": {"apple", "apricot"},
    "b": {"banana", "blueberry"},
    "c": {"cherry"},
}
pairsSlice := sequence.FromMapWithSlices(mapWithSlices)
```

## Map Operations

### Operations on Simple Maps (`pairs[K, V]`)

### `Filter(predicate predicate func(key K, value V) (bool, error))`
Filters each key-value pair based on a predicate.

```go
pairsMap := map[string]int{
	"key1":  1,
	"key2":  2,
	"key3":  3,
	"key11": 11,
}

result, _ := FromMap(pairsMap).Filter(func(key string, value int) (bool, error) {
	return strings.HasSuffix(key, "1") && value > 10, nil
}).ToMap()
// result = map["key11", 11]
```

#### `Flatten()`
Returns a sequence of all values from the map.

```go
values, _ := sequence.FromMap(map[string]int{
    "a": 1,
    "b": 2,
    "c": 3,
}).Flatten().ToSlice()
// values = [1, 2, 3] (order may vary)
```

#### `Keys()`
Returns a sequence of all keys from the map.

```go
keys, _ := sequence.FromMap(map[string]int{
    "a": 1,
    "b": 2,
    "c": 3,
}).Keys().ToSlice()
// keys = ["a", "b", "c"] (order may vary)
```

#### `ToMap()`
Converts the pairs back to a regular Go map.

```go
pairs := sequence.FromMap(map[string]int{
    "a": 1,
    "b": 2,
})
result, _ := pairs.ToMap()
// result = map["a": 1, "b": 2]
```

#### `Count()`
Returns the number of key-value pairs in the map.

```go
count, _ := sequence.FromMap(map[string]int{
    "a": 1,
    "b": 2,
    "c": 3,
}).Count()
// count = 3
```

#### `ForEach(action func(key K, value V))`
Executes a function for each key-value pair.

```go
pairs := sequence.FromMap(map[string]int{
    "a": 1,
    "b": 2,
    "c": 3,
})
err := pairs.ForEach(func(key string, value int) {
    fmt.Printf("Key: %s, Value: %d\n", key, value)
})
// Prints each key-value pair (order may vary)
```

### Operations on Maps with Slices (`pairsSlice[K, V]`)

#### `CountValues()`
Returns a simple pairs map with the count of elements for each key.

```go
result, _ := sequence.FromMapWithSlices(map[string][]string{
    "a": {"apple", "apricot", "avocado"},
    "b": {"banana", "blueberry"},
    "c": {"cherry"},
}).CountValues().ToMap()
// result = map["a": 3, "b": 2, "c": 1]
```

#### `Flatten()`
Returns a sequence of all values from all slices.

```go
values, _ := sequence.FromMapWithSlices(map[string][]int{
    "first":  {1, 2, 3},
    "second": {4, 5},
    "third":  {6},
}).Flatten().ToSlice()
// values = [1, 2, 3, 4, 5, 6] (order may vary)
```

#### `Keys()`
Returns a sequence of all keys from the map.

```go
keys, _ := sequence.FromMapWithSlices(map[string][]int{
    "first":  {1, 2, 3},
    "second": {4, 5},
    "third":  {6},
}).Keys().ToSlice()
// keys = ["first", "second", "third"] (order may vary)
```

#### `MapToFloat64(mapper func(elem V) (float64, error))`
Maps all values in all slices to float64.

```go
result, _ := sequence.FromMapWithSlices(map[string][]int{
    "group1": {1, 2, 3},
    "group2": {4, 5},
}).MapToFloat64(func(x int) (float64, error) {
    return float64(x) * 1.5, nil
}).ToMap()
// result = map["group1": [1.5, 3.0, 4.5], "group2": [6.0, 7.5]]
```

#### `Reduce(accumulator accumulator.AccumulatorSeq[V])`
Reduces each slice in the map to a single value, converting to a simple pairs map.

```go
result, _ := sequence.FromMapWithSlices(map[string][]int{
    "group1": {1, 2, 3},
    "group2": {4, 5, 6},
    "group3": {7, 8},
}).Reduce(accumulator.SumSeq).ToMap()
// result = map["group1": 6, "group2": 15, "group3": 15]
```

#### `ToMap()`
Converts the pairsSlice back to a regular Go map with slices.

```go
pairsSlice := sequence.FromMapWithSlices(map[string][]int{
    "first":  {1, 2, 3},
    "second": {4, 5, 6},
})
result, _ := pairsSlice.ToMap()
// result = map["first": [1, 2, 3], "second": [4, 5, 6]]
```

## Complete Map Examples

### Processing Simple Maps
```go
// Create a map of user scores
scores := map[string]int{
    "Alice": 95,
    "Bob":   87,
    "Carol": 92,
}

// Get all scores above 90
pairs := sequence.FromMap(scores)
highScores, _ := pairs.Flatten().
    Filter(func(score int) (bool, error) {
        return score > 90, nil
    }).
    ToSlice()
// highScores = [95, 92] (order may vary)

// Get names of high scorers
highScorerNames := make([]string, 0)
pairs.ForEach(func(name string, score int) {
    if score > 90 {
        highScorerNames = append(highScorerNames, name)
    }
})
// highScorerNames = ["Alice", "Carol"] (order may vary)
```

### Processing Maps with Slices
```go
// Create a map of students by grade
studentsByGrade := map[string][]string{
    "A": {"Alice", "Alex"},
    "B": {"Bob", "Betty", "Ben"},
    "C": {"Carol"},
}

// Count students in each grade
pairsSlice := sequence.FromMapWithSlices(studentsByGrade)
counts := pairsSlice.CountValues()
countMap, _ := counts.ToMap()
// countMap = map["A": 2, "B": 3, "C": 1]

// Get all student names
allStudents, _ := pairsSlice.Flatten().ToSlice()
// allStudents = ["Alice", "Alex", "Bob", "Betty", "Ben", "Carol"] (order may vary)

// Get all grades
grades, _ := pairsSlice.Keys().ToSlice()
// grades = ["A", "B", "C"] (order may vary)
```

### Converting Between Map Types
```go
// Start with a map of slices
groupedData := map[int][]int{
    1: {10, 20, 30},
    2: {40, 50},
    3: {60, 70, 80, 90},
}

// Reduce each group to its sum
pairsSlice := sequence.FromMapWithSlices(groupedData)
sums := pairsSlice.Reduce(accumulator.SumSeq)

// Now we have a simple map
sumMap, _ := sums.ToMap()
// sumMap = map[1: 60, 2: 90, 3: 300]

// Get the maximum sum
maxSum, _ := sums.Flatten().Reduce(accumulator.MaxSeq)
// maxSum = 300
```