package rng

import (
	"context"

	"github.com/fiffu/arisa3/app/engine"
	"github.com/fiffu/arisa3/app/types"

	dgo "github.com/bwmarrin/discordgo"
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

func (c *Cog) Name() string                                             { return "rng" }
func (c *Cog) ConfigPointer() types.StructPointer                       { return nil }
func (c *Cog) Configure(ctx context.Context, cfg types.CogConfig) error { return nil }

func (c *Cog) OnStartup(ctx context.Context, app types.IApp, rawConfig types.CogConfig) error {
	return engine.Bootstrap(ctx, app, rawConfig, c)
}

func (c *Cog) ReadyCallback(s *dgo.Session, r *dgo.Ready) error {
	err := c.commands.Register(
		s,
		c.rollCommand(),
		c.eightBallCommand(),
		c.pokiesCommand(),
	)
	if err != nil {
		return err
	}
	c.commands.BindCallbacks(s)
	return nil
}
