package metric

import (
	"fmt"

	"github.com/kodefluence/altair/core"
	"github.com/kodefluence/altair/module"
	"github.com/kodefluence/altair/plugin/metric/entity"
	"github.com/kodefluence/altair/plugin/metric/module/prometheus"
)

func version_1_0(appModule module.App, pluginBearer core.PluginBearer) error {
	var metricPlugin entity.MetricPlugin
	if err := pluginBearer.CompilePlugin("metric", &metricPlugin); err != nil {
		return err
	}

	switch metricPlugin.Config.Provider {
	case "prometheus":
		prometheus.Load(appModule)
	default:
		return fmt.Errorf("Metric plugin `%s` is currently not supported", metricPlugin.Config.Provider)
	}
	return nil
}
