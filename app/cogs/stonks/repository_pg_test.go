package stonks

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/fiffu/arisa3/app/cogs/stonks/tickers"
	"github.com/fiffu/arisa3/app/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// openTestDB returns a migrated, empty database for the *_postgres tests.
// It skips the test unless STONKS_TEST_DSN is set, for example:
//
//	docker run --rm -d --name stonks-pg -e POSTGRES_PASSWORD=pw -p 127.0.0.1:55432:5432 postgres:17-alpine
//	STONKS_TEST_DSN='postgres://postgres:pw@127.0.0.1:55432/postgres?sslmode=disable' go test ./app/cogs/stonks/
func openTestDB(t *testing.T) database.IDatabase {
	dsn := os.Getenv("STONKS_TEST_DSN")
	if dsn == "" {
		t.Skip("STONKS_TEST_DSN not set")
	}
	ctx := context.Background()

	// Reset first. The client caches applied migrations when it opens.
	reset, err := database.NewDBClient(ctx, dsn)
	require.NoError(t, err)
	_, err = reset.Exec(ctx, `
		DROP TABLE IF EXISTS stonks_tickers, stonks_guilds, stonks_wallets, stonks_coin_reactions, stonks_users;
		DELETE FROM _schema_migrations`)
	require.NoError(t, err)
	require.NoError(t, reset.Close(ctx))

	db, err := database.NewDBClient(ctx, dsn)
	require.NoError(t, err)
	t.Cleanup(func() { db.Close(ctx) })
	files, err := filepath.Glob(filepath.Join(migrationsDir, "*.sql"))
	require.NoError(t, err)
	for _, f := range files {
		schema, err := db.ParseMigration(ctx, f)
		require.NoError(t, err)
		_, err = db.Migrate(ctx, schema)
		require.NoError(t, err)
	}
	return db
}

func Test_repo_postgres(t *testing.T) {
	ctx := context.Background()
	db := openTestDB(t)

	r := NewRepository(db)
	require.NoError(t, r.EnsureTickers(ctx, []string{"ROLL", "AAA"}))
	require.NoError(t, r.EnsureTickers(ctx, []string{"ROLL"})) // idempotent

	rows, err := r.LoadTickers(ctx)
	require.NoError(t, err)
	assert.False(t, rows["AAA"].HasPrev)
	assert.JSONEq(t, `{}`, string(rows["AAA"].Meta))

	ok, err := r.AppendPrice(ctx, "AAA", 1500, 100, json.RawMessage(`{"source_ts":"x"}`))
	require.NoError(t, err)
	assert.True(t, ok)
	ok, err = r.AppendPrice(ctx, "AAA", 9999, 100, nil) // same slot
	require.NoError(t, err)
	assert.False(t, ok)
	ok, err = r.AppendPrice(ctx, "AAA", 2500, 200, nil)
	require.NoError(t, err)
	assert.True(t, ok)

	// A roll saved between LoadTickers and AppendPrice survives the meta merge.
	at := time.Date(2026, 11, 7, 9, 0, 0, 0, time.UTC)
	require.NoError(t, r.SaveRoll(ctx, tickers.LastRoll{U: 0.25, At: at}))
	_, err = r.AppendPrice(ctx, "ROLL", 100000, 100, json.RawMessage(`{"used_roll_at":"2026-11-07T05:00:00Z"}`))
	require.NoError(t, err)

	rows, err = r.LoadTickers(ctx)
	require.NoError(t, err)
	aaa := rows["AAA"]
	assert.True(t, aaa.HasPrev)
	assert.Equal(t, int64(2500), aaa.LastPrice)
	assert.Equal(t, int64(200), aaa.LastSlot)
	assert.JSONEq(t, `{"source_ts":"x"}`, string(aaa.Meta))

	var roll tickers.RollMeta
	require.NoError(t, json.Unmarshal(rows["ROLL"].Meta, &roll))
	require.NotNil(t, roll.LastRoll)
	assert.Equal(t, 0.25, roll.LastRoll.U)
	assert.True(t, at.Equal(roll.LastRoll.At))
	require.NotNil(t, roll.UsedRollAt)
	assert.True(t, time.Date(2026, 11, 7, 5, 0, 0, 0, time.UTC).Equal(*roll.UsedRollAt))
}

func Test_walletRepo_postgres(t *testing.T) {
	ctx := context.Background()
	db := openTestDB(t)
	r := NewWalletRepository(db)

	balance := func(user string) int64 {
		rows, err := db.Query(ctx, `SELECT balance FROM stonks_wallets WHERE guild_id = $1 AND user_id = $2 AND ticker = 'coin'`, guild, user)
		require.NoError(t, err)
		defer rows.Close()
		require.True(t, rows.Next(), "no coin row for %s", user)
		var b int64
		require.NoError(t, rows.Scan(&b))
		return b
	}

	// Settings: a second save replaces both keys and keeps other keys.
	require.NoError(t, r.SaveGuildSettings(ctx, guild, "1", channel))
	_, err := db.Exec(ctx, `UPDATE stonks_guilds SET data = jsonb_set(data, '{backfill}', '{"200":{"done":true}}') WHERE guild_id = $1`, guild)
	require.NoError(t, err)
	require.NoError(t, r.SaveGuildSettings(ctx, guild, coin, channel))
	emojis, err := r.GuildEmojis(ctx)
	require.NoError(t, err)
	assert.Equal(t, map[string]string{guild: coin}, emojis)
	rows, err := db.Query(ctx, `SELECT data->'backfill'->'200'->>'done' FROM stonks_guilds WHERE guild_id = $1`, guild)
	require.NoError(t, err)
	require.True(t, rows.Next())
	var done string
	require.NoError(t, rows.Scan(&done))
	require.NoError(t, rows.Close())
	assert.Equal(t, "true", done)

	cr := React{GuildID: guild, ChannelID: channel, MessageID: "600", ReactorID: reactor, RecipientID: author, Source: SourceLive}

	// An unknown author changes nothing, not even the wallet.
	_, err = r.AddCoinReaction(ctx, cr)
	assert.ErrorIs(t, err, ErrUnknownUser)
	count := func(query string, args ...any) int {
		rows, err := db.Query(ctx, query, args...)
		require.NoError(t, err)
		defer rows.Close()
		require.True(t, rows.Next())
		var n int
		require.NoError(t, rows.Scan(&n))
		return n
	}
	assert.Equal(t, 0, count(`SELECT count(*) FROM stonks_wallets`))

	// A bot author gets no wallet and no reaction row.
	require.NoError(t, r.SaveUser(ctx, "401", true))
	bot := cr
	bot.MessageID, bot.RecipientID = "602", "401"
	granted, err := r.AddCoinReaction(ctx, bot)
	require.NoError(t, err)
	assert.False(t, granted)
	assert.Equal(t, 0, count(`SELECT count(*) FROM stonks_wallets`))
	assert.Equal(t, 0, count(`SELECT count(*) FROM stonks_coin_reactions`))

	require.NoError(t, r.SaveUser(ctx, author, false))
	require.NoError(t, r.SaveUser(ctx, author, false)) // idempotent
	granted, err = r.AddCoinReaction(ctx, cr)
	require.NoError(t, err)
	assert.True(t, granted)
	granted, err = r.AddCoinReaction(ctx, cr)
	require.NoError(t, err)
	assert.False(t, granted, "an active pair grants no second coin")
	assert.Equal(t, StartingBalance+CoinPerReaction, balance(author))

	revoked, err := r.RemoveCoinReaction(ctx, cr.MessageID, cr.ReactorID)
	require.NoError(t, err)
	assert.True(t, revoked)
	revoked, err = r.RemoveCoinReaction(ctx, cr.MessageID, cr.ReactorID)
	require.NoError(t, err)
	assert.False(t, revoked)
	revoked, err = r.RemoveCoinReaction(ctx, "601", reactor)
	require.NoError(t, err)
	assert.False(t, revoked, "an unknown pair changes nothing")
	assert.Equal(t, StartingBalance, balance(author))

	// A react after an un-react activates the row again.
	cr.Source = SourceBackfill
	granted, err = r.AddCoinReaction(ctx, cr)
	require.NoError(t, err)
	assert.True(t, granted)
	assert.Equal(t, StartingBalance+CoinPerReaction, balance(author))

	// Two un-reacts at once take back one coin.
	var wg sync.WaitGroup
	results := make([]bool, 2)
	for i := range results {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ok, err := r.RemoveCoinReaction(ctx, cr.MessageID, cr.ReactorID)
			assert.NoError(t, err)
			results[i] = ok
		}()
	}
	wg.Wait()
	assert.ElementsMatch(t, []bool{true, false}, results)
	assert.Equal(t, StartingBalance, balance(author))
}
