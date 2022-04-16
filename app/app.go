package app

import (
	"arisa3/app/cogs"
	"arisa3/app/engine"
	"arisa3/app/types"
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// app implements IApp
type app struct {
	cogsConfigs map[string]interface{}
}

func (a *app) Configs() map[string]interface{} {
	return a.cogsConfigs
}
func (a *app) Debugf(ctx context.Context, m string, args ...interface{}) {
	log.Ctx(ctx).Debug().Msgf(m, args...)
}
func (a *app) Infof(ctx context.Context, m string, args ...interface{}) {
	log.Ctx(ctx).Info().Msgf(m, args...)
}
func (a *app) Warnf(ctx context.Context, m string, args ...interface{}) {
	log.Ctx(ctx).Warn().Msgf(m, args...)
}
func (a *app) Errorf(ctx context.Context, err error, m string, args ...interface{}) {
	log.Ctx(ctx).Error().Err(err).Msgf(m, args...)
}
func (a *app) ContextWithValue(ctx context.Context, key, value string) context.Context {
	l := log.Ctx(ctx)
	if l.GetLevel() == zerolog.Disabled {
		l = &log.Logger
	}
	sublog := l.With().Str(key, value).Logger() // push key to context
	return sublog.WithContext(ctx)              // push logger to context
}

func Main(configPath string) error {
	app, sess, err := newApp(configPath)
	if err != nil {
		return err
	}

	ctx := engine.StartupContext()

	if err = cogs.SetupCogs(ctx, app, sess); err != nil {
		return err
	}
	if err := sess.Open(); err != nil {
		engine.SystemLog(log.Error()).Err(err).Msg("Failed to open session")
		return err
	}

	engine.SystemLog(log.Info()).Msg("Gateway session established")
	defer sess.Close()

	engine.SystemLog(log.Info()).Msg("Press Ctrl+C to exit")
	waitUntilInterrupt()

	return err
}

func newApp(configPath string) (types.IApp, *discordgo.Session, error) {
	setupLogger()

	cfg, err := Configure(configPath)
	if err != nil {
		return nil, nil, err
	}
	cogsCfg := getCogsConfigs(cfg)

	sess, err := discordgo.New("Bot " + cfg.BotSecret)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid bot parameters: %w", err)
	}

	return &app{
		cogsConfigs: cogsCfg,
	}, sess, nil
}

func setupLogger() {
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "02-Jan-06 15:04:05.000 -0700"}
	output.FormatMessage = func(i interface{}) string {
		return fmt.Sprintf("%-5s   ", i)
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

func waitUntilInterrupt() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
}
