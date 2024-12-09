package engine

import (
	"testing"
	"time"

	"github.com/fiffu/arisa3/lib"
	"github.com/stretchr/testify/assert"
)

func Test_idempotencyKey_Before(t *testing.T) {
	now := time.Now().UTC()
	x := &idempotencyKey{"", now}
	xx := &idempotencyKey{"", now}
	y := &idempotencyKey{"", now.Add(1)}

	assert.True(t, x.Before(y))
	assert.False(t, y.Before(x))
	assert.False(t, y.Before(xx))
}

func Test_handlerIdempotency_Acquire(t *testing.T) {
	hi := newHandlerIdempotency()
	assert.True(t, hi.Acquire("abc"))
	assert.True(t, hi.Acquire("def"))
	assert.False(t, hi.Acquire("abc"))
}

func Test_handlerIdempotency_Acquire_afterExpiry(t *testing.T) {
	hi := newHandlerIdempotency()
	clock := lib.FrozenNow(t)
	hi.clock = clock.Now

	assert.True(t, hi.Acquire("abc"))
	assert.Equal(t, 1, hi.tree.Len())

	clock.Add(1 + 2*idempotencyWindow) // time passes
	assert.True(t, hi.Acquire("def"))
	assert.Equal(t, 1, hi.tree.Len()) // key 'abc' should have been deleted

	clock.Add(1)                      // later...
	assert.True(t, hi.Acquire("abc")) // since 'abc' was deleted, we can acquire it again
	assert.Equal(t, 2, hi.tree.Len()) // now there should be 2 keys
}
