package tickers

import (
	"context"
	"encoding/json"
	"math"
	"math/rand"
	"net/http"
	"testing"
	"time"

	"github.com/fiffu/arisa3/app/cogs/stonks/tickers/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var slot = time.Date(2026, 11, 7, 8, 0, 0, 0, time.UTC)

func Test_rndmStep_regimes(t *testing.T) {
	cases := []struct {
		prev Regime
		u1   float64
		want Regime
	}{
		{Flat, 0.00, Flat}, {Flat, 0.49, Flat}, {Flat, 0.50, Small}, {Flat, 0.89, Small}, {Flat, 0.90, Big},
		{Small, 0.29, Flat}, {Small, 0.30, Small}, {Small, 0.80, Big},
		{Big, 0.19, Flat}, {Big, 0.20, Small}, {Big, 0.61, Big},
		{"", 0.95, Big}, // no previous regime: odds of Flat
	}
	for _, tc := range cases {
		got, _ := rndmStep(tc.prev, rndmAnchor, tc.u1, 0.75)
		assert.Equal(t, tc.want, got, "prev=%q u1=%v", tc.prev, tc.u1)
	}
}

func Test_rndmStep_moves(t *testing.T) {
	const u1Small, u1Big = 0.6, 0.95 // from Flat

	// At the anchor, the pull is 0, so the move is exact.
	cases := []struct {
		name   string
		u1, u2 float64
		factor float64
	}{
		{"flat leaves the price alone", 0.1, 0.99, 1},
		{"small, smallest down", u1Small, 0.0, 1 / 1.01},
		{"small, smallest up", u1Small, 0.5, 1.01},
		{"small, middle up", u1Small, 0.75, 1.03},
		{"big, largest down", u1Big, 0.4999999999, 1 / 1.30},
		{"big, smallest up", u1Big, 0.5, 1.10},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, got := rndmStep(Flat, rndmAnchor, tc.u1, tc.u2)
			assert.InDelta(t, rndmAnchor*tc.factor, got, 1e-6)
		})
	}
}

func Test_rndmStep_upAndDownCancel(t *testing.T) {
	_, up := rndmStep(Flat, rndmAnchor, 0.95, 0.9) // big, ×1.26
	_, back := rndmStep(Flat, up, 0.95, 0.4)       // big, ÷1.26
	// The moves cancel. Only the second step's pull is left, from 26% above the anchor.
	want := rndmAnchor * math.Exp(-rndmPull*math.Log(1.26))
	assert.InDelta(t, want, back, 1e-9)
}

func Test_rndmStep_pullsTowardAnchor(t *testing.T) {
	_, high := rndmStep(Small, 1000, 0.5, 0.5) // smallest small move up
	assert.Less(t, high, 1000*1.01)
	_, low := rndmStep(Small, 1, 0.5, 0.4999999999) // smallest small move down... almost
	assert.Greater(t, low, 1/1.05)
}

func Test_Rndm_yearStaysInSaneRange(t *testing.T) {
	src := &Rndm{Rand: rand.New(rand.NewSource(1))}
	st := State[RndmMeta]{PrevPrice: 50, HasPrev: true}
	seen := map[Regime]int{}
	for i := 0; i < 2*365; i++ { // a year of 12-hour updates
		res, err := src.Fetch(context.Background(), st)
		require.NoError(t, err)
		require.Equal(t, OK, res.Outcome)
		seen[res.Patch.Regime]++
		st.PrevPrice, st.Meta = res.Raw, res.Patch
	}
	assert.Greater(t, st.PrevPrice, 1.0)
	assert.Less(t, st.PrevPrice, 2500.0)
	assert.Len(t, seen, 3, "every regime occurs")
}

func Test_Rndm_startsAtRandomPriceInRange(t *testing.T) {
	src := &Rndm{Rand: rand.New(rand.NewSource(1))}
	for i := 0; i < 100; i++ {
		res, err := src.Fetch(context.Background(), State[RndmMeta]{Slot: slot})
		require.NoError(t, err)
		assert.Equal(t, OK, res.Outcome)
		assert.Equal(t, RndmMeta{Regime: Flat}, res.Patch)
		assert.GreaterOrEqual(t, res.Raw, rndmStartMin)
		assert.Less(t, res.Raw, rndmStartMax)
	}
}

func Test_Roll(t *testing.T) {
	at := time.Date(2026, 11, 7, 7, 0, 0, 0, time.UTC)
	roll := &LastRoll{U: 1, At: at}

	cases := []struct {
		name    string
		st      State[RollMeta]
		raw     float64
		outcome Outcome
		patch   RollMeta
	}{
		{
			name:    "new roll moves price",
			st:      State[RollMeta]{PrevPrice: 100, HasPrev: true, Meta: RollMeta{LastRoll: roll}},
			raw:     110,
			outcome: OK,
			patch:   RollMeta{UsedRollAt: &at},
		},
		{
			name:    "used roll is not used again",
			st:      State[RollMeta]{PrevPrice: 100, HasPrev: true, Meta: RollMeta{LastRoll: roll, UsedRollAt: &at}},
			raw:     100,
			outcome: NoNewData,
		},
		{
			name:    "no roll yet",
			st:      State[RollMeta]{PrevPrice: 100, HasPrev: true},
			raw:     100,
			outcome: NoNewData,
		},
		{
			name:    "no history and no roll starts the history",
			st:      State[RollMeta]{},
			raw:     RollStartPrice,
			outcome: OK,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := (&Roll{}).Fetch(context.Background(), tc.st)
			require.NoError(t, err)
			assert.InDelta(t, tc.raw, res.Raw, 1e-9)
			assert.Equal(t, tc.outcome, res.Outcome)
			assert.Equal(t, tc.patch, res.Patch)
			assert.Nil(t, res.Patch.LastRoll, "a patch must never carry last_roll")
		})
	}
}

// Wrap decodes only the source's own keys, and its patch carries only the keys the source set.
func Test_Wrap(t *testing.T) {
	stored := json.RawMessage(`{"last_roll":{"u":1,"at":"2026-11-07T07:00:00Z"},"last_ok_at":"2026-11-07T04:00:00Z","source_ts":"x"}`)
	source := Wrap[RollMeta](&Roll{})

	res, err := source.Fetch(context.Background(), State[json.RawMessage]{PrevPrice: 100, HasPrev: true, Meta: stored})

	require.NoError(t, err)
	assert.InDelta(t, 110, res.Raw, 1e-9)
	assert.JSONEq(t, `{"used_roll_at":"2026-11-07T07:00:00Z"}`, string(res.Patch))
}

func Test_Wrap_badMetaIsAnError(t *testing.T) {
	source := Wrap[FredMeta](NewSP500(api.NewFred(http.DefaultTransport)))
	_, err := source.Fetch(context.Background(), State[json.RawMessage]{Meta: json.RawMessage(`{"source_date":5}`)})
	assert.ErrorContains(t, err, "decode meta")
}

func Test_RollFraction(t *testing.T) {
	assert.Equal(t, 0.0, RollFraction(0, 0, 99))   // default /roll, lowest
	assert.Equal(t, 1.0, RollFraction(99, 0, 99))  // default /roll, highest
	assert.Equal(t, 0.5, RollFraction(11, 2, 20))  // 2d10
	assert.Equal(t, 0.5, RollFraction(5, 5, 5))    // only one possible result
	assert.Equal(t, 1.0, RollFraction(500, 0, 99)) // clamped
}
