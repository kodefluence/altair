package cmd

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/subosito/gotenv"

	"github.com/codefluence-x/altair/core"
	"github.com/codefluence-x/altair/loader"
	"github.com/codefluence-x/journal"
)

var (
	dbConfigs    map[string]core.DatabaseConfig = map[string]core.DatabaseConfig{}
	databases    map[string]*sql.DB             = map[string]*sql.DB{}
	appConfig    core.AppConfig
	pluginBearer core.PluginBearer
	apiEngine    *gin.Engine
	RootCmd      = &cobra.Command{
		Use:   "altair",
		Short: "Light Weight and Robust API Gateway.",
	}
)

func init() {
	_ = gotenv.Load()

	cobra.OnInitialize(loadConfig)

	RootCmd.AddCommand(MigrateCmd())
	RootCmd.AddCommand(MigrateDownCmd())
	RootCmd.AddCommand(MigrateRollbackCmd())
	RootCmd.AddCommand(ServerCmd())
}

func loadConfig() {
	var err error

	loadedDBConfigs, err := loader.Database().Compile("../config/database.yml")
	if err != nil {
		journal.Error("Error loading database config", err).Log()
		os.Exit(1)
	}
	dbConfigs = loadedDBConfigs

	loadedAppConfig, err := loader.App().Compile("../config/app.yml")
	if err != nil {
		journal.Error("Error loading app config", err).Log()
		os.Exit(1)
	}
	appConfig = loadedAppConfig

	loadedPluginBearer, err := loader.Plugin().Compile("../config/plugin/")
	if err != nil {
		journal.Error("Error loading plugin config", err).Log()
		os.Exit(1)
	}
	pluginBearer = loadedPluginBearer
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
