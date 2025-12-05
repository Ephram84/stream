package common

type SortFunc[T any] func(slice []T) func(i, j int) bool

func Sort[T Ordered](arr []T) func(i, j int) bool {
	return func(i, j int) bool {
		return arr[i] < arr[j]
	}
}

func SortDesc[T Ordered](arr []T) func(i, j int) bool {
	return func(i, j int) bool {
		return arr[i] > arr[j]
	}
}
