package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_TagsSorter(t *testing.T) {
	factory := func(names ...string) []*Tag {
		tagList := make([]*Tag, len(names))
		for i, name := range names {
			tagList[i] = &Tag{ID: i, Name: name}
		}
		return tagList
	}

	testCases := []struct {
		desc        string
		input       []*Tag
		comparer    TagComparer
		expectOrder []string
	}{
		{
			desc:        "Sorted ByAlphabeticalOrder",
			input:       factory("bacon", "spam", "ham", "eggs"),
			comparer:    ByAlphabeticalOrder,
			expectOrder: []string{"bacon", "eggs", "ham", "spam"},
		},
		{
			desc:        "Sorted ByTagLength",
			input:       factory("bacon", "spam", "ham"),
			comparer:    ByTagLength,
			expectOrder: []string{"ham", "spam", "bacon"},
		},
		{
			desc:        "Sorted ByTagLength uses alphabetical order for tie-breaker",
			input:       factory("cantankerous", "candle", "cancer"),
			comparer:    ByTagLength,
			expectOrder: []string{"cancer", "candle", "cantankerous"},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			data := tc.input
			cmp := tc.comparer

			actual := TagsSorter{data, cmp}.Sorted()
			for i, expectName := range tc.expectOrder {
				assert.Equal(t, expectName, actual[i].Name)
			}
		})
	}
}
