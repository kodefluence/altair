package metric

import (
	"github.com/codefluence-x/altair/core"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type prometheusController struct{}

func NewPrometheusController() core.Controller {
	return &prometheusController{}
}

func (*prometheusController) Path() string {
	return "/metrics"
}

func (*prometheusController) Method() string {
	return "GET"
}

func (*prometheusController) Control(c *gin.Context) {
	promhttp.Handler().ServeHTTP(c.Writer, c.Request)
}
