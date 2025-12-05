# Accumulator

The `accumulator` package provides functions for aggregating values.

## Overview

There are two types of accumulators:

### 1. Batch Accumulator (`Accumulator[T any]`)

Processes all values at once as a slice. Functions are provided for summing, averaging, and determining min and max values. The prerequisite is that the values are already numbers (ints or floats). Only Avg requires float numbers. 

- **`Sum[T]`**: Calculates the sum of all numeric values
```go
numbers := slice.From([]int{1, 2, 3, 4, 5})
sum, _ := numbers.Reduce(accumulator.Sum)
// sum = 15
```

```go
numbers := slice.From([]float64{1.0, 2.0, 3.0, 4.0, 5.0})
sum, _ := numbers.Reduce(accumulator.Sum)
// sum = 15.0
```

- **`Avg`**: Calculates the average of `float64` values
```go
numbers := slice.From([]float64{1.0, 2.0, 3.0, 4.0, 5.0})
avg, _ := numbers.Reduce(accumulator.Avg)
// avg = 3.0
```

- **`Min[T]`**: Finds the minimum value
```go
numbers := slice.From([]int{1, 2, 3, 4, 5})
min, _ := numbers.Reduce(accumulator.Min)
// min = 1
```

```go
numbers := slice.From([]float64{1.0, 2.0, 3.0, 4.0, 5.0})
min, _ := numbers.Reduce(accumulator.Min)
// min = 1.0
```

- **`Max[T]`**: Finds the maximum value
```go
numbers := slice.From([]int{1, 2, 3, 4, 5})
max, _ := numbers.Reduce(accumulator.Max)
// max = 5
```

```go
numbers := slice.From([]float64{1.0, 2.0, 3.0, 4.0, 5.0})
max, _ := numbers.Reduce(accumulator.Max)
// max = 5.0
```

You can also use your own accumulator function. To do so, the function must fulfill the signature `func(values []T) (T, error)`
```go
transactions := []Transaction{
		{
			ID:          "T01",
			Amount:      10.0,
			BookingDate: 1762038000, // 2025-11-02
		},
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
	}

	oldestTransaction, _ := slice.From(transactions).Reduce(func(values []Transaction) (Transaction, error) {
		oldest := values[0]
		for _, transaction := range values[1:] {
			if transaction.BookingDate < oldest.BookingDate {
				oldest = transaction
			}
		}

		return oldest, nil
	})
    // oldestTransaction = {
    // 	    ID:          "T03",
	// 	    Amount:      30.0,
	// 	    BookingDate: 1756764000, // 2025-09-02
	// },
```

### 2. Sequential Accumulator (`AccumulatorSeq[T any]`)

Processes values sequentially with index and previous result. As with batch accumulator, functions for summing, averaging, and determining min and max are also provided for sequences.

- **`SumSeq[T]`**: Adds values sequentially
```go
numbers := sequence.From([]int{1, 2, 3, 4, 5}, nil)
sum, _ := numbers.Reduce(accumulator.SumSeq)
// sum = 15
```

```go
numbers := sequence.From([]float64{1.0, 2.0, 3.0, 4.0, 5.0}, nil)
sum, _ := numbers.Reduce(accumulator.SumSeq)
// sum = 15.0
```

- **`AvgSeq`**: Calculates the running average
```go
numbers := sequence.From([]float64{1.0, 2.0, 3.0, 4.0, 5.0})
avg, _ := numbers.Reduce(accumulator.AvgSeq)
// avg = 3.0
```

- **`MinSeq[T]`**: Determines the minimum sequentially
```go
numbers := sequence.From([]int{1, 2, 3, 4, 5}, nil)
min, _ := numbers.Reduce(accumulator.MinSeq)
// min = 1
```

```go
numbers := sequence.From([]float64{1.0, 2.0, 3.0, 4.0, 5.0}, nil)
min, _ := numbers.Reduce(accumulator.MinSeq)
// sum = 1.0
```
- **`MaxSeq[T]`**: Determines the maximum sequentially
```go
numbers := sequence.From([]int{1, 2, 3, 4, 5}, nil)
max, _ := numbers.Reduce(accumulator.MaxSeq)
// max = 15
```

```go
numbers := sequence.From([]float64{1.0, 2.0, 3.0, 4.0, 5.0}, nil)
max, _ := numbers.Reduce(accumulator.MaxSeq)
// max = 15.0
```

You can also use your own accumulator function. To do so, the function must fulfill the signature `func(index int, prev T, current T) (T, error)`

The sequential variants are useful for streaming operations where values are processed one at a time.
