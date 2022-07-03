package engine

import (
	"github.com/fiffu/arisa3/app/types"

	"github.com/rs/zerolog"
)

// AppLog contextextualizes on log messages from the base app.
func AppLog(e *zerolog.Event) *zerolog.Event {
	return e.Str(types.CtxEngine, types.CtxApp)
}

// AppLog contextextualizes on log messages from cog engine startup.
func StartupLog(e *zerolog.Event) *zerolog.Event {
	return e.Str(types.CtxEngine, types.CtxStartup)
}

// AppLog contextextualizes on log messages from the cog command registry.
func registryLog(e *zerolog.Event) *zerolog.Event {
	return e.Str(types.CtxEngine, types.CtxRegistry)
}

// AppLog contextextualizes on log messages from individual logs.
func CogLog(cog types.ICog, e *zerolog.Event) *zerolog.Event {
	return e.Str(types.CtxCog, cog.Name())
}

func CommandLog(cog types.ICog, evt types.ICommandEvent, e *zerolog.Event) *zerolog.Event {
	return e.Str(types.CtxCog, cog.Name()).Str(types.CtxCommand, evt.Command().Name())
}

func EventLog(cog types.ICog, evt types.IEvent, e *zerolog.Event) *zerolog.Event {
	return e.Str(types.CtxCog, cog.Name()).Str(types.CtxEvent, evt.Name())
}
