package general

import (
	"context"
	"fmt"
	"strings"

	dgo "github.com/bwmarrin/discordgo"
	"github.com/fiffu/arisa3/app/log"
)

const requirePermissions = int64(dgo.PermissionManageRoles)

func (c *Cog) welcome(s *dgo.Session, r *dgo.Ready) {
	rootURL := "https://discordapp.com/oauth2/authorize?"

	clientID := "client_id=" + s.State.User.ID
	scope := "scope=applications.commands%20bot"
	perms := fmt.Sprintf("permissions=%d", requirePermissions)

	inviteURL := rootURL + joinQueryParams(clientID, scope, perms)

	ctx := context.Background()
	log.Infof(ctx, "*** Bot ready:  %s#%s", s.State.User.Username, s.State.User.Discriminator)
	log.Infof(ctx, "*** Bot invite: %s", inviteURL)
	log.Infof(ctx, "*** %s", c.cfg.MOTD)
}

func joinQueryParams(qs ...string) string {
	return strings.Join(qs, "&")
}
