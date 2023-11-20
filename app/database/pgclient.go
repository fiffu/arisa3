package database

// pgclient.go implements a Postgres client satisfying IDatabase.

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/fiffu/arisa3/app/instrumentation"
	"github.com/fiffu/arisa3/app/log"
	"github.com/fiffu/arisa3/lib"
	_ "github.com/lib/pq"
	"go.opentelemetry.io/otel/trace"
)

// pgclient implements IData for database/sql + lib/pq.
type pgclient struct {
	pool               *sql.DB
	existingMigrations map[string]bool
}

func NewDBClient(ctx context.Context, dsn string) (IDatabase, error) {
	c, err := open(ctx, dsn)
	if err != nil {
		return nil, err
	}

	if err := c.seedMigration(ctx); err != nil {
		log.Errorf(ctx, err, "Seed migrations failed")
		defer c.Close(ctx)
		return nil, err
	}
	return c, err
}

func open(ctx context.Context, dsn string) (*pgclient, error) {
	ctx, span := instrumentation.SpanInContext(ctx, instrumentation.Database("sql.Open"))
	defer span.End()

	pool, err := sql.Open("postgres", dsn)
	log.Infof(ctx, "Database connection opened")
	if err != nil {
		return nil, err
	}
	return &pgclient{
		pool:               pool,
		existingMigrations: make(map[string]bool),
	}, nil
}

func newSpan(ctx context.Context, caller, sql string) (context.Context, trace.Span) {
	operation := firstWord(sql)
	ctx, span := instrumentation.SpanInContext(ctx, instrumentation.Database(sql))
	span.SetAttributes(
		instrumentation.KV.DBOperation(operation),
		instrumentation.KV.DBQuery(sql),
	)
	return ctx, span
}

func firstWord(s string) (word string) {
	word = strings.SplitN(s, " ", 2)[0]
	word = strings.TrimSuffix(word, ";")
	word = strings.ToUpper(word)
	return
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

func newOperation[T any](callable delegate[T], caller string) delegate[T] {
	return func(ctx context.Context, query string, args ...any) (T, error) {
		prettyQuery := NormalizeSQL(query)
		log.Infof(ctx, prettyQuery)
		if len(args) > 0 {
			log.Infof(ctx, "  Args: %v", args)
		}

		_, span := newSpan(ctx, caller, prettyQuery)
		defer span.End()

		return callable(ctx, query, args...)
	}
}

func (c *pgclient) Query(ctx context.Context, query string, args ...interface{}) (IRows, error) {
	op := newOperation[*sql.Rows](c.pool.QueryContext, lib.WhoCalledMe())
	rows, err := op(ctx, query, args...)
	if err == sql.ErrNoRows {
		return rows, fmt.Errorf("%w (driver: %v)", ErrNoRecords, err)
	}
	return rows, err
}

func (c *pgclient) Exec(ctx context.Context, query string, args ...interface{}) (IResult, error) {
	op := newOperation[sql.Result](c.pool.ExecContext, lib.WhoCalledMe())
	affected, err := op(ctx, query, args...)
	return affected, err
}

func (c *pgclient) Begin(ctx context.Context) (context.Context, ITransaction, error) {
	t, err := c.pool.Begin()
	if err != nil {
		return ctx, nil, err
	}

	caller := lib.WhoCalledMe()
	ctx, span := newSpan(ctx, caller, "BEGIN")
	return ctx, sqlTxnWrap{t, span}, nil
}

// sqlTxnWrap implements ITransaction for (database/sql).*Tx.
type sqlTxnWrap struct {
	*sql.Tx

	span trace.Span
}

func (txn sqlTxnWrap) Query(ctx context.Context, query string, args ...interface{}) (IRows, error) {
	op := newOperation[*sql.Rows](txn.Tx.QueryContext, lib.WhoCalledMe())
	rows, err := op(ctx, query, args...)
	if errors.Is(err, sql.ErrNoRows) {
		return rows, fmt.Errorf("%w (driver: %v)", ErrNoRecords, err)
	}
	return rows, err
}

func (txn sqlTxnWrap) Exec(ctx context.Context, query string, args ...interface{}) (IResult, error) {
	op := newOperation[sql.Result](txn.Tx.ExecContext, lib.WhoCalledMe())
	rows, err := op(ctx, query, args...)
	return rows, err
}

func (txn sqlTxnWrap) Commit(ctx context.Context) error {
	defer txn.span.End() // this would execute after commitSpan.End()

	caller := lib.WhoCalledMe()
	_, commitSpan := newSpan(ctx, caller, "COMMIT")
	defer commitSpan.End()

	return txn.Tx.Commit()
}

func (txn sqlTxnWrap) Rollback(ctx context.Context) error {
	defer txn.span.End() // this would execute after rollbackSpan.End()

	caller := lib.WhoCalledMe()
	_, rollbackSpan := newSpan(ctx, caller, "ROLLBACK")
	defer rollbackSpan.End()

	return txn.Tx.Rollback()
}
