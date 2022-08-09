package utils

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_FormatDuration(t *testing.T) {
	testCases := []struct {
		hours, mins, secs int
		expect            string
	}{
		// h
		{
			hours:  1,
			expect: "1hr",
		},
		// m
		{
			mins:   2,
			expect: "2min",
		},
		// s
		{
			secs:   3,
			expect: "less than a minute",
		},
		// h,m
		{
			hours:  1,
			mins:   2,
			expect: "1hr 2min",
		},
		// h,s
		{
			hours:  1,
			secs:   2,
			expect: "1hr",
		},
		// m,s
		{
			mins:   2,
			secs:   3,
			expect: "2min",
		},

		// h, m, s
		{
			hours:  1,
			mins:   2,
			secs:   3,
			expect: "1hr 2min",
		},
	}
	for _, tc := range testCases {
		h := time.Duration(tc.hours) * time.Hour
		m := time.Duration(tc.mins) * time.Minute
		s := time.Duration(tc.secs) * time.Second
		duration := h + m + s

		desc := fmt.Sprintf("duration %v should be formatted as %s", duration, tc.expect)

		t.Run(desc, func(t *testing.T) {
			actual := FormatDuration(duration)
			assert.Equal(t, tc.expect, actual)
		})
	}
}
