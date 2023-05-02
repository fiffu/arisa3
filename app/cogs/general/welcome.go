package general

import (
	"context"
	"fmt"
	"strings"

	"github.com/fiffu/arisa3/app/engine"

	dgo "github.com/bwmarrin/discordgo"
)

const requirePermissions = int64(dgo.PermissionManageRoles)

func (c *Cog) welcome(s *dgo.Session, r *dgo.Ready) {
	rootURL := "https://discordapp.com/oauth2/authorize?"

	clientID := "client_id=" + s.State.User.ID
	scope := "scope=applications.commands%20bot"
	perms := fmt.Sprintf("permissions=%d", requirePermissions)

	inviteURL := rootURL + joinQueryParams(clientID, scope, perms)

	ctx := context.Background()
	engine.Infof(ctx, "*** Bot ready:  %s#%s", s.State.User.Username, s.State.User.Discriminator)
	engine.Infof(ctx, "*** Bot invite: %s", inviteURL)
	engine.Infof(ctx, "*** %s", c.cfg.MOTD)
}

func joinQueryParams(qs ...string) string {
	return strings.Join(qs, "&")
}
