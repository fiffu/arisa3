package cardboard

import (
	"fmt"

	"github.com/fiffu/arisa3/app/types"
)

const (
	OptionAlias = "alias"
)

func (c *Cog) promoteCommand() *types.Command {
	return types.NewCommand("promote").ForChat().
		Desc("Indicate that posts with this tag should be prioritized over other posts.").
		Options(
			types.NewOption(OptionTag).String().Required(),
		).
		Handler(c.promote)
}

func (c *Cog) demoteCommand() *types.Command {
	return types.NewCommand("demote").ForChat().
		Desc("Indicate that posts with this tag should be de-prioritized in favour of other posts.").
		Options(
			types.NewOption(OptionTag).String().Required(),
		).
		Handler(c.demote)
}

func (c *Cog) omitCommand() *types.Command {
	return types.NewCommand("omit").ForChat().
		Desc("Indicate that posts with this tag should not be shown.").
		Options(
			types.NewOption(OptionTag).String().Required(),
		).
		Handler(c.omit)
}

func (c *Cog) aliasCommand() *types.Command {
	return types.NewCommand("alias").ForChat().
		Desc("Set an alias mapping to an actual tag.").
		Options(
			types.NewOption(OptionAlias).String().Required(),
			types.NewOption(OptionTag).String().Required(),
		).
		Handler(c.alias)
}

func (c *Cog) promote(req types.ICommandEvent) error {
	tagName, _ := req.Args().String(OptionTag)

	guildID := getGuildID(req)
	if guildID == "" {
		return req.Respond(respRequiresAdmin)
	}

	if err := c.domain.PromoteTag(tagName, guildID); err != nil {
		return err
	}
	resp := types.NewResponse().Content(fmt.Sprintf("Marked `%s` to be promoted.", tagName))
	return req.Respond(resp)
}

func (c *Cog) demote(req types.ICommandEvent) error {
	tagName, _ := req.Args().String(OptionTag)

	guildID := getGuildID(req)
	if guildID == "" {
		return req.Respond(respRequiresAdmin)
	}

	if err := c.domain.DemoteTag(tagName, guildID); err != nil {
		return err
	}
	resp := types.NewResponse().Content(fmt.Sprintf("Marked `%s` to be demoted.", tagName))
	return req.Respond(resp)
}

func (c *Cog) omit(req types.ICommandEvent) error {
	tagName, _ := req.Args().String(OptionTag)

	guildID := getGuildID(req)
	if guildID == "" {
		return req.Respond(respRequiresAdmin)
	}

	if err := c.domain.OmitTag(tagName, guildID); err != nil {
		return err
	}
	resp := types.NewResponse().Content(fmt.Sprintf("Marked `%s` to be omitted.", tagName))
	return req.Respond(resp)
}

func (c *Cog) alias(req types.ICommandEvent) error {
	actual, _ := req.Args().String(OptionTag)
	alias, _ := req.Args().String(OptionAlias)

	guildID := getGuildID(req)
	if guildID == "" {
		return req.Respond(respRequiresAdmin)
	}

	if err := c.domain.AliasTag(guildID, actual, alias); err != nil {
		return err
	}
	resp := types.NewResponse().Content(fmt.Sprintf("`%s` will be aliased as `%s`.", actual, alias))
	return req.Respond(resp)
}
