package stream

func getSize(sizes []int) int {
	if len(sizes) > 0 {
		return sizes[0]
	}
	return 1
}
