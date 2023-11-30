package engine

import (
	"context"
	"testing"

	dgo "github.com/bwmarrin/discordgo"
	"github.com/fiffu/arisa3/app/log"
	"github.com/fiffu/arisa3/app/types"
	"github.com/stretchr/testify/assert"
)

func Test_mustHandleEvent(t *testing.T) {
	ctx := context.Background()
	hdlr := func(context.Context, *dgo.Session, *dgo.Ready) { panic("testing 123") }

	msg := log.CaptureLogging(t, func() {
		mustHandleEvent(ctx, hdlr, nil, nil)
	})
	assert.Contains(t, msg, "engine.mustHandleEvent[...]")
	assert.Contains(t, msg, "testing 123")
}

func Test_mustHandleCommand(t *testing.T) {
	ctx := context.Background()
	hdlr := func(context.Context, types.ICommandEvent) error { panic("testing 123") }

	msg := log.CaptureLogging(t, func() {
		mustHandleCommand(ctx, types.NewCommand("testcommand"), hdlr, nil, nil, nil)
	})
	assert.Contains(t, msg, "engine.mustHandleCommand")
	assert.Contains(t, msg, "testing 123")
}
