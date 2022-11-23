package http_test

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/kodefluence/altair/module/apierror"
	"github.com/kodefluence/altair/module/controller"
	prometheusHttp "github.com/kodefluence/altair/plugin/metric/module/prometheus/controller/http"
	"github.com/kodefluence/altair/testhelper"
	"github.com/spf13/cobra"
	"gotest.tools/assert"
)

func TestPrometheusController(t *testing.T) {

	t.Run("Method", func(t *testing.T) {
		assert.Equal(t, "GET", prometheusHttp.NewPrometheusController().Method())
	})

	t.Run("Path", func(t *testing.T) {
		assert.Equal(t, "/metrics", prometheusHttp.NewPrometheusController().Path())
	})

	t.Run("Control", func(t *testing.T) {
		t.Run("Return metrics content", func(t *testing.T) {
			apiEngine := gin.Default()

			ctrl := prometheusHttp.NewPrometheusController()
			controller.Provide(apiEngine.Handle, apierror.Provide(), &cobra.Command{}).InjectHTTP(ctrl)
			w := testhelper.PerformRequest(apiEngine, ctrl.Method(), ctrl.Path(), nil)

			assert.Equal(t, http.StatusOK, w.Code)
		})
	})
}
