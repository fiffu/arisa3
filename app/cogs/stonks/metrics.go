package stonks

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/fiffu/arisa3/app/cogs/stonks/tickers"
	"github.com/fiffu/arisa3/app/instrumentation"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

const outcomeError = "error"

type sourceState struct {
	raw      float64
	hasRaw   bool
	saved    float64 // coins
	hasSaved bool
	lastOK   time.Time
}

type sourceMetrics struct {
	now func() time.Time

	fetches  metric.Int64Counter
	duration metric.Float64Histogram

	mu      sync.Mutex
	sources map[string]*sourceState
}

func newSourceMetrics(meter metric.Meter, now func() time.Time) (*sourceMetrics, error) {
	m := &sourceMetrics{now: now, sources: map[string]*sourceState{}}

	var err error
	if m.fetches, err = meter.Int64Counter("stonks_source_fetches",
		metric.WithDescription("Fetches, by ticker and outcome")); err != nil {
		return nil, err
	}
	if m.duration, err = meter.Float64Histogram("stonks_source_fetch_duration",
		metric.WithDescription("Time for one fetch"), metric.WithUnit("s")); err != nil {
		return nil, err
	}

	raw, err := meter.Float64ObservableGauge("stonks_source_raw_value",
		metric.WithDescription("Last value from the source, before the price floor"))
	if err != nil {
		return nil, err
	}
	saved, err := meter.Float64ObservableGauge("stonks_source_saved_value",
		metric.WithDescription("Last stored price, in coins"))
	if err != nil {
		return nil, err
	}
	age, err := meter.Float64ObservableGauge("stonks_source_age",
		metric.WithDescription("Time since the last ok fetch"), metric.WithUnit("s"))
	if err != nil {
		return nil, err
	}

	_, err = meter.RegisterCallback(func(_ context.Context, o metric.Observer) error {
		m.mu.Lock()
		defer m.mu.Unlock()
		now := m.now()
		for ticker, f := range m.sources {
			attrs := metric.WithAttributes(attribute.String("ticker", ticker))
			if f.hasRaw {
				o.ObserveFloat64(raw, f.raw, attrs)
			}
			if f.hasSaved {
				o.ObserveFloat64(saved, f.saved, attrs)
			}
			if !f.lastOK.IsZero() {
				o.ObserveFloat64(age, now.Sub(f.lastOK).Seconds(), attrs)
			}
		}
		return nil
	}, raw, saved, age)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (m *sourceMetrics) source(ticker string) *sourceState {
	f, ok := m.sources[ticker]
	if !ok {
		f = &sourceState{}
		m.sources[ticker] = f
	}
	return f
}

// Seed loads the stored state after a restart, so the gauges have values before the next fetch.
func (m *sourceMetrics) Seed(ticker string, savedMilli int64, hasSaved bool, lastOK time.Time) {
	m.mu.Lock()
	defer m.mu.Unlock()
	f := m.source(ticker)
	if hasSaved {
		f.saved, f.hasSaved = milliToCoins(savedMilli), true
	}
	f.lastOK = lastOK
}

func (m *sourceMetrics) Saved(ticker string, priceMilli int64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	f := m.source(ticker)
	f.saved, f.hasSaved = milliToCoins(priceMilli), true
}

func (m *sourceMetrics) Fetch(ctx context.Context, source tickers.Source, st tickers.State[json.RawMessage]) (tickers.Result[json.RawMessage], error) {
	ticker := source.Ticker()
	ctx, span := instrumentation.SpanInContext(ctx, instrumentation.Internal("stonks_source_fetch"))
	defer span.End()

	start := m.now()
	res, err := source.Fetch(ctx, st)
	elapsed := m.now().Sub(start)

	outcome := string(res.Outcome)
	if err != nil {
		outcome = outcomeError
		span.RecordError(err)
	}

	span.SetAttributes(
		attribute.String("ticker", ticker),
		attribute.String("outcome", outcome),
		attribute.Float64("raw_value", res.Raw),
	)
	m.fetches.Add(ctx, 1, metric.WithAttributes(attribute.String("ticker", ticker), attribute.String("outcome", outcome)))
	m.duration.Record(ctx, elapsed.Seconds(), metric.WithAttributes(attribute.String("ticker", ticker)))

	if err == nil {
		m.mu.Lock()
		f := m.source(ticker)
		if res.Outcome == tickers.OK || st.HasPrev {
			f.raw, f.hasRaw = res.Raw, true
		}
		if res.Outcome == tickers.OK {
			f.lastOK = m.now()
		}
		m.mu.Unlock()
	}
	return res, err
}

func milliToCoins(milli int64) float64 {
	return float64(milli) / 1000
}
