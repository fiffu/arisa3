package engine

import (
	"context"
	"fmt"
	"time"

	dgo "github.com/bwmarrin/discordgo"
	"github.com/fiffu/arisa3/app/instrumentation"
	"github.com/fiffu/arisa3/app/log"
)

func NewEventHandler[E any](inst instrumentation.Client, callable func(context.Context, *dgo.Session, E)) func(*dgo.Session, E) {
	return func(s *dgo.Session, evt E) {
		traceID := fmt.Sprintf("%T-%d", evt, time.Now().UTC().Unix())
		ctx := context.Background()
		ctx = log.Put(ctx, log.TraceID, traceID)

		ctx, span := instrumentation.SpanInContext(ctx, instrumentation.EventHandler(evt, callable))
		span.SetAttributes(instrumentation.KV.TraceID(traceID))
		defer span.End()

		callable(ctx, s, evt)
	}
}
