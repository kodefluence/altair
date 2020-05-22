package entity_test

import (
	"testing"

	"github.com/codefluence-x/altair/entity"
	"github.com/stretchr/testify/assert"
)

func TestAppConfig(t *testing.T) {
	t.Run("Plugins", func(t *testing.T) {
		plugins := []string{"oauth"}
		appConfig := entity.NewAppConfig(plugins)

		assert.Equal(t, plugins, appConfig.Plugins())
	})

	t.Run("PluginExists", func(t *testing.T) {
		t.Run("Not exists", func(t *testing.T) {
			t.Run("Return false", func(t *testing.T) {
				plugins := []string{}
				appConfig := entity.NewAppConfig(plugins)

				assert.False(t, appConfig.PluginExists("oauth"))
			})
		})

		t.Run("Exists", func(t *testing.T) {
			t.Run("Return true", func(t *testing.T) {
				plugins := []string{"oauth"}
				appConfig := entity.NewAppConfig(plugins)

				assert.True(t, appConfig.PluginExists("oauth"))
			})
		})
	})
}
