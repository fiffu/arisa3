package engine

import (
	"context"
	"fmt"
	"time"

	dgo "github.com/bwmarrin/discordgo"
	"github.com/fiffu/arisa3/app/instrumentation"
	"github.com/fiffu/arisa3/app/log"
	"github.com/fiffu/arisa3/app/types"
)

func NewEventHandler[E types.SupportedEvents](callable func(context.Context, *dgo.Session, E)) func(*dgo.Session, E) {
	return func(s *dgo.Session, evt E) {
		traceID := fmt.Sprintf("%T-%d", evt, time.Now().UTC().UnixMilli())
		ctx := context.Background()
		ctx = log.Put(ctx, log.TraceID, traceID)

		ctx, span := instrumentation.SpanInContext(ctx, instrumentation.EventHandler(evt, callable))
		span.SetAttributes(
			instrumentation.KV.EventName(fmt.Sprintf("%T", evt)),
			instrumentation.KV.TraceID(traceID),
		)
		defer span.End()

		mustHandleEvent(ctx, callable, s, evt)
	}
}
