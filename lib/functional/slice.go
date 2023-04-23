package functional

type Slice[T, U any] []T

func (s Slice[T, U]) Shuffle() Slice[T, U]                     { return Shuffle(s) }
func (s Slice[T, U]) Take(n int) Slice[T, U]                   { return Take(s, n) }
func (s Slice[T, U]) TakeRandom() T                            { return TakeRandom(s) }
func (s Slice[T, U]) Filter(pred func(T) bool) Slice[T, U]     { return Filter(s, pred) }
func (s Slice[T, U]) Map(mapper func(T) U) Slice[U, undefined] { return Map(s, mapper) }

type undefined any

func SliceOf[T any](elems []T) Slice[T, undefined] {
	return Slice[T, undefined](elems)
}
