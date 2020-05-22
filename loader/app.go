package loader

import (
	"io/ioutil"

	"github.com/codefluence-x/altair/core"
	"github.com/codefluence-x/altair/entity"
	"gopkg.in/yaml.v2"
)

type app struct{}

type appConfig struct {
	Plugins []string `yaml:"plugins"`
}

func App() core.AppLoader {
	return &app{}
}

func (a *app) Compile(configPath string) (core.AppConfig, error) {
	contents, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	compiledContents, err := compileTemplate(contents)
	if err != nil {
		return nil, err
	}

	var config appConfig

	if err := yaml.Unmarshal(compiledContents, &config); err != nil {
		return nil, err
	}

	return entity.NewAppConfig(config.Plugins), nil
}
