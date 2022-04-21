// package database exposes an interface for database I/O.
package database

// IDatabase describes the interface of a database client.
type IDatabase interface {
	// Close closes the database client.
	Close()

	// Query queries the database, usually a SELECT.
	Query(query string, args ...interface{}) (IRows, error)

	// Exec executes a statement on the database, usually an UPDATE/INSERT/DELETE.
	Exec(query string, args ...interface{}) (IResult, error)

	// Begin begins a transaction.
	Begin() (ITransaction, error)

	// Migrate executes a schema for database migration.
	Migrate(schema ISchema) error
}

type ITransaction interface {
	// Query queries the database, usually a SELECT.
	Query(query string, args ...interface{}) (IRows, error)

	// Exec executes a statement on the database, usually an UPDATE/INSERT/DELETE.
	Exec(query string, args ...interface{}) (IResult, error)

	// Commit commits the transaction
	Commit() error

	// Rollback rolls back the transaction
	Rollback() error
}

// IResult summarizes an executed SQL command.
type IResult interface {
	LastInsertId() (int64, error)
	RowsAffected() (int64, error)
}

type IRows interface {
	Next() bool
	Scan(dest ...interface{}) error
}

type ISchema interface {
	Version() string
	Queries() []string
	Source() string
}
