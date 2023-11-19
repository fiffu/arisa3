package database

// pgclient.go implements a Postgres client satisfying IDatabase.

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/fiffu/arisa3/app/instrumentation"
	"github.com/fiffu/arisa3/app/log"
	"github.com/fiffu/arisa3/lib"
	_ "github.com/lib/pq"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// Instrumentation attributes
const (
	attrOperation = "operation"
	attrSQL       = "query"
)

// pgclient implements IData for database/sql + lib/pq.
type pgclient struct {
	pool               *sql.DB
	existingMigrations map[string]bool
}

func NewDBClient(ctx context.Context, dsn string) (IDatabase, error) {
	pool, err := sql.Open("postgres", dsn)
	log.Infof(ctx, "Database connection opened")
	if err != nil {
		return nil, err
	}
	c := &pgclient{
		pool:               pool,
		existingMigrations: make(map[string]bool),
	}
	if err := c.seedMigration(ctx); err != nil {
		log.Errorf(ctx, err, "Seed migrations failed")
		defer c.Close(ctx)
		return nil, err
	}
	return c, err
}

func newSpan(ctx context.Context, caller, operation, sql string) (context.Context, trace.Span) {
	ctx, span := instrumentation.SpanInContext(ctx, instrumentation.Database(operation))
	span.SetAttributes(
		attribute.String(attrOperation, operation),
		attribute.String(attrSQL, sql),
	)
	return ctx, span
}

func (c *pgclient) Close(ctx context.Context) error {
	if err := c.pool.Close(); err != nil {
		log.Errorf(ctx, err, "Failed to close database connection")
		return err
	}
	log.Infof(ctx, "Database connection closed")
	return nil
}

type delegate[T any] func(ctx context.Context, query string, args ...any) (T, error)

func newOperation[T any](callable delegate[T], caller string, operation string) delegate[T] {
	return func(ctx context.Context, query string, args ...any) (T, error) {
		prettyQuery := NormalizeSQL(query)
		log.Infof(ctx, "%s: %s", operation, prettyQuery)
		if len(args) > 0 {
			log.Infof(ctx, " Args: %v", args)
		}

		_, span := newSpan(ctx, caller, operation, NormalizeSQL(operation))
		defer span.End()

		return callable(ctx, query, args...)
	}
}

func (c *pgclient) Query(ctx context.Context, query string, args ...interface{}) (IRows, error) {
	op := newOperation[*sql.Rows](c.pool.QueryContext, lib.WhoCalledMe(), "Query")
	rows, err := op(ctx, query, args...)
	if err == sql.ErrNoRows {
		return rows, fmt.Errorf("%w (driver: %v)", ErrNoRecords, err)
	}
	return rows, err
}

func (c *pgclient) Exec(ctx context.Context, query string, args ...interface{}) (IResult, error) {
	op := newOperation[sql.Result](c.pool.ExecContext, lib.WhoCalledMe(), "Exec")
	affected, err := op(ctx, query, args...)
	return affected, err
}

func (c *pgclient) Begin(ctx context.Context) (ITransaction, error) {
	t, err := c.pool.Begin()
	if err != nil {
		return nil, err
	}

	caller := lib.WhoCalledMe()
	ctx, span := newSpan(ctx, caller, "Transaction", "BEGIN")
	return sqlTxWrap{t, ctx, span}, nil
}

// sqlTxWrap implements ITransaction for (database/sql).*Tx.
type sqlTxWrap struct {
	*sql.Tx

	ctx  context.Context
	span trace.Span
}

func (t sqlTxWrap) Query(ctx context.Context, query string, args ...interface{}) (IRows, error) {
	op := newOperation[*sql.Rows](t.Tx.QueryContext, lib.WhoCalledMe(), "Query")
	rows, err := op(ctx, query, args...)
	if errors.Is(err, sql.ErrNoRows) {
		return rows, fmt.Errorf("%w (driver: %v)", ErrNoRecords, err)
	}
	return rows, err
}

func (t sqlTxWrap) Exec(ctx context.Context, query string, args ...interface{}) (IResult, error) {
	op := newOperation[sql.Result](t.Tx.ExecContext, lib.WhoCalledMe(), "Exec")
	rows, err := op(ctx, query, args...)
	return rows, err
}

func (t sqlTxWrap) Commit(ctx context.Context) error {
	defer t.span.End()
	return t.Tx.Commit()
}

func (t sqlTxWrap) Rollback(ctx context.Context) error {
	defer t.span.End()
	return t.Tx.Rollback()
}
