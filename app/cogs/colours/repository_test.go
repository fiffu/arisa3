package colours

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/fiffu/arisa3/app/database"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func newTestMember(ctrl *gomock.Controller) *MockIDomainMember {
	mem := NewMockIDomainMember(ctrl)
	mem.EXPECT().UserID().Return("123123123123").AnyTimes()
	mem.EXPECT().Username().Return("Username").AnyTimes()
	return mem
}

func Test_newRepo(t *testing.T) {
	db, _, err := database.NewMockDBClient(t)
	assert.NoError(t, err)

	repo := NewRepository(db)
	assert.NotNil(t, repo)
}

func Test_FetchUserState_whenCached_returnCachedWithoutDatabaseCall(t *testing.T) {
	ctrl := gomock.NewController(t)
	mem := newTestMember(ctrl)
	db, _, err := database.NewMockDBClient(t)
	assert.NoError(t, err)

	repo := newRepo(db)

	mutateTime := time.Now()
	repo.cachePut(&ColourState{
		UserID:     mem.UserID(),
		LastMutate: mutateTime,
		LastReroll: Never,
		LastFrozen: Never,
	})

	for reason, expect := range map[Reason]time.Time{
		Mutate: mutateTime,
		Reroll: Never,
		Freeze: Never,
	} {
		actual, err := repo.FetchUserState(context.Background(), mem, reason)
		assert.NoError(t, err)
		assert.Equal(t, expect, actual)
	}
}

func Test_FetchUserState_whenNoTimestampInDB_returnNever(t *testing.T) {
	ctrl := gomock.NewController(t)
	mem := newTestMember(ctrl)
	db, dbMock, err := database.NewMockDBClient(t)
	assert.NoError(t, err)

	repo := newRepo(db)

	for reason, expect := range map[Reason]time.Time{
		Mutate: Never,
		Reroll: Never,
		Freeze: Never,
	} {
		dbMock.ExpectQuery(`SELECT userid, tstamp, reason FROM colours WHERE userid = \$1`).
			WillReturnError(sql.ErrNoRows)
		actual, err := repo.FetchUserState(context.Background(), mem, reason)
		assert.NoError(t, err)
		assert.Equal(t, expect, actual)
	}
}

func Test_FetchUserState_cacheReadThrough(t *testing.T) {
	ctrl := gomock.NewController(t)
	mem := newTestMember(ctrl)
	db, dbMock, err := database.NewMockDBClient(t)
	assert.NoError(t, err)

	repo := newRepo(db)
	lastReroll := time.Now()

	dbMock.ExpectQuery(`SELECT userid, tstamp, reason FROM colours WHERE userid = \$1`).
		WillReturnRows(sqlmock.
			NewRows([]string{"userid", "tstamp", "reason"}).
			AddRow(mem.UserID(), lastReroll, "reroll"))

	// First fetch triggers DB call due to cache miss
	actual, err := repo.FetchUserState(context.Background(), mem, Reroll)
	assert.NoError(t, err)
	assert.Equal(t, lastReroll, actual)

	// Second fetch should not trigger DB call
	actual2, err2 := repo.FetchUserState(context.Background(), mem, Reroll)
	assert.NoError(t, err2)
	assert.Equal(t, lastReroll, actual2)
}

func Test_queryUserState_whenNoTimestampInDB_returnError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mem := newTestMember(ctrl)
	db, dbMock, err := database.NewMockDBClient(t)
	assert.NoError(t, err)

	repo := newRepo(db)

	dbMock.ExpectQuery(`SELECT userid, tstamp, reason FROM colours WHERE userid = \$1`).
		WillReturnError(sql.ErrNoRows)

	state, err := repo.queryUserState(context.Background(), mem.UserID())

	assert.ErrorIs(t, err, database.ErrNoRecords)
	assert.Nil(t, state)
}

func Test_queryUserState_whenTimestampsInDB_returnTimestamps(t *testing.T) {
	ctrl := gomock.NewController(t)
	mem := newTestMember(ctrl)
	db, dbMock, err := database.NewMockDBClient(t)
	assert.NoError(t, err)

	repo := newRepo(db)

	lastMutate := time.Now().Add(-5 * time.Hour)
	lastReroll := time.Now().Add(-1 * time.Hour)
	lastFrozen := time.Now().Add(-2 * time.Hour)

	dbMock.ExpectQuery(`SELECT userid, tstamp, reason FROM colours WHERE userid = \$1`).
		WillReturnRows(sqlmock.
			NewRows([]string{"userid", "tstamp", "reason"}).
			AddRow(mem.UserID(), lastMutate, "mutate").
			AddRow(mem.UserID(), lastReroll, "reroll").
			AddRow(mem.UserID(), lastFrozen, "freeze"))
	state, err := repo.queryUserState(context.Background(), mem.UserID())

	assert.NoError(t, err)
	assert.Equal(t, mem.UserID(), state.UserID)
	assert.Equal(t, lastMutate, state.LastMutate, "lastMutate does not match expected")
	assert.Equal(t, lastReroll, state.LastReroll, "lastReroll does not match expected")
	assert.Equal(t, lastFrozen, state.LastFrozen, "lastFrozen does not match expected")
}

func Test_queryUserState_whenSomeTimestampsInDB_returnNeverForLackingRecords(t *testing.T) {
	ctrl := gomock.NewController(t)
	mem := newTestMember(ctrl)
	db, dbMock, err := database.NewMockDBClient(t)
	assert.NoError(t, err)

	repo := newRepo(db)

	lastMutate := time.Now().Add(-5 * time.Hour)

	dbMock.ExpectQuery(`SELECT userid, tstamp, reason FROM colours WHERE userid = \$1`).
		WillReturnRows(sqlmock.
			NewRows([]string{"userid", "tstamp", "reason"}).
			AddRow(mem.UserID(), lastMutate, "mutate"))
	state, err := repo.queryUserState(context.Background(), mem.UserID())

	assert.NoError(t, err)
	assert.Equal(t, mem.UserID(), state.UserID)
	assert.Equal(t, lastMutate, state.LastMutate, "lastMutate does not match expected")
	assert.Equal(t, Never, state.LastReroll, "lastReroll does not match expected")
	assert.Equal(t, Never, state.LastFrozen, "lastFrozen does not match expected")
}

func Test_UpdateFoo(t *testing.T) {
	ctrl := gomock.NewController(t)
	mem := newTestMember(ctrl)
	db, dbMock, err := database.NewMockDBClient(t)
	assert.NoError(t, err)

	repo := newRepo(db)

	col := &Colour{1, 1, 1}
	for reason, method := range map[Reason](func() error){
		Mutate: func() error { return repo.UpdateMutate(context.Background(), mem, col) },
		Reroll: func() error { return repo.UpdateReroll(context.Background(), mem, col) },
		Freeze: func() error { return repo.UpdateFreeze(context.Background(), mem) },
	} {
		dbMock.ExpectBegin()
		dbMock.ExpectExec(`DELETE FROM colours WHERE userid = \$1 AND reason = \$2`).
			WillReturnResult(sqlmock.NewResult(1, 0))
		dbMock.ExpectExec(`INSERT INTO colours\(userid, tstamp, reason\) VALUES \(\$1, \$2, \$3\)`).
			WillReturnResult(sqlmock.NewResult(1, 0))
		dbMock.ExpectCommit()
		dbMock.ExpectExec(`INSERT INTO colours_log\(.+\) VALUES \(.+\)`).
			WillReturnResult(sqlmock.NewResult(1, 0))

		err = method()
		assert.NoError(t, err)

		tstamp, ok := repo.cachePeek(mem.UserID(), reason)
		assert.True(t, ok)
		assert.False(t, tstamp.IsZero())
	}
}

func Test_UpdateUnfreeze(t *testing.T) {
	ctrl := gomock.NewController(t)
	mem := newTestMember(ctrl)
	db, dbMock, err := database.NewMockDBClient(t)
	assert.NoError(t, err)
	dbMock.ExpectExec(`DELETE FROM colours WHERE userid = \$1 AND reason = \$2`).
		WillReturnResult(sqlmock.NewResult(1, 0))
	dbMock.ExpectExec(`INSERT INTO colours_log\(.+\) VALUES \(.+\)`).
		WillReturnResult(sqlmock.NewResult(1, 0))

	repo := newRepo(db)
	err = repo.UpdateUnfreeze(context.Background(), mem)
	assert.NoError(t, err)

	freezeTime, ok := repo.cachePeek(mem.UserID(), Freeze)
	assert.True(t, ok)
	assert.Equal(t, Never, freezeTime)
}

func Test_UpdateReroll(t *testing.T) {
	ctrl := gomock.NewController(t)
	mem := newTestMember(ctrl)
	db, dbMock, err := database.NewMockDBClient(t)
	assert.NoError(t, err)

	repo := newRepo(db)

	dbMock.ExpectBegin()
	dbMock.ExpectExec(`DELETE FROM colours WHERE userid = \$1 AND reason = \$2`).
		WillReturnResult(sqlmock.NewResult(1, 0))
	dbMock.ExpectExec(`INSERT INTO colours\(userid, tstamp, reason\) VALUES \(\$1, \$2, \$3\)`).
		WillReturnResult(sqlmock.NewResult(1, 0))
	dbMock.ExpectCommit()
	dbMock.ExpectExec(`INSERT INTO colours_log\(.+\) VALUES \(.+\)`).
		WillReturnResult(sqlmock.NewResult(1, 0))
	err = repo.UpdateReroll(context.Background(), mem, &Colour{1, 1, 1})
	assert.NoError(t, err)
}

func Test_UpdateRerollPenalty(t *testing.T) {
	ctrl := gomock.NewController(t)
	mem := newTestMember(ctrl)
	db, dbMock, err := database.NewMockDBClient(t)
	assert.NoError(t, err)

	repo := newRepo(db)

	dbMock.ExpectExec(`UPDATE colours SET tstamp=\$1 WHERE userid=\$2 AND reason=\$3`).
		WillReturnResult(sqlmock.NewResult(1, 0))

	err = repo.UpdateRerollPenalty(context.Background(), mem, time.Now())
	assert.NoError(t, err)
}
