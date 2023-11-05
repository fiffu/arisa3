package api

import (
	"context"
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

	actual, err := client.GetTags(context.Background(), []string{"capybara"})
	assert.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func Test_GetTags_Error(t *testing.T) {
	client := newClient("", "", 0)

	expectURL := "https://danbooru.donmai.us/tags.json?" +
		url.QueryEscape("search[name_comma]") +
		"=capybara"
	client.fetch = lib.StubTransportError(t, expectURL, assert.AnError)

	actual, err := client.GetTags(context.Background(), []string{"capybara"})
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

	actual, err := client.GetTagsMatching(context.Background(), "capy*")
	assert.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func Test_AutocompleteTag(t *testing.T) {
	client := newClient("username", "apikey", 0)

	expect := []*TagSuggestion{
		{
			Name:       "naidong_(artist)",
			Antecedent: "yin-ting tian",
			PostCount:  383,
			Link:       "https://danbooru.donmai.us/posts?tags=naidong_%28artist%29",
		},
		{
			Name:      "tingyun_(honkai:_star_rail)",
			PostCount: 376,
			Link:      "https://danbooru.donmai.us/posts?tags=tingyun_%28honkai%3A_star_rail%29",
		},
	}
	expectURL := "https://danbooru.donmai.us/autocomplete?" +
		url.QueryEscape("search[query]") + "=ting" +
		"&" + url.QueryEscape("search[type]") + "=tag_query" +
		"&version=1" +
		"&limit=10"

	stubResponse := `<ul>
		<li class="ui-menu-item" data-autocomplete-type="tag-word" data-autocomplete-value="naidong_(artist)">
			<div class="ui-menu-item-wrapper" tabindex="-1">
				<a class="tag-type-1" @click.prevent="" href="/posts?tags=naidong_%28artist%29">
					<span class="autocomplete-antecedent"><span>yin</span><span>-</span><b>ting</b><span> </span><span>tian</span></span>
					<span class="autocomplete-arrow">→</span> naidong (artist)
				</a>
				<span class="post-count">383</span>
			</div>
		</li>
		<li class="ui-menu-item" data-autocomplete-type="tag-word" data-autocomplete-value="tingyun_(honkai:_star_rail)">
			<div class="ui-menu-item-wrapper" tabindex="-1">
				<a class="tag-type-4" @click.prevent="" href="/posts?tags=tingyun_%28honkai%3A_star_rail%29">
					<b>ting</b><span>yun</span><span> (</span><span>honkai</span><span>: </span><span>star</span><span> </span><span>rail</span><span>)</span>
				</a>
				<span class="post-count">376</span>
			</div>
		</li>
	</ul>`
	client.fetch = lib.StubHTMLFetcher(t, expectURL, http.StatusOK, stubResponse)

	actual, err := client.AutocompleteTag(context.Background(), "ting")
	assert.NoError(t, err)
	assert.Equal(t, expect, actual)
}
