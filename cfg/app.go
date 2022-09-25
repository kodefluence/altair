package cfg

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"

	"github.com/kodefluence/altair/adapter"
	"github.com/kodefluence/altair/core"
	"github.com/kodefluence/altair/entity"
	"gopkg.in/yaml.v2"
)

type app struct{}

type baseAppConfig struct {
	Version       string   `yaml:"version"`
	Plugins       []string `yaml:"plugins"`
	Port          string   `yaml:"port"`
	ProxyHost     string   `yaml:"proxy_host"`
	Authorization struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"authorization"`
}

func App() core.AppLoader {
	return &app{}
}

func (a *app) Compile(configPath string) (core.AppConfig, error) {
	contents, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	compiledContents, err := compileTemplate(contents)
	if err != nil {
		return nil, err
	}

	var config baseAppConfig

	if err := yaml.Unmarshal(compiledContents, &config); err != nil {
		return nil, err
	}

	switch v := config.Version; v {
	case "1.0":
		var appConfigOption entity.AppConfigOption

		if config.Authorization.Username == "" {
			return nil, errors.New("config authorization `username` cannot be empty")
		}

		if config.Authorization.Password == "" {
			return nil, errors.New("config authorization `password` cannot be empty")
		}

		if config.Port == "" {
			appConfigOption.Port = 1304
		} else {
			port, err := strconv.Atoi(config.Port)
			if err != nil {
				return nil, err
			}

			appConfigOption.Port = port
		}

		if config.ProxyHost == "" {
			appConfigOption.ProxyHost = "www.local.host"
		} else {
			appConfigOption.ProxyHost = config.ProxyHost
		}

		appConfigOption.Plugins = config.Plugins
		appConfigOption.Authorization.Username = config.Authorization.Username
		appConfigOption.Authorization.Password = config.Authorization.Password

		return adapter.AppConfig(entity.NewAppConfig(appConfigOption)), nil
	default:
		return nil, fmt.Errorf("undefined template version: %s for app.yaml", v)
	}
}
