package plugin

import (
	"github.com/kodefluence/altair/core"
	"github.com/kodefluence/altair/module"
	"github.com/kodefluence/altair/plugin/metric"
	"github.com/kodefluence/altair/plugin/oauth"
)

// Load plugin for altair
// TODO: Unit test, open for contributions.
func Load(appBearer core.AppBearer, pluginBearer core.PluginBearer, dbBearer core.DatabaseBearer, apiError module.ApiError, appModule module.App) error {
	if err := metric.Load(appBearer, pluginBearer, appModule); err != nil {
		return err
	}

	if err := oauth.Load(appBearer, dbBearer, pluginBearer, apiError, appModule); err != nil {
		return err
	}

	return nil
}

// Load plugin for altair command
// TODO: Unit test, open for contributions.
func LoadCommand(appBearer core.AppBearer, pluginBearer core.PluginBearer, dbBearer core.DatabaseBearer, apiError module.ApiError, appModule module.App) error {
	if err := oauth.LoadCommand(appBearer, dbBearer, pluginBearer, apiError, appModule); err != nil {
		return err
	}

	return nil
}
