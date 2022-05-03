package cardboard

import (
	"github.com/fiffu/arisa3/app/cogs/cardboard/api"
	"github.com/fiffu/arisa3/app/types"
)

const (
	Query = "query"
)

func (c *Cog) danCommand() *types.Command {
	return types.NewCommand("dan").ForChat().
		Desc("Search the booru with the exact query, no aliases, no result filter.").
		Options(
			types.NewOption(Query).String().Required(),
		).
		Handler(c.dumbSearch)
}

func (c *Cog) cuteCommand() *types.Command {
	return types.NewCommand("cute").ForChat().
		Desc("Finds a cute picture with the given tag (spaces convert to _).").
		Options(
			types.NewOption(Query).String().Required(),
		).
		Handler(c.smartSearch(true))
}

func (c *Cog) lewdCommand() *types.Command {
	return types.NewCommand("lewd").ForChat().
		Desc("Finds a LEWD picture with the given tag (spaces convert to _).").
		Options(
			types.NewOption(Query).String().Required(),
		).
		Handler(c.smartSearch(false))
}

func (c *Cog) dumbSearch(req types.ICommandEvent) error {
	queryStr, _ := req.Args().String(Query)
	query := NewQuery(queryStr).NoMagic()

	posts, err := c.domain.PostsSearch(query)
	if err != nil {
		return err
	}

	resp, err := c.buildResponse(query, posts)
	if err != nil {
		return err
	}

	return req.Respond(resp)
}

func (c *Cog) smartSearch(safe bool) types.Handler {
	return func(req types.ICommandEvent) error {
		queryStr, _ := req.Args().String(Query)
		query := NewQuery(queryStr)
		if safe {
			query.SetSafe()
		} else {
			query.SetUnsafe()
		}

		posts, err := c.domain.PostsSearch(query)
		if err != nil {
			return err
		}

		resp, err := c.buildResponse(query, posts)
		if err != nil {
			return err
		}

		return req.Respond(resp)
	}
}

func (c *Cog) buildResponse(query IQueryPosts, posts []*api.Post) (types.ICommandResponse, error) {
	emb, err := c.domain.PostsResult(query, posts)
	if err != nil {
		return nil, err
	}

	resp := types.NewResponse().Embeds(emb)
	return resp, nil
}
