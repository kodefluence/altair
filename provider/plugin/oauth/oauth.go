package oauth

import (
	"time"

	"github.com/codefluence-x/altair/core"
	"github.com/codefluence-x/altair/provider/plugin/oauth/controller"
	"github.com/codefluence-x/altair/provider/plugin/oauth/downstream"
	"github.com/codefluence-x/altair/provider/plugin/oauth/entity"
	"github.com/codefluence-x/altair/provider/plugin/oauth/formatter"
	"github.com/codefluence-x/altair/provider/plugin/oauth/model"
	"github.com/codefluence-x/altair/provider/plugin/oauth/service"
	"github.com/codefluence-x/altair/provider/plugin/oauth/validator"
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

	// Model
	oauthApplicationModel := model.OauthApplication(db)
	oauthAccessTokenModel := model.OauthAccessToken(db)
	oauthAccessGrantModel := model.NewOauthAccessGrant(db)

	// Formatter
	oauthApplicationFormatter := formatter.OauthApplication()
	oauthModelFormatter := formatter.NewModel(accessTokenTimeout, authorizationCodeTimeout)
	oauthFormatter := formatter.Oauth()

	// Validator
	oauthValidator := validator.Oauth()

	// Service
	applicationManager := service.ApplicationManager(oauthApplicationFormatter, oauthModelFormatter, oauthApplicationModel, oauthValidator)
	authorization := service.NewAuthorization(oauthApplicationModel, oauthAccessTokenModel, oauthAccessGrantModel, oauthModelFormatter, oauthValidator, oauthFormatter)

	// DownStreamPlugin
	oauthDownStream := downstream.NewOauth(oauthAccessTokenModel)
	applicationValidationDownStream := downstream.NewApplicationValidation(oauthApplicationModel)

	// Controller of /oauth/applications
	applicationControllerDispatcher := controller.NewApplication()
	appBearer.InjectController(applicationControllerDispatcher.List(applicationManager))
	appBearer.InjectController(applicationControllerDispatcher.One(applicationManager))
	appBearer.InjectController(applicationControllerDispatcher.Create(applicationManager))
	appBearer.InjectController(applicationControllerDispatcher.Update(applicationManager))

	// Controller of /oauth/authorizations
	authorizationControllerDispatcher := controller.NewAuthorization()
	appBearer.InjectController(authorizationControllerDispatcher.Grant(authorization))
	appBearer.InjectController(authorizationControllerDispatcher.Revoke(authorization))
	appBearer.InjectController(authorizationControllerDispatcher.Token(authorization))

	appBearer.InjectDownStreamPlugin(oauthDownStream)
	appBearer.InjectDownStreamPlugin(applicationValidationDownStream)

	return nil
}
