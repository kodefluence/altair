package module

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kodefluence/monorepo/jsonapi"
	"github.com/kodefluence/monorepo/kontext"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// go:generate mockgen -source interface.go -destination mock/interface.go -package mock

type App interface {
	// TODO: Enable config via module package instead of cfg
	// Config() Config
	Controller() Controller
	// TODO: Enable plugin via module package instead of cfg
	// Plugin() Plugin
}
type Config interface {
	Port() int
	BasicAuthUsername() string
	BasicAuthPassword() string
	ProxyHost() string
	Dump() string
}

type Plugin interface {
	List() []string
	Exist(plugin string) bool
	Plugin(plugin string)
	Dump() string
}

type Controller interface {
	InjectMetric(http ...MetricController)
	InjectHTTP(http ...HttpController)
	InjectCommand(command ...CommandController)
	InjectDownstream(downstream ...DownstreamController)

	ListDownstream() []DownstreamController
	ListMetric() []MetricController
}

type MetricController interface {
	InjectCounter(metricName string, labels ...string)
	InjectHistogram(metricName string, labels ...string)
	Inc(metricName string, labels map[string]string) error
	Observe(metricName string, value float64, labels map[string]string) error
}

type HttpController interface {
	Control(ktx kontext.Context, c *gin.Context)

	// Relative path
	// /oauth/applications
	Path() string

	// GET PUT POST
	Method() string
}

type CommandController interface {
	Use() string
	Short() string
	Example() string
	Run(cmd *cobra.Command, args []string)
	ModifyFlags(flags *pflag.FlagSet)
}

type DownstreamController interface {
	Name() string
	Intervene(c *gin.Context, proxyReq *http.Request, r RouterPath) error
}

type ApiError interface {
	InternalServerError(ktx kontext.Context) jsonapi.Option
	BadRequestError(in string) jsonapi.Option
	NotFoundError(ktx kontext.Context, entityType string) jsonapi.Option
	UnauthorizedError() jsonapi.Option
	ForbiddenError(ktx kontext.Context, entityType, reason string) jsonapi.Option
	ValidationError(msg string) jsonapi.Option
}

type RouterPath interface {
	GetAuth() string
	GetScope() string
}

// type RouterCompiler interface {
// 	Compile(routesPath string) ([]entity.RouteObject, error)
// }

// type RouterForwarder interface {
// 	Generate(engine *gin.Engine, metric MetricController, routeObjects []entity.RouteObject, downStreamPlugin []DownstreamController) error
// }
