package api

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/carlmjohnson/requests"
)

const greedURL = "https://api.alternative.me/fng/"

// Greed reads the Crypto Fear & Greed Index. The index changes once a day.
type Greed struct {
	transport http.RoundTripper
}

func NewGreed(transport http.RoundTripper) *Greed {
	return &Greed{transport: transport}
}

type greedResponse struct {
	Data []struct {
		Value     string `json:"value"`
		Timestamp string `json:"timestamp"`
	} `json:"data"`
}

// Latest returns the latest index value (0 to 100), and its Unix timestamp as the API writes it.
func (g *Greed) Latest(ctx context.Context) (value float64, timestamp string, err error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var res greedResponse
	err = requests.URL(greedURL).
		Param("limit", "1").
		UserAgent(userAgent).
		Transport(g.transport).
		ToJSON(&res).
		Fetch(ctx)
	if err != nil {
		return 0, "", fmt.Errorf("fear and greed index: %w", err)
	}
	if len(res.Data) == 0 {
		return 0, "", fmt.Errorf("fear and greed index: no values in response")
	}

	v, err := strconv.ParseFloat(res.Data[0].Value, 64)
	if err != nil {
		return 0, "", fmt.Errorf("fear and greed index: bad value %q: %w", res.Data[0].Value, err)
	}
	return v, res.Data[0].Timestamp, nil
}
