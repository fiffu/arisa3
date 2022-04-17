package engine

import (
	"arisa3/app/types"

	"github.com/rs/zerolog"
)

const (
	// Keys
	CtxEngine  = "engine"
	CtxCog     = "cog"
	CtxCommand = "command"
	CtxEvent   = "event"

	// Values
	CtxApp      = "app"
	CtxRegistry = "commandsRegistry"
	CtxStartup  = "startup"
)

// AppLog contextextualizes on log messages from the base app.
func AppLog(e *zerolog.Event) *zerolog.Event {
	return e.Str(CtxEngine, CtxApp)
}

// AppLog contextextualizes on log messages from cog engine startup.
func StartupLog(e *zerolog.Event) *zerolog.Event {
	return e.Str(CtxEngine, CtxStartup)
}

// AppLog contextextualizes on log messages from the cog command registry.
func registryLog(e *zerolog.Event) *zerolog.Event {
	return e.Str(CtxEngine, CtxRegistry)
}

// AppLog contextextualizes on log messages from individual logs.
func CogLog(cog types.ICog, e *zerolog.Event) *zerolog.Event {
	return e.Str(CtxCog, cog.Name())
}

func CommandLog(cog types.ICog, evt types.ICommandEvent, e *zerolog.Event) *zerolog.Event {
	return e.Str(CtxCog, cog.Name()).Str(CtxCommand, evt.Command().Name())
}

// func EventLog(cog types.ICog, evt interface{}, e *zerolog.Event) *zerolog.Event {
// 	evtName := ParseEvent(evt)
// 	return CogLog(cog).Str(CtxEvent, evtName)
// }
