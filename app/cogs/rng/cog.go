package rng

import (
	"context"
	"time"

	"github.com/fiffu/arisa3/app/engine"
	"github.com/fiffu/arisa3/app/types"
	"github.com/fiffu/arisa3/lib"

	dgo "github.com/bwmarrin/discordgo"
)

// Cog implements ICog and IDefaultStartup
type Cog struct {
	commands    *engine.CommandsRegistry
	pokiesCache lib.ICache[*cachedEmojis, string]
}

func NewCog(a types.IApp) types.ICog {
	return &Cog{
		commands:    engine.NewCommandRegistry(),
		pokiesCache: lib.NewCache[*cachedEmojis, string](1 * time.Hour),
	}
}

func (c *Cog) Name() string                                             { return "rng" }
func (c *Cog) ConfigPointer() types.StructPointer                       { return nil }
func (c *Cog) Configure(ctx context.Context, cfg types.CogConfig) error { return nil }

func (c *Cog) OnStartup(ctx context.Context, app types.IApp, rawConfig types.CogConfig) error {
	return engine.Bootstrap(ctx, app, rawConfig, c)
}

func (c *Cog) ReadyCallback(ctx context.Context, s *dgo.Session, r *dgo.Ready) error {
	err := c.commands.Register(
		ctx,
		s,
		c.rollCommand(),
		c.bearRollCommand(),
		c.eightBallCommand(),
		c.pokiesCommand(),
	)
	if err != nil {
		return err
	}
	c.commands.BindCallbacks(s)
	return nil
}
