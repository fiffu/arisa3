package colours

import (
	"errors"
	"fmt"

	"github.com/fiffu/arisa3/app/engine"
	"github.com/fiffu/arisa3/app/types"
	"github.com/rs/zerolog/log"
)

func (c *Cog) colCommand() *types.Command {
	return types.NewCommand("col").ForChat().
		Desc("Gives you a shiny new colour").
		Handler(c.col)
}

func (c *Cog) col(req types.ICommandEvent) error {
	from := req.Interaction().Member
	if from == nil {
		return req.Respond(
			types.NewResponse().Content("You need to be in a guild to use this command."),
		)
	}

	s := NewDomainSession(req.Session())
	guildID := req.Interaction().GuildID
	userID := req.User().ID
	mem, err := s.GuildMember(guildID, userID)
	if err != nil {
		// failed to get member
		engine.CommandLog(c, req, log.Error()).Err(err).
			Msgf("Errored while retrieving member, guild=%s user=%s", guildID, userID)
		return err
	}

	// reroll here
	newColour, err := c.domain.Reroll(s, mem)
	engine.CommandLog(c, req, log.Info()).Msgf("Generated colour: #%s", newColour.ToHexcode())
	if err != nil {
		return err
	}

	r, g, b := newColour.scale255()
	hex := newColour.ToHexcode()
	title := fmt.Sprintf("#%s · rgb(%d, %d, %d)", hex, r, g, b)
	embed := types.NewEmbed().Title(title).Colour(newColour.ToDecimal())

	return req.Respond(
		types.NewResponse().Embeds(embed),
	)
}

func (c *Cog) mutate(msg types.IMessageEvent) {
	guildID := msg.GuildID()
	if guildID == "" {
		// Not a message from a guild, ignore
		return
	}

	s := NewDomainSession(msg.Event().Session())
	userID := msg.User().ID
	member, err := s.GuildMember(guildID, userID)
	if err != nil {
		engine.EventLog(c, msg.Event(), log.Error()).Err(err).
			Msgf("Errored while retrieving member, guild=%s user=%s", guildID, userID)
	}

	newColour, err := c.domain.Mutate(s, member)

	switch {
	case err == nil:
		engine.EventLog(c, msg.Event(), log.Info()).
			Msgf("Mutated colour: #%s, guild=%s user=%s", newColour.ToHexcode(), guildID, userID)

	case errors.As(err, &ErrCooldownPending):
		engine.EventLog(c, msg.Event(), log.Info()).
			Msgf("Skipped mutate due to cooldown pending, guild=%s user=%s", guildID, userID)

	default:
		engine.EventLog(c, msg.Event(), log.Error()).Err(err).
			Msgf("Errored while mutating member, guild=%s user=%s", guildID, userID)
	}
}
