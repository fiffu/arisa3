package testfixtures

import (
	_ "embed"
	"io"
	"net/http"
	"strings"
	"sync"
)

//go:embed downbooru.html
var Downbooru string

// Captured on 2026-11-08 from
// https://fred.stlouisfed.org/graph/fredgraph.csv?cosd=2026-10-24&id=SP500
//
//go:embed stonks_fred_sp500.csv
var StonksFredSP500 string

// Captured on 2026-11-08 from
// https://fred.stlouisfed.org/graph/fredgraph.csv?cosd=2026-10-24&id=DJIA
//
//go:embed stonks_fred_djia.csv
var StonksFredDJIA string

// Captured on 2026-11-08 from https://api.alternative.me/fng/?limit=1
//
//go:embed stonks_greed.json
var StonksGreed string

// Captured on 2026-11-08 from https://www.trumpstruth.org/feed.
// It holds the 100 posts from 2026-11-03 13:21 to 2026-11-07 14:42 UTC.
//
//go:embed stonks_don_feed.xml
var StonksDonFeed string

// Captured on 2026-11-08 from
// https://fred.stlouisfed.org/graph/fredgraph.csv?cosd=2025-12-22&coed=2026-01-01&id=SP500
// The last row, 2026-01-01, is a market holiday with no value.
//
//go:embed stonks_fred_sp500_holiday.csv
var StonksFredSP500Holiday string

// Transport answers every request with status 200 and Body, and records each request.
type Transport struct {
	Body string

	mu       sync.Mutex
	requests []*http.Request
}

func (tr *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	tr.mu.Lock()
	tr.requests = append(tr.requests, req)
	tr.mu.Unlock()
	return &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{},
		Body:       io.NopCloser(strings.NewReader(tr.Body)),
		Request:    req,
	}, nil
}

func (tr *Transport) Requests() []*http.Request {
	tr.mu.Lock()
	defer tr.mu.Unlock()
	return append([]*http.Request(nil), tr.requests...)
}
