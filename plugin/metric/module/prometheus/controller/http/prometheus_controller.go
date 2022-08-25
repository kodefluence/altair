package http

import (
	"github.com/gin-gonic/gin"
	"github.com/kodefluence/altair/core"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type PrometheusController struct{}

func NewPrometheusController() core.Controller {
	return &PrometheusController{}
}

func (*PrometheusController) Path() string {
	return "/metrics"
}

func (*PrometheusController) Method() string {
	return "GET"
}

func (*PrometheusController) Control(c *gin.Context) {
	promhttp.Handler().ServeHTTP(c.Writer, c.Request)
}
