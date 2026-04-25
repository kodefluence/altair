package cfg

import (
	"errors"

	"gopkg.in/yaml.v2"

	"github.com/kodefluence/altair/core"
	"github.com/kodefluence/altair/entity"
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

// DecodeConfig unmarshals the plugin's raw YAML into target. Plugin config
// structs typically wrap their fields under `Config yaml:"config"`, so the
// same raw bytes populate the inner struct while ignoring the top-level
// `plugin:`/`version:` keys.
func (p *pluginBearer) DecodeConfig(pluginName string, target interface{}) error {
	if !p.ConfigExists(pluginName) {
		return errPluginNotFound
	}

	return yaml.Unmarshal(p.plugins[pluginName].Raw, target)
}

func (p *pluginBearer) ForEach(callbackFunc func(pluginName string) error) {
	for _, plugin := range p.plugins {
		if err := callbackFunc(plugin.Plugin); err != nil {
			break
		}
	}
}
