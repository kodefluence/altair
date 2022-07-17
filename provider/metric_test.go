package provider_test

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/kodefluence/altair/adapter"
	"github.com/kodefluence/altair/cfg"
	"github.com/kodefluence/altair/entity"
	"github.com/kodefluence/altair/provider"
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

			provider.Metric(cfg.AppBearer(gin.New(), adapter.AppConfig(entity.NewAppConfig(appOption))))
		})
	})
}
