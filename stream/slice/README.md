# Eager Processing

The `stream/slice` package provides eager, in-memory operations on Go slices.

**Key Features:**
- Eager evaluation: all data is loaded and processed immediately
- Direct manipulation of Go slices
- Useful for batch operations, sorting, and when you need all elements at once
- Includes helpers to create slices from arrays, files, etc.

**Recommended Use Cases:**
- You already have a Go slice and want to perform multiple operations on it
- You want to load all data from a file or another source into memory for processing
- You need to sort, map, or filter all elements at once
- Working with datasets that fit comfortably in memory

## Creating Slices

### `From[T any](tokens []T)`
Creates a new stream from an existing slice.

```go
numbers := slice.From([]int{1, 2, 3, 4, 5})
words := slice.From([]string{"hello", "world", "go"})
```

### `FromFile(path string)`
Creates a stream of strings from a file (splits by whitespace).

```go
// Assuming a file "data.txt" contains: "apple banana cherry"
words := slice.FromFile("data.txt")
result, _ := words.ToSlice() // ["apple", "banana", "cherry"]
```

## Transformation Operations

### `Filter(filter func(elem T) (bool, error))`
Filters elements based on a predicate.

```go
numbers := slice.From([]int{1, 2, 3, 4, 5, 6})
evens := numbers.Filter(func(x int) (bool, error) {
    return x%2 == 0, nil
})
result, _ := evens.ToSlice() // [2, 4, 6]
```

### `Map(mapper func(elem T) (T, error))`
Transforms each element using a mapper function.

```go
numbers := slice.From([]int{1, 2, 3, 4})
squared := numbers.Map(func(x int) (int, error) {
    return x * x, nil
})
result, _ := squared.ToSlice() // [1, 4, 9, 16]
```

### `MapToInt(mapper func(elem T) (int, error))`
Maps elements to integers.

```go
words := slice.From([]string{"a", "bb", "ccc"})
lengths := words.MapToInt(func(s string) (int, error) {
    return len(s), nil
})
result, _ := lengths.ToSlice() // [1, 2, 3]
```

### `MapToString(mapper func(elem T) (string, error))`
Maps elements to strings.

```go
numbers := slice.From([]int{1, 2, 3})
strings := numbers.MapToString(func(x int) (string, error) {
    return fmt.Sprintf("num_%d", x), nil
})
result, _ := strings.ToSlice() // ["num_1", "num_2", "num_3"]
```

### `FlatMapSlice[T1, T2 any](s *slice[T1], mapper func(elem T1) ([]T2, error))`
Flattens nested slices.

```go
words := slice.From([]string{"hello", "world"})
chars := slice.FlatMapSlice(words, func(s string) ([]string, error) {
    return strings.Split(s, ""), nil
})
result, _ := chars.ToSlice() // ["h", "e", "l", "l", "o", "w", "o", "r", "l", "d"]
```

## Ordering Operations

### `Sort(sortFunc func(slice []T) func(i, j int) bool)`
Sorts elements using a custom comparator.

```go
numbers := slice.From([]int{5, 2, 8, 1, 9})
sorted := numbers.Sort(func(slice []int) func(i, j int) bool {
    return func(i, j int) bool {
        return slice[i] < slice[j]
    }
})
result, _ := sorted.ToSlice() // [1, 2, 5, 8, 9]
```

### `Reverse()`
Reverses the order of elements.

```go
numbers := slice.From([]int{1, 2, 3, 4, 5})
reversed := numbers.Reverse()
result, _ := reversed.ToSlice() // [5, 4, 3, 2, 1]
```

## Slicing Operations

### `Take(n int)`
Takes the first n elements.

```go
numbers := slice.From([]int{1, 2, 3, 4, 5})
first3 := numbers.Take(3)
result, _ := first3.ToSlice() // [1, 2, 3]
```

### `Skip(n int)`
Skips the first n elements.

```go
numbers := slice.From([]int{1, 2, 3, 4, 5})
afterFirst2 := numbers.Skip(2)
result, _ := afterFirst2.ToSlice() // [3, 4, 5]
```

## Set Operations

### `Distinct(eq func(a, b T) bool)`
Removes duplicate elements.

```go
numbers := slice.From([]int{1, 2, 2, 3, 3, 3, 4})
unique := numbers.Distinct(func(a, b int) bool {
    return a == b
})
result, _ := unique.ToSlice() // [1, 2, 3, 4]
```

### `Concat(slices ...*slice[T])`
Concatenates multiple streams.

```go
slice1 := slice.From([]int{1, 2, 3})
slice2 := slice.From([]int{4, 5, 6})
slice3 := slice.From([]int{7, 8, 9})
combined := slice1.Concat(slice2, slice3)
result, _ := combined.ToSlice() // [1, 2, 3, 4, 5, 6, 7, 8, 9]
```

## Grouping Operations

### `GroupByString(grouper func(elem T) string)`
Groups elements by string keys.

```go
words := slice.From([]string{"apple", "banana", "apricot", "blueberry"})
grouped := words.GroupByString(func(s string) string {
    return string(s[0]) // Group by first letter
})
// Result: map["a": ["apple", "apricot"], "b": ["banana", "blueberry"]]
```

### `PartitioningBy(predicate func(elem T) (bool, error))`
Partitions elements into two groups based on a predicate.

```go
numbers := slice.From([]int{1, 2, 3, 4, 5, 6})
partitioned := numbers.PartitioningBy(func(x int) (bool, error) {
    return x%2 == 0, nil
})
// Result: map[true: [2, 4, 6], false: [1, 3, 5]]
```

### `AssociateByString(mapper func(elem T) (string, error))`
Creates a map with string keys.

```go
words := slice.From([]string{"cat", "dog", "elephant"})
associated := words.AssociateByString(func(s string) (string, error) {
    return string(s[0]), nil
})
// Result: map["c": "cat", "d": "dog", "e": "elephant"]
```

## Terminal Operations

### `ToSlice() ([]T, error)`
Converts the stream back to a regular Go slice.

```go
stream := slice.From([]int{1, 2, 3}).Filter(func(x int) (bool, error) {
    return x > 1, nil
})
result, err := slice.ToSlice() // [2, 3], nil
```

### `Count() (int, error)`
Returns the number of elements.

```go
numbers := slice.From([]int{1, 2, 3, 4, 5})
count, _ := numbers.Count() // 5
```

### `First(orElse ...T) (*T, error)`
Returns the first element or a default value.

```go
numbers := slice.From([]int{1, 2, 3})
first, _ := numbers.First() // *1

empty := slice.From([]int{})
firstOrDefault, _ := empty.First(99) // *99
```

### `Last(orElse ...T) (*T, error)`
Returns the last element or a default value.

```go
numbers := slice.From([]int{1, 2, 3})
last, _ := numbers.Last() // *3
```

### `FirstOrNil(predicate func(elem T) (bool, error)) (*T, error)`
Returns the first element matching a predicate.

```go
numbers := slice.From([]int{1, 2, 3, 4, 5})
firstEven, _ := numbers.FirstOrNil(func(x int) (bool, error) {
    return x%2 == 0, nil
}) // *2
```

### `LastOrNil(predicate func(elem T) (bool, error)) (*T, error)`
Returns the last element matching a predicate.

```go
numbers := slice.From([]int{1, 2, 3, 4, 5})
lastEven, _ := numbers.LastOrNil(func(x int) (bool, error) {
    return x%2 == 0, nil
}) // *4
```

## Matching Operations

### `AnyMatch(predicate func(elem T) (bool, error)) (bool, error)`
Checks if any element matches the predicate.

```go
numbers := slice.From([]int{1, 3, 5, 7})
hasEven, _ := numbers.AnyMatch(func(x int) (bool, error) {
    return x%2 == 0, nil
}) // false
```

### `AllMatch(predicate func(elem T) (bool, error)) (bool, error)`
Checks if all elements match the predicate.

```go
numbers := slice.From([]int{2, 4, 6, 8})
allEven, _ := numbers.AllMatch(func(x int) (bool, error) {
    return x%2 == 0, nil
}) // true
```

### `NoneMatch(predicate func(elem T) (bool, error)) (bool, error)`
Checks if no elements match the predicate.

```go
numbers := slice.From([]int{1, 3, 5, 7})
noneEven, _ := numbers.NoneMatch(func(x int) (bool, error) {
    return x%2 == 0, nil
}) // true
```

## Reduction Operations

### `Reduce(accumulator accumulator.Accumulator[T]) (T, error)`
Reduces the stream to a single value. See [accumulator](stream/accumulator/README.md#1-batch-accumulator-accumulatort-any)

```go
numbers := slice.From([]int{1, 2, 3, 4, 5})
sum, _ := numbers.Reduce(accumulator.Sum)
// sum = 15
```

### `ForEach(consumer func(elem T) error) error`
Executes a function for each element.

```go
numbers := slice.From([]int{1, 2, 3})
err := numbers.ForEach(func(x int) error {
    fmt.Printf("Number: %d\n", x)
    return nil
})
// Prints: Number: 1, Number: 2, Number: 3
```

## Free Funcitons
Unfortunately, Go does not allow something like `func (s *slice[T]) Map[R any](mapper func(elem T) (R, error))`. Therefore, there are a few helper functions, where T becomes R

### `MapSlice[T, R any](s *slice[T], mapper func(elem T) (R, error))`
Maps a slice to another type T -> R
```go
numbers := EagerFrom([]int{1, 2, 3, 4, 5})
result, err := MapSlice(numbers, func(elem int) (string, error) {
	return "Number: " + string(rune('0'+elem)), nil
}).ToSlice()
// Result: ["Number: 1", "Number: 2", "Number: 3", "Number: 4", "Number: 5"]
```

### `ZipSlices[T1, T2 any](s1 *slice[T1], s2 *slice[T2])`
Combines two streams into tuples.

```go
numbers := slice.From([]int{1, 2, 3})
letters := slice.From([]string{"a", "b", "c"})
zipped := stream.ZipSlices(numbers, letters)
// Result: [(1,"a"), (2,"b"), (3,"c")]
```

## Utility Operations

### `Write(writer io.Writer) (int, error)`
Writes the stream as JSON to a writer.

```go
numbers := slice.From([]int{1, 2, 3})
var buf bytes.Buffer
bytesWritten, err := numbers.Write(&buf)
// buf contains: [1,2,3]
```