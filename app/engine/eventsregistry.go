package engine

import (
	"context"
	"fmt"
	"time"

	dgo "github.com/bwmarrin/discordgo"
	"github.com/fiffu/arisa3/app/instrumentation"
	"github.com/fiffu/arisa3/app/log"
	"go.opentelemetry.io/otel/attribute"
)

func NewEventHandler[E any](inst instrumentation.Client, callable func(context.Context, *dgo.Session, E)) func(*dgo.Session, E) {
	return func(s *dgo.Session, evt E) {
		evtID := fmt.Sprintf("%T-%d", evt, time.Now().UTC().Unix())
		ctx := context.Background()
		ctx = log.Put(ctx, log.TraceID, evtID)

		evtName := fmt.Sprintf("%T", evt)
		ctx, span := inst.SpanInContext(ctx, instrumentation.EventScope, evtName)
		span.SetAttributes(attribute.String(string(log.TraceID), evtID))
		defer span.End()

		callable(ctx, s, evt)
	}
}
