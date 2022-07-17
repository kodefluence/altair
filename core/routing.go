package core

import (
	"github.com/gin-gonic/gin"
	"github.com/kodefluence/altair/entity"
)

type RouteCompiler interface {
	Compile(routesPath string) ([]entity.RouteObject, error)
}

type RouteGenerator interface {
	Generate(engine *gin.Engine, metric Metric, routeObjects []entity.RouteObject, downStreamPlugin []DownStreamPlugin) error
}
