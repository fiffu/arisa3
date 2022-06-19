package general

import (
	"github.com/fiffu/arisa3/app/engine"
	"github.com/rs/zerolog/log"

	dgo "github.com/bwmarrin/discordgo"
)

func (c *Cog) welcome(s *dgo.Session, r *dgo.Ready) {
	rootURL := "https://discordapp.com/oauth2/authorize?"
	clientID := "client_id=" + s.State.User.ID
	scope := "scope=applications.commands%20bot"
	inviteURL := rootURL + clientID + "&" + scope

	engine.CogLog(c, log.Info()).Msgf(
		"*** Bot ready:  %s#%s", s.State.User.Username, s.State.User.Discriminator)
	engine.CogLog(c, log.Info()).Msgf(
		"*** Bot invite: %s", inviteURL)
	engine.CogLog(c, log.Info()).Msgf(
		"*** %s", c.cfg.MOTD)
}
