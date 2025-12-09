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

### `Map[T, R any](s *slice[T], mapper func(elem T) (R, error))`
Maps a slice to another type T -> R
```go
numbers := slice.From([]int{1, 2, 3, 4, 5})
result, _ := slice.Map(numbers, func(elem int) (string, error) {
	return "Number: " + string(rune('0'+elem)), nil
}).ToSlice()
// result = ["Number: 1", "Number: 2", "Number: 3", "Number: 4", "Number: 5"]
```

### `Zip[T1, T2 any](s1 *slice[T1], s2 *slice[T2])`
Combines two streams into tuples.

```go
numbers := slice.From([]int{1, 2, 3})
letters := slice.From([]string{"a", "b", "c"})
result, _ := slice.Zip(numbers, letters).ToSlice()
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

There are two types of map operations in the slice package:

1. **Simple Maps** (`pairs[K, V]`) - Maps with single values: `map[K]V`
2. **Maps with Slices** (`pairsSlice[K, V]`) - Maps with slice values: `map[K][]V`

### `FromMap[K common.Key, V any](m map[K]V)`
Create a new stream from an existing map with single values.

```go
// Simple map: map[string]int
simpleMap := map[string]int{
    "a": 1,
    "bb": 2,
    "ccc": 3,
}
pairs := slice.FromMap(simpleMap)
```

### `FromMapWithSlices[K common.Key, V any](m map[K][]V)`
Creates a new stream from an existing map with slice values.

```go
// Map with slices: map[string][]string
mapWithSlices := map[string][]string{
    "a": {"apple", "apricot"},
    "b": {"banana"},
}
pairsSlice := slice.FromMapWithSlices(mapWithSlices)
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
Returns a slice of all values from the map.

```go
values, _ := slice.FromMap(map[string]int{
    "a": 1,
    "b": 2,
    "c": 3,
}).Flatten().ToSlice()
// values = [1, 2, 3] (order may vary)
```

#### `Keys()`
Returns a slice of all keys from the map.

```go
keys, _ := slice.FromMap(map[string]int{
    "a": 1,
    "b": 2,
    "c": 3,
}).Keys().ToSlice()
// keys = ["a", "b", "c"] (order may vary)
```

#### `ToMap()`
Converts the pairs back to a regular Go map.

```go
pairs := slice.FromMap(map[string]int{
    "a": 1,
    "b": 2,
})
result, _ := pairs.ToMap()
// result = map["a": 1, "b": 2]
```

#### `Count()`
Returns the number of key-value pairs in the map.

```go
count, _ := slice.FromMap(map[string]int{
    "a": 1,
    "b": 2,
    "c": 3,
}).Count()
// count = 3
```

#### `ForEach(action func(key K, value V))`
Executes a function for each key-value pair.

```go
pairs := slice.FromMap(map[string]int{
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
result, _ := slice.FromMapWithSlices(map[string][]string{
    "a": {"apple", "apricot", "avocado"},
    "b": {"banana", "blueberry"},
    "c": {"cherry"},
}).CountValues().ToMap()
// result = map["a": 3, "b": 2, "c": 1]
```

#### `Flatten()`
Returns a slice of all values from all slices.

```go
values, _ := slice.FromMapWithSlices(map[string][]int{
    "first":  {1, 2, 3},
    "second": {4, 5},
    "third":  {6},
}).Flatten().ToSlice()
// result = [1, 2, 3, 4, 5, 6] (order may vary)
```

#### `Keys()`
Returns a slice of all keys from the map.

```go
keys, _ := slice.FromMapWithSlices(map[string][]int{
    "first":  {1, 2, 3},
    "second": {4, 5},
    "third":  {6},
}).Keys().ToSlice()
// keys = ["first", "second", "third"] (order may vary)
```

#### `MapToFloat64(mapper func(elem V) (float64, error))`
Maps all values in all slices to float64.

```go
result, _ := slice.FromMapWithSlices(map[string][]int{
    "group1": {1, 2, 3},
    "group2": {4, 5},
}).MapToFloat64(func(x int) (float64, error) {
    return float64(x) * 1.5, nil
}).ToMap()
// result = map["group1": [1.5, 3.0, 4.5], "group2": [6.0, 7.5]]
```

#### `Reduce(accumulator accumulator.Accumulator[V])`
Reduces each slice in the map to a single value, converting to a simple pairs map.
Note: Uses batch accumulator (`Accumulator[V]`), not sequential (`AccumulatorSeq[V]`).

```go
result, _ := slice.FromMapWithSlices(map[string][]int{
    "group1": {1, 2, 3},
    "group2": {4, 5, 6},
    "group3": {7, 8},
}).Reduce(accumulator.Sum).ToMap()
// result = map["group1": 6, "group2": 15, "group3": 15]
```

#### `ToMap()`
Converts the pairsSlice back to a regular Go map with slices.

```go
pairsSlice := slice.FromMapWithSlices(map[string][]int{
    "first":  {1, 2, 3},
    "second": {4, 5, 6},
})
result, _ := pairsSlice.ToMap()
// result = map["first": [1, 2, 3], "second": [4, 5, 6]]
```

#### `Count()`
Returns the number of keys in the map.

```go
count, _ := slice.FromMapWithSlices(map[string][]int{
    "first":  {1, 2, 3},
    "second": {4, 5},
}).Count()
// count = 2
```

#### `ForEach(action func(key K, value []V))`
Executes a function for each key-slice pair.

```go
pairsSlice := slice.FromMapWithSlices(map[string][]int{
    "group1": {1, 2, 3},
    "group2": {4, 5},
})
pairsSlice.ForEach(func(key string, values []int) {
    fmt.Printf("Key: %s, Values: %v\n", key, values)
})
// Prints each key with its slice of values (order may vary)
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
pairs := slice.FromMap(scores)
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
pairsSlice := slice.FromMapWithSlices(studentsByGrade)
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
pairsSlice := slice.FromMapWithSlices(groupedData)
sums := pairsSlice.Reduce(accumulator.Sum)

// Now we have a simple map
sumMap, _ := sums.ToMap()
// sumMap = map[1: 60, 2: 90, 3: 300]

// Get the maximum sum
maxSum, _ := sums.Flatten().Reduce(accumulator.Max)
// maxSum = 300
```

### Chaining Map and Slice Operations
```go
// Start with a map of transaction amounts by category
transactions := map[string][]float64{
    "Food":      {12.50, 8.75, 15.00},
    "Transport": {45.00, 23.50},
    "Entertainment": {50.00, 30.00, 25.00},
}

// Calculate total spending per category, then find categories over 100
pairsSlice := slice.FromMapWithSlices(transactions)
highSpendingCategories, _ := pairsSlice.Reduce(accumulator.Sum)
.Filter(func(key string, value float64) (bool, error) {
	return value > 100, nil
}).Keys().ToSlice()
// highSpendingCategories might include "Entertainment" depending on totals
```