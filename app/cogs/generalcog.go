package cogs

import (
	"arisa3/app/types"
	"context"
	"fmt"
	"math/rand"

	dgo "github.com/bwmarrin/discordgo"
)

type generalCog struct {
	util     *util
	registry *CommandsRegistry
	app      IApp
	config   *generalConfig
}
type generalConfig struct {
	Greeting string
}

func (c *generalCog) Name() string { return "general" }

func (c *generalCog) New(app IApp) ICog {
	return &generalCog{
		app:      app,
		util:     NewUtil(c, app),
		registry: NewCommandRegistry(c, app),
	}
}

func (c *generalCog) OnStartup(ctx context.Context, config CogConfig, sess *dgo.Session) error {
	cfg := &generalConfig{}
	if err := c.util.ParseConfig(config, cfg); err != nil {
		return fmt.Errorf("unable to parse cog config: %w", err)
	}
	c.config = cfg
	c.registerEvents(sess)
	return nil
}

func (c *generalCog) registerEvents(sess *dgo.Session) {
	sess.AddHandler(func(s *dgo.Session, r *dgo.Ready) {
		c.OnReady(s, r)
	})
}

func (c *generalCog) OnReady(s *dgo.Session, r *dgo.Ready) {
	ctx := c.util.EventContext(r)
	c.welcome(ctx, s)

	if err := c.registerCommands(ctx, s); err != nil {
		c.app.Errorf(ctx, err, "failed to register commands")
		return
	}
	c.registry.Finalise(ctx, s)
}

func (c *generalCog) registerCommands(ctx context.Context, s *dgo.Session) error {
	err := c.registry.Register(
		ctx,
		s,
		types.NewCommand("roll").ForChat().
			Desc("Rolls dice (supports algebraic notation, such as !roll 3d5+10)").
			Options(types.NewOption("expression").String()).
			Handler(c.roll),
	)
	return err
}

func (c *generalCog) welcome(ctx context.Context, s *dgo.Session) {
	invitePerms := dgo.PermissionUseSlashCommands
	inviteURL := fmt.Sprintf(
		"https://discordapp.com/oauth2/authorize?client_id=%s&scope=bot&permissions=%d",
		s.State.User.ID,
		invitePerms,
	)
	c.app.Infof(ctx, "*** Bot ready")
	c.app.Infof(ctx, "*** Bot user:   %s#%s", s.State.User.Username, s.State.User.Discriminator)
	c.app.Infof(ctx, "*** Bot invite: %s", inviteURL)
}

func (c *generalCog) roll(s *dgo.Session, i *dgo.InteractionCreate, cmd types.ICommand, args types.IArgs) error {
	ctx := c.util.CommandContext(cmd, i.User)
	resp := types.NewResponse()
	if expr, ok := args.String("expression"); ok {
		c.app.Infof(ctx, "expression='%s'", expr)
		lo, hi, num := roll(expr)
		resp.Content(fmt.Sprintf("Rolling %d-%d: %d", lo, hi, num))
	} else {
		c.app.Infof(ctx, "missing expression")
		resp.Content("missing expression")
	}
	return s.InteractionRespond(i.Interaction, resp.Data())
}

func roll(expr string) (int, int, int) {
	lo, hi := 0, 99
	return lo, hi, rollRand(lo, hi)
}

func rollRand(lo, hi int) int {
	return rand.Intn(lo+hi) - lo
}
