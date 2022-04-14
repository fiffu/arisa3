package cogs

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type generalCog struct {
	*util
	app    IApp
	config *generalConfig
}
type generalConfig struct {
	Greeting string
}

func (c *generalCog) Name() string { return "general" }

func (c *generalCog) New(app IApp) ICog {
	return &generalCog{
		app:  app,
		util: &util{cog: c},
	}
}

func (c *generalCog) OnStartup(ctx context.Context, config CogConfig, sess *discordgo.Session) error {
	cfg := &generalConfig{}
	if err := c.util.ParseConfig(config, cfg); err != nil {
		return fmt.Errorf("unable to parse cog config: %w", err)
	}
	c.config = cfg
	c.registerEvents(sess)
	return nil
}

func (c *generalCog) registerEvents(sess *discordgo.Session) {
	sess.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		c.OnReady(s, r)
	})
}

func (c *generalCog) OnReady(s *discordgo.Session, r *discordgo.Ready) {
	ctx := c.util.Context(c.app, r)
	c.app.Warnf(ctx, "Bot is reconnected and ready")
}
