package plugin

import (
	"github.com/kodefluence/altair/core"
	"github.com/kodefluence/altair/plugin/metric"
)

// Fabricate plugin for altair
// TODO: Unit test, open for contributions.
func Fabricate(appBearer core.AppBearer, pluginBearer core.PluginBearer) error {
	if err := metric.Provide(appBearer, pluginBearer); err != nil {
		return err
	}

	return nil
}
