package oauth

import (
	"time"

	"github.com/kodefluence/altair/module"
	"github.com/kodefluence/altair/plugin/oauth/entity"
	"github.com/kodefluence/altair/plugin/oauth/module/application"
	"github.com/kodefluence/altair/plugin/oauth/module/authorization"
	"github.com/kodefluence/altair/plugin/oauth/module/formatter"
	"github.com/kodefluence/altair/plugin/oauth/repository/mysql"
)

func loadV1_0(ctx module.PluginContext) error {
	var oauthPluginConfig entity.OauthPlugin
	if err := ctx.DecodeConfig(&oauthPluginConfig); err != nil {
		return err
	}

	sqldb, _, err := ctx.Database(oauthPluginConfig.DatabaseInstance())
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

	oauthApplicationRepo := mysql.NewOauthApplication()
	oauthAccessTokenRepo := mysql.NewOauthAccessToken()
	oauthAccessGrantRepo := mysql.NewOauthAccessGrant()
	oauthRefreshTokenRepo := mysql.NewOauthRefreshToken()

	formatter := formatter.Provide(accessTokenTimeout, authorizationCodeTimeout, refreshTokenConfig.Timeout)

	application.Load(ctx.AppModule, sqldb, oauthApplicationRepo, formatter, ctx.ApiError)
	authorization.Load(ctx.AppModule, oauthApplicationRepo, oauthAccessTokenRepo, oauthAccessGrantRepo, oauthRefreshTokenRepo, formatter, oauthPluginConfig, sqldb, ctx.ApiError)

	return nil
}

func loadCommandV1_0(ctx module.PluginContext) error {
	var oauthPluginConfig entity.OauthPlugin
	if err := ctx.DecodeConfig(&oauthPluginConfig); err != nil {
		return err
	}

	sqldb, _, err := ctx.Database(oauthPluginConfig.DatabaseInstance())
	if err != nil {
		return err
	}

	accessTokenTimeout, err := oauthPluginConfig.AccessTokenTimeout()
	if err != nil {
		return err
	}

	authorizationCodeTimeout, err := oauthPluginConfig.AuthorizationCodeTimeout()
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

	oauthApplicationRepo := mysql.NewOauthApplication()
	formatter := formatter.Provide(accessTokenTimeout, authorizationCodeTimeout, refreshTokenConfig.Timeout)

	application.LoadCommand(ctx.AppModule, sqldb, oauthApplicationRepo, formatter, ctx.ApiError)

	return nil
}
