package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/subosito/gotenv"

	"github.com/codefluence-x/journal"
	"github.com/spf13/cobra"
)

var (
	mysqlDB              *sql.DB
	mysqlConnMaxLifetime time.Duration
	mysqlMaxIdleConn     int
	mysqlMaxOpenConn     int

	migration *migrate.Migrate
)

func main() {
	_ = gotenv.Load()
	executeCommand()
}

func executeCommand() {
	rootCmd := &cobra.Command{
		Use:   "altair",
		Short: "Light Weight and Robust API Gateway.",
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}

	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Run API gateway services.",
		Run: func(cmd *cobra.Command, args []string) {
			fabricateConnection()
		},
	}

	rootCmd.AddCommand(runCmd)
	_ = rootCmd.Execute()
}

func fabricateConnection() error {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true&interpolateParams=true", os.Getenv("DATABASE_USERNAME"), os.Getenv("DATABASE_PASSWORD"), os.Getenv("DATABASE_HOST"), os.Getenv("DATABASE_NAME")))
	if err != nil {
		journal.Error(fmt.Sprintln("Fabricate connection error:", err), err).SetTags("altair", "main").Log()
		return err
	}
	db.SetConnMaxLifetime(mysqlConnMaxLifetime)
	db.SetMaxIdleConns(mysqlMaxIdleConn)
	db.SetMaxOpenConns(mysqlMaxOpenConn)

	mysqlDB = db

	journal.Info(fmt.Sprintf("Complete fabricating mysql writer connection: %s:%s@tcp(%s)/%s?", os.Getenv("DATABASE_USERNAME"), "***********", os.Getenv("DATABASE_HOST"), os.Getenv("DATABASE_NAME"))).SetTags("altair", "main").Log()

	return nil
}

func closeConnection() {
	if mysqlDB != nil {
		var err error
		for i := 0; i < 3; i++ {
			err = mysqlDB.Close()
			if err != nil {
				journal.Error(fmt.Sprintln("Error closing mysql writer because of:", err), err).SetTags("altair", "main").Log()
				continue
			}

			if err == nil {
				break
			}
		}
		if err != nil {
			journal.Info("Failed closing mysql writer and reader connection.").SetTags("altair", "main").Log()
			return
		}

		journal.Info("Success closing mysql writer and reader connection.").SetTags("altair", "main").Log()
	}
}

func fabricateMigration() error {
	driver, err := mysql.WithInstance(mysqlDB, &mysql.Config{
		MigrationsTable: "db_versions",
		DatabaseName:    os.Getenv("DATABASE_NAME"),
	})
	if err != nil {
		journal.Error(fmt.Sprintln("Fabricate migration error:", err), err).SetTags("levelup", "main").Log()
		return err
	}

	m, err := migrate.NewWithDatabaseInstance("file://migration", "mysql", driver)
	if err != nil {
		journal.Error(fmt.Sprintln("Fabricate migration error:", err), err).SetTags("levelup", "main").Log()
		return err
	}
	migration = m

	return nil
}

func closeMigration() {
	if migration != nil {
		s, err := migration.Close()
		if err != nil {
			journal.Error(fmt.Sprintln("Close migration error:", err), err).SetTags("levelup", "main").Log()
			journal.Error(fmt.Sprintln("Source:", s), s).SetTags("levelup", "main").Log()
		}
		journal.Info("Success closing migration.").SetTags("levelup", "main").Log()
	}
}
