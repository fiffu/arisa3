package colours

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/fiffu/arisa3/app/engine"
	"github.com/fiffu/arisa3/app/types"
	"github.com/fiffu/arisa3/app/utils"
	"github.com/rs/zerolog/log"
)

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
	if errors.Is(err, ErrCooldownPending) {
		engine.CommandLog(c, req, log.Info()).Err(err).
			Msgf("Blocked reroll due to cooldown pending, guild=%s user=%s", guildID, userID)

		endTime, err := c.domain.GetRerollCooldownEndTime(mem)
		if err != nil {
			return err
		}
		delta := utils.FormatDuration(time.Until(endTime))
		msg := fmt.Sprintf("You cannot reroll a new colour yet! Cooldown remaining: %s", delta)
		return req.Respond(
			types.NewResponse().Content(msg),
		)
	} else if err != nil {
		return err
	}
	engine.CommandLog(c, req, log.Info()).Msgf("Generated colour: #%s", newColour.ToHexcode())

	embed := newEmbed(newColour)
	return req.Respond(
		types.NewResponse().Embeds(embed),
	)
}

// setFreeze will freeze or unfreeze a member's colour role.
func (c *Cog) setFreeze(req types.ICommandEvent, toFrozen bool) error {
	mem, resp, err := c.fetchMember(req)
	if err != nil {
		return err
	}
	if resp != nil {
		return req.Respond(resp)
	}

	guildID := mem.Guild().ID()
	userID := mem.UserID()
	un := ""
	if !toFrozen {
		un = "un"
	}

	role := c.domain.GetColourRole(mem)
	if role == nil {
		// user has no colour role
		engine.CommandLog(c, req, log.Error()).Err(err).
			Msgf("User has no role to %sfreeze, guild=%s user=%s", un, guildID, userID)
		return req.Respond(types.NewResponse().Content("You don't even have a colour role..."))
	}

	if err := c.domain.Freeze(mem); err != nil {
		engine.CommandLog(c, req, log.Error()).Err(err).
			Msgf("Errored while freezing colour, guild=%s user=%s", guildID, userID)
		return err
	}

	emb := newEmbed(role.Colour()).Descriptionf("Your colour has been %sfrozen.", un)
	return req.Respond(types.NewResponse().Embeds(emb))
}

func (c *Cog) mutate(msg types.IMessageEvent) {
	guildID := msg.GuildID()
	if guildID == "" {
		// Not a message from a guild, ignore
		return
	}

	userID := msg.User().ID
	s := NewDomainSession(msg.Event().Session())
	member, err := s.GuildMember(guildID, userID)
	if err != nil {
		engine.EventLog(c, msg.Event(), log.Error()).Err(err).
			Msgf("Errored while retrieving member, guild=%s user=%s", guildID, userID)
		return
	}

	newColour, err := c.domain.Mutate(s, member)

	switch {
	case newColour == nil:
		// user has no colour role, do nothing
		return

	case err == nil:
		engine.EventLog(c, msg.Event(), log.Info()).
			Msgf("Mutated colour: #%s, guild=%s user=%s", newColour.ToHexcode(), guildID, userID)

	case errors.Is(err, ErrCooldownPending):
		engine.EventLog(c, msg.Event(), log.Info()).
			Msgf("Skipped mutate due to cooldown pending, guild=%s user=%s", guildID, userID)

	default:
		engine.EventLog(c, msg.Event(), log.Error()).Err(err).
			Msgf("Errored while mutating member, guild=%s user=%s", guildID, userID)
	}
}

func (c *Cog) colInfo(req types.ICommandEvent) error {
	mem, resp, err := c.fetchMember(req)
	if err != nil {
		return err
	}
	if resp != nil {
		return req.Respond(resp)
	}

	guildID := mem.Guild().ID()
	userID := mem.UserID()

	role := c.domain.GetColourRole(mem)
	if role == nil {
		engine.CommandLog(c, req, log.Error()).Err(err).
			Msgf("No colour role found, guild=%s user=%s", guildID, userID)
		return req.Respond(types.NewResponse().
			Content("You don't have a colour role. Use /col to get a random colour!"))
	}

	now := time.Now()
	desc := make([]string, 0)

	desc = append(desc, "**Reroll cooldown:**")
	rerollCDEnds, err := c.domain.GetRerollCooldownEndTime(mem)
	switch {
	case err != nil:
		engine.CommandLog(c, req, log.Error()).Err(err).
			Msgf("Errored getting cooldown end time, guild=%s user=%s", guildID, userID)
		return err
	case now.Before(rerollCDEnds):
		desc = append(desc, utils.FormatDuration(rerollCDEnds.Sub(now)))
	default:
		desc = append(desc, "_(No cooldown, reroll available)_")
	}

	desc = append(desc, "**Last mutate:**")
	lastMutateTime, ok, err := c.domain.GetLastMutate(mem)
	switch {
	case err != nil:
		engine.CommandLog(c, req, log.Error()).Err(err).
			Msgf("Errored getting last mutate time, guild=%s user=%s", guildID, userID)
		return err
	case !ok:
		desc = append(desc, "_(Never)_")
	default:
		desc = append(desc, utils.FormatDuration(now.Sub(lastMutateTime)))
	}

	lastFrozenTime, err := c.domain.GetLastFrozen(mem)
	switch {
	case err != nil:
		engine.CommandLog(c, req, log.Error()).Err(err).
			Msgf("Errored getting last frozen time, guild=%s user=%s", guildID, userID)
		return err
	case lastFrozenTime != Never:
		desc = append(desc, utils.FormatDuration(now.Sub(lastFrozenTime)))
	default:
		// do nothing
	}

	embed := newEmbed(role.Colour()).Description(strings.Join(desc, "\n"))
	return req.Respond(types.NewResponse().Embeds(embed))
}

func (c *Cog) fetchMember(req types.ICommandEvent) (IDomainMember, types.ICommandResponse, error) {
	from := req.Interaction().Member
	if from == nil {
		resp := types.NewResponse().Content("You need to be in a guild to use this command.")
		return nil, resp, nil
	}

	s := NewDomainSession(req.Session())
	guildID := req.Interaction().GuildID
	userID := req.User().ID
	mem, err := s.GuildMember(guildID, userID)
	if err != nil {
		// failed to get member
		engine.CommandLog(c, req, log.Error()).Err(err).
			Msgf("Errored while retrieving member, guild=%s user=%s", guildID, userID)
		return nil, nil, err
	}

	return mem, nil, nil
}

// newEmbed creates an embed object with title and colour defined.
func newEmbed(colour *Colour) types.IEmbed {
	r, g, b := colour.scale255()
	hex := colour.ToHexcode()
	return types.NewEmbed().
		Titlef("#%s · rgb(%d, %d, %d)", hex, r, g, b).
		Colour(colour.ToDecimal())
}
