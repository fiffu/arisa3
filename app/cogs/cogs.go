package cogs

import (
	"context"
	"errors"
	"fmt"

	"github.com/fiffu/arisa3/app/cogs/greeter"
	"github.com/fiffu/arisa3/app/cogs/rng"
	"github.com/fiffu/arisa3/app/engine"
	"github.com/fiffu/arisa3/app/types"

	dgo "github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

var (
	ErrMissingCogConfig = errors.New("missing config for cog")
)

// getCogsList maintains a list of cogs to load when the app starts.
func getCogsList(app types.IApp) []types.ICog {
	return []types.ICog{
		greeter.NewCog(app),
		rng.NewCog(app),
	}
}

// SetupCogs loads cogs.
func SetupCogs(ctx context.Context, app types.IApp, sess *dgo.Session) error {
	configs := app.Configs()

	for _, c := range getCogsList(app) {

		cfg, err := findConfig(c, configs)
		if err != nil {
			return err
		}
		if err := c.OnStartup(ctx, sess, cfg); err != nil {
			engine.StartupLog(log.Error()).
				Str(engine.CtxCog, c.Name()).
				Err(err).
				Msg("Failure to setup cog")
			return err
		}
		engine.StartupLog(log.Info()).
			Str(engine.CtxCog, c.Name()).
			Msg("Cog started")
	}
	return nil
}

// findConfig retrieves raw cog config from the app's root config.
func findConfig(cog types.ICog, cogConfigs map[string]interface{}) (types.CogConfig, error) {
	name := cog.Name()
	if cog.ConfigPointer() == nil {
		return nil, nil
	}
	if cfg, ok := cogConfigs[name]; ok {
		return cfg, nil
	}
	return nil, fmt.Errorf("%w: %s", ErrMissingCogConfig, name)
}
