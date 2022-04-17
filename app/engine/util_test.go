package engine

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_PrettifyCustomEmoji(t *testing.T) {
	type testCase struct {
		in     string
		expect string
	}
	tests := []testCase{
		{"<:birb:924875584004321361>", ":birb:"},
		{"<a:aGES_Kek:741032442814791702>", ":aGES_Kek:"},
		{":sob:", ":sob:"},
		{"123456", "123456"},
		{"<:birb:924875584004321361> <:birb:924875584004321361>", ":birb: :birb:"},
		{"testing: :sob: <:birb:924875584004321361>", "testing: :sob: :birb:"},
		{"<:    irb:924875584004321361>", "<:    irb:924875584004321361>"},
	}
	for _, tc := range tests {
		name := fmt.Sprintf(`PrettifyCustomEmoji("%s")`, tc.in)
		t.Run(name, func(t *testing.T) {
			actual := PrettifyCustomEmoji(tc.in)
			assert.Equal(t, tc.expect, actual)
		})
	}
}
