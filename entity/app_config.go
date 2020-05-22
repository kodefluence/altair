package entity

type AppConfigOption struct {
	Port          int
	ProxyHost     string
	Plugins       []string
	Authorization struct {
		Username string
		Password string
	}
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
