package colours

// repository.go contains implementation of IDomainRepository.

import (
	"errors"
	"fmt"
	"time"

	"github.com/fiffu/arisa3/app/database"
)

// ColoursRecord models table 'colours'.
type ColoursRecord struct {
	UserID string
	Reason string
	TStamp time.Time
}

// ColoursLogRecord models table 'colours_log'.
type ColoursLogRecord struct {
	UserID    string
	Username  string
	ColourHex string
	Reason    string
	TStamp    time.Time
}

// repo implements IDomainRepository.
type repo struct {
	db    database.IDatabase
	cache map[string]map[Reason]time.Time
}

func NewRepository(db database.IDatabase) IDomainRepository {
	return &repo{
		db:    db,
		cache: make(map[string]map[Reason]time.Time),
	}
}

/* Cache functions */

// cachePeek checks a value belonging to the given userID in the cache.
func (r *repo) cachePeek(userID string, reason Reason) (tstamp time.Time, ok bool) {
	state, ok := r.cache[userID]
	if !ok {
		return
	}
	tstamp, ok = state[reason]
	return
}

// cachePut upserts all values associated to the given userID into the cache.
func (r *repo) cachePut(state *ColourState) {
	r.cache[state.UserID] = map[Reason]time.Time{
		Mutate: state.LastMutate,
		Reroll: state.LastReroll,
		Freeze: state.LastFrozen,
	}
}

// cachePatch partially updates one of the given userID's values in the cache.
func (r *repo) cachePatch(userID string, reason Reason, tstamp time.Time) {
	if state, ok := r.cache[userID]; ok {
		state[reason] = tstamp
	}
	r.cache[userID] = map[Reason]time.Time{
		reason: tstamp,
	}
}

/* Exported methods for IDomainRepository, and their supporting internals. */

// FetchUserState returns the user's state for the given Reason.
func (r *repo) FetchUserState(user IDomainMember, reason Reason) (time.Time, error) {
	userID := user.UserID()
	if state, ok := r.cachePeek(userID, reason); ok {
		return state, nil
	}
	if state, err := r.queryUserState(userID); err != nil {
		if !errors.As(err, &database.ErrNoRecords) {
			return Never, nil
		}
		return Never, err
	} else {
		r.cachePut(state)
		tstamp, _ := r.cachePeek(userID, reason)
		return tstamp, nil
	}
}

func (r *repo) queryUserState(userID string) (*ColourState, error) {
	// Pull records with the given userID.
	rows, err := r.db.Query(
		"SELECT userid, tstamp, reason FROM colours WHERE userid = $1",
		userID,
	)
	if err != nil {
		return nil, err
	}

	// Parsing records.
	records := make([]ColoursRecord, 0)
	for rows.Next() {
		rec := ColoursRecord{}
		if err := rows.Scan(&rec.UserID, &rec.TStamp, &rec.Reason); err != nil {
			return nil, err
		}
		records = append(records, rec)
	}

	// Collate records into a singular state object
	// No records means no state
	if len(records) == 0 {
		return NoState, nil
	}
	state := ColourState{
		UserID:     userID,
		LastMutate: Never,
		LastReroll: Never,
		LastFrozen: Never,
	}
	for _, rec := range records {
		switch Reason(rec.Reason) {
		case Mutate:
			state.LastMutate = rec.TStamp
		case Reroll:
			state.LastReroll = rec.TStamp
		case Freeze:
			state.LastFrozen = rec.TStamp
		}
	}
	return &state, nil
}

func (r *repo) UpdateMutate(user IDomainMember, c *Colour) error {
	return r.update(user, Mutate, c, time.Now())
}
func (r *repo) UpdateReroll(user IDomainMember, c *Colour) error {
	return r.update(user, Reroll, c, time.Now())
}
func (r *repo) UpdateFreeze(user IDomainMember) error {
	return r.update(user, Freeze, nil, time.Now())
}
func (r *repo) UpdateUnfreeze(user IDomainMember) error {
	return r.update(user, Freeze, nil, Never)
}

func (r *repo) update(user IDomainMember, reason Reason, colour *Colour, tstamp time.Time) error {
	userID := user.UserID()
	// Insert the change inside a transaction
	if err := r.insert(userID, reason.String(), tstamp); err != nil {
		return err
	}

	// Log the change to the audit db
	hexcode := ""
	if colour != nil {
		hexcode = colour.ToHexcode()
	}
	if err := r.log(userID, user.Username(), reason.String(), hexcode, tstamp); err != nil {
		return err
	}
	r.cachePatch(userID, reason, tstamp)
	return nil
}

func (r *repo) insert(userID string, reason string, tstamp time.Time) error {
	rec := ColoursRecord{
		UserID: userID,
		Reason: reason,
		TStamp: tstamp,
	}
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	// Drop any records for the given reason
	if _, err := tx.Exec(
		"DELETE FROM colours WHERE userid = $1 AND reason = $2",
		rec.UserID, rec.Reason,
	); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("got another error while rolling back due to '%v': %v", err, rbErr)
		}
		return err
	}

	// Put a new record for the given reason
	if _, err := tx.Exec(
		"INSERT INTO colours(userid, tstamp, reason) VALUES ($1, $2, $3)",
		rec.UserID, rec.TStamp, rec.Reason,
	); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("got another error while rolling back due to '%v': %v", err, rbErr)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (r *repo) log(userID string, name string, reason string, hexcode string, tstamp time.Time) error {
	rec := ColoursLogRecord{
		UserID:    userID,
		Username:  name,
		Reason:    reason,
		ColourHex: hexcode,
		TStamp:    tstamp,
	}
	_, err := r.db.Exec(
		"INSERT INTO colours_logview(userid, username, colour, reason, tstamp) VALUES ($1, $2, $3, $4, $5)",
		rec.UserID, rec.Username, rec.ColourHex, rec.Reason, rec.TStamp,
	)
	return err
}

func (r *repo) UpdateRerollPenalty(user IDomainMember, tstamp time.Time) error {
	rec := ColoursRecord{
		UserID: user.UserID(),
		Reason: Reroll.String(),
		TStamp: tstamp,
	}
	_, err := r.db.Exec(
		"UPDATE colours SET tstamp=$1 WHERE userid=$2 AND reason=$3",
		rec.TStamp, rec.UserID, rec.Reason,
	)
	return err
}
