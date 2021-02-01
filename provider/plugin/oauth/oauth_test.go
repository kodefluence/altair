package oauth_test

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/codefluence-x/altair/adapter"
	coreEntity "github.com/codefluence-x/altair/entity"
	"github.com/codefluence-x/altair/loader"
	"github.com/codefluence-x/altair/mock"
	"github.com/codefluence-x/altair/provider/metric"
	"github.com/codefluence-x/altair/provider/plugin/oauth"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestOauth(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	apiEngine := gin.New()

	appOption := coreEntity.AppConfigOption{
		Port:      1304,
		ProxyHost: "www.local.host",
		Plugins:   []string{"oauth"},
	}

	appConfig := coreEntity.NewAppConfig(appOption)

	oauthDatabase := "main_database"

	plugins := map[string]coreEntity.Plugin{
		"oauth": {Plugin: "oauth", Raw: []byte(`
plugin: oauth
config:
  database: ` + oauthDatabase + `
  access_token_timeout: 24h
  authorization_code_timeout: 24h
  refresh_token:
    timeout: 24h
    active: true
`)},
	}

	MYSQLConfig := coreEntity.MYSQLDatabaseConfig{
		Database:              "altair_development",
		Username:              "some_username",
		Password:              "some_password",
		Host:                  "localhost",
		Port:                  "3306",
		ConnectionMaxLifetime: "120s",
		MaxIddleConnection:    "100",
		MaxOpenConnection:     "100",
		MigrationSource:       "file://migration",
	}

	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	t.Run("Provide", func(t *testing.T) {
		t.Run("Run gracefully", func(t *testing.T) {
			appBearer := loader.AppBearer(apiEngine, adapter.AppConfig(appConfig))
			appBearer.SetMetricProvider(metric.NewPrometheusMetric())

			dbBearer := mock.NewMockDatabaseBearer(mockCtrl)
			dbBearer.EXPECT().Database(oauthDatabase).Return(db, MYSQLConfig, nil)

			pluginBearer := loader.PluginBearer(plugins)

			assert.Nil(t, oauth.Provide(appBearer, dbBearer, pluginBearer))
		})

		t.Run("Plugin is not exists in config", func(t *testing.T) {
			appOption := coreEntity.AppConfigOption{
				Port:      1304,
				ProxyHost: "www.local.host",
				Plugins:   []string{},
			}
			appConfig := coreEntity.NewAppConfig(appOption)
			appBearer := loader.AppBearer(apiEngine, adapter.AppConfig(appConfig))

			pluginBearer := loader.PluginBearer(plugins)

			dbBearer := mock.NewMockDatabaseBearer(mockCtrl)
			dbBearer.EXPECT().Database(gomock.Any()).Times(0)

			assert.Nil(t, oauth.Provide(appBearer, dbBearer, pluginBearer))
		})

		t.Run("Compile plugin failed", func(t *testing.T) {
			appBearer := loader.AppBearer(apiEngine, adapter.AppConfig(appConfig))

			dbBearer := mock.NewMockDatabaseBearer(mockCtrl)
			dbBearer.EXPECT().Database(gomock.Any()).Times(0)

			plugins := map[string]coreEntity.Plugin{}
			pluginBearer := loader.PluginBearer(plugins)

			assert.NotNil(t, oauth.Provide(appBearer, dbBearer, pluginBearer))
		})

		t.Run("Database instance is not exists", func(t *testing.T) {
			appBearer := loader.AppBearer(apiEngine, adapter.AppConfig(appConfig))

			dbBearer := mock.NewMockDatabaseBearer(mockCtrl)
			dbBearer.EXPECT().Database(oauthDatabase).Return(nil, nil, errors.New("Database is not exists"))

			pluginBearer := loader.PluginBearer(plugins)

			assert.NotNil(t, oauth.Provide(appBearer, dbBearer, pluginBearer))
		})

		t.Run("Access token timeout wrong format", func(t *testing.T) {
			appBearer := loader.AppBearer(apiEngine, adapter.AppConfig(appConfig))

			dbBearer := mock.NewMockDatabaseBearer(mockCtrl)
			dbBearer.EXPECT().Database(oauthDatabase).Return(db, MYSQLConfig, nil)

			plugins := map[string]coreEntity.Plugin{
				"oauth": {Plugin: "oauth", Raw: []byte(`
plugin: oauth
config:
  database: ` + oauthDatabase + `
  access_token_timeout: abc // this will make it fail
  authorization_code_timeout: 24h
`)},
			}
			pluginBearer := loader.PluginBearer(plugins)

			assert.NotNil(t, oauth.Provide(appBearer, dbBearer, pluginBearer))
		})

		t.Run("Authorization code timeout wrong format", func(t *testing.T) {
			appBearer := loader.AppBearer(apiEngine, adapter.AppConfig(appConfig))

			dbBearer := mock.NewMockDatabaseBearer(mockCtrl)
			dbBearer.EXPECT().Database(oauthDatabase).Return(db, MYSQLConfig, nil)

			plugins := map[string]coreEntity.Plugin{
				"oauth": {Plugin: "oauth", Raw: []byte(`
plugin: oauth
config:
  database: ` + oauthDatabase + `
  access_token_timeout: 24h
  authorization_code_timeout: abc // this will make it fail
`)},
			}
			pluginBearer := loader.PluginBearer(plugins)

			assert.NotNil(t, oauth.Provide(appBearer, dbBearer, pluginBearer))
		})

		t.Run("Refresh token timeout in wrong format", func(t *testing.T) {
			appBearer := loader.AppBearer(apiEngine, adapter.AppConfig(appConfig))
			appBearer.SetMetricProvider(metric.NewPrometheusMetric())

			dbBearer := mock.NewMockDatabaseBearer(mockCtrl)
			dbBearer.EXPECT().Database(oauthDatabase).Return(db, MYSQLConfig, nil)

			plugins := map[string]coreEntity.Plugin{
				"oauth": {Plugin: "oauth", Raw: []byte(`
plugin: oauth
config:
  database: ` + oauthDatabase + `
  access_token_timeout: 24h
  authorization_code_timeout: 24h
  refresh_token:
    timeout: 24ssss
    active: true
`)},
			}

			pluginBearer := loader.PluginBearer(plugins)

			assert.NotNil(t, oauth.Provide(appBearer, dbBearer, pluginBearer))
		})
	})
}
