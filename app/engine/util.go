package engine

import (
	"regexp"
	"strings"
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
