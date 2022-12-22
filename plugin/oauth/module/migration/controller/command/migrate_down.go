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
	"github.com/spf13/pflag"
)

type MigrateDown struct {
	sqldb       db.DB
	sqldbconfig core.DatabaseConfig
	fs          embed.FS
}

func NewMigrateDown(sqldb db.DB, sqldbconfig core.DatabaseConfig, fs embed.FS) *MigrateDown {
	return &MigrateDown{sqldb: sqldb, sqldbconfig: sqldbconfig, fs: fs}
}

func (m *MigrateDown) Use() string {
	return "oauth/migrate:down"
}

func (m *MigrateDown) Short() string {
	return "Migrate oauth databases down"
}

func (m *MigrateDown) Example() string {
	return "altair plugin oauth/migrate:down"
}

func (m *MigrateDown) Run(cmd *cobra.Command, args []string) {
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

	if err := migrator.Down(); err != nil {
		fmt.Println("error", err)
		return
	}
}

func (m *MigrateDown) ModifyFlags(flags *pflag.FlagSet) {}
