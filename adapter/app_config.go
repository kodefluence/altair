package adapter

import (
	"github.com/kodefluence/altair/core"
	"github.com/kodefluence/altair/entity"
)

type (
	appConfig struct{ c entity.AppConfig }
)

func AppConfig(c entity.AppConfig) core.AppConfig {
	return &appConfig{c: c}
}

func (a *appConfig) Port() int                           { return a.c.Port() }
func (a *appConfig) BasicAuthUsername() string           { return a.c.BasicAuthUsername() }
func (a *appConfig) BasicAuthPassword() string           { return a.c.BasicAuthPassword() }
func (a *appConfig) ProxyHost() string                   { return a.c.ProxyHost() }
func (a *appConfig) PluginExists(pluginName string) bool { return a.c.PluginExists(pluginName) }
func (a *appConfig) Plugins() []string                   { return a.c.Plugins() }
func (a *appConfig) Dump() string                        { return a.c.Dump() }
