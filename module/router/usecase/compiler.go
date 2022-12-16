package usecase

import (
	"bytes"
	"html/template"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/kodefluence/altair/entity"
	"github.com/kodefluence/altair/util"
	"gopkg.in/yaml.v2"
)

type Compiler struct{}

func NewCompiler() *Compiler {
	return &Compiler{}
}

func (c *Compiler) Compile(routesPath string) ([]entity.RouteObject, error) {
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

func (c *Compiler) compileTemplate(b []byte) ([]byte, error) {
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

func (c *Compiler) walkAllFiles(routesPath string) ([][]byte, error) {
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
		f, _ := util.ReadFileContent(path)
		routeFiles = append(routeFiles, f)
	}

	return routeFiles, nil
}
