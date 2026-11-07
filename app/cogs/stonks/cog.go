// Package stonks implements the birbconomy: birbcoins and STONKS tickers.
package stonks

import (
	"context"
	"math/rand"
	"path/filepath"
	"time"

	"github.com/fiffu/arisa3/app/cogs/rng"
	"github.com/fiffu/arisa3/app/cogs/stonks/tickers"
	"github.com/fiffu/arisa3/app/database"
	"github.com/fiffu/arisa3/app/engine"
	"github.com/fiffu/arisa3/app/log"
	"github.com/fiffu/arisa3/app/types"
	"github.com/fiffu/arisa3/lib"
	"go.opentelemetry.io/otel"

	dgo "github.com/bwmarrin/discordgo"
)

var (
	migrationsDir = filepath.Join(lib.MustGetCallerDir(), "dbmigrations")
)

type Cog struct {
	db     database.IDatabase
	repo   IRepository
	ticker *ticker
}

func NewCog(a types.IApp) *Cog {
	return &Cog{db: a.Database()}
}

func (c *Cog) Name() string                       { return "stonks" }
func (c *Cog) ConfigPointer() types.StructPointer { return nil }
func (c *Cog) MigrationsDir() string              { return migrationsDir }

func (c *Cog) Configure(ctx context.Context, cfg types.CogConfig) error {
	c.repo = NewRepository(c.db)
	return nil
}

func (c *Cog) OnStartup(ctx context.Context, app types.IApp, rawConfig types.CogConfig) error {
	if err := engine.Bootstrap(ctx, app, rawConfig, c); err != nil {
		return err
	}
	return c.startTicker(ctx)
}

func (c *Cog) ReadyCallback(ctx context.Context, s *dgo.Session, r *dgo.Ready) error {
	return nil
}

func (c *Cog) startTicker(ctx context.Context) error {
	sources := tickers.All(rand.New(rand.NewSource(time.Now().UnixNano())))

	names := make([]string, len(sources))
	for i, s := range sources {
		names[i] = s.Ticker()
	}
	if err := c.repo.EnsureTickers(ctx, names); err != nil {
		return err
	}

	metrics, err := newSourceMetrics(otel.GetMeterProvider().Meter("arisa3/stonks"), time.Now)
	if err != nil {
		return err
	}
	c.ticker = &ticker{repo: c.repo, sources: sources, metrics: metrics, now: time.Now}
	if err := c.ticker.seedMetrics(ctx); err != nil {
		log.Errorf(ctx, err, "stonks: failed to seed source metrics")
	}

	runner := lib.NewRunner(SlotSize, c.ticker.Tick)
	runner.OnError = func(ctx context.Context, slot time.Time, err error) {
		log.Errorf(ctx, err, "stonks: update for slot %d had errors", slot.Unix())
	}
	go func() {
		// The process exits on shutdown. A tick that stops halfway finishes on the next start.
		_ = runner.Run(context.Background())
	}()
	return nil
}

// OnRoll implements rng.RollObserver.
func (c *Cog) OnRoll(ctx context.Context, roll rng.Roll) {
	if c.repo == nil || !countsForTicker(roll) {
		return
	}
	lo, hi := roll.Bounds()
	u := tickers.RollFraction(roll.Result, lo, hi)
	if err := c.repo.SaveRoll(ctx, tickers.LastRoll{U: u, At: time.Now().UTC()}); err != nil {
		log.Errorf(ctx, err, "stonks: failed to save roll")
	}
}

// Only rolls without a modifier count, such as "d50" but not "d50+2000".
// The default roll (d100-1) counts, because the user did not choose its modifier.
func countsForTicker(roll rng.Roll) bool {
	return !roll.Parsed || roll.Modifier == 0
}
