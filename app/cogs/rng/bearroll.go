package rng

import (
	"github.com/fiffu/arisa3/app/types"
	"github.com/fiffu/arisa3/app/utils"
	"github.com/fiffu/arisa3/lib"
)

var (
	bears = []string{
		":bear:",
		":polar_bear:",
		":teddy_bear:",
		utils.POOH,
	}
)

func (c *Cog) bearRollCommand() *types.Command {
	return types.NewCommand("bearroll").ForChat().
		Desc("Do a bear-er roll.").
		Handler(c.bearRoll)
}

func (c *Cog) bearRoll(req types.ICommandEvent) error {
	bear := lib.ChooseString(bears)
	resp := types.NewResponse().Content(bear)
	return req.Respond(resp)
}
