package stream

func Eq[E comparable](a, b E) bool {
	return a == b
}

func MaxInt(max, elem int) bool {
	return max < elem
}

func MinInt(min, elem int) bool {
	return min > elem
}

func MaxFloat(max, elem float64) bool {
	return max < elem
}

func MinFloat(min, elem float64) bool {
	return min > elem
}
