package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"

	_ "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/subosito/gotenv"

	"github.com/codefluence-x/altair/controller"
	"github.com/codefluence-x/altair/core"
	"github.com/codefluence-x/altair/formatter"
	"github.com/codefluence-x/altair/forwarder"
	"github.com/codefluence-x/altair/loader"
	"github.com/codefluence-x/altair/model"
	"github.com/codefluence-x/altair/plugin"
	"github.com/codefluence-x/altair/service"
	"github.com/codefluence-x/altair/validator"
	"github.com/codefluence-x/journal"
	"github.com/spf13/cobra"
)

var (
	mysqlDB              *sql.DB
	mysqlConnMaxLifetime time.Duration
	mysqlMaxIdleConn     int
	mysqlMaxOpenConn     int

	dbConfigs map[string]core.DatabaseConfig = map[string]core.DatabaseConfig{}
	databases map[string]*sql.DB             = map[string]*sql.DB{}

	appConfig core.AppConfig

	pluginBearer core.PluginBearer

	apiEngine *gin.Engine

	accessTokenTimeout time.Duration
	accessGrantTimeout time.Duration
)

func main() {
	_ = gotenv.Load()
	loadConfig()
	executeCommand()
}

func loadConfig() {
	var err error
	accessTokenTimeout, err = time.ParseDuration(os.Getenv("ACCESS_TOKEN_TIMEOUT"))
	if err != nil {
		os.Exit(1)
	}

	accessGrantTimeout, err = time.ParseDuration(os.Getenv("ACCESS_GRANT_TIMEOUT"))
	if err != nil {
		os.Exit(1)
	}

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

			runAPI()
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
		},
	}

	rootCmd.AddCommand(runCmd, migrateCmd, configCmd)
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
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&interpolateParams=true", os.Getenv("DATABASE_USERNAME"), os.Getenv("DATABASE_PASSWORD"), os.Getenv("DATABASE_HOST"), os.Getenv("DATABASE_PORT"), os.Getenv("DATABASE_NAME")))
	if err != nil {
		journal.Error(fmt.Sprintln("Fabricate connection error:", err), err).SetTags("altair", "main").Log()
		return err
	}
	db.SetConnMaxLifetime(mysqlConnMaxLifetime)
	db.SetMaxIdleConns(mysqlMaxIdleConn)
	db.SetMaxOpenConns(mysqlMaxOpenConn)

	mysqlDB = db

	journal.Info(fmt.Sprintf("Complete fabricating mysql writer connection: %s:%s@tcp(%s:%s)/%s?", os.Getenv("DATABASE_USERNAME"), "***********", os.Getenv("DATABASE_HOST"), os.Getenv("DATABASE_PORT"), os.Getenv("DATABASE_NAME"))).SetTags("altair", "main").Log()

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

// func fabricateMigration() error {
// 	driver, err := mysql.WithInstance(mysqlDB, &mysql.Config{
// 		MigrationsTable: "db_versions",
// 		DatabaseName:    os.Getenv("DATABASE_NAME"),
// 	})
// 	if err != nil {
// 		journal.Error(fmt.Sprintln("Fabricate migration error:", err), err).SetTags("altair", "main").Log()
// 		return err
// 	}

// 	m, err := migrate.NewWithDatabaseInstance("file://migration", "mysql", driver)
// 	if err != nil {
// 		journal.Error(fmt.Sprintln("Fabricate migration error:", err), err).SetTags("altair", "main").Log()
// 		return err
// 	}
// 	migration = m

// 	return nil
// }

// func closeMigration() {
// 	if migration != nil {
// 		s, err := migration.Close()
// 		if err != nil {
// 			journal.Error(fmt.Sprintln("Close migration error:", err), err).SetTags("altair", "main").Log()
// 			journal.Error(fmt.Sprintln("Source:", s), s).SetTags("altair", "main").Log()
// 		}
// 		journal.Info("Success closing migration.").SetTags("altair", "main").Log()
// 	}
// }

func runAPI() {
	gin.SetMode(gin.ReleaseMode)

	dbBearer := loader.DatabaseBearer(databases, dbConfigs)
	// TODO: DELETE DEBUG
	fmt.Println("DB Bearer", dbBearer)

	apiEngine = gin.New()
	apiEngine.GET("/health", controller.Health)

	downStreamPlugins := []core.DownStreamPlugin{}

	internalEngine := apiEngine.Group("/_plugins/", gin.BasicAuth(gin.Accounts{
		appConfig.BasicAuthUsername(): appConfig.BasicAuthPassword(),
	}))

	if appConfig.PluginExists("oauth") {
		// Model
		oauthApplicationModel := model.OauthApplication(mysqlDB)
		oauthAccessTokenModel := model.OauthAccessToken(mysqlDB)
		oauthAccessGrantModel := model.OauthAccessGrant(mysqlDB)

		// Formatter
		applicationFormatter := formatter.OauthApplication()
		modelFormatter := formatter.Model(accessTokenTimeout, accessGrantTimeout)
		oauthFormatter := formatter.Oauth()

		// Validator
		oauthValidator := validator.Oauth()

		// Service
		applicationManager := service.ApplicationManager(applicationFormatter, modelFormatter, oauthApplicationModel, oauthValidator)
		authorization := service.Authorization(oauthApplicationModel, oauthAccessTokenModel, oauthAccessGrantModel, modelFormatter, oauthValidator, oauthFormatter)

		downStreamPlugins = append(downStreamPlugins, plugin.DownStream().Oauth(oauthAccessTokenModel))

		controller.Compile(internalEngine, controller.Oauth().Application().List(applicationManager))
		controller.Compile(internalEngine, controller.Oauth().Application().One(applicationManager))
		controller.Compile(internalEngine, controller.Oauth().Application().Create(applicationManager))
		controller.Compile(internalEngine, controller.Oauth().Authorization().Grant(authorization))
		controller.Compile(internalEngine, controller.Oauth().Authorization().Revoke(authorization))
	}

	// Route Engine
	routeCompiler := forwarder.Route().Compiler()
	routeObjects, err := routeCompiler.Compile("./routes")
	if err != nil {
		journal.Error("Error compiling routes", err).
			SetTags("altair", "main").
			Log()
		os.Exit(1)
	}

	err = forwarder.Route().Generator().Generate(apiEngine, routeObjects, downStreamPlugins)
	if err != nil {
		journal.Error("Error generating routes", err).
			SetTags("altair", "main").
			Log()
		os.Exit(1)
	}

	gracefulSignal := make(chan os.Signal, 1)
	signal.Notify(gracefulSignal, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		srv := &http.Server{
			Addr:    ":" + os.Getenv("APP_PORT"),
			Handler: apiEngine,
		}

		if err := srv.ListenAndServe(); err != nil {
			journal.Error("Error running api engine", err).
				SetTags("altair", "main").
				Log()
		}
	}()

	closeSignal := <-gracefulSignal
	journal.Info(fmt.Sprintf("Receiving %s signal.... Cleaning up processes.", closeSignal.String())).Log()
}
