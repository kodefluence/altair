package oauth

import (
	"time"

	"github.com/kodefluence/altair/core"
	"github.com/kodefluence/altair/module"
	"github.com/kodefluence/altair/plugin/oauth/entity"
	"github.com/kodefluence/altair/plugin/oauth/module/application"
	"github.com/kodefluence/altair/plugin/oauth/module/formatter"
	"github.com/kodefluence/altair/plugin/oauth/repository/mysql"
)

func version_1_0(appBearer core.AppBearer, dbBearer core.DatabaseBearer, pluginBearer core.PluginBearer, apiError module.ApiError) error {
	var oauthPluginConfig entity.OauthPlugin
	if err := pluginBearer.CompilePlugin("oauth", &oauthPluginConfig); err != nil {
		return err
	}

	sqldb, _, err := dbBearer.Database(oauthPluginConfig.DatabaseInstance())
	if err != nil {
		return err
	}

	var accessTokenTimeout time.Duration
	var authorizationCodeTimeout time.Duration

	accessTokenTimeout, err = oauthPluginConfig.AccessTokenTimeout()
	if err != nil {
		return err
	}

	authorizationCodeTimeout, err = oauthPluginConfig.AuthorizationCodeTimeout()
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
	formatter := formatter.Provide(accessTokenTimeout, authorizationCodeTimeout, refreshTokenConfig.Timeout)

	// Application
	// Loading controller for oauth applications
	application.Load(appBearer, sqldb, oauthApplicationRepo, formatter, apiError)

	return nil
}
