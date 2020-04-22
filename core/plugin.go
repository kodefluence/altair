package core

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type DownStreamPlugin interface {
	Intervene(c *gin.Context, proxyReq *http.Request) error
}
