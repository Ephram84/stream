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
words, _ := slice.FromFile("data.txt").ToSlice()
// words = ["apple", "banana", "cherry"]
```

## Transformation Operations

### `Filter(filter func(elem T) (bool, error))`
Filters elements based on a predicate.

```go
evens, _ := slice.From([]int{1, 2, 3, 4, 5, 6}).Filter(func(x int) (bool, error) {
    return x%2 == 0, nil
}).ToSlice()
// evens = [2, 4, 6]
```

### `Map(mapper func(elem T) (T, error))`
Transforms each element using a mapper function.

```go
result, _ := slice.From([]int{1, 2, 3, 4}).Map(func(x int) (int, error) {
    return x * x, nil
}).ToSlice()
// result = [1, 4, 9, 16]
```

### `MapToInt(mapper func(elem T) (int, error))`
Maps elements to integers.

```go
result, _ := slice.From([]string{"a", "bb", "ccc"}).MapToInt(func(s string) (int, error) {
    return len(s), nil
}).ToSlice()
// result = [1, 2, 3]
```

### `MapToString(mapper func(elem T) (string, error))`
Maps elements to strings.

```go
result, _ := slice.From([]int{1, 2, 3}).MapToString(common.IntToString).ToSlice()
// result = ["1", "2", "3"]
```

### `FlatMapSlice[T1, T2 any](s *slice[T1], mapper func(elem T1) ([]T2, error))`
Flattens nested slices.

```go
words := slice.From([]string{"hello", "world"})
result, _ := slice.FlatMapSlice(words, func(s string) ([]string, error) {
    return strings.Split(s, ""), nil
}).ToSlice()
// result = ["h", "e", "l", "l", "o", "w", "o", "r", "l", "d"]
```

```go
    intSlice := [][]int{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	}

	result, _ := FlatMapSlice(From(intSlice), func(elem []int) ([]int, error) {
		return elem, nil
	}).ToSlice()
    // result = [1, 2, 3, 4, 5, 6, 7, 8, 9]
```

## Ordering Operations

### `Sort(sortFunc func(slice []T) func(i, j int) bool)`
Sorts elements using a custom comparator.

```go
result, _ := slice.From([]int{5, 2, 8, 1, 9}).Sort(common.Sort).ToSlice()
// result = [1, 2, 5, 8, 9]
```

Sorts desc.
```go
result, _ := slice.From([]int{5, 2, 8, 1, 9}).Sort(common.SortDesc).ToSlice()
// result = [9, 8, 5, 2, 1]
```

Custom sort function.
```go
transactions := []Transaction{
		{
			ID:          "T02",
			Amount:      20.0,
			BookingDate: 1759356000, // 2025-10-02
		},
		{
			ID:          "T03",
			Amount:      30.0,
			BookingDate: 1756764000, // 2025-09-02
		},
		{
			ID:          "T01",
			Amount:      10.0,
			BookingDate: 1762038000, // 2025-11-02
		},
	}

	sortTransaction, _ := slice.From(transactions).Sort(func(tnxs []Transaction) func(i, j int) bool {
		return func(i, j int) bool {
			return tnxs[i].BookingDate < tnxs[j].BookingDate
		}
	}).ToSlice()
    // sortTransactions = []Transaction{
	// 	{
	// 		ID:          "T01",
	// 		Amount:      10.0,
	// 		BookingDate: 1762038000, // 2025-11-02
	// 	},
	// 	{
	// 		ID:          "T02",
	// 		Amount:      20.0,
	// 		BookingDate: 1759356000, // 2025-10-02
	// 	},
	// 	{
	// 		ID:          "T03",
	// 		Amount:      30.0,
	// 		BookingDate: 1756764000, // 2025-09-02
	// 	},
	// }
```


### `Reverse()`
Reverses the order of elements.

```go
result, _ := slice.From([]int{1, 2, 3, 4, 5}).Reverse().ToSlice()
// result = [5, 4, 3, 2, 1]
```

## Slicing Operations

### `Take(n int)`
Takes the first n elements.

```go
result,  := slice.From([]int{1, 2, 3, 4, 5}).Take(3).ToSlice()
// result = [1, 2, 3]
```

### `Skip(n int)`
Skips the first n elements.

```go
result, _ := slice.From([]int{1, 2, 3, 4, 5}).Skip(2).ToSlice()
// result = [3, 4, 5]
```

## Set Operations

### `Distinct(eq func(a, b T) bool)`
Removes duplicate elements.

```go
result, _ := slice.From([]int{1, 2, 2, 3, 3, 3, 4}).Distinct(common.Eq).ToSlice()
// result = [1, 2, 3, 4]
```

### `Concat(slices ...*slice[T])`
Concatenates multiple streams.

```go
slice1 := slice.From([]int{1, 2, 3})
slice2 := slice.From([]int{4, 5, 6})
slice3 := slice.From([]int{7, 8, 9})
result, _ := slice1.Concat(slice2, slice3).ToSlice() 
// result = [1, 2, 3, 4, 5, 6, 7, 8, 9]
```

## Grouping Operations

### `GroupByString(grouper func(elem T) string)`
Groups elements by string keys.

```go
result, _ := slice.From([]string{"apple", "banana", "apricot", "blueberry"}).GroupByString(func(s string) string {
    return string(s[0]) // Group by first letter
}).ToMap()
// result = map["a": ["apple", "apricot"], "b": ["banana", "blueberry"]]
```

### `PartitioningBy(predicate func(elem T) (bool, error))`
Partitions elements into two groups based on a predicate.

```go
result, _ := slice.From([]int{1, 2, 3, 4, 5, 6}).PartitioningBy(func(x int) (bool, error) {
    return x%2 == 0, nil
}).ToMap()
// result = map[true: [2, 4, 6], false: [1, 3, 5]]
```

### `AssociateByString(mapper func(elem T) (string, error))`
Creates a map with string keys.

```go
result, _ := slice.From([]string{"cat", "dog", "elephant"}).AssociateByString(func(s string) (string, error) {
    return string(s[0]), nil
}).ToMap()
// result = map["c": "cat", "d": "dog", "e": "elephant"]
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
result, _ := slice.From([]int{1, 2, 3, 4, 5}).Count()
// result = 5
```

### `First(orElse ...T) (*T, error)`
Returns the first element or a default value.

```go
first, _ := slice.From([]int{1, 2, 3}).First()
// first = *1

firstOrDefault, _ := slice.From([]int{}).First(99)
// firstOrDefault = *99
```

### `Last(orElse ...T) (*T, error)`
Returns the last element or a default value.

```go
last, _ := slice.From([]int{1, 2, 3}).Last()
// last = *3

lastOrDefault, _ slice.From([]int{}).Last(99)
// lastOrDefault = *99
```

### `FirstOrNil(predicate func(elem T) (bool, error)) (*T, error)`
Returns the first element matching a predicate.

```go
result, _ := slice.From([]int{1, 2, 3, 4, 5}).FirstOrNil(func(x int) (bool, error) {
    return x%2 == 0, nil
})
// result = *2
```

### `LastOrNil(predicate func(elem T) (bool, error)) (*T, error)`
Returns the last element matching a predicate.

```go
result, _ := slice.From([]int{1, 2, 3, 4, 5}).LastOrNil(func(x int) (bool, error) {
    return x%2 == 0, nil
})
// result = *4
```

## Matching Operations

### `AnyMatch(predicate func(elem T) (bool, error)) (bool, error)`
Checks if any element matches the predicate.

```go
result, _ := slice.From([]int{1, 3, 5, 7}).AnyMatch(func(x int) (bool, error) {
    return x%2 == 0, nil
})
// result = false
```

### `AllMatch(predicate func(elem T) (bool, error)) (bool, error)`
Checks if all elements match the predicate.

```go
result, _ := slice.From([]int{2, 4, 6, 8}).AllMatch(func(x int) (bool, error) {
    return x%2 == 0, nil
})
// result = true
```

### `NoneMatch(predicate func(elem T) (bool, error)) (bool, error)`
Checks if no elements match the predicate.

```go
result, _ := slice.From([]int{1, 3, 5, 7}).NoneMatch(func(x int) (bool, error) {
    return x%2 == 0, nil
})
// result = true
```

## Reduction Operations

### `Reduce(accumulator accumulator.Accumulator[T]) (T, error)`
Reduces the stream to a single value. See [accumulator](stream/accumulator/README.md#1-batch-accumulator-accumulatort-any)

```go
result, _ := slice.From([]int{1, 2, 3, 4, 5}).Reduce(accumulator.Sum)
// result = 15
```

### `ForEach(consumer func(elem T) error) error`
Executes a function for each element.

```go
err := slice.From([]int{1, 2, 3}).ForEach(func(x int) error {
    fmt.Printf("Number: %d\n", x)
    return nil
})
if err != nil {
	// handle error
}
// Prints: Number: 1, Number: 2, Number: 3
```

## Free Funcitons
Unfortunately, Go does not allow something like `func (s *slice[T]) Map[R any](mapper func(elem T) (R, error))`. Therefore, there are a few helper functions, where T becomes R

### `MapSlice[T, R any](s *slice[T], mapper func(elem T) (R, error))`
Maps a slice to another type T -> R
```go
numbers := slice.From([]int{1, 2, 3, 4, 5})
result, _ := slice.MapSlice(numbers, func(elem int) (string, error) {
	return "Number: " + string(rune('0'+elem)), nil
}).ToSlice()
// result = ["Number: 1", "Number: 2", "Number: 3", "Number: 4", "Number: 5"]
```

### `ZipSlices[T1, T2 any](s1 *slice[T1], s2 *slice[T2])`
Combines two streams into tuples.

```go
numbers := slice.From([]int{1, 2, 3})
letters := slice.From([]string{"a", "b", "c"})
result, _ := slice.ZipSlices(numbers, letters).ToSlice()
// result = [(1,"a"), (2,"b"), (3,"c")]
```

## Utility Operations

### `Write(writer io.Writer) (int, error)`
Writes the stream as JSON to a writer.

```go
var buf bytes.Buffer
bytesWritten, _ := slice.From([]int{1, 2, 3}).Write(&buf)
// buf contains: [1,2,3]
```

## Creating Maps

### `FromMap[K common.Key, V any](m map[K]V)`
Create a new stream from an existing map.

```go
slice.FromMap(map[string]int{
	"a": 1,
	"bb": 2,
	"aaa": 3,
})
```

### `FromMapWithSlices[K common.Key, V any](m map[K][]V)`
Creates a new stream from an existing map with slices.

```go
slice.FromMapWithSlices(map[string][]string{
	"a": {"apple", "apricot"},
	"b": {"banana"},
})
```

## Map Operations

### `CountValues()`
Returns a map contains the numbers of elements for each key.

```go
result, _ := slice.FromMapWithSlices(map[string][]string{
	"a": {"apple", "apricot"},
	"b": {"banana"},
}).CountValues().toMap()
// result = map["a": 2, "b": 1]
```