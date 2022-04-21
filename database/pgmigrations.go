package database

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"

	"github.com/fiffu/arisa3/lib"
)

// pgmigrations.go implements migrations for pgclient

const (
	migrationInsert        = "INSERT INTO _schema_migrations (version) VALUES ($1);"
	listMigrations         = "SELECT version FROM _schema_migrations;"
	createSchemaMigrations = `
		CREATE TABLE IF NOT EXISTS "_schema_migrations" (
			version TEXT PRIMARY KEY
		);`
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
	version string
}

func (r MigrationRecord) Scan(rows IRows) error {
	return rows.Scan(&r.version)
}

// seedMigration pulls the migrations table state, or creates if it doesn't exist.
func (c *pgclient) seedMigration() error {
	if _, err := c.Exec(createSchemaMigrations); err != nil {
		return err
	}
	rows, err := c.Query(listMigrations)
	if err != nil {
		return err
	}
	for rows.Next() {
		record := MigrationRecord{}
		if err := record.Scan(rows); err != nil {
			return err
		}
		c.existingMigrations[record.version] = true
	}
	return nil
}

// Migrate executes a migration and records it in the migrations table.
func (c *pgclient) Migrate(schema ISchema) error {
	if _, ok := c.existingMigrations[schema.Version()]; ok {
		return nil
	}
	txn, err := c.pool.Begin()
	if err != nil {
		return err
	}
	for _, q := range schema.Queries() {
		if _, err := txn.Exec(q); err != nil {
			if err := txn.Rollback(); err != nil {
				return err
			}
			return err
		}
	}
	query := migrationInsert
	if _, err := txn.Exec(query, schema.Version()); err != nil {
		if err := txn.Rollback(); err != nil {
			return err
		}
		return err
	}
	return txn.Commit()
}

func (c *pgclient) ParseMigration(theFile string) (ISchema, error) {
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
