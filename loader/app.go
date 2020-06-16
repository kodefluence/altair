package loader

import (
	"errors"
	"io/ioutil"
	"strconv"

	"github.com/codefluence-x/altair/adapter"
	"github.com/codefluence-x/altair/core"
	"github.com/codefluence-x/altair/entity"
	"gopkg.in/yaml.v2"
)

type app struct{}

type appConfig struct {
	Plugins       []string `yaml:"plugins"`
	Port          string   `yaml:"port"`
	ProxyHost     string   `yaml:"proxy_host"`
	Authorization struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"authorization"`
	Metric struct {
		Interface string `yaml:"interface"`
	} `yaml:"metric"`
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

	var config appConfig

	if err := yaml.Unmarshal(compiledContents, &config); err != nil {
		return nil, err
	}

	appConfigOption, err := a.assignConfigOption(config)
	if err != nil {
		return nil, err
	}

	return adapter.AppConfig(entity.NewAppConfig(appConfigOption)), nil
}

func (a *app) assignConfigOption(config appConfig) (entity.AppConfigOption, error) {
	var appConfigOption entity.AppConfigOption

	if config.Authorization.Username == "" {
		return appConfigOption, errors.New("config authorization `username` cannot be empty")
	}

	if config.Authorization.Password == "" {
		return appConfigOption, errors.New("config authorization `password` cannot be empty")
	}

	if config.Port == "" {
		appConfigOption.Port = 1304
	} else {
		port, err := strconv.Atoi(config.Port)
		if err != nil {
			return appConfigOption, err
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
	appConfigOption.Metric.Interface = config.Metric.Interface

	return appConfigOption, nil
}
