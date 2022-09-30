package http_test

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/kodefluence/altair/module/apierror"
	"github.com/kodefluence/altair/module/controller"
	healthcheckHttp "github.com/kodefluence/altair/module/healthcheck/controller/http"
	"github.com/kodefluence/altair/plugin/metric/module/dummy/usecase"
	"github.com/kodefluence/altair/testhelper"
	"github.com/spf13/cobra"
	"gotest.tools/assert"
)

func TestHealth(t *testing.T) {
	t.Run("Health", func(t *testing.T) {
		t.Run("Return OK response", func(t *testing.T) {
			gin.SetMode(gin.ReleaseMode)
			engine := gin.New()

			controller.Provide(engine.Handle, apierror.Provide(), usecase.NewDummy(), &cobra.Command{}).InjectHTTP(healthcheckHttp.NewHealthController())
			w := testhelper.PerformRequest(engine, "GET", "/health", nil)
			assert.Equal(t, http.StatusOK, w.Code)
		})
	})
}
