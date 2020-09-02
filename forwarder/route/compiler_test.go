package route_test

import (
	"bytes"
	"os"
	"testing"
	"text/template"

	"github.com/codefluence-x/altair/entity"
	"github.com/codefluence-x/altair/forwarder/route"
	"github.com/codefluence-x/altair/mock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func TestCompiler(t *testing.T) {
	t.Run("Compile", func(t *testing.T) {
		t.Run("Given route path", func(t *testing.T) {
			t.Run("Normal scenario", func(t *testing.T) {
				routesPath := "./routes_gracefully/"

				generateAllTempTestFiles(routesPath, ExampleRoutesGracefully)

				t.Run("Return []entity.RouteObject and nil", func(t *testing.T) {
					c := route.Compiler()
					routeObjects, err := c.Compile(routesPath)

					var expectedRouteObject entity.RouteObject
					b, _ := compileTemplate([]byte(ExampleRoutesGracefully))
					_ = yaml.Unmarshal(b, &expectedRouteObject)

					assert.Nil(t, err)
					assert.Greater(t, len(routeObjects), 0)
					assert.Equal(t, expectedRouteObject, routeObjects[0])
				})

				mock.RemoveTempTestFiles(routesPath)
			})

			t.Run("Include no yaml files", func(t *testing.T) {
				routesPath := "./routes_gracefully_include_no_yaml_files/"

				generateAllTempTestFiles(routesPath, ExampleRoutesGracefully)
				mock.GenerateTempTestFiles(routesPath, "", "not_included.txt", 0666)

				t.Run("Return []entity.RouteObject and nil", func(t *testing.T) {
					c := route.Compiler()
					routeObjects, err := c.Compile(routesPath)

					var expectedRouteObject entity.RouteObject
					b, _ := compileTemplate([]byte(ExampleRoutesGracefully))
					_ = yaml.Unmarshal(b, &expectedRouteObject)

					assert.Nil(t, err)
					assert.Greater(t, len(routeObjects), 0)
					assert.Equal(t, expectedRouteObject, routeObjects[0])
				})

				mock.RemoveTempTestFiles(routesPath)
			})

			t.Run("Yaml unmarshal error", func(t *testing.T) {
				routesPath := "./routes_yaml_unmarshal_error/"

				generateAllTempTestFiles(routesPath, ExampleRoutesYamlError)

				t.Run("Return error", func(t *testing.T) {
					c := route.Compiler()
					routeObjects, err := c.Compile(routesPath)

					assert.NotNil(t, err)
					assert.Equal(t, 0, len(routeObjects))
				})

				mock.RemoveTempTestFiles(routesPath)
			})

			t.Run("No auth set", func(t *testing.T) {
				routesPath := "./routes_no_auth_set/"

				generateAllTempTestFiles(routesPath, ExampleRoutesWithNoAuth)

				t.Run("Return nil and auth will set to none", func(t *testing.T) {
					c := route.Compiler()
					routeObjects, err := c.Compile(routesPath)

					assert.Nil(t, err)
					assert.Equal(t, 1, len(routeObjects))
					assert.Equal(t, "none", routeObjects[0].Auth)
				})

				mock.RemoveTempTestFiles(routesPath)
			})

			t.Run("Template parsing error", func(t *testing.T) {
				routesPath := "./routes_template_parsing_error/"

				generateAllTempTestFiles(routesPath, ExampleTemplateParsingError)

				t.Run("Return error", func(t *testing.T) {
					c := route.Compiler()
					routeObjects, err := c.Compile(routesPath)

					assert.NotNil(t, err)
					assert.Equal(t, 0, len(routeObjects))
				})

				mock.RemoveTempTestFiles(routesPath)
			})

			t.Run("Dir is not exists", func(t *testing.T) {
				routesPath := "./routes_dir_not_exists/"
				t.Run("Return error", func(t *testing.T) {
					c := route.Compiler()
					routeObjects, err := c.Compile(routesPath)

					assert.NotNil(t, err)
					assert.Equal(t, 0, len(routeObjects))
				})
			})
		})
	})
}

func generateAllTempTestFiles(routesPath, content string) {
	err := os.Mkdir(routesPath, os.ModePerm)
	if err != nil {
		panic(err)
	}

	mock.GenerateTempTestFiles(routesPath, content, "app.yml", 0666)
}

func compileTemplate(b []byte) ([]byte, error) {
	tpl, err := template.New(uuid.New().String()).Funcs(map[string]interface{}{
		"env": os.Getenv,
	}).Parse(string(b))
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBufferString("")
	err = tpl.Execute(buf, nil)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
