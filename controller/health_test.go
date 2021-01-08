package controller_test

import (
	"net/http"
	"testing"

	"github.com/codefluence-x/altair/controller"
	"github.com/codefluence-x/altair/testhelper"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHealth(t *testing.T) {
	t.Run("Health", func(t *testing.T) {
		t.Run("Return OK response", func(t *testing.T) {
			gin.SetMode(gin.ReleaseMode)
			engine := gin.New()
			engine.GET("/health", controller.Health)
			w := testhelper.PerformRequest(engine, "GET", "/health", nil)
			assert.Equal(t, http.StatusOK, w.Code)
		})
	})
}
