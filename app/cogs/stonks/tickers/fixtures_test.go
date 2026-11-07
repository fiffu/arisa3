package tickers

import (
	"context"
	"encoding/xml"
	"testing"
	"time"

	"github.com/fiffu/arisa3/app/cogs/stonks/tickers/api"

	"github.com/fiffu/arisa3/testfixtures"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// fixtureSlot is the slot the FRED fixtures were captured for: cosd = slot − 14 days.
var fixtureSlot = time.Date(2026, 11, 7, 8, 0, 0, 0, time.UTC)

func Test_fixture_SP500(t *testing.T) {
	tr := &testfixtures.Transport{Body: testfixtures.StonksFredSP500}
	src := NewSP500(api.NewFred(tr))

	res, err := src.Fetch(context.Background(), State[FredMeta]{Slot: fixtureSlot})

	require.NoError(t, err)
	assert.Equal(t, 7353.73, res.Raw)
	assert.Equal(t, OK, res.Outcome)
	assert.Equal(t, FredMeta{SourceDate: "2026-11-05"}, res.Patch)
}

func Test_fixture_DJIA(t *testing.T) {
	tr := &testfixtures.Transport{Body: testfixtures.StonksFredDJIA}
	src := NewDJIA(api.NewFred(tr))

	res, err := src.Fetch(context.Background(), State[FredMeta]{Slot: fixtureSlot})

	require.NoError(t, err)
	assert.Equal(t, 51285.73, res.Raw)
	assert.Equal(t, FredMeta{SourceDate: "2026-11-06"}, res.Patch)
}

func Test_fixture_Greed(t *testing.T) {
	tr := &testfixtures.Transport{Body: testfixtures.StonksGreed}
	src := &Greed{Client: api.NewGreed(tr)}

	res, err := src.Fetch(context.Background(), State[GreedMeta]{Slot: fixtureSlot})

	require.NoError(t, err)
	assert.Equal(t, 25.0, res.Raw)
	assert.Equal(t, GreedMeta{SourceTS: "1794009600"}, res.Patch)
}

func Test_fixture_Don(t *testing.T) {
	at := func(day, hour int) time.Time { return time.Date(2026, 11, day, hour, 0, 0, 0, time.UTC) }

	cases := []struct {
		name     string
		from, to time.Time
		raw      float64
		outcome  Outcome
	}{
		{
			// Posts at 18:05 (1302 chars), 18:09 (202) and 18:51 (423).
			name: "three posts", from: at(6, 16), to: at(6, 20),
			raw: 1927.0 / 3, outcome: OK,
		},
		{
			// One post of 100 chars at 14:33, and two media-only posts that do not count.
			name: "media-only posts skipped", from: at(7, 12), to: at(7, 16),
			raw: 100, outcome: OK,
		},
		{
			name: "no posts", from: at(7, 16), to: at(7, 20),
			raw: 42, outcome: NoNewData,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tr := &testfixtures.Transport{Body: testfixtures.StonksDonFeed}
			src := &Don{Client: api.NewTrumpsTruth(tr)}

			res, err := src.Fetch(context.Background(), State[DonMeta]{
				Slot: tc.to, PrevSlot: tc.from, PrevPrice: 42, HasPrev: true,
			})

			require.NoError(t, err)
			assert.InDelta(t, tc.raw, res.Raw, 1e-9)
			assert.Equal(t, tc.outcome, res.Outcome)
		})
	}
}

func Test_fixture_SP500_noNewData(t *testing.T) {
	tr := &testfixtures.Transport{Body: testfixtures.StonksFredSP500}
	src := NewSP500(api.NewFred(tr))

	res, err := src.Fetch(context.Background(), State[FredMeta]{Slot: fixtureSlot, Meta: FredMeta{SourceDate: "2026-11-05"}})

	require.NoError(t, err)
	assert.Equal(t, NoNewData, res.Outcome)
}

func Test_fixture_Greed_noNewData(t *testing.T) {
	tr := &testfixtures.Transport{Body: testfixtures.StonksGreed}
	src := &Greed{Client: api.NewGreed(tr)}

	res, err := src.Fetch(context.Background(), State[GreedMeta]{Slot: fixtureSlot, Meta: GreedMeta{SourceTS: "1794009600"}})

	require.NoError(t, err)
	assert.Equal(t, NoNewData, res.Outcome)
}

// The window includes a post exactly at PrevSlot and excludes a post exactly at Slot.
func Test_fixture_Don_windowEdges(t *testing.T) {
	tr := &testfixtures.Transport{Body: testfixtures.StonksDonFeed}
	src := &Don{Client: api.NewTrumpsTruth(tr)}

	res, err := src.Fetch(context.Background(), State[DonMeta]{
		PrevSlot: time.Date(2026, 11, 6, 18, 5, 27, 0, time.UTC),  // post of 1302 chars: included
		Slot:     time.Date(2026, 11, 6, 18, 51, 30, 0, time.UTC), // post of 423 chars: excluded
		HasPrev:  true,
	})

	require.NoError(t, err)
	assert.InDelta(t, (1302.0+202.0)/2, res.Raw, 1e-9)
}

// fixturePost returns the HTML of one post in the DON fixture, by its link.
func fixturePost(t *testing.T, link string) string {
	var feed struct {
		Items []struct {
			Link        string `xml:"link"`
			Description string `xml:"description"`
		} `xml:"channel>item"`
	}
	require.NoError(t, xml.Unmarshal([]byte(testfixtures.StonksDonFeed), &feed))
	for _, item := range feed.Items {
		if item.Link == link {
			return item.Description
		}
	}
	t.Fatalf("no post %s in fixture", link)
	return ""
}

// Expected lengths come from Python's html.parser, not from this package.
func Test_fixture_postLength(t *testing.T) {
	cases := []struct {
		name, link string
		want       int
	}{
		{"media only: one View image link", "https://www.trumpstruth.org/statuses/43504", 0},
		{"quote post with br, span and blockquote", "https://www.trumpstruth.org/statuses/43503", 100},
		{"text with an &amp; entity", "https://www.trumpstruth.org/statuses/43351", 1704},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, postLength(fixturePost(t, tc.link)))
		})
	}
}
