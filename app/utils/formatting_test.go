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
		days, hours, mins, secs int
		expect                  string
	}{
		// <1 min
		{
			expect: "none",
		},
		{
			mins:   1,
			secs:   -1,
			expect: "less than a minute",
		},

		// 1 to 59.99 mins
		{
			mins:   1,
			expect: "1 min",
		},
		{
			mins:   3,
			secs:   2,
			expect: "3 mins",
		},
		{
			mins:   4,
			secs:   -1,
			expect: "3 mins",
		},

		// 1 to 23.99 hours
		{
			hours:  1,
			expect: "1 hour",
		},
		{
			hours:  3,
			secs:   4,
			expect: "3 hours",
		},
		{
			hours:  3,
			mins:   4,
			secs:   4,
			expect: "3 hours 4 mins",
		},
		{
			hours:  24,
			secs:   -1,
			expect: "23 hours 59 mins",
		},

		// at least 1 day
		{
			days:   1,
			expect: "1 day",
		},
		{
			days:   1,
			hours:  4,
			expect: "1 day, 4 hours",
		},
		{
			days:   3,
			hours:  4,
			expect: "3 days, 4 hours",
		},
		{
			days:   16,
			expect: "16 days",
		},
		{
			days:   16,
			hours:  4,
			expect: "16 days",
		},
	}
	for _, tc := range testCases {
		d := time.Duration(tc.days) * time.Hour * 24
		h := time.Duration(tc.hours) * time.Hour
		m := time.Duration(tc.mins) * time.Minute
		s := time.Duration(tc.secs) * time.Second
		duration := d + h + m + s

		desc := fmt.Sprintf("duration %v should be formatted as %s", duration, tc.expect)

		t.Run(desc, func(t *testing.T) {
			actual := FormatDuration(duration)
			assert.Equal(t, tc.expect, actual)
		})
	}
}
