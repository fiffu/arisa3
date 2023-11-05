package log

import (
	"context"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"os"
	"time"

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
	TraceID    CtxKey = "trace"
	TraceSubID CtxKey = "trace-subspan"
	User       CtxKey = "user"
	Guild      CtxKey = "guild"
	CogName    CtxKey = "cog"
)

var DoNotLogCtxKeys = []CtxKey{}

func Hash(s string) string {
	hasher := sha1.New()
	hasher.Write([]byte(s))
	return base64.StdEncoding.EncodeToString(hasher.Sum(nil))
}

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
	var traceID, subTraceID string
	if value, ok := m[TraceID]; ok {
		traceID = fmt.Sprintf("[%s] ", value)
	}
	if value, ok := m[TraceSubID]; ok {
		subTraceID = fmt.Sprintf("[%s] ", value)
	}
	for k, v := range m {
		switch k {
		case TraceID:
			traceID = v
		case TraceSubID:
			subTraceID = v
		default:
			entry = entry.Str(string(k), v)
		}
	}
	entry.Str("src", caller)
	msg = traceID + subTraceID + msg
	return msg, entry
}

func SetupLogger() {
	output := zero.ConsoleWriter{Out: os.Stdout}
	output.TimeFormat = "2006/01/02 15:04:05"
	output.FormatMessage = func(i interface{}) string {
		return fmt.Sprintf(":: %s  ", i)
	}
	zero.TimestampFunc = func() time.Time {
		return time.Now().UTC()
	}
	zerolog.Logger = zerolog.Output(output).Level(zero.InfoLevel)
}