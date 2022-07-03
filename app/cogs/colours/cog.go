package colours

import (
	"context"
	"path/filepath"

	"github.com/fiffu/arisa3/app/database"
	"github.com/fiffu/arisa3/app/engine"
	"github.com/fiffu/arisa3/app/types"
	"github.com/fiffu/arisa3/lib"
	"github.com/rs/zerolog/log"

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
	engine.CogLog(c, log.Info()).Msgf("IColoursDomain loaded")
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
	)
	if err != nil {
		return err
	}
	c.commands.BindCallbacks(s)
	return nil
}

func (c *Cog) registerEvents(sess *dgo.Session) {
	sess.AddHandler(func(s *dgo.Session, m *dgo.MessageCreate) {
		evt := types.NewMessageEvent(s, m)
		c.mutate(evt)
	})
}
