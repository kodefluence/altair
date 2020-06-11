package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"

	_ "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/subosito/gotenv"

	"github.com/codefluence-x/altair/controller"
	"github.com/codefluence-x/altair/core"
	"github.com/codefluence-x/altair/forwarder"
	"github.com/codefluence-x/altair/loader"
	"github.com/codefluence-x/altair/provider"
	"github.com/codefluence-x/journal"
	"github.com/spf13/cobra"
)

var (
	dbConfigs    map[string]core.DatabaseConfig = map[string]core.DatabaseConfig{}
	databases    map[string]*sql.DB             = map[string]*sql.DB{}
	appConfig    core.AppConfig
	pluginBearer core.PluginBearer
	apiEngine    *gin.Engine
)

func main() {
	_ = gotenv.Load()
	loadConfig()
	executeCommand()
}

func loadConfig() {
	var err error

	loadedDBConfigs, err := loader.Database().Compile("./config/database.yml")
	if err != nil {
		journal.Error("Error loading database config", err).Log()
		os.Exit(1)
	}
	dbConfigs = loadedDBConfigs

	loadedAppConfig, err := loader.App().Compile("./config/app.yml")
	if err != nil {
		journal.Error("Error loading app config", err).Log()
		os.Exit(1)
	}
	appConfig = loadedAppConfig

	loadedPluginBearer, err := loader.Plugin().Compile("./config/plugin/")
	if err != nil {
		journal.Error("Error loading plugin config", err).Log()
		os.Exit(1)
	}
	pluginBearer = loadedPluginBearer
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
			defer closeConnection()
			if err := fabricateConnection(); err != nil {
				journal.Error("Error running altair:", err).SetTags("altair", "main").Log()
				return
			}

			if err := runAPI(); err != nil {
				journal.Error("Error running altair API:", err).SetTags("altair", "main").Log()
			}
		},
	}

	configCmd := &cobra.Command{
		Use:     "config",
		Short:   "See list of configs",
		Example: "altair config app",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				fmt.Println("Invalid number of arguments, expected 1. Example `altair config [config_name]`.")
				fmt.Println("Available option:")
				fmt.Println("- all")
				fmt.Println("- app")
				fmt.Println("- db")
				return
			}

			app := func() {
				fmt.Printf("app config:\n")
				fmt.Printf("====================\n")
				fmt.Printf(appConfig.Dump())
				fmt.Printf("--------------------\n")
			}

			db := func() {
				fmt.Printf("db config:\n")
				fmt.Printf("====================\n")
				for key, config := range dbConfigs {
					fmt.Printf("instance: %s\n", key)
					fmt.Printf("driver: %s\n", config.Driver())
					fmt.Printf("--------------------\n")
					fmt.Printf(config.Dump())
				}
				fmt.Printf("--------------------\n")
			}

			switch args[0] {
			case "all":
				app()
				fmt.Println()
				db()
			case "app":
				app()
			case "db":
				db()
			default:
				fmt.Println("Invalid argument. Available: [app, db, all]")
				return
			}
		},
	}

	migrateCmd := &cobra.Command{
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

	migrateDownCmd := &cobra.Command{
		Use:   "migrate:down",
		Short: "Down the migration from current version into earlies versions.",
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

	migrateRollbackCmd := &cobra.Command{
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

	rootCmd.AddCommand(runCmd, migrateCmd, migrateDownCmd, migrateRollbackCmd, configCmd)
	_ = rootCmd.Execute()
}

func dbConnectionFabricator(dbConfig core.DatabaseConfig) (*sql.DB, error) {
	port, err := dbConfig.DBPort()
	if err != nil {
		return nil, err
	}

	maxConnLifetime, err := dbConfig.DBConnectionMaxLifetime()
	if err != nil {
		return nil, err
	}

	db, err := sql.Open(dbConfig.Driver(), fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&interpolateParams=true", dbConfig.DBUsername(), dbConfig.DBPassword(), dbConfig.DBHost(), port, dbConfig.DBDatabase()))
	if err != nil {
		return nil, err
	}

	db.SetConnMaxLifetime(maxConnLifetime)

	journal.Info(fmt.Sprintf("Complete fabricating mysql writer connection: %s:%s@tcp(%s:%d)/%s?", dbConfig.DBUsername(), "***********", dbConfig.DBHost(), port, dbConfig.DBDatabase())).SetTags("altair", "main").Log()

	return db, nil
}

func fabricateConnection() error {
	for key, config := range dbConfigs {
		sqlDB, err := dbConnectionFabricator(config)
		if err != nil {
			return err
		}

		databases[key] = sqlDB
	}

	return nil
}

func closeConnection() {
	for dbName, db := range databases {
		var err error

		for i := 0; i < 3; i++ {
			err = db.Close()
			if err != nil {
				journal.Error(fmt.Sprintln("Error closing mysql writer because of:", err), err).SetTags("altair", "main", dbName).Log()
				continue
			}

			if err == nil {
				break
			}
		}
		if err != nil {
			journal.Info("Failed closing mysql writer and reader connection.").SetTags("altair", "main", dbName).Log()
			return
		}

		journal.Info("Success closing mysql writer and reader connection.").SetTags("altair", "main", dbName).Log()
	}
}

func runAPI() error {
	gin.SetMode(gin.ReleaseMode)

	apiEngine = gin.New()
	apiEngine.GET("/health", controller.Health)

	internalEngine := apiEngine.Group("/_plugins/", gin.BasicAuth(gin.Accounts{
		appConfig.BasicAuthUsername(): appConfig.BasicAuthPassword(),
	}))

	appBearer := loader.AppBearer(internalEngine, appConfig)
	dbBearer := loader.DatabaseBearer(databases, dbConfigs)

	provider.Plugin(appBearer, dbBearer, pluginBearer)
	provider.Metric(appBearer)

	// Route Engine
	routeCompiler := forwarder.Route().Compiler()
	routeObjects, err := routeCompiler.Compile("./routes")
	if err != nil {
		journal.Error("Error compiling routes", err).
			SetTags("altair", "main").
			Log()
		return err
	}

	err = forwarder.Route().Generator().Generate(apiEngine, routeObjects, []core.DownStreamPlugin{})
	if err != nil {
		journal.Error("Error generating routes", err).
			SetTags("altair", "main").
			Log()
		return err
	}

	gracefulSignal := make(chan os.Signal, 1)
	signal.Notify(gracefulSignal, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		srv := &http.Server{
			Addr:    fmt.Sprintf(":%d", appConfig.Port()),
			Handler: apiEngine,
		}

		journal.Info(fmt.Sprintf("Running Altair in: %d", appConfig.Port())).Log()

		if err := srv.ListenAndServe(); err != nil {
			journal.Error("Error running api engine", err).
				SetTags("altair", "main").
				Log()
		}
	}()

	closeSignal := <-gracefulSignal
	journal.Info(fmt.Sprintf("Receiving %s signal.... Cleaning up processes.", closeSignal.String())).Log()
	return nil
}
