package oauth

import (
	"github.com/kodefluence/altair/core"
	"github.com/kodefluence/altair/provider/plugin/oauth/downstream"
	"github.com/kodefluence/altair/provider/plugin/oauth/entity"
	"github.com/kodefluence/altair/provider/plugin/oauth/model"
)

// Provide create new oauth plugin provider
func Provide(appBearer core.AppBearer, dbBearer core.DatabaseBearer, pluginBearer core.PluginBearer) error {
	if appBearer.Config().PluginExists("oauth") == false {
		return nil
	}

	var oauthPluginConfig entity.OauthPlugin

	if err := pluginBearer.CompilePlugin("oauth", &oauthPluginConfig); err != nil {
		return err
	}

	db, _, err := dbBearer.Database(oauthPluginConfig.DatabaseInstance())
	if err != nil {
		return err
	}

	var refreshTokenConfig entity.RefreshTokenConfig
	refreshTokenConfig.Active = oauthPluginConfig.Config.RefreshToken.Active
	if refreshTokenConfig.Active {
		refreshTokenTimeout, err := oauthPluginConfig.RefreshTokenTimeout()
		if err != nil {
			return err
		}
		refreshTokenConfig.Timeout = refreshTokenTimeout
	}

	// Model
	oauthAccessTokenModel := model.NewOauthAccessToken()

	// DownStreamPlugin
	oauthDownStream := downstream.NewOauth(oauthAccessTokenModel, db)

	appBearer.InjectDownStreamPlugin(oauthDownStream)

	return nil
}
