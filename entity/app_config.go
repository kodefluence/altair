package entity

import (
	"gopkg.in/yaml.v2"
)

type AppConfigOption struct {
	Port          int
	ProxyHost     string
	Plugins       []string
	Authorization struct {
		Username string
		Password string
	}
}

type AppConfig struct {
	plugins           []string
	pluginMap         map[string]bool
	port              int
	proxyHost         string
	basicAuthUsername string
	basicAuthPassword string
}

func NewAppConfig(option AppConfigOption) AppConfig {
	pluginMap := map[string]bool{}

	for _, p := range option.Plugins {
		pluginMap[p] = true
	}

	return AppConfig{
		plugins:           option.Plugins,
		pluginMap:         pluginMap,
		port:              option.Port,
		proxyHost:         option.ProxyHost,
		basicAuthPassword: option.Authorization.Password,
		basicAuthUsername: option.Authorization.Username,
	}
}

func (a AppConfig) PluginExists(pluginName string) bool {
	_, ok := a.pluginMap[pluginName]
	return ok
}

func (a AppConfig) Plugins() []string {
	return a.plugins
}

func (a AppConfig) Port() int {
	return a.port
}

func (a AppConfig) BasicAuthUsername() string {
	return a.basicAuthUsername
}

func (a AppConfig) BasicAuthPassword() string {
	return a.basicAuthPassword
}

func (a AppConfig) ProxyHost() string {
	return a.proxyHost
}

func (a AppConfig) Dump() string {
	appConfigOption := AppConfigOption{
		Port:      a.port,
		Plugins:   a.plugins,
		ProxyHost: a.proxyHost,
	}

	appConfigOption.Authorization.Username = a.basicAuthUsername
	appConfigOption.Authorization.Password = a.basicAuthPassword

	encodedContent, _ := yaml.Marshal(appConfigOption)
	return string(encodedContent)
}
