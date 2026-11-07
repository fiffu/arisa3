package api

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/carlmjohnson/requests"
)

const fredURL = "https://fred.stlouisfed.org/graph/fredgraph.csv"

// Fred needs no API key. FRED updates each series once per trading day.
type Fred struct {
	transport http.RoundTripper
}

func NewFred(transport http.RoundTripper) *Fred {
	return &Fred{transport: transport}
}

// LatestClose returns the last value of a series, and its date such as "2026-11-05".
// Without since, FRED returns ten years of values.
func (f *Fred) LatestClose(ctx context.Context, seriesID string, since time.Time) (date string, value float64, err error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var body []byte
	err = requests.URL(fredURL).
		Param("id", seriesID).
		Param("cosd", since.Format("2006-01-02")).
		UserAgent(userAgent).
		Transport(f.transport).
		Handle(readBody(&body)).
		Fetch(ctx)
	if err != nil {
		return "", 0, fmt.Errorf("FRED %s: %w", seriesID, err)
	}

	date, value, err = f.lastValue(body)
	if err != nil {
		return "", 0, fmt.Errorf("FRED %s: %w", seriesID, err)
	}
	return date, value, nil
}

// FRED writes "." or nothing on market holidays.
func (f *Fred) lastValue(body []byte) (date string, value float64, err error) {
	rows, err := csv.NewReader(bytes.NewReader(body)).ReadAll()
	if err != nil {
		return "", 0, err
	}
	for i := len(rows) - 1; i >= 1; i-- { // row 0 is the header
		row := rows[i]
		if len(row) < 2 {
			continue
		}
		v, err := strconv.ParseFloat(strings.TrimSpace(row[1]), 64)
		if err != nil {
			continue
		}
		if _, err := time.Parse("2006-01-02", row[0]); err != nil {
			continue
		}
		return row[0], v, nil
	}
	return "", 0, fmt.Errorf("no values in response")
}
