package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	_ "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/subosito/gotenv"

	"github.com/kodefluence/altair/cfg"
	"github.com/kodefluence/altair/core"

	"github.com/kodefluence/altair/module/apierror"
	"github.com/kodefluence/altair/module/app"
	"github.com/kodefluence/altair/module/controller"
	"github.com/kodefluence/altair/module/healthcheck"
	"github.com/kodefluence/altair/module/projectgenerator"
	"github.com/kodefluence/altair/module/router"
	"github.com/kodefluence/altair/plugin"
	"github.com/kodefluence/monorepo/db"
	"github.com/spf13/cobra"
)

var (
	dbConfigs    map[string]core.DatabaseConfig = map[string]core.DatabaseConfig{}
	databases    map[string]db.DB               = map[string]db.DB{}
	appConfig    core.AppConfig
	pluginBearer core.PluginBearer
)

func main() {
	_ = gotenv.Load()
	loadConfig()
	executeCommand()
}

func loadConfig() {
	dbConfigs, _ = cfg.Database().Compile("config/database.yml")
	appConfig, _ = cfg.App().Compile("config/app.yml")
	pluginBearer, _ = cfg.Plugin().Compile("config/plugin/")
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
			if appConfig == nil {
				fmt.Println("App config is not loaded, only run command in altair working directory")
				return
			}

			defer closeConnection()
			if err := fabricateConnection(); err != nil {
				log.Error().
					Err(err).
					Stack().
					Array("tags", zerolog.Arr().Str("altair").Str("main")).
					Msg("Error running altair")
				return
			}

			if err := runAPI(); err != nil {
				log.Error().
					Err(err).
					Stack().
					Array("tags", zerolog.Arr().Str("altair").Str("main")).
					Msg("Error running altair API")
			}
		},
	}

	configCmd := &cobra.Command{
		Use:     "config",
		Short:   "See list of configs",
		Example: "altair config app",
		Run: func(cmd *cobra.Command, args []string) {
			if appConfig == nil {
				fmt.Println("App config is not loaded, only run command in altair working directory")
				return
			}

			if len(args) < 1 {
				fmt.Println("Invalid number of arguments, expected 1. Example `altair config [config_name]`.")
				fmt.Println("Available option:")
				fmt.Println("- all")
				fmt.Println("- app")
				fmt.Println("- db")
				return
			}

			app := func() {
				fmt.Println("app config:")
				fmt.Println("====================")
				fmt.Print(appConfig.Dump())
				fmt.Println("--------------------")
			}

			db := func() {
				fmt.Println("db config:")
				fmt.Println("====================")
				for key, config := range dbConfigs {
					fmt.Printf("instance: %s\n", key)
					fmt.Printf("driver: %s\n", config.Driver())
					fmt.Println("--------------------")
					fmt.Print(config.Dump())
				}
				fmt.Println("--------------------")
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

	appController := controller.Provide(nil, nil, rootCmd)
	appModule := app.Provide(appController)
	projectgenerator.Load(appModule)

	pluginCmd := &cobra.Command{
		Use:                "plugin",
		Short:              "List of plugin commands",
		DisableFlagParsing: true,
		Run: func(cmd *cobra.Command, args []string) {
			if appConfig == nil {
				fmt.Println("App config is not loaded, only run command in altair working directory")
				return
			}

			defer closeConnection()
			if err := fabricateConnection(); err != nil {
				log.Error().
					Err(err).
					Stack().
					Array("tags", zerolog.Arr().Str("altair").Str("main")).
					Msg("Error running altair")
				return
			}

			apiError := apierror.Provide()

			pluginController := controller.Provide(nil, apiError, cmd)
			pluginModule := app.Provide(pluginController)

			appBearer := cfg.AppBearer(nil, appConfig)
			dbBearer := cfg.DatabaseBearer(databases, dbConfigs)

			if err := plugin.LoadCommand(appBearer, pluginBearer, dbBearer, apiError, pluginModule); err != nil {
				log.Error().
					Err(err).
					Stack().
					Array("tags", zerolog.Arr().Str("altair").Str("main")).
					Msg("Error generating plugins")
			}

			childCmd, _, err := cmd.Find(args)
			if err != nil || childCmd.Use == cmd.Use {
				_ = cmd.Help()
			} else {
				_ = childCmd.Execute()
			}
		},
	}

	rootCmd.AddCommand(runCmd, pluginCmd, configCmd)

	_ = rootCmd.Execute()
}

func dbConnectionFabricator(instanceName string, dbConfig core.DatabaseConfig) (db.DB, error) {
	port, err := dbConfig.DBPort()
	if err != nil {
		return nil, err
	}

	maxConnLifetime, err := dbConfig.DBConnectionMaxLifetime()
	if err != nil {
		return nil, err
	}

	maxIdleConn, err := dbConfig.DBMaxIddleConn()
	if err != nil {
		return nil, err
	}

	maxOpenConn, err := dbConfig.DBMaxOpenConn()
	if err != nil {
		return nil, err
	}

	sqldb, err := db.FabricateMySQL(instanceName, db.Config{
		Username: dbConfig.DBUsername(),
		Password: dbConfig.DBPassword(),
		Host:     dbConfig.DBHost(),
		Port:     fmt.Sprintf("%d", port),
		Name:     dbConfig.DBDatabase(),
	}, db.WithConnMaxLifetime(maxConnLifetime), db.WithMaxIdleConn(maxIdleConn), db.WithMaxOpenConn(maxOpenConn))
	if err != nil {
		return nil, err
	}

	log.Info().Msg(fmt.Sprintf("Complete fabricating mysql writer connection: %s:%s@tcp(%s:%d)/%s?", dbConfig.DBUsername(), "***********", dbConfig.DBHost(), port, dbConfig.DBDatabase()))

	return sqldb, nil
}

func fabricateConnection() error {
	for key, config := range dbConfigs {
		sqlDB, err := dbConnectionFabricator(key, config)
		if err != nil {
			return err
		}

		databases[key] = sqlDB
	}

	return nil
}

func closeConnection() {
	excs := db.CloseAll()
	for _, exc := range excs {
		log.Error().
			Err(exc).
			Stack().
			Str("detail", exc.Detail()).
			Array("tags", zerolog.Arr().Str("altair").Str("main").Str(exc.Title())).
			Msg("Error closing mysql writer")
	}
}

func runAPI() error {
	gin.SetMode(gin.ReleaseMode)

	apiEngine := gin.New()
	apiError := apierror.Provide()

	baseController := controller.Provide(apiEngine.Handle, apiError, nil)
	baseModule := app.Provide(baseController)
	healthcheck.Load(baseModule)

	pluginEngine := apiEngine.Group("/_plugins/", gin.BasicAuth(gin.Accounts{
		appConfig.BasicAuthUsername(): appConfig.BasicAuthPassword(),
	}))
	pluginController := controller.Provide(pluginEngine.Handle, apiError, nil)
	pluginModule := app.Provide(pluginController)

	appBearer := cfg.AppBearer(pluginEngine, appConfig)
	dbBearer := cfg.DatabaseBearer(databases, dbConfigs)

	if err := plugin.Load(appBearer, pluginBearer, dbBearer, apiError, pluginModule); err != nil {
		log.Error().
			Err(err).
			Stack().
			Array("tags", zerolog.Arr().Str("altair").Str("main")).
			Msg("Error generating plugins")
		return err
	}

	compiler, forwarder := router.Provide(pluginModule.Controller().ListDownstream(), pluginModule.Controller().ListMetric())
	routeObjects, err := compiler.Compile("./routes")
	if err != nil {
		log.Error().
			Err(err).
			Stack().
			Array("tags", zerolog.Arr().Str("altair").Str("main")).
			Msg("Error compiling routes")
		return err
	}

	err = forwarder.Generate(apiEngine, routeObjects)
	if err != nil {
		log.Error().
			Err(err).
			Stack().
			Array("tags", zerolog.Arr().Str("altair").Str("main")).
			Msg("Error generating routes")
		return err
	}

	gracefulSignal := make(chan os.Signal, 1)
	signal.Notify(gracefulSignal, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		srv := &http.Server{
			Addr:    fmt.Sprintf(":%d", appConfig.Port()),
			Handler: apiEngine,
		}

		log.Info().Msg(fmt.Sprintf("Running Altair in: http://127.0.0.1:%d", appConfig.Port()))

		if err := srv.ListenAndServe(); err != nil {
			log.Error().
				Err(err).
				Stack().
				Array("tags", zerolog.Arr().Str("altair").Str("main")).
				Msg("Error running api engine")
		}
	}()

	closeSignal := <-gracefulSignal
	log.Info().Array("tags", zerolog.Arr().Str("altair").Str("main")).Msg(fmt.Sprintf("Receiving %s signal.... Cleaning up processes.", closeSignal.String()))
	return nil
}
