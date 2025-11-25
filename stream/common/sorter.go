package common

func SortInts(arr []int) func(i, j int) bool {
	return func(i, j int) bool {
		return arr[i] < arr[j]
	}
}

func SortFloat(arr []float64) func(i, j int) bool {
	return func(i, j int) bool {
		return arr[i] < arr[j]
	}
}
