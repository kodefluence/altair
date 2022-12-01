package command

import (
	"embed"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/kodefluence/altair/core"
	"github.com/kodefluence/monorepo/db"
	"github.com/spf13/cobra"
)

type MigrateUp struct {
	sqldb       db.DB
	sqldbconfig core.DatabaseConfig
	fs          embed.FS
}

func NewMigrateUp(sqldb db.DB, sqldbconfig core.DatabaseConfig, fs embed.FS) *MigrateUp {
	return &MigrateUp{sqldb: sqldb, sqldbconfig: sqldbconfig, fs: fs}
}

func (m *MigrateUp) Use() string {
	return "oauth/migrate:up"
}

func (m *MigrateUp) Short() string {
	return "Migrate oauth databases"
}

func (m *MigrateUp) Example() string {
	return "altair plugin oauth/migrate:up"
}

func (m *MigrateUp) Run(cmd *cobra.Command, args []string) {
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

	if err := migrator.Up(); err != nil {
		fmt.Println("error", err)
		return
	}
}
