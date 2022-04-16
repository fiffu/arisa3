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

func SystemLog(e *zerolog.Event) *zerolog.Event {
	return e.Str(CtxEngine, CtxApp)
}

func StartupLog(e *zerolog.Event) *zerolog.Event {
	return e.Str(CtxEngine, CtxStartup)
}

func registryLog(e *zerolog.Event) *zerolog.Event {
	return e.Str(CtxEngine, CtxRegistry)
}

func CogLog(cog types.ICog, e *zerolog.Event) *zerolog.Event {
	return e.Str(CtxCog, cog.Name())
}

// func CommandLog(cog types.ICog, cmd types.ICommand, evt interface{}, e *zerolog.Event) *zerolog.Event {
// 	return e.Str(CtxCog, cog.Name()).Str(CtxCommand, cmd.Name())
// }

// func EventLog(cog types.ICog, evt interface{}, e *zerolog.Event) *zerolog.Event {
// 	evtName := ParseEvent(evt)
// 	return CogLog(cog).Str(CtxEvent, evtName)
// }
