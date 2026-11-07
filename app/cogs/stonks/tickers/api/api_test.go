package api

import (
	"context"
	"net/url"
	"testing"
	"time"

	"github.com/fiffu/arisa3/testfixtures"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// The FRED fixtures were captured with since = 2026-10-24.
var since = time.Date(2026, 10, 24, 8, 0, 0, 0, time.UTC)

func Test_Fred_LatestClose(t *testing.T) {
	cases := []struct {
		name, seriesID, fixture, date string
		value                         float64
	}{
		{"SP500", "SP500", testfixtures.StonksFredSP500, "2026-11-05", 7353.73},
		{"DJIA", "DJIA", testfixtures.StonksFredDJIA, "2026-11-06", 51285.73},
		// The last row, 2026-01-01, is a market holiday with no value.
		{"ends on a holiday", "SP500", testfixtures.StonksFredSP500Holiday, "2025-12-31", 6845.50},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tr := &testfixtures.Transport{Body: tc.fixture}

			date, value, err := NewFred(tr).LatestClose(context.Background(), tc.seriesID, since)

			require.NoError(t, err)
			assert.Equal(t, tc.date, date)
			assert.Equal(t, tc.value, value)

			reqs := tr.Requests()
			require.Len(t, reqs, 1)
			assert.Equal(t, "fred.stlouisfed.org", reqs[0].URL.Host)
			assert.Equal(t, "/graph/fredgraph.csv", reqs[0].URL.Path)
			assert.Equal(t, url.Values{"id": {tc.seriesID}, "cosd": {"2026-10-24"}}, reqs[0].URL.Query())
			assert.Equal(t, userAgent, reqs[0].Header.Get("User-Agent"))
		})
	}
}

func Test_Greed_Latest(t *testing.T) {
	tr := &testfixtures.Transport{Body: testfixtures.StonksGreed}

	value, ts, err := NewGreed(tr).Latest(context.Background())

	require.NoError(t, err)
	assert.Equal(t, 25.0, value)
	assert.Equal(t, "1794009600", ts)

	reqs := tr.Requests()
	require.Len(t, reqs, 1)
	assert.Equal(t, "https://api.alternative.me/fng/?limit=1", reqs[0].URL.String())
	assert.Equal(t, userAgent, reqs[0].Header.Get("User-Agent"))
}

func Test_TrumpsTruth_Posts(t *testing.T) {
	tr := &testfixtures.Transport{Body: testfixtures.StonksDonFeed}

	posts, err := NewTrumpsTruth(tr).Posts(context.Background())

	require.NoError(t, err)
	require.Len(t, posts, 100)
	assert.Equal(t, time.Date(2026, 11, 7, 14, 42, 52, 0, time.UTC), posts[0].At.UTC())
	assert.Equal(t, time.Date(2026, 11, 3, 13, 21, 47, 0, time.UTC), posts[99].At.UTC())
	assert.Contains(t, posts[0].HTML, "View image")

	reqs := tr.Requests()
	require.Len(t, reqs, 1)
	assert.Equal(t, "https://www.trumpstruth.org/feed", reqs[0].URL.String())
	assert.Equal(t, userAgent, reqs[0].Header.Get("User-Agent"))
}
