package core

import "database/sql"

type PluginProviderDispatcher interface{}

type PluginProvider interface {
	Controllers() []Controller
	DownStreamPlugins() []DownStreamPlugin
}

type MigrationProviderDispatcher interface {
	GoMigrate(db *sql.DB, dbConfig DatabaseConfig) (MigrationProvider, error)
}

type MigrationProvider interface {
	Migrator() (Migrator, error)
}
