package database

// pgmigrations.go implements migrations for pgclient

const (
	migrationInsert        string = "INSERT INTO _schema_migrations (version) VALUES ($1);"
	listMigrations         string = "SELECT version FROM _schema_migrations;"
	createSchemaMigrations string = `
		CREATE TABLE IF NOT EXISTS "_schema_migrations" (
			version TEXT PRIMARY KEY
		);`
)

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
