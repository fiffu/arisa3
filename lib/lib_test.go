package lib

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Atoi(t *testing.T) {
	type testCase struct {
		in  string
		out int
	}
	tests := []testCase{
		// regular stuff
		{"1", 1}, {"+1", 1}, {"-1", -1}, {"-01", -1}, {"01", 1},
		// zeroes
		{"0", 0}, {"+0", 0}, {"-0", 0}, {"000", 0},
		// edge cases - all should return 0
		{"a", 0}, {"aaaaa", 0}, {"", 0}, {"-", 0},
	}
	for _, tc := range tests {
		t.Run("Atoi(%d)", func(t *testing.T) {
			assert.Equal(t, tc.out, Atoi(tc.in))
		})
	}
}

func Test_SplitOnce(t *testing.T) {
	type testCase struct {
		input, delim, left, right string
	}
	tests := []testCase{
		{"left right", " ", "left", "right"},
		{"left ", " ", "left", ""},
		{" right", " ", "", "right"},
		{"left  ", " ", "left", " "},
		{"  right", " ", "", " right"},
		{"foo", "x", "foo", ""},
		{"foo", "", "", "foo"},
	}
	for _, tc := range tests {
		name := fmt.Sprintf(
			"SplitOnce(%s, %s) should yield '%s' and '%s'",
			tc.input, tc.delim, tc.left, tc.right,
		)
		t.Run(name, func(t *testing.T) {
			left, right := SplitOnce(tc.input, tc.delim)
			assert.Equal(t, tc.left, left)
			assert.Equal(t, tc.right, right)
		})
	}
}

func Test_FuncName(t *testing.T) {
	testCases := []struct {
		function     any
		expectSuffix string
	}{
		{
			function:     MustGetCallerDir,
			expectSuffix: "lib.MustGetCallerDir",
		},
		{
			function:     t.Log,
			expectSuffix: "testing.(*common).Log",
		},
		{
			function:     assert.True,
			expectSuffix: "assert.True",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.expectSuffix, func(t *testing.T) {
			actual := FuncName(tc.function)
			assert.Contains(t, actual, tc.expectSuffix)
			assert.True(t, strings.HasSuffix(actual, tc.expectSuffix))
		})
	}
}

func Test_MustGetCallerDir(t *testing.T) {
	here := MustGetCallerDir()
	expect := "lib"
	actual := here[len(here)-len(expect):]
	assert.Equal(t, expect, actual)
}

func Test_Clamper(t *testing.T) {
	type testCase struct {
		lo, hi, num, expect float64
	}
	for _, tc := range []testCase{
		{0, 1, 3, 1},        // +ve, overflow
		{0, 1, 0.5, 0.5},    // +ve, neutral
		{0, 1, -3, 0},       // +ve, underflow
		{-1, 0, -3, -1},     // -ve, underflow
		{-1, 0, -0.5, -0.5}, // -ve, neutral
		{-1, 0, 3, 0},       // -ve, overflow
		{-1, 1, 3, 1},       // crossing signs, overflow
		{-1, 1, 0.5, 0.5},   // crossing signs, neutral
		{-1, 1, -3, -1},     // crossing signs, underflow
	} {
		clamp := Clamper(tc.lo, tc.hi)
		actual := clamp(tc.num)
		assert.Equal(t, tc.expect, actual)
	}
}

func Test_UniformRange(t *testing.T) {
	type testCase struct {
		lo, hi float64
	}
	for _, tc := range []testCase{
		{0, 1},
		{2, 3},
		{-2, -1},
		{-1, 1},
	} {
		for i := 0; i < 10000; i++ {
			out := UniformRange(tc.lo, tc.hi)
			assert.LessOrEqual(t, tc.lo, out)
			assert.GreaterOrEqual(t, tc.hi, out)
		}
	}

	out := UniformRange(1.0, -1.0) // expect this to work even if hi and lo are swapped in inputs
	assert.LessOrEqual(t, -1.0, out)
	assert.GreaterOrEqual(t, 1.0, out)
}

func Test_CoinFlip(t *testing.T) {
	outcome := make(map[bool]bool)

	i := 0
	for len(outcome) < 2 && i < 1000 {
		choice := CoinFlip()
		outcome[choice] = true
	}
	assert.Len(t, outcome, 2)
}

func Test_IntDivmod(t *testing.T) {
	testCases := []struct {
		left, right                     int
		expectQuotient, expectRemainder int
	}{
		{
			5, 3,
			1, 2,
		},
		{
			99, 11,
			9, 0,
		},
	}
	for _, tc := range testCases {
		desc := fmt.Sprintf(
			"%d // %d should yield quotient of %d and remainder of %d",
			tc.left, tc.right, tc.expectQuotient, tc.expectRemainder,
		)
		t.Run(desc, func(t *testing.T) {
			q, r := IntDivmod(tc.left, tc.right)
			assert.Equal(t, tc.expectQuotient, q)
			assert.Equal(t, tc.expectRemainder, r)
		})
	}
}
