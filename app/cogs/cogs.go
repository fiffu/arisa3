package cogs

import (
	"arisa3/app/cogs/general"
	"arisa3/app/types"
	"context"
	"errors"
	"fmt"

	dgo "github.com/bwmarrin/discordgo"
)

var (
	ErrMissingCogConfig = errors.New("missing config for cog")
)

func getCogsList(app types.IApp) []types.ICog {
	return []types.ICog{
		general.NewCog(app),
	}
}

func SetupCogs(ctx context.Context, app types.IApp, sess *dgo.Session) error {
	configs := app.Configs()
	for _, c := range getCogsList(app) {
		ctx = app.ContextWithValue(ctx, "cog", c.Name())

		cfg, err := findConfig(c, configs)
		if err != nil {
			return err
		}
		if err := c.OnStartup(ctx, sess, cfg); err != nil {
			app.Errorf(ctx, err, "Failure to setup cog")
			return err
		}

		app.Infof(ctx, "Cog started")
	}
	return nil
}

func findConfig(cog types.ICog, cogConfigs map[string]interface{}) (types.CogConfig, error) {
	name := cog.Name()
	if cfg, ok := cogConfigs[name]; ok {
		return cfg, nil
	}
	return nil, fmt.Errorf("%w: %s", ErrMissingCogConfig, name)
}
