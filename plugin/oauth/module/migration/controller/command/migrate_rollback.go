package command

import (
	"embed"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/kodefluence/altair/core"
	"github.com/kodefluence/monorepo/db"
	"github.com/spf13/cobra"
)

type MigrateRollback struct {
	sqldb       db.DB
	sqldbconfig core.DatabaseConfig
	fs          embed.FS
}

func NewMigrateRollback(sqldb db.DB, sqldbconfig core.DatabaseConfig, fs embed.FS) *MigrateRollback {
	return &MigrateRollback{sqldb: sqldb, sqldbconfig: sqldbconfig, fs: fs}
}

func (m *MigrateRollback) Use() string {
	return "oauth/migrate:rollback"
}

func (m *MigrateRollback) Short() string {
	return "Do a migration rollback from current versions into previous versions."
}

func (m *MigrateRollback) Example() string {
	return "altair plugin oauth/migrate:rollback"
}

func (m *MigrateRollback) Run(cmd *cobra.Command, args []string) {
	dbDriver, err := mysql.WithInstance(m.sqldb.Eject(), &mysql.Config{
		MigrationsTable: "oauth_plugin_db_versions",
		DatabaseName:    m.sqldbconfig.DBDatabase(),
	})
	if err != nil {
		fmt.Println("error", err)
		return
	}

	sourceDriver, err := iofs.New(m.fs, "mysql")
	if err != nil {
		fmt.Println("error", err)
		return
	}

	migrator, err := migrate.NewWithInstance("iofs", sourceDriver, "mysql", dbDriver)
	if err != nil {
		fmt.Println("error", err)
		return
	}

	if err := migrator.Steps(-1); err != nil {
		fmt.Println("error", err)
		return
	}
}
