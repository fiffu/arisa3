package tickers

import (
	"context"
	"time"
)

type RollMeta struct {
	// Only /roll writes LastRoll. Fetch must never put it in a patch, because
	// a patch with the old value would erase a roll saved during the update.
	LastRoll   *LastRoll  `json:"last_roll,omitempty"`
	UsedRollAt *time.Time `json:"used_roll_at,omitempty"`
}

type LastRoll struct {
	// U is the roll scaled to 0..1. See RollFraction().
	U  float64   `json:"u"`
	At time.Time `json:"at"`
}

// RollStartPrice is in coins.
const RollStartPrice = 50.0

// Roll multiplies the price by 0.9 + 0.2 × u on each update, where u is the last
// /roll scaled to 0..1. A roll moves the price by −10% to +10%.
// Each roll is used once. With no new roll, the price stays the same.
//
// Users can steer ROLL on purpose, and that is by design. Only the last roll before
// an update counts, so a user can roll until the result is high, for example /roll d2
// until it shows 2, and ROLL then rises 10%.
type Roll struct{}

func (r *Roll) Ticker() string { return "ROLL" }

func (r *Roll) Fetch(_ context.Context, st State[RollMeta]) (Result[RollMeta], error) {
	prev := RollStartPrice
	if st.HasPrev {
		prev = st.PrevPrice
	}

	last, used := st.Meta.LastRoll, st.Meta.UsedRollAt
	if last == nil || (used != nil && used.Equal(last.At)) {
		if !st.HasPrev {
			// Start the price history even without a roll.
			return Result[RollMeta]{Raw: prev, Outcome: OK}, nil
		}
		return Result[RollMeta]{Raw: prev, Outcome: NoNewData}, nil
	}

	at := last.At
	return Result[RollMeta]{
		Raw:     prev * (0.9 + 0.2*clamp01(last.U)),
		Outcome: OK,
		Patch:   RollMeta{UsedRollAt: &at},
	}, nil
}

func RollFraction(result, lo, hi int) float64 {
	if hi <= lo {
		return 0.5
	}
	return clamp01(float64(result-lo) / float64(hi-lo))
}

func clamp01(f float64) float64 {
	switch {
	case f < 0:
		return 0
	case f > 1:
		return 1
	default:
		return f
	}
}
