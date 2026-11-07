package tickers

import "github.com/fiffu/arisa3/app/cogs/stonks/tickers/api"

func NewDJIA(client *api.Fred) *Fred {
	return &Fred{Client: client, Name: "DJIA", SeriesID: "DJIA"}
}
