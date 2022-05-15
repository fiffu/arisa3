package commandfilters

import (
	"fmt"

	"github.com/fiffu/arisa3/app/types"
)

type Filter func(types.ICommandEvent) bool
type CommandDecorator func(*types.Command) *types.Command

type opcode int

type Middleware struct {
	filter Filter
}

const (
	AND opcode = 0
	OR  opcode = 1
)

func NewMiddleware(f Filter) *Middleware {
	return &Middleware{f}
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

func (mw *Middleware) wrapHandler(next types.Handler, message string) types.Handler {
	return func(ev types.ICommandEvent) error {
		if ok := mw.Exec(ev); !ok {
			resp := types.NewResponse().Content(message)
			return ev.Respond(resp)
		}
		return next(ev)
	}
}

func (mw *Middleware) Asserts(message string) CommandDecorator {
	return func(cmd *types.Command) *types.Command {
		next := cmd.HandlerFunc()
		cmd.Handler(mw.wrapHandler(next, message))
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
