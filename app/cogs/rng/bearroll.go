package rng

import (
	"context"

	"github.com/fiffu/arisa3/app/types"
	"github.com/fiffu/arisa3/app/utils"
	"github.com/fiffu/arisa3/lib/functional"
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

func (c *Cog) bearRoll(ctx context.Context, req types.ICommandEvent) error {
	bear := functional.SliceOf(bears).TakeRandom()
	resp := types.NewResponse().Content(bear)
	return req.Respond(ctx, resp)
}
