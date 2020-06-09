package loader

import (
	"github.com/codefluence-x/altair/controller"
	"github.com/codefluence-x/altair/core"
	"github.com/gin-gonic/gin"
)

type appBearer struct {
	config            core.AppConfig
	downStreamPlugins []core.DownStreamPlugin
	appEngine         *gin.Engine
}

func AppBearer(appEngine *gin.Engine, config core.AppConfig) core.AppBearer {
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
	controller.Compile(a.appEngine, injectedController)
}
