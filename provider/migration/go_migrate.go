package migration

import (
	"database/sql"

	"github.com/kodefluence/altair/core"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type goMigrate struct {
	db       *sql.DB
	dbConfig core.DatabaseConfig
}

func GoMigrate(db *sql.DB, dbConfig core.DatabaseConfig) core.MigrationProvider {
	return &goMigrate{db: db, dbConfig: dbConfig}
}

func (g *goMigrate) Migrator() (core.Migrator, error) {
	driver, err := mysql.WithInstance(g.db, &mysql.Config{
		MigrationsTable: "db_versions",
		DatabaseName:    g.dbConfig.DBDatabase(),
	})
	if err != nil {
		return nil, err
	}

	// NOTES: For now only supporting mysql database. But change this if in the future support other databases.
	migrator, err := migrate.NewWithDatabaseInstance(g.dbConfig.DBMigrationSource(), "mysql", driver)
	if err != nil {
		return nil, err
	}

	return migrator, nil
}
