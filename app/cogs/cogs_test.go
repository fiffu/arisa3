package cogs

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/fiffu/arisa3/app/types"
	"github.com/stretchr/testify/assert"
)

type fakeCog struct {
	name  string
	calls *[]string
	// stop runs inside OnShutdown.
	stop func(ctx context.Context) error
}

func (f *fakeCog) Name() string                       { return f.name }
func (f *fakeCog) ConfigPointer() types.StructPointer { return nil }
func (f *fakeCog) OnStartup(context.Context, types.IApp, types.CogConfig) error {
	return nil
}

type shutdownCog struct{ fakeCog }

func (s *shutdownCog) OnShutdown(ctx context.Context) error {
	*s.calls = append(*s.calls, s.name)
	return s.stop(ctx)
}

func Test_ShutdownCogs(t *testing.T) {
	var calls []string
	ok := func(context.Context) error { return nil }
	hang := func(ctx context.Context) error { <-ctx.Done(); return ctx.Err() }
	failed := func(context.Context) error { return errors.New("boom") }

	cogs := []types.ICog{
		&shutdownCog{fakeCog{"first", &calls, ok}},
		&fakeCog{"no hook", &calls, nil},
		&shutdownCog{fakeCog{"hangs", &calls, hang}},
		&shutdownCog{fakeCog{"fails", &calls, failed}},
	}

	start := time.Now()
	ShutdownCogs(context.Background(), cogs, 50*time.Millisecond)

	assert.Equal(t, []string{"fails", "hangs", "first"}, calls, "reverse start order, and every cog still gets called")
	assert.Less(t, time.Since(start), time.Second, "the deadline bounds the hanging cog")
}
