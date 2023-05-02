package colours

import (
	"context"
	"path/filepath"

	"github.com/fiffu/arisa3/app/database"
	"github.com/fiffu/arisa3/app/engine"
	"github.com/fiffu/arisa3/app/log"
	"github.com/fiffu/arisa3/app/types"
	"github.com/fiffu/arisa3/lib"

	dgo "github.com/bwmarrin/discordgo"
)

var (
	migrationsDir = filepath.Join(lib.MustGetCallerDir(), "dbmigrations")
)

// Cog implements ICog and IDefaultStartup
type Cog struct {
	commands *engine.CommandsRegistry
	db       database.IDatabase

	cfg *Config

	domain IColoursDomain
}

type Config struct {
	MaxRoleHeightName string `mapstructure:"max_role_height_name"`

	MutateCooldownMins int `mapstructure:"mutate_cooldown_mins"`
	RerollCooldownMins int `mapstructure:"reroll_cooldown_mins"`
	RerollPenaltyMins  int `mapstructure:"reroll_penalty_mins"`
}

func NewCog(a types.IApp) types.ICog {
	return &Cog{
		commands: engine.NewCommandRegistry(),
		db:       a.Database(),
	}
}

func (c *Cog) Name() string                       { return "colours" }
func (c *Cog) ConfigPointer() types.StructPointer { return &Config{} }
func (c *Cog) Configure(ctx context.Context, cfg types.CogConfig) error {
	config, ok := cfg.(*Config)
	if !ok {
		return engine.UnexpectedConfigType(c.ConfigPointer(), cfg)
	}
	c.cfg = config
	c.domain = NewColoursDomain(
		c,
		NewRepository(c.db),
		c.cfg,
	)
	log.Infof(ctx, "IColoursDomain loaded")
	return nil
}

func (c *Cog) OnStartup(ctx context.Context, app types.IApp, rawConfig types.CogConfig) error {
	return engine.Bootstrap(ctx, app, rawConfig, c)
}

func (c *Cog) MigrationsDir() string {
	return migrationsDir
}

func (c *Cog) ReadyCallback(s *dgo.Session, r *dgo.Ready) error {
	if err := c.registerCommands(s); err != nil {
		return err
	}
	c.registerEvents(s)
	return nil
}

func (c *Cog) registerCommands(s *dgo.Session) error {
	err := c.commands.Register(
		s,
		c.colCommand(),
		c.freezeCommand(),
		c.unfreezeCommand(),
		c.colInfoCommand(),
	)
	if err != nil {
		return err
	}
	c.commands.BindCallbacks(s)
	return nil
}

func (c *Cog) registerEvents(sess *dgo.Session) {
	sess.AddHandler(engine.NewEventHandler(func(ctx context.Context, s *dgo.Session, m *dgo.MessageCreate) {
		evt := types.NewMessageEvent(s, m)
		c.onMessage(ctx, evt)
	}))
}

func (c *Cog) onMessage(ctx context.Context, evt types.IMessageEvent) {
	if evt.IsFromSelf() {
		// Ignore bot's own messages
		return
	}
	c.mutate(ctx, evt)
}

func (c *Cog) colCommand() *types.Command {
	return types.NewCommand("col").ForChat().
		Desc("Gives you a shiny new colour").
		Handler(c.col)
}

func (c *Cog) freezeCommand() *types.Command {
	return types.NewCommand("freeze").ForChat().
		Desc("Stops your colour from mutating").
		Handler(func(ctx context.Context, req types.ICommandEvent) error {
			return c.setFreeze(ctx, req, true)
		})
}

func (c *Cog) unfreezeCommand() *types.Command {
	return types.NewCommand("unfreeze").ForChat().
		Desc("Makes your colour start mutating").
		Handler(func(ctx context.Context, req types.ICommandEvent) error {
			return c.setFreeze(ctx, req, false)
		})
}
