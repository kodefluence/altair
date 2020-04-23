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

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/subosito/gotenv"

	"github.com/codefluence-x/altair/controller"
	"github.com/codefluence-x/altair/core"
	"github.com/codefluence-x/altair/formatter"
	"github.com/codefluence-x/altair/forwarder"
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

	migration *migrate.Migrate

	apiEngine *gin.Engine

	accessTokenTimeout time.Duration
	accessGrantTimeout time.Duration
)

func main() {
	_ = gotenv.Load()
	loadConfig()
	executeCommand()
	shutdownFunc()
}

func shutdownFunc() {
	closeConnection()
	closeMigration()
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
			if err := fabricateConnection(); err != nil {
				journal.Error("Error running altair:", err).SetTags("altair", "main").Log()
				return
			}

			runAPI()
		},
	}

	migrateCmd := &cobra.Command{
		Use:   "migrate",
		Short: "Do a migration.",
		Run: func(cmd *cobra.Command, args []string) {
			if err := fabricateConnection(); err != nil {
				return
			}

			if err := fabricateMigration(); err != nil {
				return
			}

			if err := migration.Up(); err != nil && err.Error() != "no change" {
				journal.Error(fmt.Sprintln("Migration error because of:", err), err).SetTags("altair", "main").Log()
				return
			}

			journal.Info("Migration migrate process is complete").SetTags("altair", "main").Log()
		},
	}

	migrateDownCmd := &cobra.Command{
		Use:   "migrate:down",
		Short: "Down the migration.",
		Run: func(cmd *cobra.Command, args []string) {
			if err := fabricateConnection(); err != nil {
				return
			}

			if err := fabricateMigration(); err != nil {
				return
			}

			if err := migration.Down(); err != nil && err.Error() != "no change" {
				journal.Error(fmt.Sprintln("Migration error because of:", err), err).SetTags("altair", "main").Log()
				return
			}

			journal.Info("Migration down process is complete").SetTags("altair", "main").Log()
		},
	}

	migrateRollbackCmd := &cobra.Command{
		Use:   "migrate:rollback",
		Short: "Rollback the migration.",
		Run: func(cmd *cobra.Command, args []string) {
			if err := fabricateConnection(); err != nil {
				return
			}

			if err := fabricateMigration(); err != nil {
				return
			}

			if err := migration.Steps(-1); err != nil && err.Error() != "no change" {
				journal.Error(fmt.Sprintln("Migration error because of:", err), err).SetTags("altair", "main").Log()
				return
			}

			journal.Info("Migration rollback one step process is complete").SetTags("altair", "main").Log()
		},
	}

	rootCmd.AddCommand(runCmd, migrateCmd, migrateDownCmd, migrateRollbackCmd)
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
		journal.Error(fmt.Sprintln("Fabricate migration error:", err), err).SetTags("altair", "main").Log()
		return err
	}

	m, err := migrate.NewWithDatabaseInstance("file://migration", "mysql", driver)
	if err != nil {
		journal.Error(fmt.Sprintln("Fabricate migration error:", err), err).SetTags("altair", "main").Log()
		return err
	}
	migration = m

	return nil
}

func closeMigration() {
	if migration != nil {
		s, err := migration.Close()
		if err != nil {
			journal.Error(fmt.Sprintln("Close migration error:", err), err).SetTags("altair", "main").Log()
			journal.Error(fmt.Sprintln("Source:", s), s).SetTags("altair", "main").Log()
		}
		journal.Info("Success closing migration.").SetTags("altair", "main").Log()
	}
}

func runAPI() {
	gin.SetMode(gin.ReleaseMode)

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

	// Route Engine
	routeCompiler := forwarder.Route().Compiler()
	routeObjects, err := routeCompiler.Compile("./routes")
	if err != nil {
		journal.Error("Error compiling routes", err).
			SetTags("altair", "main").
			Log()
		os.Exit(1)
	}

	apiEngine = gin.New()
	apiEngine.GET("/health", controller.Health)

	// DownStream Plugin
	oauthPlugin := plugin.DownStream().Oauth(oauthAccessTokenModel)

	err = forwarder.Route().Generator().Generate(apiEngine, routeObjects, []core.DownStreamPlugin{oauthPlugin})
	if err != nil {
		journal.Error("Error generating routes", err).
			SetTags("altair", "main").
			Log()
		os.Exit(1)
	}

	internalEngine := apiEngine.Group("/_plugins/", gin.BasicAuth(gin.Accounts{
		os.Getenv("BASIC_AUTH_USERNAME"): os.Getenv("BASIC_AUTH_PASSWORD"),
	}))

	controller.Compile(internalEngine, controller.Oauth().Application().List(applicationManager))
	controller.Compile(internalEngine, controller.Oauth().Application().One(applicationManager))
	controller.Compile(internalEngine, controller.Oauth().Application().Create(applicationManager))
	controller.Compile(internalEngine, controller.Oauth().Authorization().Grant(authorization))

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
