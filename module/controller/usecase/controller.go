package usecase

import (
	"github.com/gin-gonic/gin"
	"github.com/kodefluence/altair/module"
	"github.com/spf13/cobra"
)

type HttpInjector func(httpMethod, relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes

type Controller struct {
	httpController       []module.HttpController
	commandController    []module.CommandController
	downstreamController []module.DownstreamController
	metricController     []module.MetricController

	httpInjector HttpInjector
	apiError     module.ApiError

	rootCommand *cobra.Command
	commandList map[string]*cobra.Command
}

func NewController(httpInjector HttpInjector, apiError module.ApiError, rootCommand *cobra.Command) *Controller {
	return &Controller{
		httpController:       []module.HttpController{},
		commandController:    []module.CommandController{},
		downstreamController: []module.DownstreamController{},
		metricController:     []module.MetricController{},

		httpInjector: httpInjector,
		apiError:     apiError,
		rootCommand:  rootCommand,
		commandList:  map[string]*cobra.Command{},
	}
}
