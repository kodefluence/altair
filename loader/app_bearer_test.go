package loader_test

import (
	"net/http"
	"testing"

	"github.com/codefluence-x/altair/adapter"
	"github.com/codefluence-x/altair/entity"
	"github.com/codefluence-x/altair/loader"
	"github.com/codefluence-x/altair/mock"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type fakeController struct{}

func (f fakeController) Control(c *gin.Context) {

}

func (f fakeController) Path() string {
	return "/"
}

func (f fakeController) Method() string {
	return "GET"
}

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

	appBearer := loader.AppBearer(appEngine, adapter.AppConfig(appConfig))

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

	t.Run("InjectController", func(t *testing.T) {
		assert.NotPanics(t, func() {
			appBearer.InjectController(fakeController{})
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

			appBearer := loader.AppBearer(appEngine, adapter.AppConfig(appConfig))

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

		appBearer := loader.AppBearer(appEngine, adapter.AppConfig(appConfig))

		mockMetric := mock.NewMockMetric(mockCtrl)
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

				appBearer := loader.AppBearer(appEngine, adapter.AppConfig(appConfig))

				mockMetric := mock.NewMockMetric(mockCtrl)

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

				appBearer := loader.AppBearer(appEngine, adapter.AppConfig(appConfig))

				metricProvider, err := appBearer.MetricProvider()
				assert.Nil(t, metricProvider)
				assert.NotNil(t, err)
			})
		})
	})

}
