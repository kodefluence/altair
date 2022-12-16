package metric

import (
	"fmt"

	"github.com/kodefluence/altair/core"
	"github.com/kodefluence/altair/module"
	"github.com/kodefluence/altair/plugin/metric/module/dummy"
)

func Load(appBearer core.AppBearer, pluginBearer core.PluginBearer, appModule module.App) error {
	if !appBearer.Config().PluginExists("metric") {
		dummy.Load(appModule)
		return nil
	}

	version, err := pluginBearer.PluginVersion("metric")
	if err != nil {
		return err
	}

	switch version {
	case "1.0":
		return version_1_0(appModule, pluginBearer)
	default:
		return fmt.Errorf("undefined template version: %s for metric plugin", version)
	}
}
