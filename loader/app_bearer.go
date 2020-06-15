package loader

import (
	"errors"

	"github.com/codefluence-x/altair/controller"
	"github.com/codefluence-x/altair/core"
)

type appBearer struct {
	config            core.AppConfig
	downStreamPlugins []core.DownStreamPlugin
	appEngine         core.APIEngine
	metricProvider    core.Metric
}

// TODO:
// Differentiate engine with baseAPIEngine and pluginAPIEngine
// Also create injector for both baseAPIEngine and pluginAPIEngine
func AppBearer(appEngine core.APIEngine, config core.AppConfig) core.AppBearer {
	return &appBearer{
		appEngine:         appEngine,
		config:            config,
		downStreamPlugins: []core.DownStreamPlugin{},
	}
}

func (a *appBearer) Config() core.AppConfig {
	return a.config
}

func (a *appBearer) DownStreamPlugins() []core.DownStreamPlugin {
	return a.downStreamPlugins
}

func (a *appBearer) InjectDownStreamPlugin(InjectedDownStreamPlugin core.DownStreamPlugin) {
	a.downStreamPlugins = append(a.downStreamPlugins, InjectedDownStreamPlugin)
}

func (a *appBearer) InjectController(injectedController core.Controller) {
	controller.Compile(a.appEngine, a.metricProvider, injectedController)
}

func (a *appBearer) SetMetricProvider(metricProvider core.Metric) {
	if a.metricProvider == nil {
		a.metricProvider = metricProvider
	}
}

func (a *appBearer) MetricProvider() (core.Metric, error) {
	if a.metricProvider == nil {
		return nil, errors.New("Metric provider is empty")
	}

	return a.metricProvider, nil
}
