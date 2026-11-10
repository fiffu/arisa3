package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fiffu/arisa3/app/cogs"
	"github.com/fiffu/arisa3/app/database"
	"github.com/fiffu/arisa3/app/engine"
	"github.com/fiffu/arisa3/app/instrumentation"
	"github.com/fiffu/arisa3/app/log"
	"github.com/fiffu/arisa3/app/types"
	"github.com/fiffu/arisa3/app/utils"

	"github.com/bwmarrin/discordgo"
)

// IDependencyInjector is an interface for initializing injected dependencies.
type IDependencyInjector interface {
	NewDatabase(ctx context.Context, dsn string) (database.IDatabase, error)
	NewInstrumentationClient(ctx context.Context) (instrumentation.Client, error)
	Bot(token string, debugMode bool) (*discordgo.Session, error)
}

// DefaultInjector provides default methods satisfying IDependencyInjector.
type DefaultInjector struct{}

func (d DefaultInjector) NewDatabase(ctx context.Context, dsn string) (database.IDatabase, error) {
	return database.NewDBClient(ctx, dsn)
}

func (d DefaultInjector) NewInstrumentationClient(ctx context.Context) (instrumentation.Client, error) {
	return instrumentation.NewInstrumentationClient(ctx)
}

func (d DefaultInjector) Bot(token string, debugMode bool) (*discordgo.Session, error) {
	sess, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}

	sess.Debug = debugMode
	sess.Client.Transport = utils.NewInstrumentedTransport()
	return sess, nil
}

// shutdownTimeout is the time that all cogs together get to stop. It must stay well under
// the time the deploy waits after SIGTERM before it kills the container.
const shutdownTimeout = 10 * time.Second

// app implements IApp
type app struct {
	cogsConfigs map[string]interface{}
	db          database.IDatabase
	inst        instrumentation.Client
	sess        *discordgo.Session
	cogs        []types.ICog
}

func (a *app) Configs() map[string]interface{} { return a.cogsConfigs }
func (a *app) Database() database.IDatabase    { return a.db }
func (a *app) BotSession() *discordgo.Session  { return a.sess }
func (a *app) Shutdown(ctx context.Context) {
	defer a.inst.Shutdown()
	// Close the gateway first, so no new events reach the cogs while they stop.
	if err := a.sess.Close(); err != nil {
		log.Errorf(ctx, err, "Error while closing session")
		log.Stack(ctx, err)
	}
	// Cogs stop before the database closes, so their last queries do not fail.
	cogs.ShutdownCogs(ctx, a.cogs, shutdownTimeout)
	if err := a.db.Close(ctx); err != nil {
		log.Errorf(ctx, err, "Error while closing DB connection")
		log.Stack(ctx, err)
	}
}

func Main(deps IDependencyInjector, configPath string) error {
	ctx := engine.StartupContext()

	app, err := newApp(ctx, deps, configPath)
	if err != nil {
		return err
	}

	log.Infof(ctx, "Initializing cogs")
	if app.cogs, err = cogs.SetupCogs(ctx, app); err != nil {
		return err
	}

	log.Infof(ctx, "Opening gateway session")
	if err := app.BotSession().Open(); err != nil {
		log.Errorf(ctx, err, "Failed to open session")
		return err
	}

	log.Infof(ctx, "Gateway session established")
	defer app.Shutdown(ctx)

	log.Infof(ctx, "Press Ctrl+C to exit")
	waitUntilInterrupt(ctx)

	return nil
}

func newApp(ctx context.Context, deps IDependencyInjector, configPath string) (*app, error) {
	log.SetupLogger()

	cfg, err := Configure(configPath)
	if err != nil {
		return nil, err
	}
	cogsCfg := getCogsConfigs(cfg)

	db, err := deps.NewDatabase(ctx, cfg.DatabaseDSN)
	if err != nil {
		return nil, err
	}

	inst, err := instrumentation.NewInstrumentationClient(ctx)
	if err != nil {
		return nil, err
	}

	sess, err := deps.Bot(cfg.BotSecret, cfg.EnableDebug)
	if err != nil {
		return nil, fmt.Errorf("invalid bot parameters: %w", err)
	}

	return &app{
		cogsConfigs: cogsCfg,
		db:          db,
		inst:        inst,
		sess:        sess,
	}, nil
}

func getCogsConfigs(cfg *Config) map[string]interface{} {
	out := make(map[string]interface{})
	for k, v := range cfg.Cogs {
		out[k] = interface{}(v)
	}
	return out
}

func waitUntilInterrupt(ctx context.Context) {
	stop := make(chan os.Signal, 1)
	// Note: docker stop will send SIGTERM, so handle that too
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	log.Infof(ctx, "Interrupted! Shutting down...")
}
