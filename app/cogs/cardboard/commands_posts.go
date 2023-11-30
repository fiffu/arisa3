package cardboard

import (
	"context"
	"errors"

	"github.com/fiffu/arisa3/app/cogs/cardboard/api"
	"github.com/fiffu/arisa3/app/types"
)

const (
	OptionQuery = "query"
	OptionTag   = "tag"
)

func (c *Cog) danCommand() *types.Command {
	return types.NewCommand("dan").ForChat().
		Desc("Search the booru with the exact query, no aliases, no result filter.").
		Options(
			types.NewOption(OptionQuery).
				Desc("exact search query").
				String().Required(),
		).
		Handler(c.dumbSearch)
}

func (c *Cog) cuteCommand() *types.Command {
	return types.NewCommand("cute").ForChat().
		Desc("Finds a cute picture with a particular tag.").
		Options(
			types.NewOption(OptionTag).
				Desc("the tag to search (spaces convert to _)").
				String().Required(),
		).
		Handler(c.smartSearch(true))
}

func (c *Cog) lewdCommand() *types.Command {
	return types.NewCommand("lewd").ForChat().
		Desc("Finds a LEWD picture with a particular tag.").
		Options(
			types.NewOption(OptionTag).
				Desc("tag to search (spaces convert to _)").
				String().Required(),
		).
		Handler(c.smartSearch(false))
}

func (c *Cog) dumbSearch(ctx context.Context, req types.ICommandEvent) error {
	queryStr, _ := req.Args().String(OptionQuery)
	query := NewQuery(queryStr).
		WithNoMagic().
		WithGuildID(getGuildID(req))

	posts, err := c.domain.PostsSearch(ctx, query)
	if errors.Is(err, api.ErrUnderMaintenance) {
		return req.Respond(ctx, c.domain.MaintenanceResult())
	} else if err != nil {
		return err
	}

	resp, err := c.buildResponse(ctx, query, posts)
	if err != nil {
		return err
	}

	return req.Respond(ctx, resp)
}

func (c *Cog) smartSearch(safe bool) types.CommandHandler {
	return func(ctx context.Context, req types.ICommandEvent) error {
		queryStr, _ := req.Args().String(OptionTag)

		query := NewQuery(queryStr).
			WithMagic().
			WithGuildID(getGuildID(req))
		if safe {
			query.WithSafe()
		} else {
			query.WithUnsafe()
		}

		posts, err := c.domain.PostsSearch(ctx, query)
		if errors.Is(err, api.ErrUnderMaintenance) {
			return req.Respond(ctx, c.domain.MaintenanceResult())
		} else if err != nil {
			return err
		}

		resp, err := c.buildResponse(ctx, query, posts)
		if err != nil {
			return err
		}

		return req.Respond(ctx, resp)
	}
}

func (c *Cog) buildResponse(ctx context.Context, query IQueryPosts, posts []*api.Post) (types.ICommandResponse, error) {
	emb, err := c.domain.PostsResult(ctx, query, posts)
	if err != nil {
		return nil, err
	}

	resp := types.NewResponse().Embeds(emb)
	return resp, nil
}
