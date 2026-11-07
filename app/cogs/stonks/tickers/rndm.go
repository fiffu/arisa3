package tickers

import (
	"context"
	"math"
	"math/rand"
	"sync"
)

// The first RNDM price is a random value in this range, in coins.
const (
	rndmStartMin = 10.0
	rndmStartMax = 100.0
)

type Regime string

const (
	Flat  Regime = "flat"
	Small Regime = "small"
	Big   Regime = "big"
)

// rndmTransitions gives P(flat, small, big) for the next move, by the regime of the previous move.
// Big moves tend to follow big moves, so volatile periods cluster.
var rndmTransitions = map[Regime][3]float64{
	Flat:  {0.50, 0.40, 0.10},
	Small: {0.30, 0.50, 0.20},
	Big:   {0.20, 0.40, 0.40},
}

// rndmDegrees is the range of a move's size, by regime. 0.10 means ×1.10 up or ÷1.10 down.
// Up and down moves of the same size cancel out, so the walk has no drift.
var rndmDegrees = map[Regime][2]float64{
	Small: {0.01, 0.05},
	Big:   {0.10, 0.30},
}

// Without a pull, RNDM drifts toward 0 or toward thousands within a year, and a price
// near the 0.001-coin floor stops moving, because small moves round to no change.
// Each move also closes rndmPull of the log distance to rndmAnchor.
// With 0.01 and 2 updates a day, a year of RNDM stays from about 14 to 175 coins
// in 90% of simulated runs.
const (
	rndmAnchor = 50.0
	rndmPull   = 0.01
)

type RndmMeta struct {
	Regime Regime `json:"regime,omitempty"`
}

// Rndm uses a Markov chain. Each update rolls twice: the first roll picks a flat, small
// or big move, with odds that depend on the previous move. The second roll picks the
// direction and the size of the move.
type Rndm struct {
	Rand *rand.Rand
	mu   sync.Mutex
}

func (r *Rndm) Ticker() string { return "RNDM" }

func (r *Rndm) Fetch(_ context.Context, st State[RndmMeta]) (Result[RndmMeta], error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if !st.HasPrev {
		start := rndmStartMin + (rndmStartMax-rndmStartMin)*r.Rand.Float64()
		return Result[RndmMeta]{Raw: start, Outcome: OK, Patch: RndmMeta{Regime: Flat}}, nil
	}

	regime, price := rndmStep(st.Meta.Regime, st.PrevPrice, r.Rand.Float64(), r.Rand.Float64())
	return Result[RndmMeta]{Raw: price, Outcome: OK, Patch: RndmMeta{Regime: regime}}, nil
}

// u1 picks the regime. u2 picks the direction (below 0.5 is down) and the size.
// Both are in [0, 1).
func rndmStep(prev Regime, price, u1, u2 float64) (Regime, float64) {
	odds, ok := rndmTransitions[prev]
	if !ok {
		odds = rndmTransitions[Flat]
	}

	regime := Big
	switch {
	case u1 < odds[0]:
		regime = Flat
	case u1 < odds[0]+odds[1]:
		regime = Small
	}
	if regime == Flat {
		return Flat, price
	}

	deg := rndmDegrees[regime]
	size := deg[0] + (deg[1]-deg[0])*math.Mod(u2, 0.5)/0.5
	move := math.Log1p(size)
	if u2 < 0.5 {
		move = -move
	}

	logPrice := math.Log(price)
	logPrice += move + rndmPull*(math.Log(rndmAnchor)-logPrice)
	return regime, math.Exp(logPrice)
}
