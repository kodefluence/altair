package core

import "database/sql"

type PluginProviderDispatcher interface{}

type MigrationProviderDispatcher interface {
	GoMigrate(db *sql.DB, dbConfig DatabaseConfig) MigrationProvider
}

type MigrationProvider interface {
	Migrator() (Migrator, error)
}
