package engine

import (
	"arisa3/app/types"
	"context"

	dgo "github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

// IBootable describes a method set a cog should implement to be supported by Bootstrap()
type IBootable interface {
	Name() string
	ConfigPointer() types.StructPointer
	Configure(ctx context.Context, cfg interface{}) error
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