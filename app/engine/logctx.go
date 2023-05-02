package engine

import (
	"context"
	"fmt"
	"time"

	"github.com/fiffu/arisa3/app/types"
	"github.com/fiffu/arisa3/lib"

	zero "github.com/rs/zerolog"
	zerolog "github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

type internalKey string

const (
	mapKey internalKey = "_map"
)

type CtxKey string

// Context keys
const (
	TraceID CtxKey = "trace"
	User    CtxKey = "user"
	cogName CtxKey = "cog" // This should only be used by the 'engine' package during dependency injection
)

func Put(parent context.Context, key CtxKey, value any) context.Context {
	ctx, m := GetMap(parent)
	m[key] = stringify(value)
	return ctx
}

func Get(ctx context.Context, key CtxKey) string {
	_, m := GetMap(ctx)
	return m[key]
}

func GetMap(parent context.Context) (context.Context, map[CtxKey]string) {
	m := parent.Value(mapKey)
	if m != nil {
		if m, ok := m.(map[CtxKey]string); ok {
			return parent, m
		}
	}
	newMap := make(map[CtxKey]string)
	ctx := context.WithValue(parent, mapKey, newMap)
	return ctx, newMap
}

func stringify(value any) (str string) {
	type stringer interface{ String() string }

	switch v := value.(type) {
	case string:
		str = v
	case time.Time:
		str = v.UTC().Format(time.RFC3339)
	case stringer:
		str = value.(stringer).String()
	default:
		str = fmt.Sprint(value)
	}
	return
}

func Debugf(ctx context.Context, msg string, args ...any) {
	msg, entry := newEntry(ctx, zerolog.Debug(), lib.WhoCalledMe(), msg)
	entry.Msgf(msg, args...)
}

func Infof(ctx context.Context, msg string, args ...any) {
	msg, entry := newEntry(ctx, zerolog.Info(), lib.WhoCalledMe(), msg)
	entry.Msgf(msg, args...)
}

func Warnf(ctx context.Context, msg string, args ...any) {
	msg, entry := newEntry(ctx, zerolog.Warn(), lib.WhoCalledMe(), msg)
	entry.Msgf(msg, args...)
}

func Errorf(ctx context.Context, err error, msg string, args ...any) {
	msg, entry := newEntry(ctx, zerolog.Error(), lib.WhoCalledMe(), msg)
	entry.Err(err).Stack().Msgf(msg, args...)
}

func Stack(ctx context.Context, err error) {
	msg := "Stack dump"
	msg, entry := newEntry(ctx, zerolog.Error(), lib.WhoCalledMe(), msg)
	zero.ErrorStackMarshaler = pkgerrors.MarshalStack
	entry.Stack().Err(err).Msg(msg)
}

func newEntry(ctx context.Context, entry *zero.Event, caller, msg string) (string, *zero.Event) {
	_, m := GetMap(ctx)
	for k, v := range m {
		if k == TraceID {
			msg = "[" + v + "]"
		}
		entry = entry.Str(string(k), v)
	}
	entry.Str("src", caller)
	return msg, entry
}

func CommandLog(cog types.ICog, evt types.ICommandEvent, e *zero.Event) *zero.Event {
	return e.Str(types.CtxCog, cog.Name()).Str(types.CtxCommand, evt.Command().Name())
}

func EventLog(cog types.ICog, evt types.IEvent, e *zero.Event) *zero.Event {
	return e.Str(types.CtxCog, cog.Name()).Str(types.CtxEvent, evt.Name())
}
