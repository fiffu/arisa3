package functional

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type S struct{ int }

func (s *S) String() string { return fmt.Sprint(s.int) }

func Test_Map(t *testing.T) {
	assert.Equal(t,
		[]string{"1", "2", "3"},
		Map(
			[]*S{{1}, {2}, {3}},
			func(s *S) string { return s.String() },
		),
	)
}

func Test_Zip(t *testing.T) {
	testCases := []struct {
		desc   string
		left   []int
		right  []int
		expect []Tuple[int, int]
	}{
		{
			desc:  "equal length",
			left:  []int{1, 2, 3},
			right: []int{1, 4, 9},
			expect: []Tuple[int, int]{
				{1, 1}, {2, 4}, {3, 9},
			},
		},
		{
			desc:  "left longer than right",
			left:  []int{1, 2, 3, 4, 5, 6},
			right: []int{1, 4, 9},
			expect: []Tuple[int, int]{
				{1, 1}, {2, 4}, {3, 9},
			},
		},
		{
			desc:  "right longer than left",
			left:  []int{1, 2, 3},
			right: []int{1, 4, 9, 16, 25, 36},
			expect: []Tuple[int, int]{
				{1, 1}, {2, 4}, {3, 9},
			},
		},
		{
			desc:   "left empty",
			left:   []int{},
			right:  []int{1, 4, 9},
			expect: []Tuple[int, int]{},
		},
		{
			desc:   "right empty",
			left:   []int{1, 2, 3},
			right:  []int{},
			expect: []Tuple[int, int]{},
		},
		{
			desc:   "both empty",
			left:   []int{},
			right:  []int{},
			expect: []Tuple[int, int]{},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			actual := Zip(tc.left, tc.right)
			assert.ElementsMatch(t, tc.expect, actual)
		})
	}
}

func Test_Zip_MixedTypes(t *testing.T) {
	left := []int{1, 2, 3}
	right := []string{"one", "two", "three"}

	z := Zip(left, right)
	assert.Equal(t, z[0].Left, 1)
	assert.Equal(t, z[0].Right, "one")
}

func Test_Shuffle(t *testing.T) {
	original := []int{
		1, 2, 3, 4, 5, 6, 7, 8, 9, 10,
		11, 12, 13, 14, 15, 16, 17, 18, 19, 20,
	}

	buf := make([]int, len(original))
	test := func() (fail bool) {
		copy(buf, original)
		shuffled := Shuffle(buf)
		return Equal(original, shuffled)
	}

	testCount := 1000
	identicalCount := 0
	for i := 0; i < testCount; i++ {
		if fail := test(); fail {
			identicalCount += 1
		}
	}

	assert.Less(t, identicalCount, testCount, "shuffled outcome matches the original input for all of %d", testCount)
}

func Test_Take(t *testing.T) {
	testCases := []struct {
		desc   string
		input  []int
		take   int
		expect []int
	}{
		{
			desc:   "take less than length",
			input:  []int{1, 2, 3},
			take:   2,
			expect: []int{1, 2},
		},
		{
			desc:   "take more than length",
			input:  []int{1, 2, 3},
			take:   100,
			expect: []int{1, 2, 3},
		},
		{
			desc:   "take none",
			input:  []int{1, 2, 3},
			take:   0,
			expect: []int{},
		},
		{
			desc:   "take from empty input",
			input:  []int{},
			take:   5,
			expect: []int{},
		},
		{
			desc:   "take negative number",
			input:  []int{1, 2, 3},
			take:   -5,
			expect: []int{},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			actual := Take(tc.input, tc.take)
			assert.ElementsMatch(t, tc.expect, actual)
		})
	}
}

func Test_Contains(t *testing.T) {
	abc := []string{"a", "b", "c"}
	assert.True(t, Contains(abc, "a"))
	assert.True(t, Contains(abc, "b"))
	assert.True(t, Contains(abc, "c"))
	assert.False(t, Contains(abc, "z"))

	empty := []string{}
	assert.False(t, Contains(empty, "a"))
	assert.False(t, Contains(empty, "b"))
	assert.False(t, Contains(empty, "c"))
	assert.False(t, Contains(empty, "z"))
}

func Test_Deref(t *testing.T) {
	one, two, three := 1, 2, 3
	abc := []*int{
		&one, &two, &three,
	}
	assert.ElementsMatch(
		t,
		Deref(abc),
		[]int{1, 2, 3},
	)
}
