package loader_test

import (
	"fmt"
	"testing"

	"github.com/codefluence-x/altair/entity"
	"github.com/codefluence-x/altair/loader"
	"github.com/stretchr/testify/assert"
)

func TestApp(t *testing.T) {

	t.Run("Compile", func(t *testing.T) {
		t.Run("Given config path", func(t *testing.T) {
			t.Run("Normal scenario", func(t *testing.T) {
				t.Run("Return app config", func(t *testing.T) {
					configPath := "./app_normal/"
					fileName := "app.yml"

					generateTempTestFiles(configPath, AppConfigNormal, fileName, 0666)

					expectedAppConfig := entity.NewAppConfig([]string{"oauth"})

					appConfig, err := loader.App().Compile(fmt.Sprintf("%s%s", configPath, fileName))
					assert.Nil(t, err)

					assert.Equal(t, expectedAppConfig.Plugins(), appConfig.Plugins())

					removeTempTestFiles(configPath)
				})
			})

			t.Run("File not found", func(t *testing.T) {
				t.Run("Return error", func(t *testing.T) {
					configPath := "./app_not_found/"
					fileName := "app.yml"

					generateTempTestFiles(configPath, AppConfigNormal, fileName, 0666)

					appConfig, err := loader.App().Compile(fmt.Sprintf("%s%s", configPath, "should_be_not_found_yml"))
					assert.NotNil(t, err)
					assert.Nil(t, appConfig)

					removeTempTestFiles(configPath)
				})
			})

			t.Run("Template error", func(t *testing.T) {
				t.Run("Return error", func(t *testing.T) {
					configPath := "./app_template_error/"
					fileName := "app.yml"

					generateTempTestFiles(configPath, AppConfigTemplateError, fileName, 0666)

					appConfig, err := loader.App().Compile(fmt.Sprintf("%s%s", configPath, fileName))
					assert.NotNil(t, err)
					assert.Nil(t, appConfig)

					removeTempTestFiles(configPath)
				})
			})

			t.Run("Unmarshal failed", func(t *testing.T) {
				t.Run("Return error", func(t *testing.T) {
					configPath := "./app_unmarshal_failed/"
					fileName := "app.yml"

					generateTempTestFiles(configPath, AppConfigUnmarshalError, fileName, 0666)

					appConfig, err := loader.App().Compile(fmt.Sprintf("%s%s", configPath, fileName))
					assert.NotNil(t, err)
					assert.Nil(t, appConfig)

					removeTempTestFiles(configPath)
				})
			})
		})
	})
}
