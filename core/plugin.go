package core

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type DownStreamPlugin interface {
	Name() string
	Intervene(c *gin.Context, proxyReq *http.Request) error
}
