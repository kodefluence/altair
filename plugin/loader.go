package plugin

import (
	"github.com/kodefluence/altair/core"
	"github.com/kodefluence/altair/plugin/metric"
)

// Load plugin for altair
// TODO: Unit test, open for contributions.
func Load(appBearer core.AppBearer, pluginBearer core.PluginBearer) error {
	if err := metric.Load(appBearer, pluginBearer); err != nil {
		return err
	}

	return nil
}
