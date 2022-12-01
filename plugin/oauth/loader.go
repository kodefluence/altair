package oauth

import (
	"fmt"

	"github.com/kodefluence/altair/core"
	"github.com/kodefluence/altair/module"
)

// Provide create new oauth plugin provider
func Load(appBearer core.AppBearer, dbBearer core.DatabaseBearer, pluginBearer core.PluginBearer, apiError module.ApiError, appModule module.App) error {
	if appBearer.Config().PluginExists("oauth") == false {
		return nil
	}

	version, err := pluginBearer.PluginVersion("oauth")
	if err != nil {
		return err
	}

	switch version {
	case "1.0":
		return version_1_0(dbBearer, pluginBearer, apiError, appModule)
	default:
		return fmt.Errorf("undefined template version: %s for metric plugin", version)
	}
}

// Provide create new oauth plugin provider
func LoadCommand(appBearer core.AppBearer, dbBearer core.DatabaseBearer, pluginBearer core.PluginBearer, appModule module.App) error {
	if appBearer.Config().PluginExists("oauth") == false {
		return nil
	}

	version, err := pluginBearer.PluginVersion("oauth")
	if err != nil {
		return err
	}

	switch version {
	case "1.0":
		return version_1_0_command(dbBearer, pluginBearer, appModule)
	default:
		return fmt.Errorf("undefined template version: %s for metric plugin", version)
	}
}
