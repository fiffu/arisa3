package database

// pgclient.go implements a Postgres client satisfying IDatabase.

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

// pgclient implements IData for database/sql + lib/pq.
type pgclient struct {
	pool               *sql.DB
	existingMigrations map[string]bool
}

func NewDBClient(dsn string) (IDatabase, error) {
	pool, err := sql.Open("postgres", dsn)
	log.Info().Msgf("Database connection opened")
	if err != nil {
		return nil, err
	}
	c := &pgclient{
		pool:               pool,
		existingMigrations: make(map[string]bool),
	}
	if err := c.seedMigration(); err != nil {
		log.Error().Msgf("Seed migrations failed")
		defer c.Close()
		return nil, err
	}
	return c, err
}

func (c *pgclient) Close() error {
	return c.pool.Close()
}

func (c *pgclient) Query(query string, args ...interface{}) (IRows, error) {
	rows, err := c.pool.Query(query, args...)
	if errors.Is(err, sql.ErrNoRows) {
		return rows, fmt.Errorf("%w (driver: %v)", ErrNoRecords, err)
	}
	return rows, err
}

func (c *pgclient) Exec(query string, args ...interface{}) (IResult, error) {
	affected, err := c.pool.Exec(query, args...)
	return affected, err
}

func (c *pgclient) Begin() (ITransaction, error) {
	t, err := c.pool.Begin()
	if err != nil {
		return nil, err
	}
	return pgtxn{t}, nil
}

// pgtxn implements ITransaction for database/sql + lib/pq.
type pgtxn struct {
	*sql.Tx
}

func (t pgtxn) Query(query string, args ...interface{}) (IRows, error) {
	rows, err := t.Tx.Query(query, args...)
	if errors.Is(err, sql.ErrNoRows) {
		return rows, fmt.Errorf("%w (driver: %v)", ErrNoRecords, err)
	}
	return rows, err
}

func (t pgtxn) Exec(query string, args ...interface{}) (IResult, error) {
	return t.Tx.Exec(query, args...)
}
