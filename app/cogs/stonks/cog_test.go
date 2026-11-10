package stonks

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/fiffu/arisa3/app/cogs/rng"
	"github.com/fiffu/arisa3/app/cogs/stonks/tickers"
	"github.com/fiffu/arisa3/lib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_OnRoll_savesOnlyRollsWithoutModifier(t *testing.T) {
	cases := []struct {
		name  string
		roll  rng.Roll
		saved bool
	}{
		{"default roll, d100-1", rng.Roll{Result: 50, Count: 1, Sides: 100, Modifier: -1, Parsed: false}, true},
		{"d50", rng.Roll{Result: 25, Count: 1, Sides: 50, Parsed: true}, true},
		{"d50+2000", rng.Roll{Result: 2025, Count: 1, Sides: 50, Modifier: 2000, Parsed: true}, false},
		{"3d6-1", rng.Roll{Result: 9, Count: 3, Sides: 6, Modifier: -1, Parsed: true}, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			repo := newFakeRepo()
			require.NoError(t, repo.EnsureTickers(context.Background(), []string{"ROLL"}))
			cog := &Cog{repo: repo}

			cog.OnRoll(context.Background(), tc.roll)

			_, saved := repo.meta["ROLL"]["last_roll"]
			assert.Equal(t, tc.saved, saved)
		})
	}
}

// blockingSource fails only when ctx ends, like an HTTP fetch cut off by shutdown.
type blockingSource struct {
	started chan struct{}
}

func (b *blockingSource) Ticker() string { return "BLOCK" }
func (b *blockingSource) Fetch(ctx context.Context, _ tickers.State[json.RawMessage]) (tickers.Result[json.RawMessage], error) {
	close(b.started)
	<-ctx.Done()
	return tickers.Result[json.RawMessage]{}, ctx.Err()
}

func Test_OnShutdown_stopsUpdateWithoutPricingTheSlot(t *testing.T) {
	ctx := context.Background()
	repo := newFakeRepo()
	require.NoError(t, repo.EnsureTickers(ctx, []string{"BLOCK"}))
	prevSlot := lib.SlotOf(time.Now(), SlotSize).Add(-SlotSize).Unix()
	_, err := repo.AppendPrice(ctx, "BLOCK", 5000, prevSlot, nil)
	require.NoError(t, err)

	src := &blockingSource{started: make(chan struct{})}
	cog := &Cog{repo: repo}
	require.NoError(t, cog.startTicker(ctx, []tickers.Source{src}))
	<-src.started

	shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	require.NoError(t, cog.OnShutdown(shutdownCtx))

	repo.mu.Lock()
	defer repo.mu.Unlock()
	assert.Equal(t, []int64{prevSlot}, repo.slots["BLOCK"], "the stopped slot must stay unpriced")
}
