package core

import (
	"github.com/codefluence-x/altair/entity"
	"github.com/gin-gonic/gin"
)

type RouteCompiler interface {
	Compile(path string) []entity.RouteObject
}

type RouteGenerator interface {
	Generate(c *gin.Engine, routeObjects []entity.RouteObject)
}
