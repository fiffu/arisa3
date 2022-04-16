package general

import (
	"arisa3/app/engine"
	"arisa3/app/types"
	"context"
	"fmt"

	dgo "github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

// Cog implements types.ICog
type Cog struct {
}

func (c *Cog) Name() string { return "general" }

func NewCog(a types.IApp) types.ICog {
	return &Cog{}
}

func (c *Cog) OnStartup(ctx context.Context, sess *dgo.Session, rawConfig types.CogConfig) error {
	return engine.DefaultStartup(ctx, sess, rawConfig, c)
}

func (c *Cog) ConfigType() interface{} {
	return nil
}

func (c *Cog) Configure(ctx context.Context, cfg interface{}) error {
	return nil
}

func (c *Cog) ReadyCallback(s *dgo.Session, r *dgo.Ready) {
	welcome(s)
}

func welcome(s *dgo.Session) {
	invitePerms := dgo.PermissionUseSlashCommands
	inviteURL := fmt.Sprintf(
		"https://discordapp.com/oauth2/authorize?client_id=%s&scope=bot&permissions=%d",
		s.State.User.ID,
		invitePerms,
	)
	log.Info().Msgf("*** Bot ready")
	log.Info().Msgf("*** Bot user:   %s#%s", s.State.User.Username, s.State.User.Discriminator)
	log.Info().Msgf("*** Bot invite: %s", inviteURL)
}
