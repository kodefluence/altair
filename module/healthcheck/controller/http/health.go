package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kodefluence/monorepo/kontext"

	"github.com/kodefluence/altair/module"
)

type HealthController struct{}

func NewHealthController() module.HttpController {
	return &HealthController{}
}

func (*HealthController) Path() string {
	return "/health"
}

func (*HealthController) Method() string {
	return "GET"
}

func (*HealthController) Control(ktx kontext.Context, c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "OK",
	})
}
