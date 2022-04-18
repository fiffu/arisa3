package engine

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/fiffu/arisa3/app/types"
	"github.com/fiffu/arisa3/lib/envconfig"

	dgo "github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

// Errors
var (
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

// StartupContext creates a runtime context for the app startup sequence.
func StartupContext() context.Context {
	// we can inject timeouts etc here
	type contextKey string
	ctx := context.Background()
	return context.WithValue(ctx, contextKey(CtxEngine), "startup")
}

// Bootstrap parses config and pushes it to cog, and sets up a handler for discordgo.Ready event.
func Bootstrap(ctx context.Context, sess *dgo.Session, rawConfig types.CogConfig, cog IBootable) error {
	cfg := cog.ConfigPointer()
	if err := ParseConfig(rawConfig, cfg); err != nil {
		return err
	}
	if replaced, err := envconfig.MergeEnvVars(cfg, EnvKeyPrefix(cog)); err != nil {
		return err
	} else if len(replaced) > 0 {
		for envKey, fld := range replaced {
			registryLog(log.Info()).Msgf(
				"replaced %v with environment var %s",
				fld,
				envKey,
			)
		}
	}
	if err := cog.Configure(ctx, cfg); err != nil {
		return err
	}

	sess.AddHandler(func(s *dgo.Session, r *dgo.Ready) {
		if err := cog.ReadyCallback(s, r); err != nil {
			log.Error().
				Str(CtxEngine, "ReadyCallback").
				Str(CtxCog, cog.Name()).
				Err(err).Msg("error in ReadyCallback")
		}
	})
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
