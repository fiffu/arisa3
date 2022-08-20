package functional

import (
	"math/rand"
)

func min(i, j int) int {
	if i < j {
		return i
	}
	return j
}

type Tuple[T, U any] struct {
	Left  T
	Right U
}

// Shuffle returns a copy of input, with elements randomly reordered.
func Shuffle[T any](input []T) []T {
	elems := make([]T, len(input))
	copy(elems, input)
	rand.Shuffle(len(elems), func(i, j int) {
		elems[i], elems[j] = elems[j], elems[i]
	})
	return elems
}

// Take returns a slice containing the first n elements of the input slice.
// The returned elements are not removed from the input slice.
func Take[T any](input []T, n int) []T {
	ts := make([]T, 0)
	for i, t := range input {
		if i >= n {
			break
		}
		ts = append(ts, t)
	}
	return ts
}

// TakeRandom returns a random element from the input slice.
// The taken element is not removed from the input slice.
func TakeRandom[T any](elems []T) T {
	n := rand.Intn(len(elems))
	return elems[n]
}

// Contains returns whether the given elem is found in the input slice of elems.
func Contains[T comparable](elems []T, elem T) bool {
	for _, e := range elems {
		if e == elem {
			return true
		}
	}
	return false
}

// Zip combines both input slice into a slice of tuples.
// Each nth tuple holds the nth element from both lists.
// The two input slices can hold different types.
func Zip[T, U any](left []T, right []U) []Tuple[T, U] {
	var result []Tuple[T, U]

	size := min(len(left), len(right))
	for i, t := range left {
		if i == size {
			break
		}
		u := right[i]
		result = append(result, Tuple[T, U]{Left: t, Right: u})
	}
	return result
}

// Map applies the given mapper function to every element of the input slice.
// The result is a slice holding each result returned by the mapper function.
func Map[T, U any](elems []T, mapper func(T) U) []U {
	var result []U
	for _, t := range elems {
		result = append(result, mapper(t))
	}
	return result
}

// Deref takes an input slice of pointers and returns a slice storing their
// dereferenced values.
func Deref[T any](elems []*T) []T {
	return Map(elems, func(in *T) T {
		return *in
	})
}
