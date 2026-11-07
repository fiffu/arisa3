package api

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/carlmjohnson/requests"
)

const trumpsTruthURL = "https://www.trumpstruth.org/feed"

// TrumpsTruth reads a third-party archive of Donald Trump's Truth Social posts.
// The feed has the latest 100 posts.
type TrumpsTruth struct {
	transport http.RoundTripper
}

func NewTrumpsTruth(transport http.RoundTripper) *TrumpsTruth {
	return &TrumpsTruth{transport: transport}
}

type Post struct {
	At   time.Time
	HTML string
}

type rssFeed struct {
	Items []struct {
		PubDate     string `xml:"pubDate"`
		Description string `xml:"description"`
	} `xml:"channel>item"`
}

// Posts skips items with a date it cannot parse.
func (c *TrumpsTruth) Posts(ctx context.Context) ([]Post, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var body []byte
	err := requests.URL(trumpsTruthURL).
		UserAgent(userAgent).
		Transport(c.transport).
		Handle(readBody(&body)).
		Fetch(ctx)
	if err != nil {
		return nil, fmt.Errorf("trumpstruth feed: %w", err)
	}

	var feed rssFeed
	if err := xml.Unmarshal(body, &feed); err != nil {
		return nil, fmt.Errorf("trumpstruth feed: %w", err)
	}
	if len(feed.Items) == 0 {
		return nil, fmt.Errorf("trumpstruth feed: no items")
	}

	posts := make([]Post, 0, len(feed.Items))
	for _, item := range feed.Items {
		at, err := time.Parse(time.RFC1123Z, strings.TrimSpace(item.PubDate))
		if err != nil {
			continue
		}
		posts = append(posts, Post{At: at, HTML: item.Description})
	}
	return posts, nil
}
