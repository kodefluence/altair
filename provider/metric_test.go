package provider_test

import (
	"testing"

	"github.com/codefluence-x/altair/adapter"
	"github.com/codefluence-x/altair/entity"
	"github.com/codefluence-x/altair/loader"
	"github.com/codefluence-x/altair/provider"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestMetric(t *testing.T) {

	t.Run("Metric", func(t *testing.T) {
		assert.NotPanics(t, func() {
			appOption := entity.AppConfigOption{
				Port:      1304,
				ProxyHost: "www.local.host",
				Plugins:   []string{"oauth"},
			}

			provider.Metric(loader.AppBearer(gin.New(), adapter.AppConfig(entity.NewAppConfig(appOption))))
		})
	})
}
