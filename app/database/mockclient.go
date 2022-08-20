package database

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func NewMockDBClient(t *testing.T) (IDatabase, sqlmock.Sqlmock, error) {
	db, dbMock, err := sqlmock.New()
	return &mockClient{db}, dbMock, err
}

type mockClient struct {
	db *sql.DB
}

func (c *mockClient) Close() error {
	return c.db.Close()
}

func (c *mockClient) Query(query string, args ...interface{}) (IRows, error) {
	rows, err := c.db.Query(query, args...)
	if err == sql.ErrNoRows {
		return rows, fmt.Errorf("%w (driver: %v)", ErrNoRecords, err)
	}
	return rows, err
}

func (c *mockClient) Exec(query string, args ...interface{}) (IResult, error) {
	return c.db.Exec(query, args...)
}

func (c *mockClient) Begin() (ITransaction, error) {
	t, err := c.db.Begin()
	if err != nil {
		return nil, err
	}
	return sqlTxWrap{t}, nil
}

func (c *mockClient) Migrate(ISchema) (executed bool, err error) {
	panic("not implemented")
}

func (c *mockClient) ParseMigration(filepath string) (ISchema, error) {
	panic("not implemented")
}
