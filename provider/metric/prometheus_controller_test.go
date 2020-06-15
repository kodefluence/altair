package metric_test

import (
	"net/http"
	"testing"

	"github.com/codefluence-x/altair/provider/metric"
	"github.com/codefluence-x/altair/testhelper"
	"github.com/gin-gonic/gin"
	"gotest.tools/assert"
)

func TestPrometheusController(t *testing.T) {

	t.Run("Method", func(t *testing.T) {
		assert.Equal(t, "GET", metric.NewPrometheusController().Method())
	})

	t.Run("Path", func(t *testing.T) {
		assert.Equal(t, "/metrics", metric.NewPrometheusController().Path())
	})

	t.Run("Control", func(t *testing.T) {
		t.Run("Return metrics content", func(t *testing.T) {
			apiEngine := gin.Default()

			ctrl := metric.NewPrometheusController()
			apiEngine.Handle(ctrl.Method(), ctrl.Path(), ctrl.Control)

			w := testhelper.PerformRequest(apiEngine, ctrl.Method(), ctrl.Path(), nil)

			assert.Equal(t, http.StatusOK, w.Code)
		})
	})
}
