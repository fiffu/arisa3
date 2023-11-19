package database

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/fiffu/arisa3/app/instrumentation"
)

func NewMockDBClient(t *testing.T) (IDatabase, sqlmock.Sqlmock, error) {
	db, dbMock, err := sqlmock.New()
	return &mockClient{db}, dbMock, err
}

type mockClient struct {
	db *sql.DB
}

func (c *mockClient) Close(ctx context.Context) error {
	return c.db.Close()
}

func (c *mockClient) Query(ctx context.Context, query string, args ...interface{}) (IRows, error) {
	rows, err := c.db.Query(query, args...)
	if err == sql.ErrNoRows {
		return rows, fmt.Errorf("%w (driver: %v)", ErrNoRecords, err)
	}
	return rows, err
}

func (c *mockClient) Exec(ctx context.Context, query string, args ...interface{}) (IResult, error) {
	return c.db.Exec(query, args...)
}

func (c *mockClient) Begin(ctx context.Context) (context.Context, ITransaction, error) {
	t, err := c.db.Begin()
	if err != nil {
		return nil, nil, err
	}
	ctx, span := instrumentation.SpanInContext(ctx, instrumentation.Database("Transaction"))
	return ctx, sqlTxWrap{t, span}, nil
}

func (c *mockClient) Migrate(ctx context.Context, schema ISchema) (executed bool, err error) {
	panic("not implemented")
}

func (c *mockClient) ParseMigration(ctx context.Context, filepath string) (ISchema, error) {
	panic("not implemented")
}
