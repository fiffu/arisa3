package general

import (
	"context"

	dgo "github.com/bwmarrin/discordgo"
	"github.com/fiffu/arisa3/app/engine"
	"github.com/fiffu/arisa3/app/types"
)

// Cog implements ICog and IDefaultStartup
type Cog struct {
	cfg      *Config
	commands *engine.CommandsRegistry
}
type Config struct {
	MOTD            string `mapstructure:"motd" envvar:"motd"`
	RepoName        string `mapstructure:"repo_name"`
	RepoWebURL      string `mapstructure:"repo_web_url"`
	RepoIssuesURL   string `mapstructure:"repo_issues_url"`
	RepoGitCloneURL string `mapstructure:"repo_gitclone_url"`
}

func NewCog(a types.IApp) types.ICog {
	return &Cog{
		commands: engine.NewCommandRegistry(),
	}
}

func (c *Cog) Name() string                       { return "general" }
func (c *Cog) ConfigPointer() types.StructPointer { return &Config{} }
func (c *Cog) Configure(ctx context.Context, cfg types.CogConfig) error {
	if config, ok := cfg.(*Config); ok {
		c.cfg = config
		return nil
	}
	return engine.UnexpectedConfigType(c.ConfigPointer(), cfg)
}

func (c *Cog) OnStartup(ctx context.Context, app types.IApp, rawConfig types.CogConfig) error {
	return engine.Bootstrap(ctx, app, rawConfig, c)
}

func (c *Cog) ReadyCallback(ctx context.Context, s *dgo.Session, r *dgo.Ready) error {
	c.welcome(s, r)
	err := c.commands.Register(
		ctx,
		s,
		c.gitCommand(),
	)
	if err != nil {
		return err
	}
	c.commands.BindCallbacks(s)
	return nil
}
