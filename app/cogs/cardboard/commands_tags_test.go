package cardboard

import (
	"testing"

	"github.com/fiffu/arisa3/app/cogs/cardboard/api"
	"github.com/stretchr/testify/assert"
)

func Test_formatSuggestion(t *testing.T) {
	testCases := map[string]*api.TagSuggestion{
		"yin-ting tian → [`naidong_(artist)`](https://danbooru.donmai.us/posts?tags=naidong_%28artist%29) (383)": {
			Name:       "naidong_(artist)",
			Antecedent: "yin-ting tian",
			PostCount:  383,
			Link:       "https://danbooru.donmai.us/posts?tags=naidong_%28artist%29",
		},
		"[`tingyun_(honkai:_star_rail)`](https://danbooru.donmai.us/posts?tags=tingyun_%28honkai%3A_star_rail%29) (376)": {
			Name:      "tingyun_(honkai:_star_rail)",
			PostCount: 376,
			Link:      "https://danbooru.donmai.us/posts?tags=tingyun_%28honkai%3A_star_rail%29",
		},
	}
	for expect, input := range testCases {
		t.Run(expect, func(t *testing.T) {
			actual := formatSuggestion(input)
			assert.Equal(t, expect, actual)
		})
	}
}
