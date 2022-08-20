package api

import (
	"context"
	"net/http"
	"testing"

	"github.com/carlmjohnson/requests"
	"github.com/fiffu/arisa3/lib"
	"github.com/stretchr/testify/assert"
)

var samplePost = &Post{
	ID:            5605442,
	MD5:           "e10c930adeb9d9bb8232ad1f9c7185de",
	FileExt:       "jpg",
	AllTags:       "1girl :> capybara crab digitan_(porforever) original porforever",
	GeneralTags:   "1girl :> capybara crab",
	CharacterTags: "digitan_(porforever)",
	CopyrightTags: "original",
	ArtistTags:    "porforever",
	File:          "https://cdn.donmai.us/original/e1/0c/e10c930adeb9d9bb8232ad1f9c7185de.jpg",
	FileLarge:     "https://cdn.donmai.us/sample/e1/0c/sample-e10c930adeb9d9bb8232ad1f9c7185de.jpg",
	FilePreview:   "https://cdn.donmai.us/preview/e1/0c/e10c930adeb9d9bb8232ad1f9c7185de.jpg",
}

func Test_GetPosts_Error(t *testing.T) {
	client := &client{}

	client.fetch = func(ctx context.Context, builder *requests.Builder) error {
		return assert.AnError
	}

	posts, err := client.GetPosts([]string{""})
	assert.Nil(t, posts)
	assert.Error(t, err)
}

func Test_GetPosts(t *testing.T) {
	client := newClient("username", "apikey", 0)

	expect := []*Post{samplePost}
	expectURL := "https://danbooru.donmai.us/posts.json?limit=100&tags=capybara"

	stubJSON := `[
		{
			"id": 5605442,
			"md5": "e10c930adeb9d9bb8232ad1f9c7185de",
			"file_ext": "jpg",
			"tag_string": "1girl :> capybara crab digitan_(porforever) original porforever",
			"tag_string_general": "1girl :> capybara crab",
			    "tag_string_character": "digitan_(porforever)",
			    "tag_string_copyright": "original",
			    "tag_string_artist": "porforever",
			"file_url": "https://cdn.donmai.us/original/e1/0c/e10c930adeb9d9bb8232ad1f9c7185de.jpg",
			"large_file_url": "https://cdn.donmai.us/sample/e1/0c/sample-e10c930adeb9d9bb8232ad1f9c7185de.jpg",
			"preview_file_url": "https://cdn.donmai.us/preview/e1/0c/e10c930adeb9d9bb8232ad1f9c7185de.jpg"
		}
	]`
	client.fetch = lib.StubJSONFetcher(t, expectURL, http.StatusOK, stubJSON)

	actual, err := client.GetPosts([]string{"capybara"})

	assert.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func Test_TagsList(t *testing.T) {
	expect := []string{"1girl", ":>", "capybara", "crab", "digitan_(porforever)", "original", "porforever"}
	actual := samplePost.TagsList()
	assert.ElementsMatch(t, expect, actual)
}

func Test_GetFileURL(t *testing.T) {
	testCases := []struct {
		desc                         string
		File, FileLarge, FilePreview string
		expect                       string
	}{
		{
			desc:        "Default to File first",
			File:        "https://cdn.donmai.us/original/e1/0c/e10c930adeb9d9bb8232ad1f9c7185de.jpg",
			FileLarge:   "https://cdn.donmai.us/sample/e1/0c/sample-e10c930adeb9d9bb8232ad1f9c7185de.jpg",
			FilePreview: "https://cdn.donmai.us/preview/e1/0c/e10c930adeb9d9bb8232ad1f9c7185de.jpg",
			expect:      "https://cdn.donmai.us/original/e1/0c/e10c930adeb9d9bb8232ad1f9c7185de.jpg",
		},
		{
			desc:        "If File missing, fallback to FileLarge",
			File:        "",
			FileLarge:   "https://cdn.donmai.us/sample/e1/0c/sample-e10c930adeb9d9bb8232ad1f9c7185de.jpg",
			FilePreview: "https://cdn.donmai.us/preview/e1/0c/e10c930adeb9d9bb8232ad1f9c7185de.jpg",
			expect:      "https://cdn.donmai.us/sample/e1/0c/sample-e10c930adeb9d9bb8232ad1f9c7185de.jpg",
		},
		{
			desc:        "If File and FileLarge missing, fallback to FilePreview",
			File:        "",
			FileLarge:   "",
			FilePreview: "https://cdn.donmai.us/preview/e1/0c/e10c930adeb9d9bb8232ad1f9c7185de.jpg",
			expect:      "https://cdn.donmai.us/preview/e1/0c/e10c930adeb9d9bb8232ad1f9c7185de.jpg",
		},
		{
			desc:   "If all missing, expect empty string",
			File:   "",
			expect: "",
		},
		{
			desc:   "If host is missing from URL, it should be fixed",
			File:   "/data/foo.jpg",
			expect: "https://danbooru.donmai.us/data/foo.jpg",
		},
		{
			desc:   "If //data is in url, it should be fixed",
			File:   "https://cdn.donmai.us//data/foo.jpg",
			expect: "https://cdn.donmai.us/data/foo.jpg",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			post := &Post{
				File:        tc.File,
				FileLarge:   tc.FileLarge,
				FilePreview: tc.FilePreview,
			}
			actual := post.GetFileURL()
			assert.Equal(t, tc.expect, actual)
		})
	}
}
