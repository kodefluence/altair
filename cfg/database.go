package cfg

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/kodefluence/altair/core"
	"github.com/kodefluence/altair/entity"
	"gopkg.in/yaml.v2"
)

type database struct{}

func Database() core.DatabaseLoader {
	return &database{}
}

func (d *database) Compile(configPath string) (map[string]core.DatabaseConfig, error) {
	f, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}

	contents, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	compiledContents, err := compileTemplate(contents)
	if err != nil {
		return nil, err
	}

	driverConfigs, err := d.unmarshalDriver(compiledContents)
	if err != nil {
		return nil, err
	}

	databaseConfigs, err := d.assignConfig(driverConfigs)
	if err != nil {
		return nil, err
	}

	return databaseConfigs, nil
}

func (d *database) unmarshalDriver(compiledContents []byte) (map[string]map[string]string, error) {
	var driverConfigs map[string]map[string]string

	err := yaml.Unmarshal(compiledContents, &driverConfigs)
	if err != nil {
		return nil, err
	}

	return driverConfigs, nil
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

	if max_iddle_connection, ok := config["max_iddle_connection"]; ok && max_iddle_connection != "" {
		mysqlConfig.MaxIddleConnection = max_iddle_connection
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
