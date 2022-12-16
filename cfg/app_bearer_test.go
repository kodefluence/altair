package cfg_test

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/kodefluence/altair/adapter"
	"github.com/kodefluence/altair/cfg"
	"github.com/kodefluence/altair/entity"
	"github.com/kodefluence/altair/plugin/metric/module/dummy/controller/metric"
	"github.com/stretchr/testify/assert"
)

type fakeDownStreamPlugin struct{}

func (fakeDownStreamPlugin) Name() string {
	return "fakeDownStreamPlugin"
}

func (fakeDownStreamPlugin) Intervene(c *gin.Context, proxyReq *http.Request, r entity.RouterPath) error {
	return nil
}

func TestAppBearer(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	appOption := entity.AppConfigOption{
		Port:      1304,
		ProxyHost: "www.local.host",
		Plugins:   []string{"oauth"},
	}

	appOption.Authorization.Username = "altair"
	appOption.Authorization.Password = "secret"

	appConfig := entity.NewAppConfig(appOption)
	appEngine := gin.Default()

	appBearer := cfg.AppBearer(appEngine, adapter.AppConfig(appConfig))
	appBearer.SetMetricProvider(metric.NewDummy())

	t.Run("Config", func(t *testing.T) {
		t.Run("Return AppConfig", func(t *testing.T) {
			assert.Equal(t, adapter.AppConfig(appConfig), appBearer.Config())
		})
	})

	t.Run("InjectDownStreamPlugin", func(t *testing.T) {
		assert.NotPanics(t, func() {
			appBearer.InjectDownStreamPlugin(fakeDownStreamPlugin{})
		})
	})

	t.Run("DownStreamPlugins", func(t *testing.T) {
		t.Run("return DownStreamPlugins", func(t *testing.T) {
			appOption := entity.AppConfigOption{
				Port:      1304,
				ProxyHost: "www.local.host",
				Plugins:   []string{"oauth"},
			}

			appOption.Authorization.Username = "altair"
			appOption.Authorization.Password = "secret"

			appConfig := entity.NewAppConfig(appOption)
			appEngine := gin.Default()

			appBearer := cfg.AppBearer(appEngine, adapter.AppConfig(appConfig))

			appBearer.InjectDownStreamPlugin(fakeDownStreamPlugin{})
			appBearer.InjectDownStreamPlugin(fakeDownStreamPlugin{})

			assert.Equal(t, 2, len(appBearer.DownStreamPlugins()))
		})
	})

	t.Run("SetMetricProvider", func(t *testing.T) {
		appOption := entity.AppConfigOption{
			Port:      1304,
			ProxyHost: "www.local.host",
			Plugins:   []string{"oauth"},
		}

		appOption.Authorization.Username = "altair"
		appOption.Authorization.Password = "secret"

		appConfig := entity.NewAppConfig(appOption)
		appEngine := gin.Default()

		appBearer := cfg.AppBearer(appEngine, adapter.AppConfig(appConfig))

		mockMetric := metric.NewDummy()
		assert.NotPanics(t, func() {
			appBearer.SetMetricProvider(mockMetric)
		})
	})

	t.Run("SetMetricProvider", func(t *testing.T) {
		t.Run("Has metric provider", func(t *testing.T) {
			t.Run("Return metric provider", func(t *testing.T) {
				appOption := entity.AppConfigOption{
					Port:      1304,
					ProxyHost: "www.local.host",
					Plugins:   []string{"oauth"},
				}

				appOption.Authorization.Username = "altair"
				appOption.Authorization.Password = "secret"

				appConfig := entity.NewAppConfig(appOption)
				appEngine := gin.Default()

				appBearer := cfg.AppBearer(appEngine, adapter.AppConfig(appConfig))

				mockMetric := metric.NewDummy()

				appBearer.SetMetricProvider(mockMetric)
				metricProvider, err := appBearer.MetricProvider()
				assert.Nil(t, err)
				assert.Equal(t, mockMetric, metricProvider)
			})
		})

		t.Run("No metric provider", func(t *testing.T) {
			t.Run("Return error", func(t *testing.T) {
				appOption := entity.AppConfigOption{
					Port:      1304,
					ProxyHost: "www.local.host",
					Plugins:   []string{"oauth"},
				}

				appOption.Authorization.Username = "altair"
				appOption.Authorization.Password = "secret"

				appConfig := entity.NewAppConfig(appOption)
				appEngine := gin.Default()

				appBearer := cfg.AppBearer(appEngine, adapter.AppConfig(appConfig))

				metricProvider, err := appBearer.MetricProvider()
				assert.Nil(t, metricProvider)
				assert.NotNil(t, err)
			})
		})
	})

}
