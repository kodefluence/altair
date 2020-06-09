package loader_test

import (
	"net/http"
	"testing"

	"github.com/codefluence-x/altair/entity"
	"github.com/codefluence-x/altair/loader"
	"github.com/gin-gonic/gin"
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
	appOption := entity.AppConfigOption{
		Port:      1304,
		ProxyHost: "www.local.host",
		Plugins:   []string{"oauth"},
	}

	appOption.Authorization.Username = "altair"
	appOption.Authorization.Password = "secret"

	appConfig := entity.NewAppConfig(appOption)
	appEngine := gin.Default()

	appBearer := loader.AppBearer(appEngine, appConfig)

	t.Run("Config", func(t *testing.T) {
		t.Run("Return AppConfig", func(t *testing.T) {
			assert.Equal(t, appConfig, appBearer.Config())
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

			appBearer := loader.AppBearer(appEngine, appConfig)

			appBearer.InjectDownStreamPlugin(fakeDownStreamPlugin{})
			appBearer.InjectDownStreamPlugin(fakeDownStreamPlugin{})

			assert.Equal(t, 2, len(appBearer.DownStreamPlugins()))
		})
	})
}
