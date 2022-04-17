package rng

import (
	"arisa3/app/engine"
	"arisa3/app/types"
	"context"

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

func (c *Cog) Name() string                                         { return "rng" }
func (c *Cog) ConfigPointer() types.StructPointer                   { return nil }
func (c *Cog) Configure(ctx context.Context, cfg interface{}) error { return nil }

func (c *Cog) OnStartup(ctx context.Context, sess *dgo.Session, rawConfig types.CogConfig) error {
	return engine.Bootstrap(ctx, sess, rawConfig, c)
}

func (c *Cog) ReadyCallback(s *dgo.Session, r *dgo.Ready) error {
	err := c.commands.Register(
		s,
		c.rollCommand(),
	)
	if err != nil {
		return err
	}
	c.commands.BindCallbacks(s)
	return nil
}
