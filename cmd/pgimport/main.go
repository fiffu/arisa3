// pgimport copies arisa3's data from the old Postgres database into a new SQLite file.
//
// It creates the SQLite file with the app's own migrations, then copies every table in one transaction,
// rewriting timestamps into database.TimeFormat. It's a one-off for the Postgres -> SQLite move; delete it
// (and the lib/pq dependency) once the migration is done.
//
// Run from the repo root (it reads app/cogs/*/dbmigrations):
//
//	PG_DSN='postgres://user:pass@host:5432/db?sslmode=disable' go run ./cmd/pgimport -out arisa3.db -source-tz UTC
//
// The DSN is read from the environment so the password stays off the command line.
package main

import (
	"context"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/fiffu/arisa3/app/database"
	_ "github.com/lib/pq"
)

type table struct {
	name       string
	columns    []string
	timestamps map[string]bool // columns of Postgres type TIMESTAMP
}

// The tables the migrations create, minus _schema_migrations (the new file gets its own).
// colours_log is the partitioned parent, so selecting from it reads every partition.
var tables = []table{
	{"colours", []string{"userid", "tstamp", "reason"}, map[string]bool{"tstamp": true}},
	{"colours_log", []string{"userid", "username", "colour", "reason", "tstamp"}, map[string]bool{"tstamp": true}},
	{"aliases", []string{"alias", "actual", "guildid"}, nil},
	{"tag_promote", []string{"tag", "guildid"}, nil},
	{"tag_demote", []string{"tag", "guildid"}, nil},
	{"tag_omit", []string{"tag", "guildid"}, nil},
}

func main() {
	out := flag.String("out", "", "SQLite file to create (must not exist)")
	sourceTZ := flag.String("source-tz", "", "timezone the bot ran in when it wrote to Postgres, e.g. UTC or Asia/Singapore")
	flag.Parse()
	if *out == "" || *sourceTZ == "" {
		flag.Usage()
		os.Exit(2)
	}
	if err := run(*out, *sourceTZ, os.Getenv("PG_DSN")); err != nil {
		log.Fatalf("pgimport: %v", err)
	}
}

func run(out, sourceTZ, pgDSN string) error {
	ctx := context.Background()
	if pgDSN == "" {
		return errors.New("PG_DSN is not set")
	}
	// Postgres TIMESTAMP has no zone: it holds the wall-clock time of whoever wrote it.
	loc, err := time.LoadLocation(sourceTZ)
	if err != nil {
		return fmt.Errorf("bad -source-tz: %w", err)
	}
	if _, err := os.Stat(out); err == nil {
		return fmt.Errorf("%s already exists; refusing to overwrite", out)
	}

	pg, err := sql.Open("postgres", pgDSN)
	if err != nil {
		return err
	}
	defer pg.Close()
	if err := pg.PingContext(ctx); err != nil {
		return fmt.Errorf("connecting to postgres: %w", err)
	}

	if err := migrate(ctx, out); err != nil {
		return fmt.Errorf("creating schema: %w", err)
	}

	// Raw handle for the copy: the app's client logs every statement, which is too noisy here.
	lite, err := sql.Open("sqlite", "file:"+out+"?_pragma=busy_timeout(5000)&_txlock=immediate")
	if err != nil {
		return err
	}
	defer lite.Close()

	tx, err := lite.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	counts := map[string]int{}
	for _, t := range tables {
		n, err := copyTable(ctx, pg, tx, t, loc)
		if err != nil {
			return fmt.Errorf("copying %s: %w", t.name, err)
		}
		counts[t.name] = n
	}
	if err := tx.Commit(); err != nil {
		return err
	}

	// Check that row counts match on both sides.
	ok := true
	for _, t := range tables {
		var pgN, liteN int
		if err := pg.QueryRowContext(ctx, `SELECT count(*) FROM "`+t.name+`"`).Scan(&pgN); err != nil {
			return err
		}
		if err := lite.QueryRowContext(ctx, `SELECT count(*) FROM "`+t.name+`"`).Scan(&liteN); err != nil {
			return err
		}
		mark := "ok"
		if pgN != liteN {
			// Postgres can gain rows mid-copy if the bot is still running; that's expected drift.
			mark = "DIFFERS (rows written to postgres during the copy?)"
			ok = false
		}
		fmt.Printf("%-12s postgres=%-8d sqlite=%-8d %s\n", t.name, pgN, liteN, mark)
	}

	var integrity string
	if err := lite.QueryRowContext(ctx, "PRAGMA integrity_check").Scan(&integrity); err != nil {
		return err
	}
	fmt.Println("integrity_check:", integrity)
	if _, err := lite.ExecContext(ctx, "VACUUM"); err != nil {
		return err
	}
	if !ok {
		fmt.Println("row counts differ; see above")
	}
	fmt.Println("wrote", out)
	return nil
}

// migrate creates the file and applies every cog's migrations through the app's own client, so
// _schema_migrations is filled in the same way as on a normal start.
func migrate(ctx context.Context, path string) error {
	db, err := database.NewDBClient(ctx, path)
	if err != nil {
		return err
	}
	defer db.Close(ctx)

	files, err := filepath.Glob("app/cogs/*/dbmigrations/*.sql")
	if err != nil {
		return err
	}
	if len(files) == 0 {
		return errors.New("no migrations found; run this from the repo root")
	}
	// Each cog's migrations only depend on that cog's earlier ones; sort by directory, then file name.
	sort.Strings(files)
	for _, f := range files {
		schema, err := db.ParseMigration(ctx, f)
		if err != nil {
			return err
		}
		if _, err := db.Migrate(ctx, schema); err != nil {
			return fmt.Errorf("%s: %w", f, err)
		}
	}
	return nil
}

func copyTable(ctx context.Context, pg *sql.DB, tx *sql.Tx, t table, loc *time.Location) (int, error) {
	cols := strings.Join(t.columns, ", ")
	rows, err := pg.QueryContext(ctx, fmt.Sprintf(`SELECT %s FROM "%s"`, cols, t.name))
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	placeholders := make([]string, len(t.columns))
	for i := range placeholders {
		placeholders[i] = "?"
	}
	ins, err := tx.PrepareContext(ctx, fmt.Sprintf(`INSERT INTO "%s" (%s) VALUES (%s)`,
		t.name, cols, strings.Join(placeholders, ", ")))
	if err != nil {
		return 0, err
	}
	defer ins.Close()

	n := 0
	for rows.Next() {
		vals := make([]any, len(t.columns))
		ptrs := make([]any, len(t.columns))
		for i, c := range t.columns {
			if t.timestamps[c] {
				ptrs[i] = new(sql.NullTime)
			} else {
				ptrs[i] = new(sql.NullString)
			}
		}
		if err := rows.Scan(ptrs...); err != nil {
			return n, err
		}
		for i, p := range ptrs {
			switch v := p.(type) {
			case *sql.NullTime:
				if v.Valid {
					// Keep the wall clock, attach the real zone, then store as UTC.
					w := v.Time
					vals[i] = database.FormatTime(time.Date(w.Year(), w.Month(), w.Day(),
						w.Hour(), w.Minute(), w.Second(), w.Nanosecond(), loc))
				}
			case *sql.NullString:
				if v.Valid {
					vals[i] = v.String
				}
			}
		}
		if _, err := ins.ExecContext(ctx, vals...); err != nil {
			return n, err
		}
		n++
	}
	return n, rows.Err()
}
