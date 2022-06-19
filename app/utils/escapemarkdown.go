package utils

import (
	"strings"
)

func EscapeMarkdown(text string) string {
	markdownSymbols := "\\*>_:~`|"
	for _, char := range markdownSymbols {
		s := string(char)
		text = strings.ReplaceAll(text, s, "\\"+s)
	}
	return text
}
