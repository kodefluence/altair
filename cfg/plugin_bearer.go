package cfg

import (
	"errors"

	"github.com/kodefluence/altair/core"
	"github.com/kodefluence/altair/entity"
	"gopkg.in/yaml.v2"
)

var errPluginNotFound = errors.New("Plugin is not exists")

type pluginBearer struct {
	plugins map[string]entity.Plugin
}

func PluginBearer(plugins map[string]entity.Plugin) core.PluginBearer {
	return &pluginBearer{plugins: plugins}
}

func (p *pluginBearer) Length() int {
	return len(p.plugins)
}

func (p *pluginBearer) ConfigExists(pluginName string) bool {
	_, err := p.PluginVersion(pluginName)
	return err == nil
}

func (p *pluginBearer) PluginVersion(pluginName string) (string, error) {
	plugin, ok := p.plugins[pluginName]
	if !ok {
		return "", errPluginNotFound
	}
	return plugin.Version, nil
}

func (p *pluginBearer) CompilePlugin(pluginName string, injectedStruct interface{}) error {
	if p.ConfigExists(pluginName) == false {
		return errPluginNotFound
	}

	return yaml.Unmarshal(p.plugins[pluginName].Raw, injectedStruct)
}

func (p *pluginBearer) ForEach(callbackFunc func(pluginName string) error) {
	for _, plugin := range p.plugins {
		if err := callbackFunc(plugin.Plugin); err != nil {
			break
		}
	}
}
