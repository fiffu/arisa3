// Package tickers has no database access. The caller stores prices and handles failures.
package tickers

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/fiffu/arisa3/app/cogs/stonks/tickers/api"
)

type Outcome string

const (
	OK        Outcome = "ok"
	NoNewData Outcome = "no_new_data"
)

type State[T any] struct {
	Slot time.Time
	// PrevSlot is zero if there is no previous price.
	PrevSlot time.Time
	// PrevPrice is in coins. It is valid only if HasPrev is true.
	PrevPrice float64
	HasPrev   bool
	Meta      T
}

type Result[T any] struct {
	Raw     float64
	Outcome Outcome
	// Patch is merged into the stored meta, key by key. Tag every meta field
	// omitempty, so a patch writes only the fields the source set.
	Patch T
}

// T is the source's own meta type. It holds only the keys that the source reads.
type source[T any] interface {
	Ticker() string
	Fetch(ctx context.Context, st State[T]) (Result[T], error)
}

// Source carries meta as raw JSON, so sources with different meta types fit in one slice.
type Source interface {
	Ticker() string
	Fetch(ctx context.Context, st State[json.RawMessage]) (Result[json.RawMessage], error)
}

// NoMeta is the meta type of a source that keeps no state.
type NoMeta struct{}

func Wrap[T any](src source[T]) Source {
	return wrapped[T]{src: src}
}

type wrapped[T any] struct {
	src source[T]
}

func (w wrapped[T]) Ticker() string { return w.src.Ticker() }

func (w wrapped[T]) Fetch(ctx context.Context, st State[json.RawMessage]) (Result[json.RawMessage], error) {
	// Unmarshal ignores keys that T does not have, such as the keys of the update loop.
	var meta T
	if len(st.Meta) > 0 {
		if err := json.Unmarshal(st.Meta, &meta); err != nil {
			return Result[json.RawMessage]{}, fmt.Errorf("decode meta: %w", err)
		}
	}

	res, err := w.src.Fetch(ctx, State[T]{
		Slot:      st.Slot,
		PrevSlot:  st.PrevSlot,
		PrevPrice: st.PrevPrice,
		HasPrev:   st.HasPrev,
		Meta:      meta,
	})
	out := Result[json.RawMessage]{Raw: res.Raw, Outcome: res.Outcome}
	if err != nil {
		return out, err
	}

	patch, err := json.Marshal(res.Patch)
	if err != nil {
		return out, fmt.Errorf("encode meta patch: %w", err)
	}
	out.Patch = patch
	return out, nil
}

func Sources(transport http.RoundTripper, rnd *rand.Rand) []Source {
	fred := api.NewFred(transport)
	return []Source{
		Wrap[RndmMeta](&Rndm{Rand: rnd}),
		Wrap[RollMeta](&Roll{}),
		Wrap[FredMeta](NewSP500(fred)),
		Wrap[FredMeta](NewDJIA(fred)),
		Wrap[GreedMeta](&Greed{Client: api.NewGreed(transport)}),
		Wrap[DonMeta](&Don{Client: api.NewTrumpsTruth(transport)}),
	}
}
