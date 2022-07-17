package cfg

import (
	"errors"

	"github.com/kodefluence/altair/core"
	"github.com/kodefluence/altair/entity"
	"gopkg.in/yaml.v2"
)

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
	_, ok := p.plugins[pluginName]
	return ok
}

func (p *pluginBearer) CompilePlugin(pluginName string, injectedStruct interface{}) error {
	if p.ConfigExists(pluginName) == false {
		return errors.New("Plugin is not exists")
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
