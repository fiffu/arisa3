package colours

import (
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
	mem, err := s.GuildMember(
		req.Interaction().GuildID,
		req.User().ID,
	)
	if err != nil {
		// failed to get member
		return err
	}

	// reroll here
	newColour, err := c.domain.Reroll(s, mem)
	engine.CommandLog(c, req, log.Info()).Msgf("Generated colour: %+v", newColour)
	if err != nil {
		return err
	}

	r, g, b := newColour.scale255()
	hex := newColour.ToHexcode()
	title := fmt.Sprintf("#%s · rgb(%d, %d, %d)", hex, r, g, b)
	embed := types.NewEmbed().Title(title).Color(newColour.ToDecimal())

	return req.Respond(
		types.NewResponse().Embeds(embed),
	)
}
