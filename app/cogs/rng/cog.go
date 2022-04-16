package rng

import (
	"arisa3/app/engine"
	"arisa3/app/types"
	"context"
	"fmt"
	"math/rand"

	dgo "github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

// Cog implements ICog and IDefaultStartup
type Cog struct {
	commands *engine.CommandsRegistry
}

func NewCog(a types.IApp) types.ICog {
	return &Cog{
		commands: engine.NewCommandRegistry(),
	}
}

func (c *Cog) Name() string                                         { return "rng" }
func (c *Cog) ConfigPointer() types.StructPointer                   { return nil }
func (c *Cog) Configure(ctx context.Context, cfg interface{}) error { return nil }

func (c *Cog) OnStartup(ctx context.Context, sess *dgo.Session, rawConfig types.CogConfig) error {
	return engine.Bootstrap(ctx, sess, rawConfig, c)
}

func (c *Cog) ReadyCallback(s *dgo.Session, r *dgo.Ready) error {
	err := c.commands.Register(
		s,
		types.NewCommand("roll").ForChat().
			Desc("Rolls dice (supports algebraic notation, such as !roll 3d5+10)").
			Options(types.NewOption("expression").String()).
			Handler(c.roll),
	)
	if err != nil {
		return err
	}
	c.commands.BindCallbacks(s)
	return nil
}

func (c *Cog) roll(s *dgo.Session, i *dgo.InteractionCreate, cmd types.ICommand, args types.IArgs) error {
	resp := types.NewResponse()
	if expr, ok := args.String("expression"); ok {
		log.Info().Msgf("expression='%s'", expr)
		lo, hi, num := roll(expr)
		resp.Content(fmt.Sprintf("Rolling %d-%d: %d", lo, hi, num))
	} else {
		log.Info().Msgf("missing expression")
		resp.Content("missing expression")
	}
	return s.InteractionRespond(i.Interaction, resp.Data())
}

func roll(expr string) (int, int, int) {
	lo, hi := 0, 99
	choice := rand.Intn(lo+hi) - lo
	return lo, hi, choice
}
