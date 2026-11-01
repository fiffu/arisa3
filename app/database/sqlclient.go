package database

// sqlclient.go implements a SQLite client satisfying IDatabase.

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/fiffu/arisa3/app/instrumentation"
	"github.com/fiffu/arisa3/app/log"
	"github.com/fiffu/arisa3/lib"
	"go.opentelemetry.io/otel/trace"
	_ "modernc.org/sqlite"
)

// TimeFormat is how timestamps are stored. SQLite has no timestamp type, so they are TEXT and compared as
// strings. A fixed width in UTC keeps string order the same as time order.
const TimeFormat = "2006-01-02T15:04:05.000000000Z"

// FormatTime converts t to its stored form.
func FormatTime(t time.Time) string {
	return t.UTC().Format(TimeFormat)
}

// sqlitePragmas are applied to every connection. WAL lets reads run alongside a write, busy_timeout waits
// for a lock instead of failing at once, and _txlock=immediate takes the write lock at BEGIN so two
// transactions can't deadlock upgrading from read to write.
var sqlitePragmas = []string{
	"_pragma=busy_timeout(5000)",
	"_pragma=journal_mode(WAL)",
	"_pragma=synchronous(NORMAL)",
	"_txlock=immediate",
}

// sqliteDSN turns a file path (or a "file:" URI) into a DSN with sqlitePragmas added.
func sqliteDSN(dsn string) (string, error) {
	if dsn == "" {
		return "", errors.New("database_dsn is empty; set it to a SQLite file path")
	}
	if strings.Contains(dsn, "://") {
		return "", fmt.Errorf("database_dsn looks like a URL (%q...), expected a SQLite file path", strings.SplitN(dsn, "://", 2)[0])
	}
	if !strings.HasPrefix(dsn, "file:") {
		dsn = "file:" + dsn
	}
	sep := "?"
	if strings.Contains(dsn, "?") {
		sep = "&"
	}
	return dsn + sep + strings.Join(sqlitePragmas, "&"), nil
}

// normalizeArgs stores every time.Time in TimeFormat, whatever its location.
func normalizeArgs(args []any) []any {
	out := make([]any, len(args))
	for i, a := range args {
		if t, ok := a.(time.Time); ok {
			a = FormatTime(t)
		}
		out[i] = a
	}
	return out
}

// sqlclient implements IDatabase for database/sql + modernc.org/sqlite.
type sqlclient struct {
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

func open(ctx context.Context, dsn string) (*sqlclient, error) {
	ctx, span := instrumentation.SpanInContext(ctx, instrumentation.Database("sql.Open"))
	defer span.End()

	dsn, err := sqliteDSN(dsn)
	if err != nil {
		return nil, err
	}
	pool, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}
	// sql.Open doesn't connect; fail now if the file can't be opened.
	if err := pool.PingContext(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("opening sqlite database: %w", err)
	}
	log.Infof(ctx, "Database connection opened")
	return &sqlclient{
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

func (c *sqlclient) Close(ctx context.Context) error {
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

		return callable(ctx, query, normalizeArgs(args)...)
	}
}

func (c *sqlclient) Query(ctx context.Context, query string, args ...interface{}) (IRows, error) {
	op := newOperation[*sql.Rows](c.pool.QueryContext, lib.WhoCalledMe())
	rows, err := op(ctx, query, args...)
	if err == sql.ErrNoRows {
		return rows, fmt.Errorf("%w (driver: %v)", ErrNoRecords, err)
	}
	return rows, err
}

func (c *sqlclient) Exec(ctx context.Context, query string, args ...interface{}) (IResult, error) {
	op := newOperation[sql.Result](c.pool.ExecContext, lib.WhoCalledMe())
	affected, err := op(ctx, query, args...)
	return affected, err
}

func (c *sqlclient) Begin(ctx context.Context) (context.Context, ITransaction, error) {
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
