// Package common provides shared types, utilities, and interfaces used across the stream package.
package common

// Key is a constraint that allows types that can be used as map keys.
// It includes int, int64, float64, string, and bool types.
type Key interface {
	int | int64 | float64 | string | bool
}

// Numbers is a constraint for numeric types that support arithmetic operations.
// It includes int, int64, and float64 types with their underlying types.
type Numbers interface {
	~int | ~int64 | ~float64
}

// Ordered is a constraint for types that support ordering comparisons.
// It includes int, int64, float64, and string types with their underlying types.
type Ordered interface {
	~int | ~int64 | ~float64 | ~string
}
