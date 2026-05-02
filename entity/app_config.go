package entity

import (
	"time"

	"gopkg.in/yaml.v2"
)

type AppConfigOption struct {
	Port               int
	ProxyHost          string
	UpstreamTimeout    time.Duration
	MaxRequestBodySize int64
	Plugins            []string
	AutoMigrate        bool
	Authorization      struct {
		Username string
		Password string
	}
}

type AppConfig struct {
	plugins            []string
	pluginMap          map[string]bool
	port               int
	proxyHost          string
	upstreamTimeout    time.Duration
	maxRequestBodySize int64
	basicAuthUsername  string
	basicAuthPassword  string
	autoMigrate        bool
}

func NewAppConfig(option AppConfigOption) AppConfig {
	pluginMap := map[string]bool{}

	for _, p := range option.Plugins {
		pluginMap[p] = true
	}

	return AppConfig{
		plugins:            option.Plugins,
		pluginMap:          pluginMap,
		port:               option.Port,
		proxyHost:          option.ProxyHost,
		upstreamTimeout:    option.UpstreamTimeout,
		maxRequestBodySize: option.MaxRequestBodySize,
		basicAuthPassword:  option.Authorization.Password,
		basicAuthUsername:  option.Authorization.Username,
		autoMigrate:        option.AutoMigrate,
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

func (a AppConfig) UpstreamTimeout() time.Duration {
	return a.upstreamTimeout
}

func (a AppConfig) MaxRequestBodySize() int64 {
	return a.maxRequestBodySize
}

func (a AppConfig) AutoMigrate() bool {
	return a.autoMigrate
}

func (a AppConfig) Dump() string {
	appConfigOption := AppConfigOption{
		Port:               a.port,
		Plugins:            a.plugins,
		ProxyHost:          a.proxyHost,
		UpstreamTimeout:    a.upstreamTimeout,
		MaxRequestBodySize: a.maxRequestBodySize,
		AutoMigrate:        a.autoMigrate,
	}

	appConfigOption.Authorization.Username = a.basicAuthUsername
	appConfigOption.Authorization.Password = a.basicAuthPassword

	encodedContent, _ := yaml.Marshal(appConfigOption)
	return string(encodedContent)
}
