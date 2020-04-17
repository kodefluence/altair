package route

import (
	"github.com/codefluence-x/altair/core"
	"github.com/codefluence-x/altair/entity"
	"github.com/gin-gonic/gin"
)

type generator struct{}

func Generator() core.RouteGenerator {
	return &generator{}
}

func (g *generator) Generate(c *gin.Engine, routeObjects []entity.RouteObject) {

}
