package database

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"

	"github.com/fiffu/arisa3/app/log"
	"github.com/fiffu/arisa3/lib"
)

// pgmigrations.go implements migrations for pgclient

const (
	createSchemaMigrations = `CREATE TABLE IF NOT EXISTS "_schema_migrations" (version TEXT PRIMARY KEY);`
)

var (
	filenamePattern   = regexp.MustCompile(`\d+_.+\.sql`)
	ErrParseMigration = errors.New("failed to parse migration")
)

// sqlSchema implements ISchema for .sql files.
type sqlSchema struct {
	source  string
	version string
	queries []string
}

func (s sqlSchema) Source() string    { return s.source }
func (s sqlSchema) Version() string   { return s.version }
func (s sqlSchema) Queries() []string { return s.queries }

type MigrationRecord struct {
	Version string
}

func (r *MigrationRecord) Scan(rows IRows) error {
	return rows.Scan(&r.Version)
}

// seedMigration pulls the migrations table state, or creates if it doesn't exist.
func (c *pgclient) seedMigration(ctx context.Context) error {
	log.Infof(ctx, "Creating schema migrations table")
	if _, err := c.Exec(ctx, createSchemaMigrations); err != nil {
		log.Errorf(ctx, err, "Failed to creating seed migrations table")
		return err
	}
	rows, err := c.Query(ctx, "SELECT version FROM _schema_migrations;")
	if err != nil {
		return err
	}

	for rows.Next() {
		rec := &MigrationRecord{}
		if err := rec.Scan(rows); err != nil {
			log.Errorf(ctx, err, "Failed parsing migration record: %v", rows)
			return err
		}
		c.existingMigrations[rec.Version] = true
	}
	log.Infof(ctx, "Loaded schema migrations (noted %d migration records)", len(c.existingMigrations))
	return nil
}

// Migrate executes a migration and records it in the migrations table.
func (c *pgclient) Migrate(ctx context.Context, schema ISchema) (bool, error) {
	if _, ok := c.existingMigrations[schema.Version()]; ok {
		return false, nil
	}
	log.Infof(ctx, "Executing migration %s (%s)", schema.Version(), schema.Source())
	txn, err := c.pool.Begin()
	if err != nil {
		return false, err
	}
	for _, q := range schema.Queries() {
		if _, err := txn.Exec(q); err != nil {
			if err := txn.Rollback(); err != nil {
				return false, err
			}
			return false, err
		}
	}
	query := "INSERT INTO _schema_migrations (version) VALUES ($1);"
	if _, err := txn.Exec(query, schema.Version()); err != nil {
		if err := txn.Rollback(); err != nil {
			return false, err
		}
		return false, err
	}
	if err := txn.Commit(); err != nil {
		return false, err
	}
	return true, nil
}

// ParseMigration implements parsing of files into sqlSchema.
func (c *pgclient) ParseMigration(ctx context.Context, theFile string) (ISchema, error) {
	name, err := validateFileName(theFile)
	if err != nil {
		return nil, err
	}

	bytes, err := ioutil.ReadFile(theFile)
	if err != nil {
		return nil, err
	}

	s := string(bytes)
	version, _ := lib.SplitOnce(name, "_")

	return sqlSchema{
		source:  theFile,
		version: version,
		queries: []string{s},
	}, nil
}

func validateFileName(theFile string) (string, error) {
	name := filepath.Base(theFile)
	if !filenamePattern.MatchString(name) {
		return "", fmt.Errorf(
			"%w: file should have pattern %s, got '%s'",
			ErrParseMigration,
			filenamePattern.String(), name,
		)
	}
	return name, nil
}
