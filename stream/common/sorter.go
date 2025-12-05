package common

// SortFunc is a function type that takes a slice and returns a comparison function
// for sorting. The returned function takes two indices and returns true if the element
// at index i should come before the element at index j.
type SortFunc[T any] func(slice []T) func(i, j int) bool

// Sort returns a comparison function for sorting a slice in ascending order.
// It works with any type that implements the Ordered constraint.
func Sort[T Ordered](arr []T) func(i, j int) bool {
	return func(i, j int) bool {
		return arr[i] < arr[j]
	}
}

// SortDesc returns a comparison function for sorting a slice in descending order.
// It works with any type that implements the Ordered constraint.
func SortDesc[T Ordered](arr []T) func(i, j int) bool {
	return func(i, j int) bool {
		return arr[i] > arr[j]
	}
}
