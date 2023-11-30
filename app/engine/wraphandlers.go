package engine

import (
	"context"
	"errors"
	"fmt"

	dgo "github.com/bwmarrin/discordgo"
	"github.com/fiffu/arisa3/app/instrumentation"
	"github.com/fiffu/arisa3/app/log"
	"github.com/fiffu/arisa3/app/types"
	"github.com/fiffu/arisa3/lib"
)

var (
	errPanic = errors.New("panic while executing handler")
)

func newErrPanic(recovered any) error {
	return fmt.Errorf("%w: %v", errPanic, recovered)
}

// mustHandleCommand executes a command's handler, trapping and logging any panics.
func mustHandleCommand(
	ctx context.Context,
	cmd types.ICommand,
	handler types.CommandHandler,
	args types.IArgs, s *dgo.Session, i *dgo.InteractionCreate,
) (returnErr error) {
	defer func() {
		if r := recover(); r != nil {
			instrumentation.EmitErrorf(ctx, "command %s panic: %v", cmd.Name(), r)
			returnErr = newErrPanic(r)
			log.Stack(ctx, returnErr)
		}
	}()

	returnErr = handler(ctx, types.NewCommandEvent(s, i, cmd, args))
	return
}

// mustHandleEvent executes an event handler, trapping and logging any panics.
func mustHandleEvent[E types.SupportedEvents](
	ctx context.Context,
	handler types.EventHandler[E],
	s *dgo.Session,
	evt E,
) {
	defer func() {
		if r := recover(); r != nil {
			instrumentation.EmitErrorf(ctx, "event handler %T panic: %v", lib.FuncName(handler), r)
			log.Stack(ctx, newErrPanic(r))
		}
	}()

	handler(ctx, s, evt)
}
