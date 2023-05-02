// package database exposes an interface for database I/O.
package database

import (
	"context"
	"errors"
)

var (
	// ErrNoRecords indicates 0 results were returned for a query.
	ErrNoRecords = errors.New("no records found")
)

// IDatabase describes the interface of a database client.
//go:generate mockgen -source=database.go -destination=./database_mock.go -package=database
type IDatabase interface {
	// Close closes the database client.
	Close(ctx context.Context) error

	// Query queries the database, usually a SELECT.
	Query(ctx context.Context, query string, args ...interface{}) (IRows, error)

	// Exec executes a statement on the database, usually an UPDATE/INSERT/DELETE.
	Exec(ctx context.Context, query string, args ...interface{}) (IResult, error)

	// Begin begins a transaction.
	Begin(ctx context.Context) (ITransaction, error)

	// Migrate executes a schema for database migration.
	Migrate(ctx context.Context, schema ISchema) (executed bool, err error)

	// ParseMigration is a helper function for reading migrations.
	ParseMigration(ctx context.Context, filepath string) (ISchema, error)
}

// ITransaction describes an interface of a database transaction.
type ITransaction interface {
	// Query queries the database, usually a SELECT.
	Query(ctx context.Context, query string, args ...interface{}) (IRows, error)

	// Exec executes a statement on the database, usually an UPDATE/INSERT/DELETE.
	Exec(ctx context.Context, query string, args ...interface{}) (IResult, error)

	// Commit commits the transaction
	Commit(ctx context.Context) error

	// Rollback rolls back the transaction
	Rollback(ctx context.Context) error
}

// IResult summarizes an executed SQL command.
type IResult interface {
	LastInsertId() (int64, error)
	RowsAffected() (int64, error)
}

// IRows represents an iterable cursor over items returned by a database query.
type IRows interface {
	Next() bool
	Scan(dest ...interface{}) error
}

// ISchema represents a schema used in database migrations.
type ISchema interface {
	Version() string
	Source() string
	Queries() []string
}
