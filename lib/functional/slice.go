package functional

type Slice[T any] []T

func (s Slice[T]) Shuffle() Slice[T]           { return Shuffle(s) }
func (s Slice[T]) Take(n int) Slice[T]         { return Take(s, n) }
func (s Slice[T]) TakeRandom() T               { return TakeRandom(s) }
func (s Slice[T]) Zip(other []T) []Tuple[T, T] { return Zip(s, other) }

func SliceOf[T any](elems []T) Slice[T] {
	return Slice[T](elems)
}
