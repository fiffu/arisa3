package cogs

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/fiffu/arisa3/app/cogs/cardboard"
	"github.com/fiffu/arisa3/app/cogs/colours"
	"github.com/fiffu/arisa3/app/cogs/general"
	"github.com/fiffu/arisa3/app/cogs/rng"
	"github.com/fiffu/arisa3/app/cogs/stonks"
	"github.com/fiffu/arisa3/app/log"
	"github.com/fiffu/arisa3/app/types"
)

var (
	ErrMissingCogConfig = errors.New("missing config for cog")
)

// getCogsList maintains a list of cogs to load when the app starts.
func getCogsList(app types.IApp) []types.ICog {
	stonksCog := stonks.NewCog(app)
	return []types.ICog{
		general.NewCog(app),
		rng.NewCog(app, stonksCog),
		colours.NewCog(app),
		cardboard.NewCog(app),
		stonksCog,
	}
}

// SetupCogs loads cogs. It returns the cogs in start order, for ShutdownCogs.
func SetupCogs(ctx context.Context, app types.IApp) ([]types.ICog, error) {
	configs := app.Configs()

	cogs := getCogsList(app)
	for _, c := range cogs {

		cfg, err := findConfig(c, configs)
		if err != nil {
			return nil, err
		}
		if err := c.OnStartup(ctx, app, cfg); err != nil {
			log.Errorf(ctx, err, "Failed to setup cog: %s", c.Name())
			return nil, err
		}
		log.Infof(ctx, "%s cog init complete ⚙️", c.Name())
	}
	return cogs, nil
}

// ShutdownCogs calls OnShutdown on each cog that implements types.IShutdown, in reverse start order.
// All calls share one deadline. A cog that misses it is logged, and the next cog still gets called.
func ShutdownCogs(ctx context.Context, cogs []types.ICog, timeout time.Duration) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	for i := len(cogs) - 1; i >= 0; i-- {
		c, ok := cogs[i].(types.IShutdown)
		if !ok {
			continue
		}
		if err := c.OnShutdown(ctx); err != nil {
			log.Errorf(ctx, err, "Failed to shut down cog: %s", cogs[i].Name())
			continue
		}
		log.Infof(ctx, "%s cog shut down", cogs[i].Name())
	}
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
