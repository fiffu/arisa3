package colours

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/fiffu/arisa3/app/log"
	"github.com/fiffu/arisa3/app/types"
	"github.com/fiffu/arisa3/app/utils"
)

func (c *Cog) col(ctx context.Context, req types.ICommandEvent) error {
	from := req.Interaction().Member
	if from == nil {
		return req.Respond(ctx, types.NewResponse().Content("You need to be in a guild to use this command."))
	}

	s := NewDomainSession(req.Session())
	guildID := req.Interaction().GuildID
	userID := req.User().ID
	mem, err := s.GuildMember(ctx, guildID, userID)
	if err != nil {
		// failed to get member
		log.Errorf(ctx, err, "Errored while retrieving member, guild=%s user=%s", guildID, userID)
		return err
	}

	// reroll here
	newColour, err := c.domain.Reroll(ctx, s, mem)
	if errors.Is(err, ErrRerollCooldownPending) {
		log.Errorf(ctx, err, "Blocked reroll due to cooldown pending, guild=%s user=%s", guildID, userID)

		endTime, err := c.domain.GetRerollCooldownEndTime(ctx, mem)
		if err != nil {
			return err
		}
		delta := utils.FormatDuration(time.Until(endTime))
		msg := fmt.Sprintf("You cannot reroll a new colour yet! Cooldown remaining: %s", delta)
		return req.Respond(ctx, types.NewResponse().Content(msg))
	} else if err != nil {
		return err
	}
	log.Infof(ctx, "Generated colour: #%s", newColour.ToHexcode())

	embed := newEmbed(newColour)
	return req.Respond(ctx, types.NewResponse().Embeds(embed))
}

// setFreeze will freeze or unfreeze a member's colour role.
func (c *Cog) setFreeze(ctx context.Context, req types.ICommandEvent, toFrozen bool) error {
	mem, resp, err := c.fetchMember(ctx, req)
	if err != nil {
		return err
	}
	if resp != nil {
		return req.Respond(ctx, resp)
	}

	guildID := mem.Guild().ID()
	userID := mem.UserID()

	action := c.domain.Freeze
	un := ""
	if !toFrozen {
		action = c.domain.Unfreeze
		un = "un"
	}

	role := c.domain.GetColourRole(ctx, mem)
	if role == nil {
		// user has no colour role
		log.Warnf(ctx, "User has no role to %sfreeze, guild=%s user=%s", un, guildID, userID)
		return req.Respond(ctx, types.NewResponse().Content("You don't even have a colour role..."))
	}

	if err := action(ctx, mem); err != nil {
		log.Errorf(ctx, err, "Errored while %sfreezing colour, guild=%s user=%s", un, guildID, userID)
		return err
	}

	emb := newEmbed(role.Colour()).Descriptionf("Your colour has been %sfrozen.", un)
	return req.Respond(ctx, types.NewResponse().Embeds(emb))
}

func (c *Cog) mutate(ctx context.Context, msg types.IMessageEvent) {
	guildID := msg.GuildID()
	if guildID == "" {
		// Not a message from a guild, ignore
		return
	}

	userID := msg.User().ID
	s := NewDomainSession(msg.Event().Session())
	member, err := s.GuildMember(ctx, guildID, userID)
	if err != nil {
		log.Errorf(ctx, err, "Errored while retrieving member, guild=%s user=%s", guildID, userID)
		return
	}

	newColour, err := c.domain.Mutate(ctx, s, member)

	switch {
	case newColour == nil:
		// user has no colour role, do nothing
		return

	case err == nil:
		log.Infof(ctx, "Mutated colour: #%s, guild=%s user=%s", newColour.ToHexcode(), guildID, userID)

	case errors.Is(err, ErrMutateCooldownPending):
		log.Infof(ctx, "Skipped mutate due to cooldown pending, guild=%s user=%s", guildID, userID)

	default:
		log.Errorf(ctx, err, "Errored while mutating member, guild=%s user=%s", guildID, userID)
	}
}

func (c *Cog) fetchMember(ctx context.Context, req types.ICommandEvent) (IDomainMember, types.ICommandResponse, error) {
	from := req.Interaction().Member
	if from == nil {
		resp := types.NewResponse().Content("You need to be in a guild to use this command.")
		return nil, resp, nil
	}

	s := NewDomainSession(req.Session())
	guildID := req.Interaction().GuildID
	userID := req.User().ID
	mem, err := s.GuildMember(ctx, guildID, userID)
	if err != nil {
		// failed to get member
		log.Errorf(ctx, err, "Errored while retrieving member, guild=%s user=%s", guildID, userID)
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
