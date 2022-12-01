package migration

import (
	"embed"

	"github.com/kodefluence/altair/core"
	"github.com/kodefluence/altair/module"
	"github.com/kodefluence/altair/plugin/oauth/module/migration/controller/command"
	"github.com/kodefluence/monorepo/db"
)

//go:embed mysql/*.sql
var fs embed.FS

func LoadCommand(sqldb db.DB, sqldbconfig core.DatabaseConfig, appModule module.App) error {
	appModule.Controller().InjectCommand(
		command.NewMigrateUp(sqldb, sqldbconfig, fs),
		command.NewMigrateDown(sqldb, sqldbconfig, fs),
		command.NewMigrateRollback(sqldb, sqldbconfig, fs),
	)
	return nil
}
