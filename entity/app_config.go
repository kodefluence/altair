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
	Metric struct {
		Interface string
	}
}

type AppConfig struct {
	plugins           []string
	pluginMap         map[string]bool
	port              int
	proxyHost         string
	basicAuthUsername string
	basicAuthPassword string
	metricConfig      *MetricConfig
}

type MetricConfig struct {
	metricInterface string
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
		metricConfig: &MetricConfig{
			metricInterface: option.Metric.Interface,
		},
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
	appConfigOption.Metric.Interface = a.metricConfig.metricInterface

	encodedContent, _ := yaml.Marshal(appConfigOption)
	return string(encodedContent)
}

func (a AppConfig) Metric() *MetricConfig {
	return a.metricConfig
}

func (m *MetricConfig) Interface() string {
	return m.metricInterface
}
