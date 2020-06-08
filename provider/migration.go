package provider

import (
	"database/sql"

	"github.com/codefluence-x/altair/core"
	mp "github.com/codefluence-x/altair/provider/migration"
)

type migration struct{}

func Migration() core.MigrationProviderDispatcher {
	return &migration{}
}

func (m *migration) GoMigrate(db *sql.DB, dbConfig core.DatabaseConfig) core.MigrationProvider {
	return mp.GoMigrate(db, dbConfig)
}
