package database

// pgclient.go implements a Postgres client satisfying IDatabase.

import (
	"context"
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

func NewDBClient(ctx context.Context, dsn string) (IDatabase, error) {
	pool, err := sql.Open("postgres", dsn)
	log.Info().Msgf("Database connection opened")
	if err != nil {
		return nil, err
	}
	c := &pgclient{
		pool:               pool,
		existingMigrations: make(map[string]bool),
	}
	if err := c.seedMigration(ctx); err != nil {
		log.Error().Msgf("Seed migrations failed")
		defer c.Close(ctx)
		return nil, err
	}
	return c, err
}

func (c *pgclient) Close(ctx context.Context) error {
	if err := c.pool.Close(); err != nil {
		log.Error().Err(err).Msgf("Failed to close database connection")
		return err
	}
	log.Info().Msgf("Database connection closed")
	return nil
}

func (c *pgclient) Query(ctx context.Context, query string, args ...interface{}) (IRows, error) {
	log.Info().Msgf("Query: %s", query)
	if len(args) > 0 {
		log.Info().Msgf(" Args: %v", args)
	}
	rows, err := c.pool.Query(query, args...)
	if err == sql.ErrNoRows {
		return rows, fmt.Errorf("%w (driver: %v)", ErrNoRecords, err)
	}
	return rows, err
}

func (c *pgclient) Exec(ctx context.Context, query string, args ...interface{}) (IResult, error) {
	log.Info().Msgf(" Exec: %s", query)
	if len(args) > 0 {
		log.Info().Msgf(" Args: %v", args)
	}
	affected, err := c.pool.Exec(query, args...)
	return affected, err
}

func (c *pgclient) Begin(ctx context.Context) (ITransaction, error) {
	t, err := c.pool.Begin()
	if err != nil {
		return nil, err
	}
	return sqlTxWrap{t}, nil
}

// sqlTxWrap implements ITransaction for (database/sql).*Tx.
type sqlTxWrap struct {
	*sql.Tx
}

func (t sqlTxWrap) Query(ctx context.Context, query string, args ...interface{}) (IRows, error) {
	log.Info().Msgf("Query: %s", query)
	rows, err := t.Tx.Query(query, args...)
	if errors.Is(err, sql.ErrNoRows) {
		return rows, fmt.Errorf("%w (driver: %v)", ErrNoRecords, err)
	}
	return rows, err
}

func (t sqlTxWrap) Exec(ctx context.Context, query string, args ...interface{}) (IResult, error) {
	log.Info().Msgf("Exec: %s", query)
	return t.Tx.Exec(query, args...)
}

func (t sqlTxWrap) Commit(ctx context.Context) error {
	return t.Tx.Commit()
}

func (t sqlTxWrap) Rollback(ctx context.Context) error {
	return t.Tx.Rollback()
}
