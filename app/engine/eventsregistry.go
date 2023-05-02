package engine

import (
	"context"
	"fmt"
	"time"

	dgo "github.com/bwmarrin/discordgo"
)

func NewEventHandler[E any](callable func(context.Context, *dgo.Session, E)) func(*dgo.Session, E) {
	return func(s *dgo.Session, evt E) {
		evtID := fmt.Sprintf("%T@%d", evt, time.Now().UTC().Unix())
		ctx := context.Background()
		ctx = Put(ctx, traceID, evtID)

		callable(ctx, s, evt)
	}
}
