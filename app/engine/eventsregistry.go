package engine

import (
	"context"
	"fmt"
	"time"

	dgo "github.com/bwmarrin/discordgo"
	"github.com/fiffu/arisa3/app/log"
)

func NewEventHandler[E any](callable func(context.Context, *dgo.Session, E)) func(*dgo.Session, E) {
	return func(s *dgo.Session, evt E) {
		evtID := fmt.Sprintf("%T-%d", evt, time.Now().UTC().Unix())
		ctx := context.Background()
		ctx = log.Put(ctx, log.TraceID, evtID)

		callable(ctx, s, evt)
	}
}
