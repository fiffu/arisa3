package cardboard

import (
	"context"
	"path/filepath"

	"github.com/fiffu/arisa3/app/commandfilters"
	"github.com/fiffu/arisa3/app/database"
	"github.com/fiffu/arisa3/app/engine"
	"github.com/fiffu/arisa3/app/types"
	"github.com/fiffu/arisa3/lib"

	dgo "github.com/bwmarrin/discordgo"
)

var (
	migrationsDir = filepath.Join(lib.MustGetCallerDir(), "dbmigrations")

	respRequiresGuild = types.NewResponse().Content("This command can only be used from a server.")
	respRequiresAdmin = types.NewResponse().Content("This command can only be used from a server by a server admin.")
)

// Cog implements ICog and IDefaultStartup
type Cog struct {
	commands *engine.CommandsRegistry
	db       database.IDatabase

	cfg    *Config
	domain IDomain
}

type Config struct {
	User           string `mapstructure:"user" envvar:"DANBOORU_USER"`
	APIKey         string `mapstructure:"api_key" envvar:"DANBOORU_API_KEY"`
	APITimeoutSecs int    `mapstructure:"api_timeout_secs"`
}

func NewCog(a types.IApp) types.ICog {
	return &Cog{
		commands: engine.NewCommandRegistry(),
		db:       a.Database(),
	}
}

func (c *Cog) Name() string                       { return "cardboard" }
func (c *Cog) ConfigPointer() types.StructPointer { return &Config{} }
func (c *Cog) Configure(ctx context.Context, cfg types.CogConfig) error {
	config, ok := cfg.(*Config)
	if !ok {
		return engine.UnexpectedConfigType(c.ConfigPointer(), cfg)
	}
	c.cfg = config

	c.domain = NewDomain(
		c.db,
		c.cfg,
	)
	return nil
}

func (c *Cog) OnStartup(ctx context.Context, app types.IApp, rawConfig types.CogConfig) error {
	return engine.Bootstrap(ctx, app, rawConfig, c)
}

func (c *Cog) MigrationsDir() string {
	return migrationsDir
}

func (c *Cog) ReadyCallback(s *dgo.Session, r *dgo.Ready) error {
	guildOnly := commandfilters.NewMiddleware(commandfilters.IsFromGuild).
		FailureResponse(respRequiresGuild).
		CommandDecorator()
	adminOnly := commandfilters.NewMiddleware(commandfilters.IsGuildAdmin).
		FailureResponse(respRequiresAdmin).
		CommandDecorator()

	err := c.commands.Register(
		s,
		// commands to fetch posts
		c.danCommand(),
		c.cuteCommand(),
		c.lewdCommand(),

		// commands to set tag ops
		adminOnly(c.promoteCommand()),
		adminOnly(c.demoteCommand()),
		adminOnly(c.omitCommand()),
		adminOnly(c.aliasCommand()),
		guildOnly(c.aliasesCommand()),
	)
	if err != nil {
		return err
	}
	c.commands.BindCallbacks(s)
	return nil
}
