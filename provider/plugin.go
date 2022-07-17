package provider

import (
	"github.com/kodefluence/altair/core"
	"github.com/kodefluence/altair/provider/plugin"
)

func Plugin(appBearer core.AppBearer, dbBearer core.DatabaseBearer, pluginBearer core.PluginBearer) error {
	if err := plugin.Oauth(appBearer, dbBearer, pluginBearer); err != nil {
		return err
	}

	return nil
}
