package app

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/fiffu/arisa3/app/cogs"
	"github.com/fiffu/arisa3/app/engine"
	"github.com/fiffu/arisa3/app/types"

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
		engine.AppLog(log.Error()).Err(err).Msg("Failed to open session")
		return err
	}

	engine.AppLog(log.Info()).Msg("Gateway session established")
	defer sess.Close()

	engine.AppLog(log.Info()).Msg("Press Ctrl+C to exit")
	waitUntilInterrupt()

	return nil
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
	output := zerolog.ConsoleWriter{Out: os.Stdout}
	output.TimeFormat = "02-Jan-06 15:04:05 -0700"
	output.FormatMessage = func(i interface{}) string {
		return fmt.Sprintf("%-5s   ", i)
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

func waitUntilInterrupt() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	engine.AppLog(log.Info()).Msg("Interrupted! Shutting down...")
}
