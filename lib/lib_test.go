package lib

import (
	"fmt"
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

func Test_MustGetCallerDir(t *testing.T) {
	here := MustGetCallerDir()
	expect := "lib"
	actual := here[len(here)-len(expect):]
	assert.Equal(t, expect, actual)
}

func Test_ChooseString(t *testing.T) {
	choices := []string{"a", "b"}
	outcome := make(map[string]bool)

	i := 0
	for len(outcome) < 2 && i < 1000 {
		choice := ChooseString(choices)
		outcome[choice] = true
	}
	assert.Len(t, outcome, len(choices))
}

func Test_ChooseBool(t *testing.T) {
	outcome := make(map[bool]bool)

	i := 0
	for len(outcome) < 2 && i < 1000 {
		choice := ChooseBool()
		outcome[choice] = true
	}
	assert.Len(t, outcome, 2)
}

func Test_DecimalToRGB(t *testing.T) {
	// ffff00
	num := 255<<16 + 255<<8 + 0
	r, g, b := DecimalToRGB(num)
	assert.Equal(t, 1.0, r)
	assert.Equal(t, 1.0, g)
	assert.Equal(t, 0.0, b)
}
