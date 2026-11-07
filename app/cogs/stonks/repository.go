package stonks

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/fiffu/arisa3/app/cogs/stonks/tickers"
	"github.com/fiffu/arisa3/app/database"
)

type tickerRow struct {
	Ticker string
	// LastPrice is in thousandths of a coin. It is valid only if HasPrev is true.
	LastPrice int64
	// LastSlot is in Unix seconds. It is valid only if HasPrev is true.
	LastSlot int64
	HasPrev  bool
	Meta     json.RawMessage
}

type IRepository interface {
	EnsureTickers(ctx context.Context, names []string) error
	LoadTickers(ctx context.Context) (map[string]tickerRow, error)
	// AppendPrice does nothing if the ticker already has a price for this slot or a later one.
	AppendPrice(ctx context.Context, ticker string, priceMilli int64, slot int64, patch json.RawMessage) (appended bool, err error)
	SaveRoll(ctx context.Context, roll tickers.LastRoll) error
}

type repo struct {
	db database.IDatabase
}

func NewRepository(db database.IDatabase) IRepository {
	return &repo{db: db}
}

func (r *repo) EnsureTickers(ctx context.Context, names []string) error {
	for _, name := range names {
		_, err := r.db.Exec(ctx, `INSERT INTO stonks_tickers (ticker) VALUES ($1) ON CONFLICT DO NOTHING`, name)
		if err != nil {
			return fmt.Errorf("ensure ticker %s: %w", name, err)
		}
	}
	return nil
}

func (r *repo) LoadTickers(ctx context.Context) (map[string]tickerRow, error) {
	rows, err := r.db.Query(ctx, `
		SELECT ticker, prices[cardinality(prices)], slots[cardinality(slots)], meta
		FROM stonks_tickers`)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	out := map[string]tickerRow{}
	for rows.Next() {
		var (
			row   tickerRow
			price sql.NullInt64
			slot  sql.NullInt64
			meta  []byte
		)
		if err := rows.Scan(&row.Ticker, &price, &slot, &meta); err != nil {
			return nil, err
		}
		row.HasPrev = price.Valid && slot.Valid
		row.LastPrice = price.Int64
		row.LastSlot = slot.Int64
		row.Meta = meta
		out[row.Ticker] = row
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

func (r *repo) AppendPrice(ctx context.Context, ticker string, priceMilli int64, slot int64, patch json.RawMessage) (bool, error) {
	if len(patch) == 0 {
		patch = json.RawMessage("{}")
	}

	// `meta || patch` changes only the keys in patch. It never overwrites a
	// roll that /roll saved after LoadTickers ran.
	res, err := r.db.Exec(ctx, `
		UPDATE stonks_tickers
		SET prices = array_append(prices, $2::bigint),
		    slots  = array_append(slots, $3::bigint),
		    meta   = meta || $4::jsonb
		WHERE ticker = $1
		  AND (cardinality(slots) = 0 OR slots[cardinality(slots)] < $3::bigint)`,
		ticker, priceMilli, slot, string(patch))
	if err != nil {
		return false, err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return false, err
	}
	return n > 0, nil
}

func (r *repo) SaveRoll(ctx context.Context, roll tickers.LastRoll) error {
	patch, err := json.Marshal(tickers.RollMeta{LastRoll: &roll})
	if err != nil {
		return err
	}
	_, err = r.db.Exec(ctx, `UPDATE stonks_tickers SET meta = meta || $2::jsonb WHERE ticker = $1`,
		"ROLL", string(patch))
	return err
}
