package rng

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_parseExpr(t *testing.T) {
	type testCase struct {
		name                string
		in                  []string
		count, sides, modif int
		parsed              bool
	}
	tests := []testCase{
		{
			name:   "weird inputs",
			in:     []string{"", "0d", "d", " "},
			parsed: false,
		},
		{
			"zeroes, implicit count", []string{"0", "000", "1d0", "d0"},
			1, 0, 0, true,
		},
		{
			"zeroes, explicit count", []string{"0d0", "00d00", "000d000-000"},
			0, 0, 0, true,
		},
		{
			"whitespace", []string{" 1d2+3", "\t1d2+3", "\t 1d2+3", "1d2+3 ", "1d2+3\t", "1d2+3 \t"},
			1, 2, 3, true,
		},
		{
			"rolling d99", []string{"44", "d44", "D44"},
			1, 44, 0, true,
		},
		{
			"rolling 3d5", []string{"3d5", "3D5", "3d5+0", "3d5-0"},
			3, 5, 0, true,
		},
		{
			"rolling d5+10", []string{"d5+10", "D5+10", "1d5+10", "1D5+10"},
			1, 5, 10, true,
		},
		{
			"rolling 3d5-10", []string{"d5-10", "D5-10", "1d5-10", "1D5-10"},
			1, 5, -10, true,
		},
	}
	for _, tc := range tests {
		for _, s := range tc.in {
			t.Run(tc.name+" with "+s, func(t *testing.T) {
				actual := parseExpr(s)
				if !tc.parsed {
					assert.False(t, actual.parsed)
				} else {
					expect := dice{tc.count, tc.sides, tc.modif, true}
					assert.Equal(t, expect, actual)
				}
			})
		}
	}
}

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

func Test_toss(t *testing.T) {
	t.Run("testing count", func(t *testing.T) {
		d := dice{count: 33, sides: 1, modif: 0}
		expect := 33
		actual := toss(d)
		assert.Equal(t, expect, actual)
	})
	t.Run("testing sides", func(t *testing.T) {
		d := dice{count: 9, sides: 0, modif: -1}
		expect := -1
		actual := toss(d)
		assert.Equal(t, expect, actual)
	})
	t.Run("testing modif", func(t *testing.T) {
		d := dice{count: 1, sides: 1, modif: 1000}
		expect := 1001
		actual := toss(d)
		assert.Equal(t, expect, actual)
	})
}

func Test_throwDie(t *testing.T) {
	sides := 6
	t.Run("monte carlo testing", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			actual := throwDie(sides)
			assert.LessOrEqual(t, 1, actual)
			assert.GreaterOrEqual(t, 6, actual)
		}
	})
}

func Test_formatDice(t *testing.T) {
	type testCase struct {
		name   string
		d      dice
		expect string
	}
	tests := []testCase{
		{
			name:   "parse is negative",
			d:      dice{parsed: false},
			expect: "0-99",
		},
		{
			d:      dice{count: 0, sides: 0, modif: 0, parsed: true},
			expect: "0d0",
		},
		{
			d:      dice{count: 0, sides: 0, modif: -1, parsed: true},
			expect: "0d0-1",
		},
		{
			d:      dice{count: 3, sides: 5, modif: 0, parsed: true},
			expect: "3d5",
		},
		{
			d:      dice{count: 3, sides: 5, modif: 1, parsed: true},
			expect: "3d5+1",
		},
	}
	for _, tc := range tests {
		n := tc.name
		if n == "" {
			n = tc.expect
		}
		t.Run(n, func(t *testing.T) {
			actual := formatDice(tc.d)
			assert.Equal(t, tc.expect, actual)
		})
	}
}
