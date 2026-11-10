package engine

import (
	"context"
	"errors"
	"testing"

	"github.com/fiffu/arisa3/app/types"
	"github.com/stretchr/testify/assert"

	dgo "github.com/bwmarrin/discordgo"
)

type readyCog struct {
	calls int
	errs  []error
}

func (c *readyCog) Name() string                                     { return "test" }
func (c *readyCog) ConfigPointer() types.StructPointer               { return nil }
func (c *readyCog) Configure(context.Context, types.CogConfig) error { return nil }
func (c *readyCog) ReadyCallback(context.Context, *dgo.Session, *dgo.Ready) error {
	i := c.calls
	c.calls++
	if i < len(c.errs) {
		return c.errs[i]
	}
	return nil
}

func Test_readyOnce(t *testing.T) {
	ctx := context.Background()

	t.Run("later Ready events are ignored after a success", func(t *testing.T) {
		cog := &readyCog{}
		ready := readyOnce(cog)
		ready(ctx, nil, &dgo.Ready{})
		ready(ctx, nil, &dgo.Ready{})
		assert.Equal(t, 1, cog.calls)
	})

	t.Run("a failed callback runs again on the next Ready", func(t *testing.T) {
		cog := &readyCog{errs: []error{errors.New("register failed")}}
		ready := readyOnce(cog)
		ready(ctx, nil, &dgo.Ready{})
		ready(ctx, nil, &dgo.Ready{})
		ready(ctx, nil, &dgo.Ready{})
		assert.Equal(t, 2, cog.calls)
	})
}
