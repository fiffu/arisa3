package colours

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/fiffu/arisa3/app/database"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newSQLiteRepo opens a real SQLite file and applies this cog's migrations to it.
func newSQLiteRepo(t *testing.T) (*repo, database.IDatabase) {
	ctx := context.Background()
	db, err := database.NewDBClient(ctx, filepath.Join(t.TempDir(), "test.db"))
	require.NoError(t, err)
	t.Cleanup(func() { db.Close(ctx) })

	files, err := os.ReadDir(migrationsDir)
	require.NoError(t, err)
	for _, f := range files {
		schema, err := db.ParseMigration(ctx, filepath.Join(migrationsDir, f.Name()))
		require.NoError(t, err)
		_, err = db.Migrate(ctx, schema)
		require.NoError(t, err, f.Name())
	}
	return newRepo(db), db
}

func Test_sqlite_updateThenFetch(t *testing.T) {
	ctx := context.Background()
	mem := newTestMember(gomock.NewController(t))
	r, _ := newSQLiteRepo(t)

	colour := &Colour{}
	require.NoError(t, r.UpdateMutate(ctx, mem, colour))
	require.NoError(t, r.UpdateReroll(ctx, mem, colour))
	require.NoError(t, r.UpdateFreeze(ctx, mem))

	// Bypass the cache to read what was stored.
	state, err := r.queryUserState(ctx, mem.UserID())
	require.NoError(t, err)
	for _, ts := range []time.Time{state.LastMutate, state.LastReroll, state.LastFrozen} {
		assert.NotEqual(t, Never, ts)
		assert.WithinDuration(t, time.Now(), ts, 5*time.Second)
	}

	require.NoError(t, r.UpdateUnfreeze(ctx, mem))
	state, err = r.queryUserState(ctx, mem.UserID())
	require.NoError(t, err)
	assert.Equal(t, Never, state.LastFrozen)
	assert.NotEqual(t, Never, state.LastMutate)
}

func Test_sqlite_rerollPenaltyKeepsTimezoneAndPrecision(t *testing.T) {
	ctx := context.Background()
	mem := newTestMember(gomock.NewController(t))
	r, _ := newSQLiteRepo(t)

	require.NoError(t, r.UpdateReroll(ctx, mem, &Colour{}))

	// A non-UTC time must come back as the same instant.
	sgt := time.FixedZone("SGT", 8*60*60)
	penalty := time.Date(2030, 1, 2, 3, 4, 5, 600, sgt)
	require.NoError(t, r.UpdateRerollPenalty(ctx, mem, penalty))

	state, err := r.queryUserState(ctx, mem.UserID())
	require.NoError(t, err)
	assert.True(t, penalty.Equal(state.LastReroll), "got %v, want %v", state.LastReroll, penalty)
}

func Test_sqlite_historyRespectsSinceAndWindow(t *testing.T) {
	ctx := context.Background()
	mem := newTestMember(gomock.NewController(t))
	r, db := newSQLiteRepo(t)

	now := time.Now()
	insert := func(colour string, ts time.Time) {
		_, err := db.Exec(ctx,
			"INSERT INTO colours_log(userid, username, colour, reason, tstamp) VALUES ($1, $2, $3, $4, $5)",
			mem.UserID(), mem.Username(), colour, "mutate", ts)
		require.NoError(t, err)
	}
	insert("#old", now.AddDate(-3, 0, 0))                                // outside the 2-year view
	insert("#a", now.Add(-48*time.Hour))                                 // before `since`
	insert("#b", now.Add(-time.Hour).In(time.FixedZone("X", -10*60*60))) // inside, odd timezone
	insert("#c", now.Add(-time.Minute).Truncate(time.Second))            // inside, no fraction
	insert("", now.Add(-30*time.Second))                                 // empty colour is skipped

	history, err := r.FetchUserHistory(ctx, mem, now.Add(-24*time.Hour))
	require.NoError(t, err)

	var got []string
	for _, h := range history {
		got = append(got, h.ColourHex)
	}
	assert.Equal(t, []string{"#b", "#c"}, got)
}
