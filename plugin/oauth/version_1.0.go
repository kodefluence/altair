package oauth

import (
	"time"

	"github.com/kodefluence/altair/core"
	"github.com/kodefluence/altair/module"
	"github.com/kodefluence/altair/plugin/oauth/entity"
	"github.com/kodefluence/altair/plugin/oauth/module/application"
	"github.com/kodefluence/altair/plugin/oauth/module/authorization"
	"github.com/kodefluence/altair/plugin/oauth/module/formatter"
	"github.com/kodefluence/altair/plugin/oauth/module/migration"
	"github.com/kodefluence/altair/plugin/oauth/repository/mysql"
)

func version_1_0(dbBearer core.DatabaseBearer, pluginBearer core.PluginBearer, apiError module.ApiError, appModule module.App) error {
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
	oauthAccessTokenRepo := mysql.NewOauthAccessToken()
	oauthAccessGrantRepo := mysql.NewOauthAccessGrant()
	oauthRefreshTokenRepo := mysql.NewOauthRefreshToken()

	// Formatter
	formatter := formatter.Provide(accessTokenTimeout, authorizationCodeTimeout, refreshTokenConfig.Timeout)

	// Application
	// Loading controller for oauth applications and downstream
	application.Load(appModule, sqldb, oauthApplicationRepo, formatter, apiError)

	// Authorization
	// Loading controller for authorization and downstream
	authorization.Load(appModule, oauthApplicationRepo, oauthAccessTokenRepo, oauthAccessGrantRepo, oauthRefreshTokenRepo, formatter, oauthPluginConfig, sqldb, apiError)

	return nil
}

func version_1_0_command(dbBearer core.DatabaseBearer, pluginBearer core.PluginBearer, apiError module.ApiError, appModule module.App) error {
	var oauthPluginConfig entity.OauthPlugin
	if err := pluginBearer.CompilePlugin("oauth", &oauthPluginConfig); err != nil {
		return err
	}

	sqldb, sqldbconfig, err := dbBearer.Database(oauthPluginConfig.DatabaseInstance())
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

	// Formatter
	formatter := formatter.Provide(accessTokenTimeout, authorizationCodeTimeout, refreshTokenConfig.Timeout)

	// Migration
	// Set up migration for oauth plugin
	migration.LoadCommand(sqldb, sqldbconfig, appModule)
	application.LoadCommand(appModule, sqldb, oauthApplicationRepo, formatter, nil)

	return nil
}
