package route_test

import (
	"bytes"
	"fmt"
	"os"
	"testing"
	"text/template"

	"github.com/codefluence-x/altair/entity"
	"github.com/codefluence-x/altair/forwarder/route"
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

				removeAllTemplateTestFiles(routesPath)
			})

			t.Run("Include no yaml files", func(t *testing.T) {
				routesPath := "./routes_gracefully_include_no_yaml_files/"

				generateAllTempTestFiles(routesPath, ExampleRoutesGracefully)
				generateFiles(routesPath, "", "not_included.txt", 0666)

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

				removeAllTemplateTestFiles(routesPath)
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

				removeAllTemplateTestFiles(routesPath)
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

				removeAllTemplateTestFiles(routesPath)
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

	generateFiles(routesPath, content, "app.yml", 0666)
}

func generateFiles(routesPath, content, fileName string, mode os.FileMode) {
	f, err := os.OpenFile(fmt.Sprintf("%s%s", routesPath, fileName), os.O_RDWR|os.O_CREATE|os.O_TRUNC, mode)
	if err != nil {
		panic(err)
	}

	_, err = f.WriteString(content)
	if err != nil {
		panic(err)
	}
}

func removeAllTemplateTestFiles(routesPath string) {
	err := os.RemoveAll(routesPath)
	if err != nil {
		panic(err)
	}
}

func compileTemplate(b []byte) ([]byte, error) {
	tpl, err := template.New(uuid.New().String()).Parse(string(b))
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBufferString("")
	tpl = tpl.Funcs(map[string]interface{}{
		"env": os.Getenv,
	})
	err = tpl.Execute(buf, nil)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
