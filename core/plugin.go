package core

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kodefluence/altair/entity"
)

type DownStreamPlugin interface {
	Name() string
	Intervene(c *gin.Context, proxyReq *http.Request, r entity.RouterPath) error
}
