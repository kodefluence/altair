package entity

type AppConfig struct {
	plugins   []string
	pluginMap map[string]bool
}

func NewAppConfig(plugins []string) AppConfig {
	pluginMap := map[string]bool{}

	for _, p := range plugins {
		pluginMap[p] = true
	}

	return AppConfig{
		plugins:   plugins,
		pluginMap: pluginMap,
	}
}

func (a AppConfig) PluginExists(pluginName string) bool {
	_, ok := a.pluginMap[pluginName]
	return ok
}

func (a AppConfig) Plugins() []string {
	return a.plugins
}
