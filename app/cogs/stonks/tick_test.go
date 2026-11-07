package stonks

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/fiffu/arisa3/app/cogs/stonks/tickers"
	"github.com/fiffu/arisa3/app/cogs/stonks/tickers/api"
	"github.com/fiffu/arisa3/testfixtures"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/metric/noop"
)

// fakeRepo keeps stonks_tickers in memory, with the same append and merge rules as the SQL.
type fakeRepo struct {
	mu     sync.Mutex
	prices map[string][]int64
	slots  map[string][]int64
	meta   map[string]map[string]json.RawMessage
}

func newFakeRepo() *fakeRepo {
	return &fakeRepo{
		prices: map[string][]int64{},
		slots:  map[string][]int64{},
		meta:   map[string]map[string]json.RawMessage{},
	}
}

func (f *fakeRepo) EnsureTickers(_ context.Context, names []string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	for _, n := range names {
		if _, ok := f.meta[n]; !ok {
			f.meta[n] = map[string]json.RawMessage{}
		}
	}
	return nil
}

func (f *fakeRepo) LoadTickers(context.Context) (map[string]tickerRow, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	out := map[string]tickerRow{}
	for n, m := range f.meta {
		meta, err := json.Marshal(m)
		if err != nil {
			return nil, err
		}
		row := tickerRow{Ticker: n, Meta: meta}
		if k := len(f.slots[n]); k > 0 {
			row.HasPrev, row.LastPrice, row.LastSlot = true, f.prices[n][k-1], f.slots[n][k-1]
		}
		out[n] = row
	}
	return out, nil
}

func (f *fakeRepo) AppendPrice(_ context.Context, n string, price, slot int64, patch json.RawMessage) (bool, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if k := len(f.slots[n]); k > 0 && f.slots[n][k-1] >= slot {
		return false, nil
	}
	f.prices[n] = append(f.prices[n], price)
	f.slots[n] = append(f.slots[n], slot)
	return true, f.merge(n, patch)
}

func (f *fakeRepo) SaveRoll(_ context.Context, roll tickers.LastRoll) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	patch, err := json.Marshal(tickers.RollMeta{LastRoll: &roll})
	if err != nil {
		return err
	}
	return f.merge("ROLL", patch)
}

// merge needs f.mu held.
func (f *fakeRepo) merge(n string, patch json.RawMessage) error {
	if len(patch) == 0 {
		return nil
	}
	var m map[string]json.RawMessage
	if err := json.Unmarshal(patch, &m); err != nil {
		return err
	}
	for k, v := range m {
		f.meta[n][k] = v
	}
	return nil
}

// stubSource returns fixed results in order.
type stubSource struct {
	name    string
	results []tickers.Result[json.RawMessage]
	errs    []error
	calls   int
}

func (s *stubSource) Ticker() string { return s.name }
func (s *stubSource) Fetch(context.Context, tickers.State[json.RawMessage]) (tickers.Result[json.RawMessage], error) {
	i := s.calls
	s.calls++
	var err error
	if i < len(s.errs) {
		err = s.errs[i]
	}
	return s.results[i], err
}

func newTestTicker(t *testing.T, repo *fakeRepo, srcs ...tickers.Source) *ticker {
	m, err := newSourceMetrics(noop.NewMeterProvider().Meter("test"), time.Now)
	require.NoError(t, err)
	names := []string{}
	for _, s := range srcs {
		names = append(names, s.Ticker())
	}
	require.NoError(t, repo.EnsureTickers(context.Background(), names))
	return &ticker{repo: repo, sources: srcs, metrics: m, now: time.Now}
}

var (
	slot0 = time.Unix(1794052800, 0).UTC() // 2026-11-07 12:00 UTC
	slot1 = slot0.Add(SlotSize)
	slot2 = slot1.Add(SlotSize)
)

func Test_Tick_appendsAndIsIdempotent(t *testing.T) {
	repo := newFakeRepo()
	src := &stubSource{name: "AAA", results: []tickers.Result[json.RawMessage]{
		{Raw: 12.3456, Outcome: tickers.OK, Patch: json.RawMessage(`{"source_ts":"1"}`)},
	}}
	tk := newTestTicker(t, repo, src)

	require.NoError(t, tk.Tick(context.Background(), slot0))
	require.NoError(t, tk.Tick(context.Background(), slot0)) // rerun: no fetch, no append

	assert.Equal(t, 1, src.calls)
	assert.Equal(t, []int64{12346}, repo.prices["AAA"])
	assert.Equal(t, []int64{slot0.Unix()}, repo.slots["AAA"])
	assert.JSONEq(t, `"1"`, string(repo.meta["AAA"]["source_ts"]))
	assert.NotEmpty(t, repo.meta["AAA"]["last_ok_at"])
}

func Test_Tick_noNewDataAndErrorsKeepPreviousPrice(t *testing.T) {
	repo := newFakeRepo()
	src := &stubSource{
		name: "AAA",
		results: []tickers.Result[json.RawMessage]{
			{Raw: 5, Outcome: tickers.OK},
			{Raw: 999, Outcome: tickers.NoNewData},
			{},
		},
		errs: []error{nil, nil, errors.New("feed down")},
	}
	tk := newTestTicker(t, repo, src)

	require.NoError(t, tk.Tick(context.Background(), slot0))
	require.NoError(t, tk.Tick(context.Background(), slot1))
	err := tk.Tick(context.Background(), slot2)

	assert.ErrorContains(t, err, "feed down")
	assert.Equal(t, []int64{5000, 5000, 5000}, repo.prices["AAA"])
	assert.Equal(t, []int64{slot0.Unix(), slot1.Unix(), slot2.Unix()}, repo.slots["AAA"])
}

func Test_Tick_noHistorySkipsUntilFirstValue(t *testing.T) {
	repo := newFakeRepo()
	src := &stubSource{
		name:    "AAA",
		results: []tickers.Result[json.RawMessage]{{}, {Outcome: tickers.NoNewData}, {Raw: 7, Outcome: tickers.OK}},
		errs:    []error{errors.New("feed down")},
	}
	tk := newTestTicker(t, repo, src)

	assert.Error(t, tk.Tick(context.Background(), slot0))
	assert.NoError(t, tk.Tick(context.Background(), slot1))
	assert.NoError(t, tk.Tick(context.Background(), slot2))

	assert.Equal(t, []int64{7000}, repo.prices["AAA"])
	assert.Equal(t, []int64{slot2.Unix()}, repo.slots["AAA"])
}

func Test_Tick_oneFailingFeedDoesNotBlockOthers(t *testing.T) {
	repo := newFakeRepo()
	bad := &stubSource{name: "BAD", results: []tickers.Result[json.RawMessage]{{}}, errs: []error{errors.New("down")}}
	good := &stubSource{name: "GOOD", results: []tickers.Result[json.RawMessage]{{Raw: 1, Outcome: tickers.OK}}}
	tk := newTestTicker(t, repo, bad, good)

	err := tk.Tick(context.Background(), slot0)

	assert.ErrorContains(t, err, "BAD: down")
	assert.Equal(t, []int64{1000}, repo.prices["GOOD"])
}

func Test_toMilli_appliesFloor(t *testing.T) {
	assert.Equal(t, int64(6000120), toMilli(6000.12))
	assert.Equal(t, PriceFloor, toMilli(0))
	assert.Equal(t, PriceFloor, toMilli(-3))
	assert.Equal(t, PriceFloor, toMilli(0.0001))
}

func Test_Tick_rollUsesEachRollOnce(t *testing.T) {
	repo := newFakeRepo()
	tk := newTestTicker(t, repo, tickers.Wrap[tickers.RollMeta](&tickers.Roll{}))
	ctx := context.Background()

	require.NoError(t, tk.Tick(ctx, slot0)) // starts at 50
	require.NoError(t, repo.SaveRoll(ctx, tickers.LastRoll{U: 1, At: slot0.Add(time.Hour)}))
	require.NoError(t, tk.Tick(ctx, slot1)) // +10%
	require.NoError(t, tk.Tick(ctx, slot2)) // same roll: no change

	assert.Equal(t, []int64{50000, 55000, 55000}, repo.prices["ROLL"])
	assert.NotEmpty(t, repo.meta["ROLL"]["last_roll"], "the update must not erase the saved roll")
}

func Test_mergePatches(t *testing.T) {
	got, err := mergePatches(json.RawMessage(`{"a":1,"b":2}`), nil, json.RawMessage(`{"b":3,"c":4}`))
	require.NoError(t, err)
	assert.JSONEq(t, `{"a":1,"b":3,"c":4}`, string(got))
}

// barrierSource returns only after every barrierSource in the test has started its fetch.
type barrierSource struct {
	name    string
	started *sync.WaitGroup
}

func (s *barrierSource) Ticker() string { return s.name }
func (s *barrierSource) Fetch(context.Context, tickers.State[json.RawMessage]) (tickers.Result[json.RawMessage], error) {
	s.started.Done()
	all := make(chan struct{})
	go func() { s.started.Wait(); close(all) }()
	select {
	case <-all:
		return tickers.Result[json.RawMessage]{Raw: 1, Outcome: tickers.OK}, nil
	case <-time.After(2 * time.Second):
		return tickers.Result[json.RawMessage]{}, errors.New("fetches ran one at a time")
	}
}

func Test_Tick_fetchesSourcesConcurrently(t *testing.T) {
	repo := newFakeRepo()
	var started sync.WaitGroup
	started.Add(3)
	tk := newTestTicker(t, repo,
		&barrierSource{name: "A", started: &started},
		&barrierSource{name: "B", started: &started},
		&barrierSource{name: "C", started: &started},
	)

	require.NoError(t, tk.Tick(context.Background(), slot0))

	for _, n := range []string{"A", "B", "C"} {
		assert.Equal(t, []int64{1000}, repo.prices[n], n)
	}
}

func Test_Tick_missingRowIsAnErrorNotANewTicker(t *testing.T) {
	repo := newFakeRepo()
	tk := newTestTicker(t, repo, tickers.Wrap[tickers.RollMeta](&tickers.Roll{}))
	ctx := context.Background()
	require.NoError(t, tk.Tick(ctx, slot0))

	// A partial read returns no row for ROLL.
	saved := repo.meta["ROLL"]
	delete(repo.meta, "ROLL")
	err := tk.Tick(ctx, slot1)
	repo.meta["ROLL"] = saved

	assert.ErrorContains(t, err, "ROLL: no row in stonks_tickers")
	assert.Equal(t, []int64{50000}, repo.prices["ROLL"], "the price must not reset to the start price")
}

// flakyTransport fails while fail is true, and serves a fixture otherwise.
type flakyTransport struct {
	fail    bool
	fixture *testfixtures.Transport
}

func (f *flakyTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if f.fail {
		return nil, errors.New("connection refused")
	}
	return f.fixture.RoundTrip(req)
}

func Test_Tick_donCountsPostsFromAFailedWindow(t *testing.T) {
	repo := newFakeRepo()
	tr := &flakyTransport{fixture: &testfixtures.Transport{Body: testfixtures.StonksDonFeed}}
	don := tickers.Wrap[tickers.DonMeta](&tickers.Don{Client: api.NewTrumpsTruth(tr)})
	tk := newTestTicker(t, repo, don)
	ctx := context.Background()
	at := func(day, hour int) time.Time { return time.Date(2026, 11, day, hour, 0, 0, 0, time.UTC) }

	// Window Nov 6 00:00–12:00: 3 posts, 300 chars.
	require.NoError(t, tk.Tick(ctx, at(6, 12)))

	tr.fail = true
	assert.Error(t, tk.Tick(ctx, at(7, 0)))

	// Window Nov 6 12:00–Nov 7 12:00 covers the failed window too: 5 posts, 2272 chars.
	// Without CoveredUntil, it would start at Nov 7 00:00: 1 post, 100 chars.
	tr.fail = false
	require.NoError(t, tk.Tick(ctx, at(7, 12)))

	assert.Equal(t, []int64{100000, 100000, 454400}, repo.prices["DON"])
}
