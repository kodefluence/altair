package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/codefluence-x/altair/loader"
	"github.com/codefluence-x/altair/provider"
	"github.com/codefluence-x/journal"
)

func MigrateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "migrate",
		Short: "Do a migration from current version into latest versions.",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				fmt.Println("Invalid number of arguments, expected 1. Example `altair migrate [database_instance_name]`.")
				fmt.Println("To see available database instance, use `altair config db`.")
				return
			}

			defer closeConnection()
			if err := fabricateConnection(); err != nil {
				return
			}

			dbBearer := loader.DatabaseBearer(databases, dbConfigs)
			db, config, err := dbBearer.Database(args[0])
			if err != nil {
				journal.Error(fmt.Sprintf("Error loading database instance of: `%s`", args[0]), err).SetTags("altair", "main").Log()
				return
			}

			migrationProvider := provider.Migration().GoMigrate(db, config)
			migrator, err := migrationProvider.Migrator()
			if err != nil {
				journal.Error("Error providing migrator", err).SetTags("altair", "main").Log()
				return
			}
			defer migrator.Close()

			if err := migrator.Up(); err != nil && err.Error() != "no change" {
				journal.Error("Error doing database migration", err).SetTags("altair", "main").Log()
				return
			}

			journal.Info("Successfully migrating databases").SetTags("altair", "main").Log()
		},
	}
}

func MigrateDownCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "migrate:down",
		Short: "Down the migration from current version into earliest versions.",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				fmt.Println("Invalid number of arguments, expected 1. Example `altair migrate:down [database_instance_name]`.")
				fmt.Println("To see available database instance, use `altair config db`.")
				return
			}

			defer closeConnection()
			if err := fabricateConnection(); err != nil {
				return
			}

			dbBearer := loader.DatabaseBearer(databases, dbConfigs)
			db, config, err := dbBearer.Database(args[0])
			if err != nil {
				journal.Error(fmt.Sprintf("Error loading database instance of: `%s`", args[0]), err).SetTags("altair", "main").Log()
				return
			}

			migrationProvider := provider.Migration().GoMigrate(db, config)
			migrator, err := migrationProvider.Migrator()
			if err != nil {
				journal.Error("Error providing migrator", err).SetTags("altair", "main").Log()
				return
			}
			defer migrator.Close()

			if err := migrator.Down(); err != nil && err.Error() != "no change" {
				journal.Error("Error doing database migration", err).SetTags("altair", "main").Log()
				return
			}

			journal.Info("Successfully migrating databases").SetTags("altair", "main").Log()
		},
	}
}

func MigrateRollbackCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "migrate:rollback",
		Short: "Do a migration rollback from current versions into previous versions.",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				fmt.Println("Invalid number of arguments, expected 1. Example `altair migrate:rollback [database_instance_name]`.")
				fmt.Println("To see available database instance, use `altair config db`.")
				return
			}

			defer closeConnection()
			if err := fabricateConnection(); err != nil {
				return
			}

			dbBearer := loader.DatabaseBearer(databases, dbConfigs)
			db, config, err := dbBearer.Database(args[0])
			if err != nil {
				journal.Error(fmt.Sprintf("Error loading database instance of: `%s`", args[0]), err).SetTags("altair", "main").Log()
				return
			}

			migrationProvider := provider.Migration().GoMigrate(db, config)
			migrator, err := migrationProvider.Migrator()
			if err != nil {
				journal.Error("Error providing migrator", err).SetTags("altair", "main").Log()
				return
			}
			defer migrator.Close()

			if err := migrator.Steps(-1); err != nil && err.Error() != "no change" {
				journal.Error("Error doing database migration", err).SetTags("altair", "main").Log()
				return
			}

			journal.Info("Successfully migrating databases").SetTags("altair", "main").Log()
		},
	}
}
