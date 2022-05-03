package api

import (
	"context"
	"strings"
	"time"

	"github.com/carlmjohnson/requests"
	"github.com/rs/zerolog/log"
)

const (
	WildcardCharacter = "*"
	apiHost           = "danbooru.donmai.us"
	apiHostHTTP       = "https://" + apiHost
)

var (
	MediaFileExts = []string{"png", "jpg", "jpeg", "gif"}
)

type client struct {
	auth        bool
	username    string
	apiKey      string
	host        string
	timeoutSecs int
}

func NewClient(username, apiKey string, timeoutSecs int) IClient {
	auth := false
	switch {
	case username != "" && apiKey != "":
		auth = true
	case username != "" && apiKey == "":
		log.Warn().Msg("Client username was provided, but API key was not. Skipping auth...")
	case username == "" && apiKey != "":
		log.Warn().Msg("Client API key was provided, but username was not. Skipping auth...")
	}

	if timeoutSecs == 0 {
		timeoutSecs = 2
	}
	return &client{
		auth,
		username,
		apiKey,
		apiHost,
		timeoutSecs,
	}
}

func commaJoin(strs []string) string {
	return strings.Join(strs, ",")
}

func (c *client) context() (context.Context, context.CancelFunc) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Duration(c.timeoutSecs)*time.Second)
	return ctx, cancel
}

func (c *client) baseRequest() *requests.Builder {
	b := &requests.Builder{}
	b.Host(c.host)
	if c.auth {
		b.BasicAuth(c.username, c.apiKey)
	}
	return b
}

func (c *client) requestPosts() *requests.Builder {
	return c.baseRequest().Path("posts.json")
}

func (c *client) requestTags() *requests.Builder {
	return c.baseRequest().Path("tags.json")
}
