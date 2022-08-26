package cfg_test

import (
	"fmt"
	"testing"

	"github.com/kodefluence/altair/cfg"
	"github.com/kodefluence/altair/entity"
	"github.com/kodefluence/altair/testhelper"
	"github.com/stretchr/testify/assert"
)

func TestApp(t *testing.T) {
	appConfigOption := entity.AppConfigOption{
		Plugins:   []string{"oauth"},
		Port:      1304,
		ProxyHost: "www.local.host",
	}

	appConfigOption.Authorization.Username = "altair"
	appConfigOption.Authorization.Password = "secret"
	appConfigOption.Metric.Interface = "prometheus"

	t.Run("Compile", func(t *testing.T) {
		t.Run("Given config path", func(t *testing.T) {
			t.Run("Normal scenario", func(t *testing.T) {
				t.Run("Return app config", func(t *testing.T) {
					configPath := "./app_normal/"
					fileName := "app.yml"

					testhelper.GenerateTempTestFiles(configPath, AppConfigNormal, fileName, 0666)

					expectedAppConfig := entity.NewAppConfig(appConfigOption)

					appConfig, err := cfg.App().Compile(fmt.Sprintf("%s%s", configPath, fileName))
					assert.Nil(t, err)

					assert.Equal(t, expectedAppConfig.Plugins(), appConfig.Plugins())
					assert.Equal(t, expectedAppConfig.Port(), appConfig.Port())
					assert.Equal(t, expectedAppConfig.ProxyHost(), appConfig.ProxyHost())
					assert.Equal(t, expectedAppConfig.BasicAuthPassword(), appConfig.BasicAuthPassword())
					assert.Equal(t, expectedAppConfig.BasicAuthUsername(), appConfig.BasicAuthUsername())
					assert.Equal(t, expectedAppConfig.Metric().Interface(), appConfig.Metric().Interface())

					testhelper.RemoveTempTestFiles(configPath)
				})
			})

			t.Run("With custom port", func(t *testing.T) {
				appConfigOption := entity.AppConfigOption{
					Plugins:   []string{"oauth"},
					Port:      7001,
					ProxyHost: "www.local.host",
				}

				appConfigOption.Authorization.Username = "altair"
				appConfigOption.Authorization.Password = "secret"

				t.Run("Return app config", func(t *testing.T) {
					configPath := "./app_with_custom_port/"
					fileName := "app.yml"

					testhelper.GenerateTempTestFiles(configPath, AppConfigWithCustomPort, fileName, 0666)

					expectedAppConfig := entity.NewAppConfig(appConfigOption)

					appConfig, err := cfg.App().Compile(fmt.Sprintf("%s%s", configPath, fileName))
					assert.Nil(t, err)

					assert.Equal(t, expectedAppConfig.Plugins(), appConfig.Plugins())
					assert.Equal(t, expectedAppConfig.Port(), appConfig.Port())
					assert.Equal(t, expectedAppConfig.ProxyHost(), appConfig.ProxyHost())
					assert.Equal(t, expectedAppConfig.BasicAuthPassword(), appConfig.BasicAuthPassword())
					assert.Equal(t, expectedAppConfig.BasicAuthUsername(), appConfig.BasicAuthUsername())

					testhelper.RemoveTempTestFiles(configPath)
				})
			})

			t.Run("With custom proxy host", func(t *testing.T) {
				appConfigOption := entity.AppConfigOption{
					Plugins:   []string{"oauth"},
					Port:      1304,
					ProxyHost: "www.altair.id",
				}

				appConfigOption.Authorization.Username = "altair"
				appConfigOption.Authorization.Password = "secret"

				t.Run("Return app config", func(t *testing.T) {
					configPath := "./app_with_custom_proxy_host/"
					fileName := "app.yml"

					testhelper.GenerateTempTestFiles(configPath, AppConfigWithCustomProxyHost, fileName, 0666)

					expectedAppConfig := entity.NewAppConfig(appConfigOption)

					appConfig, err := cfg.App().Compile(fmt.Sprintf("%s%s", configPath, fileName))
					assert.Nil(t, err)

					assert.Equal(t, expectedAppConfig.Plugins(), appConfig.Plugins())
					assert.Equal(t, expectedAppConfig.Port(), appConfig.Port())
					assert.Equal(t, expectedAppConfig.ProxyHost(), appConfig.ProxyHost())
					assert.Equal(t, expectedAppConfig.BasicAuthPassword(), appConfig.BasicAuthPassword())
					assert.Equal(t, expectedAppConfig.BasicAuthUsername(), appConfig.BasicAuthUsername())

					testhelper.RemoveTempTestFiles(configPath)
				})
			})

			t.Run("Empty authorization username", func(t *testing.T) {
				t.Run("Return error", func(t *testing.T) {
					configPath := "./app_empty_username/"
					fileName := "app.yml"

					testhelper.GenerateTempTestFiles(configPath, AppConfigAuthUsernameEmpty, fileName, 0666)

					appConfig, err := cfg.App().Compile(fmt.Sprintf("%s%s", configPath, fileName))
					assert.NotNil(t, err)
					assert.Nil(t, appConfig)

					testhelper.RemoveTempTestFiles(configPath)
				})
			})

			t.Run("Empty authorization password", func(t *testing.T) {
				t.Run("Return error", func(t *testing.T) {
					configPath := "./app_empty_password/"
					fileName := "app.yml"

					testhelper.GenerateTempTestFiles(configPath, AppConfigAuthPasswordEmpty, fileName, 0666)

					appConfig, err := cfg.App().Compile(fmt.Sprintf("%s%s", configPath, fileName))
					assert.NotNil(t, err)
					assert.Nil(t, appConfig)

					testhelper.RemoveTempTestFiles(configPath)
				})
			})

			t.Run("Invalid custom port", func(t *testing.T) {
				t.Run("Return error", func(t *testing.T) {
					configPath := "./app_invalid_custom_port/"
					fileName := "app.yml"

					testhelper.GenerateTempTestFiles(configPath, AppConfigWithInvalidCustomPort, fileName, 0666)

					appConfig, err := cfg.App().Compile(fmt.Sprintf("%s%s", configPath, fileName))
					assert.NotNil(t, err)
					assert.Nil(t, appConfig)

					testhelper.RemoveTempTestFiles(configPath)
				})
			})

			t.Run("File not found", func(t *testing.T) {
				t.Run("Return error", func(t *testing.T) {
					configPath := "./app_not_found/"
					fileName := "app.yml"

					testhelper.GenerateTempTestFiles(configPath, AppConfigNormal, fileName, 0666)

					appConfig, err := cfg.App().Compile(fmt.Sprintf("%s%s", configPath, "should_be_not_found_yml"))
					assert.NotNil(t, err)
					assert.Nil(t, appConfig)

					testhelper.RemoveTempTestFiles(configPath)
				})
			})

			t.Run("Template error", func(t *testing.T) {
				t.Run("Return error", func(t *testing.T) {
					configPath := "./app_template_error/"
					fileName := "app.yml"

					testhelper.GenerateTempTestFiles(configPath, AppConfigTemplateError, fileName, 0666)

					appConfig, err := cfg.App().Compile(fmt.Sprintf("%s%s", configPath, fileName))
					assert.NotNil(t, err)
					assert.Nil(t, appConfig)

					testhelper.RemoveTempTestFiles(configPath)
				})
			})

			t.Run("Unmarshal failed", func(t *testing.T) {
				t.Run("Return error", func(t *testing.T) {
					configPath := "./app_unmarshal_failed/"
					fileName := "app.yml"

					testhelper.GenerateTempTestFiles(configPath, AppConfigUnmarshalError, fileName, 0666)

					appConfig, err := cfg.App().Compile(fmt.Sprintf("%s%s", configPath, fileName))
					assert.NotNil(t, err)
					assert.Nil(t, appConfig)

					testhelper.RemoveTempTestFiles(configPath)
				})
			})
		})
	})
}
