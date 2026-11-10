package stonks

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/fiffu/arisa3/app/cogs/stonks/tickers"
	"github.com/fiffu/arisa3/app/log"
)

const SlotSize = tickers.DefaultSlotSize

// tickMeta holds the meta keys that the update loop owns, not any source.
type tickMeta struct {
	// LastOKAt lets the age gauge survive a restart.
	LastOKAt *time.Time `json:"last_ok_at,omitempty"`
}

// PriceFloor is the lowest price, in thousandths of a coin. A price of 0 would make buys free.
const PriceFloor int64 = 1

// toMilli applies PriceFloor.
func toMilli(raw float64) int64 {
	if math.IsNaN(raw) || math.IsInf(raw, 0) {
		return PriceFloor
	}
	p := int64(math.Round(raw * 1000))
	if p < PriceFloor {
		return PriceFloor
	}
	return p
}

type ticker struct {
	repo    IRepository
	sources []tickers.Source
	metrics *sourceMetrics
	now     func() time.Time
}

// Each write in Tick is idempotent, so a rerun is safe.
func (t *ticker) Tick(ctx context.Context, slot time.Time) error {
	rows, err := t.repo.LoadTickers(ctx)
	if err != nil {
		return fmt.Errorf("load tickers: %w", err)
	}

	// Each source writes only its own row, so all sources run at once.
	// errs is indexed like t.sources, so errors keep a stable order.
	errs := make([]error, len(t.sources))
	var wg sync.WaitGroup
	for i, source := range t.sources {
		// EnsureTickers() creates every row at startup, so a missing row means a bad read.
		// Treating it as a new ticker would reset its price.
		row, ok := rows[source.Ticker()]
		if !ok {
			errs[i] = fmt.Errorf("%s: no row in stonks_tickers", source.Ticker())
			continue
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := t.tickOne(ctx, source, row, slot); err != nil {
				errs[i] = fmt.Errorf("%s: %w", source.Ticker(), err)
			}
		}()
	}
	wg.Wait()
	return errors.Join(errs...)
}

func (t *ticker) tickOne(ctx context.Context, source tickers.Source, row tickerRow, slot time.Time) error {
	name := source.Ticker()
	if row.HasPrev && row.LastSlot >= slot.Unix() {
		return nil // already priced
	}

	st := tickers.State[json.RawMessage]{Slot: slot, Meta: row.Meta, HasPrev: row.HasPrev}
	if row.HasPrev {
		st.PrevSlot = time.Unix(row.LastSlot, 0).UTC()
		st.PrevPrice = milliToCoins(row.LastPrice)
	}

	res, fetchErr := t.metrics.Fetch(ctx, source, st)

	var (
		price int64
		patch json.RawMessage
	)
	switch {
	case fetchErr != nil:
		// A fetch stopped by shutdown did not fail. Storing the previous price would mark
		// the slot as priced, and the next start would not fetch it again.
		if !row.HasPrev || ctx.Err() != nil {
			return fetchErr
		}
		// Keep the previous price, so every ticker has a value for every slot.
		price = row.LastPrice
	case res.Outcome == tickers.NoNewData:
		if !row.HasPrev {
			return nil
		}
		price, patch = row.LastPrice, res.Patch
	default:
		now := t.now().UTC()
		own, err := json.Marshal(tickMeta{LastOKAt: &now})
		if err != nil {
			return err
		}
		if patch, err = mergePatches(res.Patch, own); err != nil {
			return fmt.Errorf("merge meta patches: %w", err)
		}
		price = toMilli(res.Raw)
	}

	appended, err := t.repo.AppendPrice(ctx, name, price, slot.Unix(), patch)
	if err != nil {
		return errors.Join(fetchErr, fmt.Errorf("append price: %w", err))
	}
	if appended {
		t.metrics.Saved(name, price)
		log.Infof(ctx, "stonks: %s = %d milli-coins at slot %d", name, price, slot.Unix())
	}
	return fetchErr
}

func (t *ticker) seedMetrics(ctx context.Context) error {
	rows, err := t.repo.LoadTickers(ctx)
	if err != nil {
		return err
	}
	for name, row := range rows {
		var meta tickMeta
		var lastOK time.Time
		if len(row.Meta) > 0 && json.Unmarshal(row.Meta, &meta) == nil && meta.LastOKAt != nil {
			lastOK = *meta.LastOKAt
		}
		t.metrics.Seed(name, row.LastPrice, row.HasPrev, lastOK)
	}
	return nil
}

// Later patches win when two patches set the same key.
func mergePatches(patches ...json.RawMessage) (json.RawMessage, error) {
	merged := map[string]json.RawMessage{}
	for _, p := range patches {
		if len(p) == 0 {
			continue
		}
		var m map[string]json.RawMessage
		if err := json.Unmarshal(p, &m); err != nil {
			return nil, err
		}
		for k, v := range m {
			merged[k] = v
		}
	}
	return json.Marshal(merged)
}
