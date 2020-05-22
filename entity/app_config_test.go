package entity_test

import (
	"testing"

	"github.com/codefluence-x/altair/entity"
	"github.com/stretchr/testify/assert"
)

func TestAppConfig(t *testing.T) {
	appOption := entity.AppConfigOption{
		Port:      1304,
		ProxyHost: "www.local.host",
		Plugins:   []string{"oauth"},
	}

	appOption.Authorization.Username = "altair"
	appOption.Authorization.Password = "secret"

	t.Run("Plugins", func(t *testing.T) {
		appConfig := entity.NewAppConfig(appOption)

		assert.Equal(t, appOption.Plugins, appConfig.Plugins())
	})

	t.Run("Port", func(t *testing.T) {
		appConfig := entity.NewAppConfig(appOption)

		assert.Equal(t, appOption.Port, appConfig.Port())
	})

	t.Run("BasicAuthUsername", func(t *testing.T) {
		appConfig := entity.NewAppConfig(appOption)

		assert.Equal(t, appOption.Authorization.Username, appConfig.BasicAuthUsername())
	})

	t.Run("BasicAuthPassword", func(t *testing.T) {
		appConfig := entity.NewAppConfig(appOption)

		assert.Equal(t, appOption.Authorization.Password, appConfig.BasicAuthPassword())
	})

	t.Run("ProxyHost", func(t *testing.T) {
		appConfig := entity.NewAppConfig(appOption)

		assert.Equal(t, appOption.ProxyHost, appConfig.ProxyHost())
	})
	t.Run("PluginExists", func(t *testing.T) {
		t.Run("Not exists", func(t *testing.T) {
			appOption := entity.AppConfigOption{
				Port:      1304,
				ProxyHost: "www.local.host",
				Plugins:   []string{},
			}

			appOption.Authorization.Username = "altair"
			appOption.Authorization.Password = "secret"

			t.Run("Return false", func(t *testing.T) {
				appConfig := entity.NewAppConfig(appOption)

				assert.False(t, appConfig.PluginExists("oauth"))
			})
		})

		t.Run("Exists", func(t *testing.T) {
			t.Run("Return true", func(t *testing.T) {
				appConfig := entity.NewAppConfig(appOption)

				assert.True(t, appConfig.PluginExists("oauth"))
			})
		})
	})
}
