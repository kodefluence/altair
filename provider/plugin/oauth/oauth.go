package oauth

import (
	"time"

	"github.com/kodefluence/altair/core"
	"github.com/kodefluence/altair/provider/plugin/oauth/controller"
	"github.com/kodefluence/altair/provider/plugin/oauth/downstream"
	"github.com/kodefluence/altair/provider/plugin/oauth/entity"
	"github.com/kodefluence/altair/provider/plugin/oauth/formatter"
	"github.com/kodefluence/altair/provider/plugin/oauth/model"
	"github.com/kodefluence/altair/provider/plugin/oauth/service"
	"github.com/kodefluence/altair/provider/plugin/oauth/validator"
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

	// Model
	oauthApplicationModel := model.NewOauthApplication()
	oauthAccessTokenModel := model.NewOauthAccessToken()
	oauthAccessGrantModel := model.NewOauthAccessGrant()
	oauthRefreshTokenModel := model.NewOauthRefreshToken()

	// // Formatter
	oauthModelFormatter := formatter.NewModel(accessTokenTimeout, authorizationCodeTimeout, refreshTokenConfig.Timeout)
	oauthFormatter := formatter.Oauth()

	// Validator
	oauthValidator := validator.NewOauth(refreshTokenConfig.Active)

	// Service
	authorization := service.NewAuthorization(oauthApplicationModel, oauthAccessTokenModel, oauthAccessGrantModel, oauthRefreshTokenModel, oauthModelFormatter, oauthValidator, oauthFormatter, refreshTokenConfig.Active, db)

	// DownStreamPlugin
	oauthDownStream := downstream.NewOauth(oauthAccessTokenModel, db)

	// Controller of /oauth/authorizations
	authorizationControllerDispatcher := controller.NewAuthorization()
	appBearer.InjectController(authorizationControllerDispatcher.Grant(authorization))
	appBearer.InjectController(authorizationControllerDispatcher.Revoke(authorization))
	appBearer.InjectController(authorizationControllerDispatcher.Token(authorization))

	appBearer.InjectDownStreamPlugin(oauthDownStream)

	return nil
}
