package tickers

import (
	"context"

	"github.com/fiffu/arisa3/app/cogs/stonks/tickers/api"
)

type FredMeta struct {
	// SourceDate is the date of the last value read, such as "2026-11-05".
	SourceDate string `json:"source_date,omitempty"`
}

// Fred prices a ticker at the latest daily close of a FRED series, such as SP500 or DJIA.
type Fred struct {
	Client   *api.Fred
	Name     string
	SeriesID string
}

func (f *Fred) Ticker() string { return f.Name }

func (f *Fred) Fetch(ctx context.Context, st State[FredMeta]) (Result[FredMeta], error) {
	// Two weeks covers any market closure, such as a holiday weekend.
	date, value, err := f.Client.LatestClose(ctx, f.SeriesID, st.Slot.AddDate(0, 0, -14))
	if err != nil {
		return Result[FredMeta]{}, err
	}

	outcome := OK
	if date == st.Meta.SourceDate {
		outcome = NoNewData
	}
	return Result[FredMeta]{Raw: value, Outcome: outcome, Patch: FredMeta{SourceDate: date}}, nil
}
