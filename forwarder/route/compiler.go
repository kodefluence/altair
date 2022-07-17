package route

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"

	"github.com/google/uuid"
	"github.com/kodefluence/altair/core"
	"github.com/kodefluence/altair/entity"
	"gopkg.in/yaml.v2"
)

type compiler struct{}

// Compiler for all routes yaml file
func Compiler() core.RouteCompiler {
	return &compiler{}
}

func (c *compiler) Compile(routesPath string) ([]entity.RouteObject, error) {
	var routeObjects []entity.RouteObject

	listOfBytes, err := c.walkAllFiles(routesPath)
	if err != nil {
		return routeObjects, err
	}

	for _, b := range listOfBytes {

		compiledBytes, err := c.compileTemplate(b)
		if err != nil {
			return routeObjects, err
		}

		var routeObject entity.RouteObject
		if err := yaml.Unmarshal(compiledBytes, &routeObject); err != nil {
			return routeObjects, err
		}

		if routeObject.Auth == "" {
			routeObject.Auth = "none"
		}

		routeObjects = append(routeObjects, routeObject)
	}

	return routeObjects, nil
}

func (c *compiler) compileTemplate(b []byte) ([]byte, error) {
	tpl, err := template.New(uuid.New().String()).Funcs(template.FuncMap{
		"env": os.Getenv,
	}).Parse(string(b))
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBufferString("")
	err = tpl.Execute(buf, nil)
	return buf.Bytes(), err
}

func (c *compiler) walkAllFiles(routesPath string) ([][]byte, error) {
	var files []string
	var routeFiles [][]byte

	err := filepath.Walk(routesPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || (filepath.Ext(path) != ".yaml" && filepath.Ext(path) != ".yml") {
			return nil
		}
		files = append(files, path)
		return nil
	})
	if err != nil {
		return routeFiles, err
	}

	for _, path := range files {
		f, _ := ioutil.ReadFile(path)
		routeFiles = append(routeFiles, f)
	}

	return routeFiles, nil
}
