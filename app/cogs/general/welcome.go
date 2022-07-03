package general

import (
	"fmt"
	"strings"

	"github.com/fiffu/arisa3/app/engine"
	"github.com/rs/zerolog/log"

	"github.com/bwmarrin/discordgo"
	dgo "github.com/bwmarrin/discordgo"
)

const requirePermissions = int64(discordgo.PermissionManageRoles)

func (c *Cog) welcome(s *dgo.Session, r *dgo.Ready) {
	rootURL := "https://discordapp.com/oauth2/authorize?"

	clientID := "client_id=" + s.State.User.ID
	scope := "scope=applications.commands%20bot"
	perms := fmt.Sprintf("permissions=%d", requirePermissions)

	inviteURL := rootURL + joinQueryParams(clientID, scope, perms)

	engine.CogLog(c, log.Info()).Msgf(
		"*** Bot ready:  %s#%s", s.State.User.Username, s.State.User.Discriminator)
	engine.CogLog(c, log.Info()).Msgf(
		"*** Bot invite: %s", inviteURL)
	engine.CogLog(c, log.Info()).Msgf(
		"*** %s", c.cfg.MOTD)
}

func joinQueryParams(qs ...string) string {
	return strings.Join(qs, "&")
}
