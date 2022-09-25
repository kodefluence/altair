package metric

import (
	"fmt"

	"github.com/kodefluence/altair/core"
	"github.com/kodefluence/altair/plugin/metric/module/dummy"
)

func Load(appBearer core.AppBearer, pluginBearer core.PluginBearer) error {
	if appBearer.Config().PluginExists("metric") == false {
		dummy.Load(appBearer)
		return nil
	}

	version, err := pluginBearer.PluginVersion("metric")
	if err != nil {
		return err
	}

	switch version {
	case "1.0":
		return version_1_0(appBearer, pluginBearer)
	default:
		return fmt.Errorf("undefined template version: %s for metric plugin", version)
	}
}
