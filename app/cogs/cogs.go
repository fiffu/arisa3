package cogs

import (
	"context"
	"errors"
	"fmt"

	"github.com/fiffu/arisa3/app/cogs/cardboard"
	"github.com/fiffu/arisa3/app/cogs/colours"
	"github.com/fiffu/arisa3/app/cogs/general"
	"github.com/fiffu/arisa3/app/cogs/rng"
	"github.com/fiffu/arisa3/app/engine"
	"github.com/fiffu/arisa3/app/types"
)

var (
	ErrMissingCogConfig = errors.New("missing config for cog")
)

// getCogsList maintains a list of cogs to load when the app starts.
func getCogsList(app types.IApp) []types.ICog {
	return []types.ICog{
		general.NewCog(app),
		rng.NewCog(app),
		colours.NewCog(app),
		cardboard.NewCog(app),
	}
}

// SetupCogs loads cogs.
func SetupCogs(ctx context.Context, app types.IApp) error {
	configs := app.Configs()

	for _, c := range getCogsList(app) {

		cfg, err := findConfig(c, configs)
		if err != nil {
			return err
		}
		if err := c.OnStartup(ctx, app, cfg); err != nil {
			engine.Errorf(ctx, err, "Failed to setup cog: %s", c.Name())
			return err
		}
		engine.Infof(ctx, "%s cog init complete ⚙️", c.Name())
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
