package oauth

import (
	"github.com/kodefluence/altair/core"
	"github.com/kodefluence/altair/plugin/oauth/entity"
	"github.com/kodefluence/altair/plugin/oauth/repository/mysql"
)

// Provide create new oauth plugin provider
func Load(appBearer core.AppBearer, dbBearer core.DatabaseBearer, pluginBearer core.PluginBearer) error {
	if appBearer.Config().PluginExists("oauth") == false {
		return nil
	}

	var oauthPluginConfig entity.OauthPlugin
	if err := pluginBearer.CompilePlugin("oauth", &oauthPluginConfig); err != nil {
		return err
	}

	_, _, err := dbBearer.Database(oauthPluginConfig.DatabaseInstance())
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
	_ = mysql.NewOauthApplication()
	_ = mysql.NewOauthAccessToken()
	_ = mysql.NewOauthAccessGrant()
	_ = mysql.NewOauthRefreshToken()

	// Formatter
	// oauthApplicationFormatter := formatter.OauthApplication()
	// oauthModelFormatter := formatter.NewModel(accessTokenTimeout, authorizationCodeTimeout, refreshTokenConfig.Timeout)
	// oauthFormatter := formatter.Oauth()

	// Validator
	// oauthValidator := validator.NewOauth(refreshTokenConfig.Active)

	// Service
	// applicationManager := service.NewApplicationManager(oauthApplicationFormatter, oauthModelFormatter, oauthApplicationModel, oauthValidator, db)
	// authorization := service.NewAuthorization(oauthApplicationModel, oauthAccessTokenModel, oauthAccessGrantModel, oauthRefreshTokenModel, oauthModelFormatter, oauthValidator, oauthFormatter, refreshTokenConfig.Active, db)

	// DownStreamPlugin
	// oauthDownStream := downstream.NewOauth(oauthAccessTokenModel, db)
	// applicationValidationDownStream := downstream.NewApplicationValidation(oauthApplicationModel, db)

	// Controller of /oauth/applications
	// applicationControllerDispatcher := controller.NewApplication()
	// appBearer.InjectController(applicationControllerDispatcher.List(applicationManager))
	// appBearer.InjectController(applicationControllerDispatcher.One(applicationManager))
	// appBearer.InjectController(applicationControllerDispatcher.Create(applicationManager))
	// appBearer.InjectController(applicationControllerDispatcher.Update(applicationManager))

	// Controller of /oauth/authorizations
	// authorizationControllerDispatcher := controller.NewAuthorization()
	// appBearer.InjectController(authorizationControllerDispatcher.Grant(authorization))
	// appBearer.InjectController(authorizationControllerDispatcher.Revoke(authorization))
	// appBearer.InjectController(authorizationControllerDispatcher.Token(authorization))

	// appBearer.InjectDownStreamPlugin(oauthDownStream)
	// appBearer.InjectDownStreamPlugin(applicationValidationDownStream)

	return nil
}