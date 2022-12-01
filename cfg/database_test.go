package cfg_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/kodefluence/altair/cfg"
	"github.com/kodefluence/altair/entity"
	"github.com/kodefluence/altair/testhelper"
	"github.com/stretchr/testify/assert"
)

func TestDatabase(t *testing.T) {

	t.Run("Compile", func(t *testing.T) {
		t.Run("Given config path as parameter", func(t *testing.T) {
			t.Run("General scenario", func(t *testing.T) {
				t.Run("Normal scenario", func(t *testing.T) {
					t.Run("Run gracefully", func(t *testing.T) {
						dbName := "db_1"
						dbUsername := "username_db1"
						dbPassword := "password_db1"
						dbHost := "127.0.0.1"
						dbPort := "3306"

						os.Setenv("DATABASE_NAME_DB_CONFIG_NORMAL_SCENARIO", dbName)
						os.Setenv("DATABASE_USERNAME_DB_CONFIG_NORMAL_SCENARIO", dbUsername)
						os.Setenv("DATABASE_PASSWORD_DB_CONFIG_NORMAL_SCENARIO", dbPassword)
						os.Setenv("DATABASE_HOST_DB_CONFIG_NORMAL_SCENARIO", dbHost)
						os.Setenv("DATABASE_PORT_DB_CONFIG_NORMAL_SCENARIO", dbPort)

						configPath := "./db_normal_config/"
						fileName := "database.yml"

						expectedMYSQLConfig := entity.MYSQLDatabaseConfig{
							Database:              dbName,
							Username:              dbUsername,
							Password:              dbPassword,
							Host:                  dbHost,
							Port:                  dbPort,
							ConnectionMaxLifetime: "120s",
							MaxIddleConnection:    "100",
							MaxOpenConnection:     "100",
						}

						testhelper.GenerateTempTestFiles(configPath, DatabaseConfigNormalScenario, fileName, 0666)

						dbConfigs, err := cfg.Database().Compile(fmt.Sprintf("%s%s", configPath, fileName))
						assert.Nil(t, err)

						c, ok := dbConfigs["oauth_database"]

						assert.True(t, ok)
						assert.Equal(t, expectedMYSQLConfig.Driver(), c.Driver())
						assert.Equal(t, expectedMYSQLConfig.DBHost(), c.DBHost())

						expectedDBPort, _ := expectedMYSQLConfig.DBPort()
						actualDBPort, err := c.DBPort()
						assert.Nil(t, err)
						assert.Equal(t, expectedDBPort, actualDBPort)

						assert.Equal(t, expectedMYSQLConfig.DBUsername(), c.DBUsername())
						assert.Equal(t, expectedMYSQLConfig.DBPassword(), c.DBPassword())
						assert.Equal(t, expectedMYSQLConfig.DBDatabase(), c.DBDatabase())

						expectedMaxConnLifetime, _ := expectedMYSQLConfig.DBConnectionMaxLifetime()
						actualMaxConnLifetime, err := c.DBConnectionMaxLifetime()

						assert.Equal(t, expectedMaxConnLifetime, actualMaxConnLifetime)

						expectedMaxIddleConn, _ := expectedMYSQLConfig.DBMaxIddleConn()
						actualMaxIddleConn, err := c.DBMaxIddleConn()
						assert.Nil(t, err)
						assert.Equal(t, expectedMaxIddleConn, actualMaxIddleConn)

						expectedMaxOpenConn, _ := expectedMYSQLConfig.DBMaxOpenConn()
						actualMaxOpenConn, err := c.DBMaxOpenConn()
						assert.Nil(t, err)
						assert.Equal(t, expectedMaxOpenConn, actualMaxOpenConn)

						testhelper.RemoveTempTestFiles(configPath)
					})
				})

				t.Run("Some env not found", func(t *testing.T) {
					t.Run("Run gracefully", func(t *testing.T) {
						dbName := "db_1"
						dbUsername := "username_db1"
						dbPassword := "password_db1"
						dbHost := "127.0.0.1"
						dbPort := "3306"

						os.Setenv("DATABASE_NAME_DB_CONFIG_ENV_CONFIG_NOT_FOUND", dbName)
						os.Setenv("DATABASE_USERNAME_DB_CONFIG_ENV_CONFIG_NOT_FOUND", dbUsername)
						os.Setenv("DATABASE_PASSWORD_DB_CONFIG_ENV_CONFIG_NOT_FOUND", dbPassword)
						os.Setenv("DATABASE_HOST_DB_CONFIG_ENV_CONFIG_NOT_FOUND", dbHost)
						os.Setenv("DATABASE_PORT_DB_CONFIG_ENV_CONFIG_NOT_FOUND", dbPort)

						configPath := "./db_normal_but_some_env_not_found/"
						fileName := "database.yml"

						expectedMYSQLConfig := entity.MYSQLDatabaseConfig{
							Database:              dbName,
							Username:              dbUsername,
							Password:              dbPassword,
							Host:                  dbHost,
							Port:                  dbPort,
							ConnectionMaxLifetime: "0",
							MaxIddleConnection:    "100",
							MaxOpenConnection:     "100",
						}

						testhelper.GenerateTempTestFiles(configPath, DatabaseConfigWithNotFoundENV, fileName, 0666)

						dbConfigs, err := cfg.Database().Compile(fmt.Sprintf("%s%s", configPath, fileName))
						assert.Nil(t, err)

						c, ok := dbConfigs["oauth_database"]

						assert.True(t, ok)
						assert.Equal(t, expectedMYSQLConfig.Driver(), c.Driver())
						assert.Equal(t, expectedMYSQLConfig.DBHost(), c.DBHost())

						expectedDBPort, _ := expectedMYSQLConfig.DBPort()
						actualDBPort, err := c.DBPort()
						assert.Nil(t, err)
						assert.Equal(t, expectedDBPort, actualDBPort)

						assert.Equal(t, expectedMYSQLConfig.DBUsername(), c.DBUsername())
						assert.Equal(t, expectedMYSQLConfig.DBPassword(), c.DBPassword())
						assert.Equal(t, expectedMYSQLConfig.DBDatabase(), c.DBDatabase())

						expectedMaxConnLifetime, _ := expectedMYSQLConfig.DBConnectionMaxLifetime()
						actualMaxConnLifetime, err := c.DBConnectionMaxLifetime()

						assert.Equal(t, expectedMaxConnLifetime, actualMaxConnLifetime)

						expectedMaxIddleConn, _ := expectedMYSQLConfig.DBMaxIddleConn()
						actualMaxIddleConn, err := c.DBMaxIddleConn()
						assert.Nil(t, err)
						assert.Equal(t, expectedMaxIddleConn, actualMaxIddleConn)

						expectedMaxOpenConn, _ := expectedMYSQLConfig.DBMaxOpenConn()
						actualMaxOpenConn, err := c.DBMaxOpenConn()
						assert.Nil(t, err)
						assert.Equal(t, expectedMaxOpenConn, actualMaxOpenConn)

						testhelper.RemoveTempTestFiles(configPath)
					})
				})

				t.Run("Normal scenario with 2 config", func(t *testing.T) {
					t.Run("Run gracefully", func(t *testing.T) {
						dbName1 := "db_1"
						dbUsername1 := "username_db1"
						dbPassword1 := "password_db1"
						dbHost1 := "127.0.0.1"
						dbPort1 := "3306"

						os.Setenv("DATABASE_NAME_TWO_SCENARIO_1", dbName1)
						os.Setenv("DATABASE_USERNAME_TWO_SCENARIO_1", dbUsername1)
						os.Setenv("DATABASE_PASSWORD_TWO_SCENARIO_1", dbPassword1)
						os.Setenv("DATABASE_HOST_TWO_SCENARIO_1", dbHost1)
						os.Setenv("DATABASE_PORT_TWO_SCENARIO_1", dbPort1)

						expectedMYSQLConfig1 := entity.MYSQLDatabaseConfig{
							Database:              dbName1,
							Username:              dbUsername1,
							Password:              dbPassword1,
							Host:                  dbHost1,
							Port:                  dbPort1,
							ConnectionMaxLifetime: "120s",
							MaxIddleConnection:    "100",
							MaxOpenConnection:     "100",
						}

						dbName2 := "db_1"
						dbUsername2 := "username_db1"
						dbPassword2 := "password_db1"
						dbHost2 := "127.0.0.1"
						dbPort2 := "3306"

						os.Setenv("DATABASE_NAME_TWO_SCENARIO_2", dbName2)
						os.Setenv("DATABASE_USERNAME_TWO_SCENARIO_2", dbUsername2)
						os.Setenv("DATABASE_PASSWORD_TWO_SCENARIO_2", dbPassword2)
						os.Setenv("DATABASE_HOST_TWO_SCENARIO_2", dbHost2)
						os.Setenv("DATABASE_PORT_TWO_SCENARIO_2", dbPort2)

						expectedMYSQLConfig2 := entity.MYSQLDatabaseConfig{
							Database:              dbName2,
							Username:              dbUsername2,
							Password:              dbPassword2,
							Host:                  dbHost2,
							Port:                  dbPort2,
							ConnectionMaxLifetime: "120s",
							MaxIddleConnection:    "100",
							MaxOpenConnection:     "100",
						}

						configPath := "./db_normal_config_2/"
						fileName := "database.yml"

						testhelper.GenerateTempTestFiles(configPath, DatabaseConfigNormalScenarioWithTwoValue, fileName, 0666)

						dbConfigs, err := cfg.Database().Compile(fmt.Sprintf("%s%s", configPath, fileName))
						assert.Nil(t, err)

						c1, ok := dbConfigs["oauth_database"]

						assert.True(t, ok)
						assert.Equal(t, expectedMYSQLConfig1.Driver(), c1.Driver())
						assert.Equal(t, expectedMYSQLConfig1.DBHost(), c1.DBHost())

						expectedDBPort1, _ := expectedMYSQLConfig1.DBPort()
						actualDBPort1, err := c1.DBPort()
						assert.Nil(t, err)
						assert.Equal(t, expectedDBPort1, actualDBPort1)

						assert.Equal(t, expectedMYSQLConfig1.DBUsername(), c1.DBUsername())
						assert.Equal(t, expectedMYSQLConfig1.DBPassword(), c1.DBPassword())
						assert.Equal(t, expectedMYSQLConfig1.DBDatabase(), c1.DBDatabase())

						expectedMaxConnLifetime1, _ := expectedMYSQLConfig1.DBConnectionMaxLifetime()
						actualMaxConnLifetime1, err := c1.DBConnectionMaxLifetime()

						assert.Equal(t, expectedMaxConnLifetime1, actualMaxConnLifetime1)

						expectedMaxIddleConn1, _ := expectedMYSQLConfig1.DBMaxIddleConn()
						actualMaxIddleConn1, err := c1.DBMaxIddleConn()
						assert.Nil(t, err)
						assert.Equal(t, expectedMaxIddleConn1, actualMaxIddleConn1)

						expectedMaxOpenConn1, _ := expectedMYSQLConfig1.DBMaxOpenConn()
						actualMaxOpenConn1, err := c1.DBMaxOpenConn()
						assert.Nil(t, err)
						assert.Equal(t, expectedMaxOpenConn1, actualMaxOpenConn1)

						c2, ok := dbConfigs["other_database"]

						assert.True(t, ok)
						assert.Equal(t, expectedMYSQLConfig2.Driver(), c2.Driver())
						assert.Equal(t, expectedMYSQLConfig2.DBHost(), c2.DBHost())

						expectedDBPort2, _ := expectedMYSQLConfig2.DBPort()
						actualDBPort2, err := c2.DBPort()
						assert.Nil(t, err)
						assert.Equal(t, expectedDBPort2, actualDBPort2)

						assert.Equal(t, expectedMYSQLConfig2.DBUsername(), c2.DBUsername())
						assert.Equal(t, expectedMYSQLConfig2.DBPassword(), c2.DBPassword())
						assert.Equal(t, expectedMYSQLConfig2.DBDatabase(), c2.DBDatabase())

						expectedMaxConnLifetime2, _ := expectedMYSQLConfig2.DBConnectionMaxLifetime()
						actualMaxConnLifetime2, err := c2.DBConnectionMaxLifetime()

						assert.Equal(t, expectedMaxConnLifetime2, actualMaxConnLifetime2)

						expectedMaxIddleConn2, _ := expectedMYSQLConfig2.DBMaxIddleConn()
						actualMaxIddleConn2, err := c2.DBMaxIddleConn()
						assert.Nil(t, err)
						assert.Equal(t, expectedMaxIddleConn2, actualMaxIddleConn2)

						expectedMaxOpenConn2, _ := expectedMYSQLConfig2.DBMaxOpenConn()
						actualMaxOpenConn2, err := c2.DBMaxOpenConn()
						assert.Nil(t, err)
						assert.Equal(t, expectedMaxOpenConn2, actualMaxOpenConn2)

						testhelper.RemoveTempTestFiles(configPath)
					})
				})

				t.Run("Invalid driver", func(t *testing.T) {
					t.Run("Return error", func(t *testing.T) {
						configPath := "./db_invalid_driver/"
						fileName := "database.yml"

						testhelper.GenerateTempTestFiles(configPath, DatabaseConfigInvalidDriver, fileName, 0666)

						dbConfigs, err := cfg.Database().Compile(fmt.Sprintf("%s%s", configPath, fileName))
						assert.NotNil(t, err)
						assert.Nil(t, dbConfigs)

						testhelper.RemoveTempTestFiles(configPath)
					})
				})

				t.Run("File not found", func(t *testing.T) {
					t.Run("Return error", func(t *testing.T) {
						configPath := "./db_file_not_found/"
						fileName := "database.yml"

						testhelper.GenerateTempTestFiles(configPath, DatabaseConfigInvalidDriver, fileName, 0666)

						dbConfigs, err := cfg.Database().Compile(fmt.Sprintf("%s%s", configPath, "xxx.yml"))
						assert.NotNil(t, err)
						assert.Nil(t, dbConfigs)

						testhelper.RemoveTempTestFiles(configPath)
					})
				})

				t.Run("Invalid yaml content", func(t *testing.T) {
					t.Run("Return error", func(t *testing.T) {
						configPath := "./db_invalid_yaml/"
						fileName := "database.yml"

						testhelper.GenerateTempTestFiles(configPath, DatabaseConfigInvalidYaml, fileName, 0666)

						dbConfigs, err := cfg.Database().Compile(fmt.Sprintf("%s%s", configPath, fileName))
						assert.NotNil(t, err)
						assert.Nil(t, dbConfigs)

						testhelper.RemoveTempTestFiles(configPath)
					})
				})

				t.Run("Invalid go text template", func(t *testing.T) {
					t.Run("Return error", func(t *testing.T) {
						configPath := "./db_invalid_go_text_template/"
						fileName := "database.yml"

						testhelper.GenerateTempTestFiles(configPath, DatabaseConfigInvalidTemplateFormatting, fileName, 0666)

						dbConfigs, err := cfg.Database().Compile(fmt.Sprintf("%s%s", configPath, fileName))
						assert.NotNil(t, err)
						assert.Nil(t, dbConfigs)

						testhelper.RemoveTempTestFiles(configPath)
					})
				})

				t.Run("Missing driver", func(t *testing.T) {
					t.Run("Return error", func(t *testing.T) {
						configPath := "./db_missing_driver/"
						fileName := "database.yml"

						testhelper.GenerateTempTestFiles(configPath, DatabaseConfigMissingDriver, fileName, 0666)

						dbConfigs, err := cfg.Database().Compile(fmt.Sprintf("%s%s", configPath, fileName))
						assert.NotNil(t, err)
						assert.Nil(t, dbConfigs)

						testhelper.RemoveTempTestFiles(configPath)
					})
				})
			})

			t.Run("Mysql assign scenario", func(t *testing.T) {
				t.Run("Empty database", func(t *testing.T) {
					t.Run("Return error", func(t *testing.T) {
						configPath := "./db_mysql_empty_database/"
						fileName := "database.yml"

						testhelper.GenerateTempTestFiles(configPath, DatabaseConfigMYSQLEmptyDatabase, fileName, 0666)

						dbConfigs, err := cfg.Database().Compile(fmt.Sprintf("%s%s", configPath, fileName))
						assert.NotNil(t, err)
						assert.Nil(t, dbConfigs)

						testhelper.RemoveTempTestFiles(configPath)
					})
				})

				t.Run("Empty host", func(t *testing.T) {
					t.Run("Return error", func(t *testing.T) {
						configPath := "./db_mysql_empty_host/"
						fileName := "database.yml"

						testhelper.GenerateTempTestFiles(configPath, DatabaseConfigMYSQLEmptyHost, fileName, 0666)

						dbConfigs, err := cfg.Database().Compile(fmt.Sprintf("%s%s", configPath, fileName))
						assert.NotNil(t, err)
						assert.Nil(t, dbConfigs)

						testhelper.RemoveTempTestFiles(configPath)
					})
				})

				t.Run("Empty username", func(t *testing.T) {
					t.Run("Return error", func(t *testing.T) {
						configPath := "./db_mysql_empty_username/"
						fileName := "database.yml"

						testhelper.GenerateTempTestFiles(configPath, DatabaseConfigMYSQLEmptyUsername, fileName, 0666)

						dbConfigs, err := cfg.Database().Compile(fmt.Sprintf("%s%s", configPath, fileName))
						assert.NotNil(t, err)
						assert.Nil(t, dbConfigs)

						testhelper.RemoveTempTestFiles(configPath)
					})
				})

				t.Run("Empty port", func(t *testing.T) {
					t.Run("Fallback to 3306", func(t *testing.T) {
						configPath := "./db_mysql_empty_port/"
						fileName := "database.yml"

						testhelper.GenerateTempTestFiles(configPath, DatabaseConfigMYSQLEmptyPort, fileName, 0666)

						dbConfigs, err := cfg.Database().Compile(fmt.Sprintf("%s%s", configPath, fileName))
						assert.Nil(t, err)

						config, ok := dbConfigs["oauth_database"]
						assert.True(t, ok)

						port, err := config.DBPort()
						assert.Nil(t, err)
						assert.Equal(t, 3306, port)

						testhelper.RemoveTempTestFiles(configPath)
					})
				})

				t.Run("Empty connection_max_lifetime", func(t *testing.T) {
					t.Run("Fallback to 0", func(t *testing.T) {
						configPath := "./db_mysql_empty_conn_max_lifetime/"
						fileName := "database.yml"

						testhelper.GenerateTempTestFiles(configPath, DatabaseConfigMYSQLEmptyConnectionMaxLifetime, fileName, 0666)

						dbConfigs, err := cfg.Database().Compile(fmt.Sprintf("%s%s", configPath, fileName))
						assert.Nil(t, err)

						config, ok := dbConfigs["oauth_database"]
						assert.True(t, ok)

						connMaxLifetime, err := config.DBConnectionMaxLifetime()
						assert.Nil(t, err)
						assert.Equal(t, time.Duration(0), connMaxLifetime)

						testhelper.RemoveTempTestFiles(configPath)
					})
				})

				t.Run("Empty max_iddle_connection", func(t *testing.T) {
					t.Run("Fallback to 0", func(t *testing.T) {
						configPath := "./db_mysql_empty_max_iddle_conn/"
						fileName := "database.yml"

						testhelper.GenerateTempTestFiles(configPath, DatabaseConfigMYSQLEmptyMaxIddleConnection, fileName, 0666)

						dbConfigs, err := cfg.Database().Compile(fmt.Sprintf("%s%s", configPath, fileName))
						assert.Nil(t, err)

						config, ok := dbConfigs["oauth_database"]
						assert.True(t, ok)

						connMaxIddleConn, err := config.DBMaxIddleConn()
						assert.Nil(t, err)
						assert.Equal(t, 0, connMaxIddleConn)

						testhelper.RemoveTempTestFiles(configPath)
					})
				})

				t.Run("Empty max_open_connection", func(t *testing.T) {
					t.Run("Fallback to 0", func(t *testing.T) {
						configPath := "./db_mysql_empty_max_open_conn/"
						fileName := "database.yml"

						testhelper.GenerateTempTestFiles(configPath, DatabaseConfigMYSQLEmptyMaxOpenConnection, fileName, 0666)

						dbConfigs, err := cfg.Database().Compile(fmt.Sprintf("%s%s", configPath, fileName))
						assert.Nil(t, err)

						config, ok := dbConfigs["oauth_database"]
						assert.True(t, ok)

						maxOpenConn, err := config.DBMaxOpenConn()
						assert.Nil(t, err)
						assert.Equal(t, 0, maxOpenConn)

						testhelper.RemoveTempTestFiles(configPath)
					})
				})
			})
		})
	})
}
