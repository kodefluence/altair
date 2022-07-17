package cfg_test

import (
	"testing"

	"github.com/kodefluence/altair/cfg"
	"github.com/kodefluence/altair/testhelper"
	"github.com/stretchr/testify/assert"
)

func TestPlugin(t *testing.T) {

	t.Run("Compile", func(t *testing.T) {
		t.Run("Given plugin path", func(t *testing.T) {
			t.Run("Normal scenario", func(t *testing.T) {
				t.Run("Return map string of entity.Plugin", func(t *testing.T) {
					pluginPath := "./normal_scenario_plugin_path/"

					testhelper.GenerateTempTestFiles(pluginPath, PluginConfigNormal1, "oauth.yaml", 0666)
					testhelper.GenerateTempTestFiles(pluginPath, PluginConfigNormal2, "cache.yaml", 0666)

					pluginBearer, err := cfg.Plugin().Compile(pluginPath)

					assert.Nil(t, err)
					assert.Equal(t, 2, pluginBearer.Length())
					assert.True(t, pluginBearer.ConfigExists("oauth"))
					assert.True(t, pluginBearer.ConfigExists("cache"))

					testhelper.RemoveTempTestFiles(pluginPath)
				})
			})

			t.Run("Plugin already defined", func(t *testing.T) {
				t.Run("Return error", func(t *testing.T) {
					pluginPath := "./plugin_already_defined/"

					testhelper.GenerateTempTestFiles(pluginPath, PluginConfigNormal1, "oauth.yaml", 0666)
					testhelper.GenerateTempTestFiles(pluginPath, PluginConfigNormal1, "oauth_2.yaml", 0666)

					pluginBearer, err := cfg.Plugin().Compile(pluginPath)

					assert.NotNil(t, err)
					assert.Nil(t, pluginBearer)

					testhelper.RemoveTempTestFiles(pluginPath)
				})
			})

			t.Run("Yaml unmarshal error", func(t *testing.T) {
				t.Run("Return error", func(t *testing.T) {
					pluginPath := "./plugin_config_yaml_unmarshal_error/"

					testhelper.GenerateTempTestFiles(pluginPath, PluginConfigYamlUnmarshalError, "oauth.yaml", 0666)

					pluginBearer, err := cfg.Plugin().Compile(pluginPath)

					assert.NotNil(t, err)
					assert.Nil(t, pluginBearer)

					testhelper.RemoveTempTestFiles(pluginPath)
				})
			})

			t.Run("Template parsing error error", func(t *testing.T) {
				t.Run("Return error", func(t *testing.T) {
					pluginPath := "./plugin_config_template_parsing_error/"

					testhelper.GenerateTempTestFiles(pluginPath, PluginConfigTemplateParsingError, "oauth.yaml", 0666)

					pluginBearer, err := cfg.Plugin().Compile(pluginPath)

					assert.NotNil(t, err)
					assert.Nil(t, pluginBearer)

					testhelper.RemoveTempTestFiles(pluginPath)
				})
			})

			t.Run("Dir is not exists", func(t *testing.T) {
				t.Run("Return error", func(t *testing.T) {
					pluginPath := "./plugin_config_dir_not_exists/"

					pluginBearer, err := cfg.Plugin().Compile(pluginPath)

					assert.NotNil(t, err)
					assert.Nil(t, pluginBearer)
				})
			})
		})
	})
}
