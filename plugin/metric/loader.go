package metric

import (
	"fmt"

	"github.com/kodefluence/altair/core"
	"github.com/kodefluence/altair/plugin/metric/entity"
	"github.com/kodefluence/altair/plugin/metric/module/dummy"
	"github.com/kodefluence/altair/plugin/metric/module/prometheus"
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
		var metricPlugin entity.MetricPlugin
		if err := pluginBearer.CompilePlugin("metric", &metricPlugin); err != nil {
			return err
		}

		switch metricPlugin.Config.Provider {
		case "prometheus":
			prometheus.Load(appBearer)
		default:
			return fmt.Errorf("Metric plugin `%s` is currently not supported", metricPlugin.Config.Provider)
		}
		return nil
	default:
		return fmt.Errorf("undefined template version: %s for metric plugin", version)
	}
}
