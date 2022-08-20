package utils

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_EscapeMarkdown(t *testing.T) {
	testCases := []struct {
		input  string
		expect string
	}{
		{
			input:  "",
			expect: "",
		},
		{
			input:  `\*escaped input`,
			expect: `\\\*escaped input`,
		},
		{
			input:  "`inline monospace`",
			expect: "\\`inline monospace\\`",
		},
		{
			input:  "```code fence```",
			expect: "\\`\\`\\`code fence\\`\\`\\`",
		},
		{
			input:  `*italics*`,
			expect: `\*italics\*`,
		},
		{
			input:  `**bold**`,
			expect: `\*\*bold\*\*`,
		},
		{
			input:  `_italics_`,
			expect: `\_italics\_`,
		},
		{
			input:  `__underline__`,
			expect: `\_\_underline\_\_`,
		},
		{
			input:  `> quotation`,
			expect: `\> quotation`,
		},

		// discord-specific markup
		{
			input:  `:emoji:`,
			expect: `\:emoji\:`,
		},
		{
			input:  `~strikethrough~`,
			expect: `\~strikethrough\~`,
		},
		{
			input:  `||spoiler||`,
			expect: `\|\|spoiler\|\|`,
		},
	}
	for _, tc := range testCases {
		desc := fmt.Sprintf(`EscapeMarkdown("%s") == %s`, tc.input, tc.expect)
		t.Run(desc, func(t *testing.T) {
			actual := EscapeMarkdown(tc.input)
			assert.Equal(t, tc.expect, actual)
		})
	}
}

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

		// none
		{
			expect: "none",
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
