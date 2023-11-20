package utils

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/fiffu/arisa3/lib"
)

// Misc values
var (
	CustomEmojiRegex *regexp.Regexp = regexp.MustCompile(`(<a?)?:\w+:(\d{18}>)?`)
)

// PrettifyCustomEmoji converts "<:birb:924875584004321361>" into ":birb:"
func PrettifyCustomEmoji(str string) string {
	replacer := func(s string) string {
		start := strings.Index(s, ":")
		end := strings.LastIndex(s, ":")
		return s[start : end+1]
	}
	return CustomEmojiRegex.ReplaceAllStringFunc(str, replacer)
}

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
	oneMinSecs := 60
	oneHourSecs := 60 * oneMinSecs
	oneDaySecs := 24 * oneHourSecs

	days, remainder := lib.IntDivmod(deltaSeconds, oneDaySecs)
	hours, remainder := lib.IntDivmod(remainder, 60*60)
	mins, secs := lib.IntDivmod(remainder, 60)

	return formatDuration(days, hours, mins, secs)
}

func formatDuration(days, hours, mins, secs int) string {
	fmtD := func(n int) string { return formatUnit(n, "1 day", "%d days") }
	fmtH := func(n int) string { return formatUnit(n, "1 hour", "%d hours") }
	fmtM := func(n int) string { return formatUnit(n, "1 min", "%d mins") }

	var output []string
	switch {
	case days >= 14:
		return fmtD(days)

	case days > 0:
		output = append(output, fmtD(days))
		if h := fmtH(hours); h != "" {
			output = append(output, h)
		}
		return strings.Join(output, ", ")

	case hours > 0:
		output = append(output, fmtH(hours))
		if m := fmtM(mins); m != "" {
			output = append(output, m)
		}
		return strings.Join(output, " ")

	case mins > 0:
		return fmtM(mins)

	case secs > 0:
		return "less than a minute"

	default:
		return "none"
	}
}

func formatUnit(n int, one, many string) string {
	switch n {
	case 0:
		return ""
	case 1:
		return one
	default:
		return fmt.Sprintf(many, n)
	}
}
