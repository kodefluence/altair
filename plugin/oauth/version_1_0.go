package oauth

import (
	"errors"
	"time"

	"github.com/kodefluence/altair/module"
	"github.com/kodefluence/altair/plugin/oauth/entity"
	"github.com/kodefluence/altair/plugin/oauth/module/application"
	"github.com/kodefluence/altair/plugin/oauth/module/authorization"
	"github.com/kodefluence/altair/plugin/oauth/module/formatter"
	"github.com/kodefluence/altair/plugin/oauth/repository/mysql"
)

// errMissingDecodeConfig and errMissingDatabase guard against PluginContext
// values constructed outside of plugin.runner.buildContext (which always
// populates both closures). Production code never trips these — they make
// plugins safe to call from external test fixtures.
var (
	errMissingDecodeConfig = errors.New("oauth plugin: PluginContext.DecodeConfig is nil")
	errMissingDatabase     = errors.New("oauth plugin: PluginContext.Database is nil")
)

func loadV1_0(ctx module.PluginContext) error {
	if ctx.DecodeConfig == nil {
		return errMissingDecodeConfig
	}
	if ctx.Database == nil {
		return errMissingDatabase
	}
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
	if ctx.DecodeConfig == nil {
		return errMissingDecodeConfig
	}
	if ctx.Database == nil {
		return errMissingDatabase
	}
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
