package common

type Key interface {
	int | int64 | float64 | string | bool
}

type Numbers interface {
	~int | ~int64 | ~float64
}
