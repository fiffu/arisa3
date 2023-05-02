package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/fiffu/arisa3/app/cogs"
	"github.com/fiffu/arisa3/app/database"
	"github.com/fiffu/arisa3/app/engine"
	"github.com/fiffu/arisa3/app/types"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// IDependencyInjector is an interface for initializing injected dependencies.
type IDependencyInjector interface {
	NewDatabase(dsn string) (database.IDatabase, error)
	Bot(token string, debugMode bool) (*discordgo.Session, error)
}

// DefaultInjector provides default methods satisfying IDependencyInjector.
type DefaultInjector struct{}

func (d DefaultInjector) NewDatabase(dsn string) (database.IDatabase, error) {
	return database.NewDBClient(dsn)
}

func (d DefaultInjector) Bot(token string, debugMode bool) (*discordgo.Session, error) {
	sess, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}

	sess.Debug = debugMode
	return sess, nil
}

// app implements IApp
type app struct {
	cogsConfigs map[string]interface{}
	db          database.IDatabase
	sess        *discordgo.Session
}

func (a *app) Configs() map[string]interface{} { return a.cogsConfigs }
func (a *app) Database() database.IDatabase    { return a.db }
func (a *app) BotSession() *discordgo.Session  { return a.sess }
func (a *app) Shutdown(ctx context.Context) {
	if err := a.sess.Close(); err != nil {
		engine.Errorf(ctx, err, "Error while closing session")
		engine.Stack(ctx, err)
	}
	if err := a.db.Close(); err != nil {
		engine.Errorf(ctx, err, "Error while closing DB connection")
		engine.Stack(ctx, err)
	}
}

func Main(deps IDependencyInjector, configPath string) error {
	app, err := newApp(deps, configPath)
	if err != nil {
		return err
	}

	ctx := engine.StartupContext()
	ctx = engine.Put(ctx, engine.FromEngine, "main")

	engine.Infof(ctx, "Initializing cogs")
	if err = cogs.SetupCogs(ctx, app); err != nil {
		return err
	}

	engine.Infof(ctx, "Opening gateway session")
	if err := app.BotSession().Open(); err != nil {
		engine.Errorf(ctx, err, "Failed to open session")
		return err
	}

	engine.Infof(ctx, "Gateway session established")
	defer app.Shutdown(ctx)

	engine.Infof(ctx, "Press Ctrl+C to exit")
	waitUntilInterrupt(ctx)

	return nil
}

func newApp(deps IDependencyInjector, configPath string) (types.IApp, error) {
	setupLogger()

	cfg, err := Configure(configPath)
	if err != nil {
		return nil, err
	}
	cogsCfg := getCogsConfigs(cfg)

	db, err := deps.NewDatabase(cfg.DatabaseDSN)
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
		sess:        sess,
	}, nil
}

func setupLogger() {
	output := zerolog.ConsoleWriter{Out: os.Stdout}
	output.TimeFormat = "2006/01/02 15:04:05"
	output.FormatMessage = func(i interface{}) string {
		return fmt.Sprintf(":: %s  ", i)
	}
	zerolog.TimestampFunc = func() time.Time {
		return time.Now().UTC()
	}
	log.Logger = log.Output(output).Level(zerolog.InfoLevel)
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
	signal.Notify(stop, os.Interrupt)
	<-stop
	engine.Infof(ctx, "Interrupted! Shutting down...")
}
