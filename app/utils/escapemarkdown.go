package utils

import (
	"strings"
)

func EscapeMarkdown(s string) string {
	for _, c := range "\\*>_:~`|" {
		ch := string(c)
		s = strings.ReplaceAll(s, ch, "\\"+ch)
	}
	return s
}
