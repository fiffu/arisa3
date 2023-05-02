package cardboard

import (
	"context"
	"fmt"
	"strings"

	"github.com/fiffu/arisa3/app/types"
)

const (
	OptionAlias = "alias"
)

func (c *Cog) promoteCommand() *types.Command {
	return types.NewCommand("promote").ForChat().
		Desc("(Admins only) Indicate that posts with this tag should be prioritized over other posts.").
		Options(
			types.NewOption(OptionTag).
				Desc("tag to be promoted").
				String().Required(),
		).
		Handler(c.promote)
}

func (c *Cog) demoteCommand() *types.Command {
	return types.NewCommand("demote").ForChat().
		Desc("(Admins only) Indicate that posts with this tag should be de-prioritized in favour of other posts.").
		Options(
			types.NewOption(OptionTag).
				Desc("tag to be demoted").
				String().Required(),
		).
		Handler(c.demote)
}

func (c *Cog) omitCommand() *types.Command {
	return types.NewCommand("omit").ForChat().
		Desc("(Admins only) Indicate that posts with this tag should not be shown.").
		Options(
			types.NewOption(OptionTag).
				Desc("tag to be omitted").
				String().Required(),
		).
		Handler(c.omit)
}

func (c *Cog) aliasCommand() *types.Command {
	return types.NewCommand("alias").ForChat().
		Desc("(Admins only) Set an alias mapping to an actual tag.").
		Options(
			types.NewOption(OptionAlias).
				Desc("alias name to be created").
				String().Required(),
			types.NewOption(OptionTag).
				Desc("actual tag").
				String().Required(),
		).
		Handler(c.alias)
}

func (c *Cog) aliasesCommand() *types.Command {
	return types.NewCommand("aliases").ForChat().
		Desc("List available aliases.").
		Handler(c.listAliases)
}

func (c *Cog) promote(ctx context.Context, req types.ICommandEvent) error {
	tagName, _ := req.Args().String(OptionTag)

	guildID := getGuildID(req)
	if guildID == "" {
		return req.Respond(ctx, respRequiresAdmin)
	}

	if err := c.domain.SetPromote(tagName, guildID); err != nil {
		return err
	}
	resp := types.NewResponse().Content(fmt.Sprintf("Marked `%s` to be promoted.", tagName))
	return req.Respond(ctx, resp)
}

func (c *Cog) demote(ctx context.Context, req types.ICommandEvent) error {
	tagName, _ := req.Args().String(OptionTag)

	guildID := getGuildID(req)
	if guildID == "" {
		return req.Respond(ctx, respRequiresAdmin)
	}

	if err := c.domain.SetDemote(tagName, guildID); err != nil {
		return err
	}
	resp := types.NewResponse().Content(fmt.Sprintf("Marked `%s` to be demoted.", tagName))
	return req.Respond(ctx, resp)
}

func (c *Cog) omit(ctx context.Context, req types.ICommandEvent) error {
	tagName, _ := req.Args().String(OptionTag)

	guildID := getGuildID(req)
	if guildID == "" {
		return req.Respond(ctx, respRequiresAdmin)
	}

	if err := c.domain.SetOmit(tagName, guildID); err != nil {
		return err
	}
	resp := types.NewResponse().Content(fmt.Sprintf("Marked `%s` to be omitted.", tagName))
	return req.Respond(ctx, resp)
}

func (c *Cog) alias(ctx context.Context, req types.ICommandEvent) error {
	actual, _ := req.Args().String(OptionTag)
	alias, _ := req.Args().String(OptionAlias)

	guildID := getGuildID(req)
	if guildID == "" {
		return req.Respond(ctx, respRequiresAdmin)
	}

	if err := c.domain.SetAlias(guildID, Alias(alias), Actual(actual)); err != nil {
		return err
	}
	resp := types.NewResponse().Content(fmt.Sprintf("`%s` will be aliased as `%s`.", actual, alias))
	return req.Respond(ctx, resp)
}

func (c *Cog) listAliases(ctx context.Context, req types.ICommandEvent) error {
	guildID := getGuildID(req)
	if guildID == "" {
		return req.Respond(ctx, respRequiresAdmin)
	}

	aliasMap, err := c.domain.GetAliases(guildID)
	if err != nil {
		return err
	}

	list := []string{
		"Alias → Actual Tag",
		"==================",
	}
	for ali, act := range aliasMap {
		list = append(list, fmt.Sprintf("%s → %s", string(ali), string(act)))
	}

	resp := respNoAliases
	if len(list) > 0 {
		message := fmt.Sprintf("```\n" + strings.Join(list, "\n") + "```")
		resp.Content(message)
	}
	return req.Respond(ctx, resp)
}
