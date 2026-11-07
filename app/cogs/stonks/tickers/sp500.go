package tickers

import "github.com/fiffu/arisa3/app/cogs/stonks/tickers/api"

func NewSP500(client *api.Fred) *Fred {
	return &Fred{Client: client, Name: "SP500", SeriesID: "SP500"}
}
