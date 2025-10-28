package tupel

type Tupel[A any, B any] struct {
	First  A
	Second B
}

func New[A any, B any](first A, second B) *Tupel[A, B] {
	return &Tupel[A, B]{
		First:  first,
		Second: second,
	}
}
