package cfg

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"

	"github.com/kodefluence/altair/core"
	"github.com/kodefluence/altair/entity"
)

type database struct{}

// baseDatabaseConfig wraps the v1.0 envelope. Both Version and Instances are
// empty for legacy (pre-version) files; the loader detects that and falls
// back to the flat map shape below.
type baseDatabaseConfig struct {
	Version   string                       `yaml:"version"`
	Instances map[string]map[string]string `yaml:"instances"`
}

func Database() core.DatabaseLoader {
	return &database{}
}

func (d *database) Compile(configPath string) (map[string]core.DatabaseConfig, error) {
	f, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	contents, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	compiledContents, err := compileTemplate(contents)
	if err != nil {
		return nil, err
	}

	driverConfigs, err := d.selectDriverConfigs(compiledContents)
	if err != nil {
		return nil, err
	}

	return d.assignConfig(driverConfigs)
}

// selectDriverConfigs returns the flat driver-config map for either envelope.
//
// v1.0:
//
//	version: "1.0"
//	instances:
//	  main_database: {driver: mysql, ...}
//
// legacy (no version field):
//
//	main_database: {driver: mysql, ...}
//
// Legacy emits a one-shot deprecation warning so deployments keep working
// while operators migrate.
func (d *database) selectDriverConfigs(compiledContents []byte) (map[string]map[string]string, error) {
	var envelope baseDatabaseConfig
	if err := yaml.Unmarshal(compiledContents, &envelope); err == nil && envelope.Version != "" {
		switch envelope.Version {
		case "1.0":
			if envelope.Instances == nil {
				return nil, errors.New("database config version `1.0` requires an `instances` map")
			}
			return envelope.Instances, nil
		default:
			return nil, fmt.Errorf("undefined template version: %s for database.yml", envelope.Version)
		}
	}

	var legacy map[string]map[string]string
	if err := yaml.Unmarshal(compiledContents, &legacy); err != nil {
		return nil, err
	}

	if legacy != nil {
		log.Warn().
			Array("tags", zerolog.Arr().Str("altair").Str("cfg").Str("database").Str("deprecation")).
			Msg("database.yml is missing a `version` field; add `version: \"1.0\"` and wrap instances under `instances:`. Legacy format will be removed in a future release.")
	}

	return legacy, nil
}

func (d *database) assignConfig(driverConfigs map[string]map[string]string) (map[string]core.DatabaseConfig, error) {
	databaseConfigs := map[string]core.DatabaseConfig{}

	for key, config := range driverConfigs {
		driver, ok := config["driver"]
		if !ok {
			return nil, errors.New("database driver is not specified")
		}

		switch driver {
		case "mysql":
			c, err := d.assignMysqlConfig(config)
			if err != nil {
				return nil, err
			}
			databaseConfigs[key] = c
		default:
			return nil, fmt.Errorf("database driver:  `%s` is not supported", driver)
		}
	}

	return databaseConfigs, nil
}

func (d *database) assignMysqlConfig(config map[string]string) (core.DatabaseConfig, error) {
	var mysqlConfig entity.MYSQLDatabaseConfig

	if database, ok := config["database"]; ok && database != "" {
		mysqlConfig.Database = database
	} else {
		return nil, errors.New("Config `database` cannot be empty for mysql driver")
	}

	if username, ok := config["username"]; ok && username != "" {
		mysqlConfig.Username = username
	} else {
		return nil, errors.New("Config `username` cannot be empty for mysql driver")
	}

	if password, ok := config["password"]; ok && password != "" {
		mysqlConfig.Password = password
	}

	if host, ok := config["host"]; ok && host != "" {
		mysqlConfig.Host = host
	} else {
		return nil, errors.New("Config `host` cannot be empty for mysql driver")
	}

	if port, ok := config["port"]; ok && port != "" {
		mysqlConfig.Port = port
	} else {
		mysqlConfig.Port = "3306"
	}

	if connection_max_lifetime, ok := config["connection_max_lifetime"]; ok && connection_max_lifetime != "" {
		mysqlConfig.ConnectionMaxLifetime = connection_max_lifetime
	} else {
		mysqlConfig.ConnectionMaxLifetime = "0"
	}

	// Accept both the canonical `max_idle_connection` and the legacy typo
	// `max_iddle_connection`. Prefer the canonical spelling when both are set;
	// emit a one-shot deprecation warning when only the legacy key is present.
	if maxIdle, ok := config["max_idle_connection"]; ok && maxIdle != "" {
		mysqlConfig.MaxIddleConnection = maxIdle
	} else if maxIdleLegacy, ok := config["max_iddle_connection"]; ok && maxIdleLegacy != "" {
		log.Warn().
			Array("tags", zerolog.Arr().Str("altair").Str("cfg").Str("database").Str("deprecation")).
			Msg("`max_iddle_connection` is a misspelling; use `max_idle_connection` instead. The legacy key will be removed in a future release.")
		mysqlConfig.MaxIddleConnection = maxIdleLegacy
	} else {
		mysqlConfig.MaxIddleConnection = "0"
	}

	if max_open_connection, ok := config["max_open_connection"]; ok && max_open_connection != "" {
		mysqlConfig.MaxOpenConnection = max_open_connection
	} else {
		mysqlConfig.MaxOpenConnection = "0"
	}

	return mysqlConfig, nil
}
