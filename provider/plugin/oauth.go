package plugin

import (
	"github.com/codefluence-x/altair/core"
	"github.com/codefluence-x/altair/provider/plugin/oauth"
)

func Oauth(appBearer core.AppBearer, dbBearer core.DatabaseBearer, pluginBearer core.PluginBearer) error {
	return oauth.Provide(appBearer, dbBearer, pluginBearer)
}
