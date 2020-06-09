package loader_test

import (
	"errors"
	"testing"

	"github.com/codefluence-x/altair/entity"
	"github.com/codefluence-x/altair/loader"
	"github.com/stretchr/testify/assert"
)

func TestPluginBearer(t *testing.T) {

	// We use oauth plugin as an example of struct that will be compiled by plugin bearer
	type OauthPlugin struct {
		Config struct {
			Database string `yaml:"database"`

			AccessTokenTimeoutRaw       string `yaml:"access_token_timeout"`
			AuthorizationCodeTimeoutRaw string `yaml:"authorization_code_timeout"`
		} `yaml:"config"`
	}

	t.Run("ConfigExists", func(t *testing.T) {
		t.Run("Given plugin name", func(t *testing.T) {
			plugins := map[string]entity.Plugin{
				"oauth": {Plugin: "oauth"},
			}
			pluginBearer := loader.PluginBearer(plugins)

			t.Run("Plugin exists", func(t *testing.T) {
				t.Run("Return true", func(t *testing.T) {
					assert.True(t, pluginBearer.ConfigExists("oauth"))
				})
			})

			t.Run("Plugin is not exists", func(t *testing.T) {
				t.Run("Return false", func(t *testing.T) {
					assert.False(t, pluginBearer.ConfigExists("cache"))
				})
			})
		})
	})

	t.Run("Length", func(t *testing.T) {
		t.Run("Return the length of plugin list", func(t *testing.T) {
			plugins := map[string]entity.Plugin{
				"oauth": {Plugin: "oauth"},
				"cache": {Plugin: "cache"},
			}
			pluginBearer := loader.PluginBearer(plugins)

			assert.Equal(t, len(plugins), pluginBearer.Length())
		})
	})

	t.Run("CompilePlugin", func(t *testing.T) {
		t.Run("Given Plugin and Injected Struct", func(t *testing.T) {
			t.Run("Injected struct has been injected with yaml value", func(t *testing.T) {
				t.Run("Return nil", func(t *testing.T) {
					plugins := map[string]entity.Plugin{
						"oauth": {Plugin: "oauth", Raw: []byte(PluginConfigNormal1)},
					}
					pluginBearer := loader.PluginBearer(plugins)

					oauthPlugins := OauthPlugin{}

					err := pluginBearer.CompilePlugin("oauth", &oauthPlugins)
					assert.Nil(t, err)
					assert.NotEqual(t, "", oauthPlugins.Config.Database)
					assert.NotEqual(t, "", oauthPlugins.Config.AccessTokenTimeoutRaw)
					assert.NotEqual(t, "", oauthPlugins.Config.AuthorizationCodeTimeoutRaw)
				})
			})

			t.Run("Plugin is not exists", func(t *testing.T) {
				t.Run("Return error", func(t *testing.T) {
					plugins := map[string]entity.Plugin{}
					pluginBearer := loader.PluginBearer(plugins)

					oauthPlugins := OauthPlugin{}
					err := pluginBearer.CompilePlugin("oauth", &oauthPlugins)
					assert.NotNil(t, err)
				})
			})

			t.Run("Unmarshal failed because of injected struct is not struct ", func(t *testing.T) {
				t.Run("Return error", func(t *testing.T) {
					plugins := map[string]entity.Plugin{
						"oauth": {Plugin: "oauth", Raw: []byte(PluginConfigNormal1)},
					}
					pluginBearer := loader.PluginBearer(plugins)

					err := pluginBearer.CompilePlugin("oauth", "not struct, should be failed")
					assert.NotNil(t, err)
				})
			})
		})
	})

	t.Run("ForEach", func(t *testing.T) {
		t.Run("Given callback function", func(t *testing.T) {
			t.Run("Run gracefully", func(t *testing.T) {
				plugins := map[string]entity.Plugin{
					"oauth": {Plugin: "oauth"},
					"cache": {Plugin: "cache"},
				}

				pluginBearer := loader.PluginBearer(plugins)
				count := 0
				pluginBearer.ForEach(func(pluginName string) error {
					count++
					return nil
				})

				assert.Equal(t, count, len(plugins))
			})

			t.Run("There is error in callback function", func(t *testing.T) {
				t.Run("Break the iteration", func(t *testing.T) {
					plugins := map[string]entity.Plugin{
						"oauth": {Plugin: "oauth"},
						"cache": {Plugin: "cache"},
					}

					pluginBearer := loader.PluginBearer(plugins)
					count := 0
					pluginBearer.ForEach(func(pluginName string) error {
						if count == 1 {
							return errors.New("Stop the iteration")
						}
						count++
						return nil
					})

					assert.Equal(t, count, 1)
				})
			})
		})
	})
}
