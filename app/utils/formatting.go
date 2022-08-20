package utils

import (
	"fmt"
	"strings"
	"time"

	"github.com/fiffu/arisa3/lib"
)

func EscapeMarkdown(text string) string {
	markdownSymbols := "\\*>_:~`|"
	for _, char := range markdownSymbols {
		s := string(char)
		text = strings.ReplaceAll(text, s, "\\"+s)
	}
	return text
}

func FormatDuration(delta time.Duration) string {
	deltaSeconds := int(delta.Seconds())
	if deltaSeconds <= 0 {
		return "none"
	}
	hours, remainder := lib.IntDivmod(deltaSeconds, 60*60)
	mins, secs := lib.IntDivmod(remainder, 60)

	output := make([]string, 0)
	var h, m, s string
	if hours > 0 {
		h = fmt.Sprintf("%dhr", hours)
		output = append(output, h)
	}
	if mins > 0 {
		m = fmt.Sprintf("%dmin", mins)
		output = append(output, m)
	}
	if secs > 0 && h == "" && m == "" {
		s = "less than a minute"
		output = append(output, s)
	}
	return strings.Join(output, " ")
}
