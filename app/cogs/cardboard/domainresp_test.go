package cardboard

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/fiffu/arisa3/app/cogs/cardboard/api"
	"github.com/fiffu/arisa3/app/database"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func newTestPost() *api.Post {
	return &api.Post{
		FileExt:       "png",
		AllTags:       "general character copyright artist",
		GeneralTags:   "general",
		CharacterTags: "character",
		CopyrightTags: "copyright",
		ArtistTags:    "artist",
		File:          "https://foo.png",
	}
}

func Test_formatZeroResults(t *testing.T) {
	ctrl := gomock.NewController(t)
	db := database.NewMockIDatabase(ctrl)
	d := NewDomain(db, &Config{})

	term := "testing"
	query := NewMockIQueryPosts(ctrl)
	query.EXPECT().Term().Return(term)

	embed := d.formatZeroResults(query)
	assert.Contains(t, embed.Data().Description, term)
}

func Test_formatResult(t *testing.T) {
	ctrl := gomock.NewController(t)
	db := database.NewMockIDatabase(ctrl)
	client := api.NewMockIClient(ctrl)
	client.EXPECT().
		GetTags(gomock.Any(), gomock.Any()).
		AnyTimes().
		DoAndReturn(func(_ context.Context, _ []string) (map[string]*api.Tag, error) {
			return map[string]*api.Tag{
				"general":   {ID: 0, Name: "general", PostCount: 999},
				"artist":    {ID: 1, Name: "artist", PostCount: 999},
				"character": {ID: 2, Name: "character", PostCount: 999},
				"copyright": {ID: 3, Name: "copyright", PostCount: 999},
			}, nil
		})

	d := NewDomain(db, &Config{})
	d.client = client

	testTerm := "testing"
	query := NewQuery(testTerm)

	testPost := newTestPost()
	posts := []*api.Post{testPost}

	embed, err := d.formatResult(context.Background(), query, posts)
	assert.NoError(t, err)
	assert.Equal(t, embed.Data().Color, embedColour)
	assert.Contains(t, embed.Data().Title, testPost.CharacterTags)
	assert.Contains(t, embed.Data().Image.URL, testPost.File)
	assert.Contains(t, embed.Data().Fields[0].Value, testPost.ArtistTags)
	assert.Contains(t, embed.Data().Fields[1].Value, testPost.CopyrightTags)
	assert.Contains(t, embed.Data().Fields[2].Value, testPost.File)
	assert.Contains(t, embed.Data().Footer.Text, testTerm)
}

func Test_formatResult_shouldOnlyGetTagsForTermArtistCopyright(t *testing.T) {
	ctrl := gomock.NewController(t)
	db := database.NewMockIDatabase(ctrl)
	client := api.NewMockIClient(ctrl)
	client.EXPECT().
		GetTags(gomock.Any(), []string{"testing", "artist", "copyright"}).
		AnyTimes().
		DoAndReturn(func(_ context.Context, _ []string) (map[string]*api.Tag, error) {
			return map[string]*api.Tag{
				"general":   {ID: 0, Name: "general", PostCount: 999},
				"artist":    {ID: 1, Name: "artist", PostCount: 999},
				"character": {ID: 2, Name: "character", PostCount: 999},
				"copyright": {ID: 3, Name: "copyright", PostCount: 999},
			}, nil
		})

	d := NewDomain(db, &Config{})
	d.client = client

	testTerm := "testing"
	query := NewQuery(testTerm)

	testPost := newTestPost()
	posts := []*api.Post{testPost}

	embed, err := d.formatResult(context.Background(), query, posts)
	assert.NoError(t, err)
	assert.Contains(t, embed.Data().Fields[0].Value, testPost.ArtistTags)
	assert.Contains(t, embed.Data().Fields[0].Value, "999")
	assert.Contains(t, embed.Data().Fields[1].Value, testPost.CopyrightTags)
	assert.Contains(t, embed.Data().Fields[1].Value, "999")
}

func Test_formatResult_shouldNotFailIfGetTagsErrors(t *testing.T) {
	ctrl := gomock.NewController(t)
	db := database.NewMockIDatabase(ctrl)
	client := api.NewMockIClient(ctrl)
	client.EXPECT().
		GetTags(gomock.Any(), []string{"testing", "artist", "copyright"}).
		AnyTimes().
		DoAndReturn(func(_ context.Context, _ []string) (map[string]*api.Tag, error) {
			return nil, assert.AnError
		})

	d := NewDomain(db, &Config{})
	d.client = client

	testTerm := "testing"
	query := NewQuery(testTerm)

	testPost := newTestPost()
	posts := []*api.Post{testPost}

	embed, err := d.formatResult(context.Background(), query, posts)
	assert.NoError(t, err)
	assert.Contains(t, embed.Data().Fields[0].Value, testPost.ArtistTags)
	assert.NotContains(t, embed.Data().Fields[0].Value, "999")
	assert.Contains(t, embed.Data().Fields[1].Value, testPost.CopyrightTags)
	assert.NotContains(t, embed.Data().Fields[1].Value, "999")
}

func Test_embedTitle(t *testing.T) {
	type testCase struct {
		name              string
		characters        string
		artists           string
		expectContains    []string
		expectNotContains []string
	}

	tests := []testCase{
		{
			name:           "regular case where characters and artists fit within the length limit",
			characters:     "character",
			artists:        "artist",
			expectContains: []string{"character drawn by artist"},
		},
		{
			name:              "too many characters should fallback to artists",
			characters:        strings.Repeat("character", 255),
			artists:           "artist",
			expectContains:    []string{"artist"},
			expectNotContains: []string{"character"},
		},
		{
			name:              "no artist should omit 'drawn by' in title",
			characters:        "character",
			artists:           "",
			expectContains:    []string{"character"},
			expectNotContains: []string{"drawn by"},
		},
		{
			name:              "artist name too long should omit 'drawn by' in title",
			characters:        "character",
			artists:           strings.Repeat("artist", 255),
			expectContains:    []string{"character"},
			expectNotContains: []string{"drawn by"},
		},
		{
			name:           "no character, no artist",
			characters:     "",
			artists:        "",
			expectContains: []string{"Picture #"},
		},
		{
			name:           "https://github.com/fiffu/arisa3/issues/151",
			characters:     "shun_(blue_archive) shun_(small)_(blue_archive)",
			artists:        "zhnyy3",
			expectContains: []string{"shun (blue archive) and shun (small) (blue archive) drawn by zhnyy3"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			post := &api.Post{
				CharacterTags: tc.characters,
				ArtistTags:    tc.artists,
			}
			actual := embedTitle(post)
			for _, x := range tc.expectContains {
				assert.Contains(t, actual, x)
			}
			for _, x := range tc.expectNotContains {
				assert.NotContains(t, actual, x)
			}
		})
	}
}

func Test_joinWithTail(t *testing.T) {
	testCases := []struct {
		input  []string
		expect string
	}{
		{
			input:  []string{"a"},
			expect: "a",
		},
		{
			input:  []string{"a", "b"},
			expect: "a and b",
		},
		{
			input:  []string{"a", "b", "c"},
			expect: "a, b and c",
		},
		{
			input:  []string{"a", "b", "c", "d"},
			expect: "a, b, c and d",
		},
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprint(tc.input), func(t *testing.T) {
			actual := joinWithTail(
				tc.input,
				", ",
				" and ",
			)
			assert.Equal(t, tc.expect, actual)
		})
	}
}

func Test_embedTitleArtists(t *testing.T) {
	testCases := []struct {
		desc    string
		artists string
		expect  string
	}{
		{
			desc:    "no artists",
			artists: "",
			expect:  "",
		},
		{
			desc:    "1 artist",
			artists: "artist1",
			expect:  " drawn by artist1",
		},
		{
			desc:    "multiple artists",
			artists: "artist1 artist2 artist3",
			expect:  " drawn by artist1 and others",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			post := &api.Post{ArtistTags: tc.artists}
			actual := embedTitleArtists(post)
			assert.Equal(t, tc.expect, actual)
		})
	}
}

func Test_embedFieldTags(t *testing.T) {
	testCases := []struct {
		desc           string
		tags           []string
		expectContains []string
	}{
		{
			desc:           "Regular case with tag that has a count",
			tags:           []string{"foo"},
			expectContains: []string{"foo", "https", "999"},
		},
		{
			desc:           "Tag not in tagData",
			tags:           []string{"banana"},
			expectContains: []string{"banana", "https"},
		},
		{
			desc:           "No tags",
			tags:           []string{},
			expectContains: []string{"(none)"},
		},
	}

	tagData := map[string]*api.Tag{
		"foo": {PostCount: 999},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			actual := embedFieldTags(tc.tags, tagData)
			for _, x := range tc.expectContains {
				assert.Contains(t, actual, x)
			}
		})
	}
}

func Test_splitTags(t *testing.T) {
	testCases := []struct {
		desc   string
		str    string
		expect []string
	}{
		{
			desc:   "Single tag",
			str:    "tag1",
			expect: []string{"tag1"},
		},
		{
			desc:   "Multiple tags",
			str:    "tag1 tag2",
			expect: []string{"tag1", "tag2"},
		},
		{
			desc:   "No tags",
			str:    "",
			expect: []string{},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			actual := splitTags(tc.str)
			assert.Equal(t, len(tc.expect), len(actual))
			assert.ElementsMatch(t, tc.expect, actual)
		})
	}
}
