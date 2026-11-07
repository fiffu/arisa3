package stonks

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/fiffu/arisa3/app/cogs/stonks/tickers"
	"github.com/fiffu/arisa3/app/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test_repo_postgres runs the SQL against a real Postgres.
// It runs only if STONKS_TEST_DSN is set, for example:
//
//	docker run --rm -d --name stonks-pg -e POSTGRES_PASSWORD=pw -p 127.0.0.1:55432:5432 postgres:17-alpine
//	STONKS_TEST_DSN='postgres://postgres:pw@127.0.0.1:55432/postgres?sslmode=disable' go test ./app/cogs/stonks/
func Test_repo_postgres(t *testing.T) {
	dsn := os.Getenv("STONKS_TEST_DSN")
	if dsn == "" {
		t.Skip("STONKS_TEST_DSN not set")
	}
	ctx := context.Background()

	// Reset first. The client caches applied migrations when it opens.
	reset, err := database.NewDBClient(ctx, dsn)
	require.NoError(t, err)
	_, err = reset.Exec(ctx, `DROP TABLE IF EXISTS stonks_tickers; DELETE FROM _schema_migrations`)
	require.NoError(t, err)
	require.NoError(t, reset.Close(ctx))

	db, err := database.NewDBClient(ctx, dsn)
	require.NoError(t, err)
	defer db.Close(ctx)
	files, err := filepath.Glob(filepath.Join(migrationsDir, "*.sql"))
	require.NoError(t, err)
	for _, f := range files {
		schema, err := db.ParseMigration(ctx, f)
		require.NoError(t, err)
		_, err = db.Migrate(ctx, schema)
		require.NoError(t, err)
	}

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
