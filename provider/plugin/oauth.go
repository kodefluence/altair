package plugin

import (
	"github.com/kodefluence/altair/core"
	"github.com/kodefluence/altair/provider/plugin/oauth"
)

func Oauth(appBearer core.AppBearer, dbBearer core.DatabaseBearer, pluginBearer core.PluginBearer) error {
	return oauth.Provide(appBearer, dbBearer, pluginBearer)
}
