package frames

import (
	"context"

	dgo "github.com/bwmarrin/discordgo"
	"github.com/fiffu/arisa3/app/engine"
	"github.com/fiffu/arisa3/app/types"
	"github.com/rs/zerolog/log"
)

// Cog implements ICog and IDefaultStartup
type Cog struct {
	commands *engine.CommandsRegistry
}

func NewCog(a types.IApp) types.ICog {
	return &Cog{
		commands: engine.NewCommandRegistry(),
	}
}

func (c *Cog) Name() string                       { return "general" }
func (c *Cog) ConfigPointer() types.StructPointer { return nil }

func (c *Cog) OnStartup(ctx context.Context, app types.IApp, rawConfig types.CogConfig) error {
	sess := app.BotSession()
	sess.AddHandler(func(s *dgo.Session, r *dgo.Ready) {
		err := c.commands.Register(
			s,
			c.tweetCommand(),
		)
		if err != nil {
			engine.CogLog(c, log.Error()).
				Err(err).
				Msg("Error while registering commands")
		}
	})
	return nil
}
