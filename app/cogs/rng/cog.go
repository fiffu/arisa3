package rng

import (
	"context"
	"time"

	"github.com/fiffu/arisa3/app/engine"
	"github.com/fiffu/arisa3/app/types"
	"github.com/fiffu/arisa3/lib"

	dgo "github.com/bwmarrin/discordgo"
)

// RollObserver gets every /roll. The observer decides which rolls it uses.
type RollObserver interface {
	OnRoll(ctx context.Context, roll Roll)
}

type Roll struct {
	Result   int
	Count    int
	Sides    int
	Modifier int
	// Parsed is false when the user gave no dice expression. The roll then
	// uses the default d100-1, whose modifier the user did not choose.
	Parsed bool
}

// Bounds returns the lowest and highest possible results.
func (r Roll) Bounds() (lo, hi int) {
	if r.Sides <= 0 {
		return r.Modifier, r.Modifier
	}
	return r.Count + r.Modifier, r.Count*r.Sides + r.Modifier
}

// Cog implements ICog and IDefaultStartup
type Cog struct {
	commands      *engine.CommandsRegistry
	pokiesCache   lib.ICache[*cachedEmojis, string]
	rollObservers []RollObserver
}

func NewCog(a types.IApp, observers ...RollObserver) types.ICog {
	return &Cog{
		commands:      engine.NewCommandRegistry(),
		pokiesCache:   lib.NewCache[*cachedEmojis, string](1 * time.Hour),
		rollObservers: observers,
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
