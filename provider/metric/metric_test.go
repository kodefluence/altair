package metric_test

import (
	"testing"

	"github.com/codefluence-x/altair/adapter"
	"github.com/codefluence-x/altair/entity"
	"github.com/codefluence-x/altair/loader"
	"github.com/codefluence-x/altair/provider/metric"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestProvide(t *testing.T) {

	t.Run("Provide", func(t *testing.T) {
		assert.NotPanics(t, func() {
			appOption := entity.AppConfigOption{
				Port:      1304,
				ProxyHost: "www.local.host",
				Plugins:   []string{"oauth"},
			}

			metric.Provide(loader.AppBearer(gin.New(), adapter.AppConfig(entity.NewAppConfig(appOption))))
		})
	})
}
