package cogs

import (
	"context"
	"errors"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

var (
	ErrMissingCogConfig = errors.New("missing config for cog")
)
var CogsList = []ICog{
	&generalCog{},
}

type CogConfig interface{}

type IApp interface {
	Configs() map[string]interface{}
	Debugf(context.Context, string, ...interface{})
	Infof(context.Context, string, ...interface{})
	Warnf(context.Context, string, ...interface{})
	Errorf(context.Context, error, string, ...interface{})
	ContextWithValue(ctx context.Context, key, value string) context.Context
}

type ICog interface {
	New(IApp) ICog
	Name() string
	OnStartup(ctx context.Context, config CogConfig, sess *discordgo.Session) error
}

func SetupCogs(ctx context.Context, a IApp, sess *discordgo.Session) error {
	configs := a.Configs()
	for _, c := range CogsList {
		cog := c.New(a)
		ctx = a.ContextWithValue(ctx, "cog", cog.Name())

		cfg, err := findConfig(cog, configs)
		if err != nil {
			return err
		}
		if err := cog.OnStartup(ctx, cfg, sess); err != nil {
			a.Errorf(ctx, err, "Failure to setup cog")
			return err
		}

		a.Infof(ctx, "Cog started")
	}
	return nil
}

func findConfig(cog ICog, cogConfigs map[string]interface{}) (CogConfig, error) {
	name := cog.Name()
	if cfg, ok := cogConfigs[name]; ok {
		return cfg, nil
	}
	return nil, fmt.Errorf("%w: %s", ErrMissingCogConfig, name)
}
