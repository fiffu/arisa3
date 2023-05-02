package engine

import (
	"context"
	"fmt"
	"time"

	dgo "github.com/bwmarrin/discordgo"
)

type eventHandler[E any] func(*dgo.Session, E)

func NewEventHandler[E any](callable func(context.Context, *dgo.Session, E)) eventHandler[E] {
	return func(s *dgo.Session, evt E) {
		evtID := fmt.Sprintf("%T@%d", evt, time.Now().UTC().Unix())
		ctx := context.Background()
		ctx = Put(ctx, TraceID, evtID)

		callable(ctx, s, evt)
	}
}
