package commandfilters

import (
	"context"
	"fmt"

	"github.com/fiffu/arisa3/app/types"
)

type Filter func(types.ICommandEvent) bool
type CommandDecorator func(*types.Command) *types.Command
type FailureResponse types.ICommandResponse

type opcode int

type Middleware struct {
	filter                   Filter
	assertionFailureResponse FailureResponse
}

const (
	AND opcode = 0
	OR  opcode = 1
)

func NewMiddleware(filter Filter) *Middleware {
	defaultResponse := types.NewResponse().Content("Command failed! Try again later?")
	return &Middleware{
		filter,
		defaultResponse,
	}
}

func (mw *Middleware) And(f Filter) *Middleware {
	return mw.cmp(AND, f)
}

func (mw *Middleware) Or(f Filter) *Middleware {
	return mw.cmp(OR, f)
}

func (mw *Middleware) cmp(op opcode, f Filter) *Middleware {
	prev := mw.filter
	mw.filter = func(ev types.ICommandEvent) bool {
		return compare(op, ev, prev, f)
	}
	return mw
}

func (mw *Middleware) Exec(ev types.ICommandEvent) bool {
	return mw.filter(ev)
}

func (mw *Middleware) FailureResponse(resp FailureResponse) *Middleware {
	mw.assertionFailureResponse = resp
	return mw
}

func (mw *Middleware) CommandDecorator() CommandDecorator {
	return func(cmd *types.Command) *types.Command {
		// next is the subsequent handler to be called after this handler
		next := cmd.HandlerFunc()
		// assertionHandler calls next() only if Exec() returns true
		assertionHandler := func(ctx context.Context, ev types.ICommandEvent) error {
			if !mw.Exec(ev) {
				return ev.Respond(mw.assertionFailureResponse)
			}
			return next(ctx, ev)
		}
		// Overwrite command's handler with the assertionHandler
		cmd.Handler(assertionHandler)
		return cmd
	}
}

// compare executes comparison. We MUST pass in `left` and `right` as callables and only
// invoke them inline with comparison operators to achieve short-circuiting.
func compare(op opcode, ev types.ICommandEvent, left, right Filter) bool {
	switch op {
	case AND:
		return left(ev) && right(ev)
	case OR:
		return left(ev) || right(ev)
	}
	panic("unknown opcode: " + fmt.Sprint(op))
}
