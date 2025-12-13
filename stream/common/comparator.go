package common

// Eq checks if two comparable values are equal.
func Eq[E comparable](a, b E) bool {
	return a == b
}

// Min returns the smaller of two numeric values.
func Min[T Numbers](a, b T) T {
	if a < b {
		return a
	}
	return b
}

// Max returns the larger of two numeric values.
func Max[T Numbers](a, b T) T {
	if a > b {
		return a
	}
	return b
}
