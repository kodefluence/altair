package core

import (
	"github.com/codefluence-x/altair/entity"
	"github.com/gin-gonic/gin"
)

type RouteCompiler interface {
	Compile(routesPath string) ([]entity.RouteObject, error)
}

type RouteGenerator interface {
	Generate(engine *gin.Engine, routeObjects []entity.RouteObject, downStreamPlugin []DownStreamPlugin) error
}
