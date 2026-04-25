package cfg

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"

	"github.com/kodefluence/altair/core"
	"github.com/kodefluence/altair/entity"
)

type plugin struct{}

func Plugin() core.PluginLoader {
	return &plugin{}
}

func (p *plugin) Compile(pluginPath string) (core.PluginBearer, error) {
	var pluginList = map[string]entity.Plugin{}

	listOfBytes, err := p.walkAllFiles(pluginPath)
	if err != nil {
		return nil, err
	}

	for _, b := range listOfBytes {
		var plugin entity.Plugin

		compiledBytes, err := compileTemplate(b)
		if err != nil {
			return nil, err
		}

		err = yaml.Unmarshal(compiledBytes, &plugin)
		if err != nil {
			return nil, err
		}

		if plugin.Plugin == "" {
			return nil, errors.New("plugin config missing required `plugin` field")
		}
		if plugin.Version == "" {
			return nil, fmt.Errorf("plugin `%s` config missing required `version` field", plugin.Plugin)
		}

		if _, ok := pluginList[plugin.Plugin]; ok {
			return nil, fmt.Errorf("Plugin `%s` already defined", plugin.Plugin)
		}

		plugin.Raw = compiledBytes
		pluginList[plugin.Plugin] = plugin
	}

	return PluginBearer(pluginList), nil
}

func (p *plugin) walkAllFiles(pluginPath string) ([][]byte, error) {
	var files []string
	var routeFiles [][]byte

	err := filepath.Walk(pluginPath, func(path string, info os.FileInfo, err error) error {
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
		content, err := readPluginFile(path)
		if err != nil {
			return nil, err
		}
		routeFiles = append(routeFiles, content)
	}

	return routeFiles, nil
}

func readPluginFile(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return io.ReadAll(f)
}
