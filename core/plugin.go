package core

import (
	"net/http"

	"github.com/codefluence-x/altair/entity"
	"github.com/gin-gonic/gin"
)

type PluginProviderDispatcher interface {
	Oauth() PluginProvider
}

type PluginProvider interface {
	Controllers() []Controller
	DownStreamPlugins() []DownStreamPlugin
}

type DownStreamPlugin interface {
	Name() string
	Intervene(c *gin.Context, proxyReq *http.Request, r entity.RouterPath) error
}
