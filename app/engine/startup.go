package engine

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/fiffu/arisa3/app/database"
	"github.com/fiffu/arisa3/app/types"
	"github.com/fiffu/arisa3/lib/envconfig"

	dgo "github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

// Errors
var (
	ErrBootstrap             = errors.New("bootstrap error")
	ErrCogNotBootable        = errors.New("cog does not implement IBootable")
	ErrCogParseConfig        = errors.New("unable to parse cog config")
	ErrUnexpectedConfigValue = errors.New("config type assert failed")
)

// IBootable describes a method set a cog should implement to be supported by Bootstrap()
type IBootable interface {
	Name() string
	ConfigPointer() types.StructPointer
	Configure(ctx context.Context, cfg types.CogConfig) error
	ReadyCallback(s *dgo.Session, r *dgo.Ready) error
}

type IRepository interface {
	Name() string
	MigrationsDir() string
}

// StartupContext creates a runtime context for the app startup sequence.
func StartupContext() context.Context {
	// we can inject timeouts etc here
	type contextKey string
	ctx := context.Background()
	return context.WithValue(ctx, contextKey(types.CtxEngine), "startup")
}

// Bootstrap parses config and pushes it to cog, and sets up a handler for discordgo.Ready event.
func Bootstrap(ctx context.Context, app types.IApp, rawConfig types.CogConfig, c interface{}) error {
	bootError := func(e error) error {
		return fmt.Errorf("%w: %v", ErrBootstrap, e)
	}

	// Assert interface satisfies IBootable
	cog, ok := c.(IBootable)
	if !ok {
		return bootError(ErrCogNotBootable)
	}

	// Parse config
	cfg := cog.ConfigPointer()
	if err := ParseConfig(rawConfig, cfg); err != nil {
		return bootError(err)
	}
	// Merge config from env vars
	if replaced, err := envconfig.MergeEnvVars(cfg, EnvKeyPrefix(cog)); err != nil {
		return bootError(err)
	} else if len(replaced) > 0 {
		for envKey, fld := range replaced {
			registryLog(log.Info()).Msgf(
				"replaced %v with environment var %s",
				fld,
				envKey,
			)
		}
	}
	// Assign config
	if err := cog.Configure(ctx, cfg); err != nil {
		return bootError(err)
	}

	// Setup repo migrations
	if rcog, ok := c.(IRepository); ok {
		db := app.Database()
		registryLog(log.Info()).Str(types.CtxCog, cog.Name()).Msgf(
			"Running migrations",
		)
		if err := runMigrations(rcog, db); err != nil {
			registryLog(log.Error()).Err(err).Str(types.CtxCog, cog.Name()).Msgf(
				"Migrations failed",
			)
			if closeErr := db.Close(); closeErr != nil {
				return bootError(fmt.Errorf(
					"failed to close DB connection (%v) during teardown due to "+
						"migration error (%v)",
					closeErr, err,
				))
			}
			return bootError(err)
		}
	} else {
		registryLog(log.Info()).Str(types.CtxCog, cog.Name()).Msgf(
			"Skipping migrations (interface assert failed)",
		)
	}

	// Bind ready callback after boot sequence is ready
	sess := app.BotSession()
	sess.AddHandler(func(s *dgo.Session, r *dgo.Ready) {
		if err := cog.ReadyCallback(s, r); err != nil {
			log.Error().
				Str(types.CtxEngine, "ReadyCallback").
				Str(types.CtxCog, cog.Name()).
				Err(err).Msg("error in ReadyCallback")
		}
	})
	return nil
}

func runMigrations(cog IRepository, db database.IDatabase) error {
	dir := cog.MigrationsDir()
	files, err := ioutil.ReadDir(dir)
	registryLog(log.Info()).Str(types.CtxCog, cog.Name()).Msgf(
		"Looking migrations in: %s (found %d files)", dir, len(files),
	)
	if err != nil {
		return err
	}
	for _, file := range files {
		path := filepath.Join(dir, file.Name())
		schema, err := db.ParseMigration(path)
		if err != nil {
			return err
		}
		if err := db.Migrate(schema); err != nil {
			return err
		}
	}
	return nil
}

// EnvKeyPrefix derives a prefix for environment keys from an IBootable cog.
func EnvKeyPrefix(cog IBootable) string {
	return fmt.Sprintf("ARISA3_%sCOG_", cog.Name())
}

// ParseConfig decodes input data and assigns it to output.
func ParseConfig(in types.CogConfig, out types.StructPointer) error {
	if out == nil {
		return nil
	}
	bytes, err := json.Marshal(in)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrCogParseConfig, err)
	}
	err = json.Unmarshal(bytes, out)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrCogParseConfig, err)
	}
	return nil
}

// UnexpectedConfigType is shorthand to create an error based on ErrUnexpectedConfigValue.
func UnexpectedConfigType(wanted interface{}, got interface{}) error {
	return fmt.Errorf("%w, wanted: %T, got: %#v", ErrUnexpectedConfigValue, wanted, got)
}
