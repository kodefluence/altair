package oauth

import (
	"github.com/kodefluence/altair/core"
	"github.com/kodefluence/altair/module"
	"github.com/kodefluence/altair/plugin/oauth/entity"
	"github.com/kodefluence/altair/plugin/oauth/module/application"
	"github.com/kodefluence/altair/plugin/oauth/module/formatter"
	"github.com/kodefluence/altair/plugin/oauth/repository/mysql"
)

// Provide create new oauth plugin provider
func Load(appBearer core.AppBearer, dbBearer core.DatabaseBearer, pluginBearer core.PluginBearer, apiError module.ApiError) error {
	if appBearer.Config().PluginExists("oauth") == false {
		return nil
	}

	var oauthPluginConfig entity.OauthPlugin
	if err := pluginBearer.CompilePlugin("oauth", &oauthPluginConfig); err != nil {
		return err
	}

	sqldb, _, err := dbBearer.Database(oauthPluginConfig.DatabaseInstance())
	if err != nil {
		return err
	}

	_, err = oauthPluginConfig.AccessTokenTimeout()
	if err != nil {
		return err
	}

	_, err = oauthPluginConfig.AuthorizationCodeTimeout()
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

	// Repository
	oauthApplicationRepo := mysql.NewOauthApplication()
	_ = mysql.NewOauthAccessToken()
	_ = mysql.NewOauthAccessGrant()
	_ = mysql.NewOauthRefreshToken()

	// Formatter
	formatter := formatter.Provide()

	// Application
	// Loading controller for oauth applications
	application.Load(appBearer, sqldb, oauthApplicationRepo, formatter, apiError)

	return nil
}
