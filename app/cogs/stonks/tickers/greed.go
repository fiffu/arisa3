package tickers

import (
	"context"

	"github.com/fiffu/arisa3/app/cogs/stonks/tickers/api"
)

type GreedMeta struct {
	// SourceTS is the Unix timestamp of the last value read, as the API returns it.
	SourceTS string `json:"source_ts,omitempty"`
}

// Greed uses the Crypto Fear & Greed Index (0 to 100) as the price. It changes once a day.
type Greed struct {
	Client *api.Greed
}

func (g *Greed) Ticker() string { return "GREED" }

func (g *Greed) Fetch(ctx context.Context, st State[GreedMeta]) (Result[GreedMeta], error) {
	v, ts, err := g.Client.Latest(ctx)
	if err != nil {
		return Result[GreedMeta]{}, err
	}

	outcome := OK
	if ts == st.Meta.SourceTS {
		outcome = NoNewData
	}
	return Result[GreedMeta]{Raw: v, Outcome: outcome, Patch: GreedMeta{SourceTS: ts}}, nil
}
