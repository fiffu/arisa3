package functional

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Slice_ChainingCalls(t *testing.T) {
	s := SliceOf([]int{99})

	assert.ElementsMatch(
		t,
		[]int{99},
		s.Take(1),
	)
	assert.ElementsMatch(
		t,
		[]int{99},
		s.Take(1).Shuffle(),
	)
	assert.Equal(
		t,
		99,
		s.Take(1).Shuffle().TakeRandom(),
	)
}
