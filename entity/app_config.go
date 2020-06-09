package entity

import "gopkg.in/yaml.v2"

type AppConfigOption struct {
	Port          int      `yaml:"port"`
	ProxyHost     string   `yaml:"proxy_host"`
	Plugins       []string `yaml:"plugins"`
	Authorization struct {
		Username string `yaml:"username"`
		Password string `yaml:"pasword"`
	} `yaml:"authorization"`
}

type appConfig struct {
	plugins           []string
	pluginMap         map[string]bool
	port              int
	proxyHost         string
	basicAuthUsername string
	basicAuthPassword string
}

func NewAppConfig(option AppConfigOption) appConfig {
	pluginMap := map[string]bool{}

	for _, p := range option.Plugins {
		pluginMap[p] = true
	}

	return appConfig{
		plugins:           option.Plugins,
		pluginMap:         pluginMap,
		port:              option.Port,
		proxyHost:         option.ProxyHost,
		basicAuthPassword: option.Authorization.Password,
		basicAuthUsername: option.Authorization.Username,
	}
}

func (a appConfig) PluginExists(pluginName string) bool {
	_, ok := a.pluginMap[pluginName]
	return ok
}

func (a appConfig) Plugins() []string {
	return a.plugins
}

func (a appConfig) Port() int {
	return a.port
}

func (a appConfig) BasicAuthUsername() string {
	return a.basicAuthUsername
}

func (a appConfig) BasicAuthPassword() string {
	return a.basicAuthPassword
}

func (a appConfig) ProxyHost() string {
	return a.proxyHost
}

func (a appConfig) Dump() string {
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
