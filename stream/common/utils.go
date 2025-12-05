package common

func Min[T Numbers](a, b T) T {
	if a < b {
		return a
	}
	return b
}

func Max[T Numbers](a, b T) T {
	if a > b {
		return a
	}
	return b
}
