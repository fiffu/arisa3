package api

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/fiffu/arisa3/lib"
	"github.com/stretchr/testify/assert"
)

func Test_GetTags(t *testing.T) {
	client := newClient("username", "apikey", 0)

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

	actual, err := client.GetTags([]string{"capybara"})
	assert.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func Test_GetTags_Error(t *testing.T) {
	client := newClient("", "", 0)

	expectURL := "https://danbooru.donmai.us/tags.json?" +
		url.QueryEscape("search[name_comma]") +
		"=capybara"
	client.fetch = lib.StubTransportError(t, expectURL, assert.AnError)

	actual, err := client.GetTags([]string{"capybara"})
	assert.Error(t, err)
	assert.Nil(t, actual)
}

func Test_GetTags_multipleTags(t *testing.T) {
	client := newClient("", "", 0)

	expectURL := "https://danbooru.donmai.us/tags.json?" +
		url.QueryEscape("search[name_comma]") +
		"=" +
		url.QueryEscape("capybara,asdfqwe")
	client.fetch = lib.StubTransportError(t, expectURL, assert.AnError)

	actual, err := client.GetTags([]string{"capybara", "asdfqwe"})
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
	actual := indexTagsByName(tags)
	assert.Equal(t, expect, actual)
}

func Test_GetTagsMatching(t *testing.T) {
	client := newClient("username", "apikey", 0)

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

	actual, err := client.GetTagsMatching("capy*")
	assert.NoError(t, err)
	assert.Equal(t, expect, actual)
}
