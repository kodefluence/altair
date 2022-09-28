package usecase

import (
	"github.com/gin-gonic/gin"
	"github.com/kodefluence/altair/module"
	"github.com/spf13/cobra"
)

type HttpInjector func(httpMethod, relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes

type Metric interface {
	InjectCounter(metricName string, labels ...string)
	InjectHistogram(metricName string, labels ...string)
	Inc(metricName string, labels map[string]string) error
	Observe(metricName string, value float64, labels map[string]string) error
}

type Controller struct {
	httpController       []module.HttpController
	commandController    []module.CommandController
	downstreamController []module.DownstreamController

	httpInjector HttpInjector
	apiError     module.ApiError
	metric       Metric
	rootCommand  *cobra.Command
}

func NewController(httpInjector HttpInjector, apiError module.ApiError, metric Metric, rootCommand *cobra.Command) *Controller {
	return &Controller{
		httpController:       []module.HttpController{},
		commandController:    []module.CommandController{},
		downstreamController: []module.DownstreamController{},

		httpInjector: httpInjector,
		apiError:     apiError,
		metric:       metric,
		rootCommand:  rootCommand,
	}
}
