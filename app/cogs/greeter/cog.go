package greeter

import (
	"context"
	"fmt"

	"github.com/fiffu/arisa3/app/engine"
	"github.com/fiffu/arisa3/app/types"

	dgo "github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

// Cog implements ICog and IDefaultStartup
type Cog struct {
	cfg *Config
}
type Config struct {
	MOTD string `env:"motd"`
}

func NewCog(a types.IApp) types.ICog { return &Cog{} }

func (c *Cog) Name() string                       { return "greeter" }
func (c *Cog) ConfigPointer() types.StructPointer { return &Config{} }
func (c *Cog) Configure(ctx context.Context, cfg types.CogConfig) error {
	if config, ok := cfg.(*Config); ok {
		c.cfg = config
		return nil
	}
	return engine.UnexpectedConfigType(c.ConfigPointer(), cfg)
}
func (c *Cog) RunMigrations() {

}
func (c *Cog) OnStartup(ctx context.Context, app types.IApp, rawConfig types.CogConfig) error {
	return engine.Bootstrap(ctx, app, rawConfig, c)
}

func (c *Cog) ReadyCallback(s *dgo.Session, r *dgo.Ready) error {
	c.welcome(s, r)
	return nil
}

func (c *Cog) welcome(s *dgo.Session, r *dgo.Ready) {
	invitePerms := dgo.PermissionUseSlashCommands
	inviteURL := fmt.Sprintf(
		"https://discordapp.com/oauth2/authorize?client_id=%s&scope=bot&permissions=%d",
		s.State.User.ID,
		invitePerms,
	)
	engine.CogLog(c, log.Info()).Msgf(
		"*** Bot ready:  %s#%s", s.State.User.Username, s.State.User.Discriminator)
	engine.CogLog(c, log.Info()).Msgf(
		"*** Bot invite: %s", inviteURL)
	engine.CogLog(c, log.Info()).Msgf(
		"*** %s", c.cfg.MOTD)

}
