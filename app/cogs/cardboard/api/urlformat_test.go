package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GetPostURL(t *testing.T) {
	post := &Post{ID: 12345}

	assert.Equal(
		t,
		"https://danbooru.donmai.us/posts/12345",
		GetPostURL(post),
	)
}

func Test_GetSearchURL(t *testing.T) {
	assert.Equal(
		t,
		"https://danbooru.donmai.us/posts?utf8=%E2%9C%93&tags=foo%20bar",
		GetSearchURL("foo bar"),
	)
}
