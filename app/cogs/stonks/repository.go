package stonks

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
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

const (
	// StartingBalance is in thousandths of a coin. A wallet gets it when it is created.
	StartingBalance int64 = 50_000
	// CoinPerReaction is in thousandths of a coin.
	CoinPerReaction int64 = 1_000
)

type ReactionSource string

const (
	SourceLive     ReactionSource = "live"
	SourceBackfill ReactionSource = "backfill"
)

// React holds Discord snowflake IDs as decimal strings.
type React struct {
	GuildID     string
	ChannelID   string
	MessageID   string
	ReactorID   string
	RecipientID string
	Source      ReactionSource
	// EmojiID is not stored. It is "" for a built-in emoji.
	EmojiID string
}

// ErrUnknownUser means the recipient is not in stonks_users yet. Save it with SaveUser, then retry.
var ErrUnknownUser = errors.New("stonks: user not in stonks_users")

type IWalletRepository interface {
	// GuildEmojis maps each guild ID to the ID of its coin emoji.
	GuildEmojis(ctx context.Context) (map[string]string, error)
	SaveGuildSettings(ctx context.Context, guildID, emojiID, marketChannelID string) error
	SaveUser(ctx context.Context, userID string, isBot bool) error
	// AddCoinReaction grants the recipient one coin, unless the pair is already active
	// or the recipient is a bot. It returns ErrUnknownUser and changes nothing
	// if the recipient is not in stonks_users.
	AddCoinReaction(ctx context.Context, r React) (granted bool, err error)
	// RemoveCoinReaction takes the coin back, unless the pair is not active.
	RemoveCoinReaction(ctx context.Context, messageID, reactorID string) (revoked bool, err error)
}

func NewWalletRepository(db database.IDatabase) IWalletRepository {
	return &repo{db: db}
}

func (r *repo) GuildEmojis(ctx context.Context) (map[string]string, error) {
	rows, err := r.db.Query(ctx, `
		SELECT guild_id, data->>'emoji_id' FROM stonks_guilds
		WHERE data->>'emoji_id' IS NOT NULL`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := map[string]string{}
	for rows.Next() {
		var guildID, emojiID string
		if err := rows.Scan(&guildID, &emojiID); err != nil {
			return nil, err
		}
		out[guildID] = emojiID
	}
	return out, rows.Err()
}

func (r *repo) SaveGuildSettings(ctx context.Context, guildID, emojiID, marketChannelID string) error {
	// || replaces only these two keys, so it keeps backfill progress written at the same time.
	_, err := r.db.Exec(ctx, `
		INSERT INTO stonks_guilds (guild_id, data)
		VALUES ($1, jsonb_build_object('emoji_id', $2::text, 'market_channel_id', $3::text))
		ON CONFLICT (guild_id) DO UPDATE
		SET data = stonks_guilds.data || EXCLUDED.data, updated_at = now()`,
		guildID, emojiID, marketChannelID)
	return err
}

func (r *repo) SaveUser(ctx context.Context, userID string, isBot bool) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO stonks_users (user_id, is_bot) VALUES ($1, $2)
		ON CONFLICT (user_id) DO UPDATE SET is_bot = EXCLUDED.is_bot, updated_at = now()`,
		userID, isBot)
	return err
}

func (r *repo) AddCoinReaction(ctx context.Context, cr React) (bool, error) {
	return r.inTx(ctx, func(ctx context.Context, tx database.ITransaction) (bool, error) {
		rows, err := tx.Query(ctx, `SELECT is_bot FROM stonks_users WHERE user_id = $1`, cr.RecipientID)
		if err != nil {
			return false, err
		}
		var isBot bool
		found := rows.Next()
		if found {
			err = rows.Scan(&isBot)
		}
		if err = errors.Join(err, rows.Err(), rows.Close()); err != nil {
			return false, err
		}
		if !found {
			return false, ErrUnknownUser
		}
		if isBot {
			return false, nil
		}

		if _, err := tx.Exec(ctx, `
			INSERT INTO stonks_wallets (guild_id, user_id, ticker, balance)
			VALUES ($1, $2, 'coin', $3)
			ON CONFLICT DO NOTHING`,
			cr.GuildID, cr.RecipientID, StartingBalance); err != nil {
			return false, fmt.Errorf("create wallet: %w", err)
		}

		// The WHERE on the conflict branch makes a react on an active pair change nothing,
		// so RETURNING gives a row only when this call activated the pair.
		changed, err := execReturnsRow(ctx, tx, `
			INSERT INTO stonks_coin_reactions
				(guild_id, channel_id, message_id, reactor_id, recipient_id, active, source)
			VALUES ($1, $2, $3, $4, $5, true, $6)
			ON CONFLICT (message_id, reactor_id) DO UPDATE
			SET active = true, updated_at = now()
			WHERE stonks_coin_reactions.active = false
			RETURNING 1`,
			cr.GuildID, cr.ChannelID, cr.MessageID, cr.ReactorID, cr.RecipientID, string(cr.Source))
		if err != nil || !changed {
			return false, err
		}

		if _, err := tx.Exec(ctx, `
			UPDATE stonks_wallets SET balance = balance + $3, updated_at = now()
			WHERE guild_id = $1 AND user_id = $2 AND ticker = 'coin'`,
			cr.GuildID, cr.RecipientID, CoinPerReaction); err != nil {
			return false, fmt.Errorf("credit coin: %w", err)
		}
		return true, nil
	})
}

func (r *repo) RemoveCoinReaction(ctx context.Context, messageID, reactorID string) (bool, error) {
	return r.inTx(ctx, func(ctx context.Context, tx database.ITransaction) (bool, error) {
		rows, err := tx.Query(ctx, `
			UPDATE stonks_coin_reactions SET active = false, updated_at = now()
			WHERE message_id = $1 AND reactor_id = $2 AND active
			RETURNING guild_id, recipient_id`,
			messageID, reactorID)
		if err != nil {
			return false, err
		}
		var guildID, recipientID string
		found := rows.Next()
		if found {
			err = rows.Scan(&guildID, &recipientID)
		}
		// The rows must be closed before the next statement in the transaction.
		err = errors.Join(err, rows.Err(), rows.Close())
		if err != nil || !found {
			return false, err
		}

		// The balance can go negative, for example after the recipient spent the coin.
		if _, err := tx.Exec(ctx, `
			UPDATE stonks_wallets SET balance = balance - $3, updated_at = now()
			WHERE guild_id = $1 AND user_id = $2 AND ticker = 'coin'`,
			guildID, recipientID, CoinPerReaction); err != nil {
			return false, fmt.Errorf("debit coin: %w", err)
		}
		return true, nil
	})
}

// inTx commits if fn returns no error, and rolls back otherwise.
func (r *repo) inTx(ctx context.Context, fn func(context.Context, database.ITransaction) (bool, error)) (bool, error) {
	ctx, tx, err := r.db.Begin(ctx)
	if err != nil {
		return false, err
	}
	ok, err := fn(ctx, tx)
	if err != nil {
		return false, errors.Join(err, tx.Rollback(ctx))
	}
	if err := tx.Commit(ctx); err != nil {
		return false, err
	}
	return ok, nil
}

func execReturnsRow(ctx context.Context, tx database.ITransaction, query string, args ...any) (bool, error) {
	rows, err := tx.Query(ctx, query, args...)
	if err != nil {
		return false, err
	}
	found := rows.Next()
	return found, errors.Join(rows.Err(), rows.Close())
}
