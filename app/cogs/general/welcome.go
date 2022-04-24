package general

import (
	"fmt"

	"github.com/fiffu/arisa3/app/engine"
	"github.com/rs/zerolog/log"

	dgo "github.com/bwmarrin/discordgo"
)

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
