package api

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/fiffu/arisa3/lib"
	"github.com/fiffu/arisa3/testfixtures"
	"github.com/stretchr/testify/assert"
)

func Test_GetTags(t *testing.T) {
	client := newClient("username", "apikey", "123456", 0)

	expect := map[string]*Tag{
		"capybara": {
			ID:        406745,
			Name:      "capybara",
			PostCount: 201,
		},
	}
	expectURL := "https://danbooru.donmai.us/tags.json?" +
		url.QueryEscape("search[name_comma]") +
		"=capybara"

	stubJSON := `[
		{
			"id": 406745,
			"name": "capybara",
			"post_count": 201,
			"category": 0,
			"created_at": "2013-02-28T02:29:28.204-05:00",
			"updated_at": "2019-09-02T09:55:08.730-04:00",
			"is_locked": false,
			"is_deprecated": false
		}
	]`
	client.fetch = lib.StubJSONFetcher(t, expectURL, http.StatusOK, stubJSON)

	actual, err := client.GetTags(context.Background(), []string{"capybara"})
	assert.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func Test_GetTags_Error(t *testing.T) {
	client := newClient("", "", "", 0)

	expectURL := "https://danbooru.donmai.us/tags.json?" +
		url.QueryEscape("search[name_comma]") +
		"=capybara"
	client.fetch = lib.StubTransportError(t, expectURL, assert.AnError)

	actual, err := client.GetTags(context.Background(), []string{"capybara"})
	assert.Error(t, err)
	assert.Nil(t, actual)
}

func Test_GetTags_multipleTags(t *testing.T) {
	client := newClient("", "", "", 0)

	expectURL := "https://danbooru.donmai.us/tags.json?" +
		url.QueryEscape("search[name_comma]") +
		"=" +
		url.QueryEscape("capybara,asdfqwe")
	client.fetch = lib.StubTransportError(t, expectURL, assert.AnError)

	actual, err := client.GetTags(context.Background(), []string{"capybara", "asdfqwe"})
	assert.Error(t, err)
	assert.Nil(t, actual)
}

func Test_indexTagsByName(t *testing.T) {
	tags := []*Tag{
		{ID: 406745, Name: "capybara", PostCount: 1},
		{ID: 406745, Name: "capybara", PostCount: 1}, // same name, same ID should be ignored
		{ID: 1, Name: "capybara", PostCount: 333},    // same name, diff ID should overwrite
	}
	expect := map[string]*Tag{
		"capybara": {
			ID:        1,
			Name:      "capybara",
			PostCount: 333,
		},
	}
	actual := indexTagsByName(context.Background(), tags)
	assert.Equal(t, expect, actual)
}

func Test_GetTagsMatching(t *testing.T) {
	client := newClient("username", "apikey", "123456", 0)

	expect := []*Tag{{
		ID:        406745,
		Name:      "capybara",
		PostCount: 201,
	}}
	expectURL := "https://danbooru.donmai.us/tags.json?" +
		url.QueryEscape("search[name_matches]") +
		"=" +
		url.QueryEscape("capy*")

	stubJSON := `[
		{
			"id": 406745,
			"name": "capybara",
			"post_count": 201,
			"category": 0,
			"created_at": "2013-02-28T02:29:28.204-05:00",
			"updated_at": "2019-09-02T09:55:08.730-04:00",
			"is_locked": false,
			"is_deprecated": false
		}
	]`
	client.fetch = lib.StubJSONFetcher(t, expectURL, http.StatusOK, stubJSON)

	actual, err := client.GetTagsMatching(context.Background(), "capy*")
	assert.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func Test_AutocompleteTag(t *testing.T) {
	client := newClient("username", "apikey", "123456", 0)

	expect := []*TagSuggestion{
		{
			Name:       "naidong_(artist)",
			Antecedent: "yin-ting_tian",
			PostCount:  383,
			Link:       "https://danbooru.donmai.us/posts?tags=naidong_%28artist%29",
		},
		{
			Name:      "tingyun_(honkai:_star_rail)",
			PostCount: 376,
			Link:      "https://danbooru.donmai.us/posts?tags=tingyun_%28honkai%3A_star_rail%29",
		},
	}
	expectURL := "https://danbooru.donmai.us/autocomplete.json?" +
		url.QueryEscape("search[query]") + "=ting" +
		"&" + url.QueryEscape("search[type]") + "=tag_query" +
		"&limit=10"

	stubJSON := `[
		{"type":"tag-word","label":"naidong (artist)","value":"naidong_(artist)","category":1,"post_count":383,"antecedent":"yin-ting_tian","tag":{"id":1,"name":"naidong_(artist)"}},
		{"type":"tag-word","label":"tingyun (honkai: star rail)","value":"tingyun_(honkai:_star_rail)","category":4,"post_count":376,"tag":{"id":2,"name":"tingyun_(honkai:_star_rail)"}}
	]`
	client.fetch = lib.StubJSONFetcher(t, expectURL, http.StatusOK, stubJSON)

	actual, err := client.AutocompleteTag(context.Background(), "ting")
	assert.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func Test_AutocompleteTag_maintenance(t *testing.T) {
	client := newClient("username", "apikey", "123456", 0)
	client.fetch = lib.StubHTMLFetcher(t,
		"https://danbooru.donmai.us/autocomplete.json?limit=10&search%5Bquery%5D=ting&search%5Btype%5D=tag_query",
		http.StatusServiceUnavailable,
		testfixtures.Downbooru,
	)

	_, err := client.AutocompleteTag(context.Background(), "ting")
	assert.Equal(t, ErrUnderMaintenance, err)
}

func Test_AutocompleteTag_blocked(t *testing.T) {
	client := newClient("username", "apikey", "123456", 0)
	client.fetch = lib.StubHTMLFetcher(t,
		"https://danbooru.donmai.us/autocomplete.json?limit=10&search%5Bquery%5D=ting&search%5Btype%5D=tag_query",
		http.StatusForbidden,
		"<html><title>Just a moment...</title></html>",
	)

	_, err := client.AutocompleteTag(context.Background(), "ting")
	assert.Error(t, err)
}
