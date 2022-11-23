package authorization

import (
	"github.com/kodefluence/altair/module"
	"github.com/kodefluence/altair/plugin/oauth/entity"
	"github.com/kodefluence/altair/plugin/oauth/module/authorization/controller/http"
	"github.com/kodefluence/altair/plugin/oauth/module/authorization/usecase"
	"github.com/kodefluence/monorepo/db"
)

func Load(
	appModule module.App,
	oauthApplicationRepo usecase.OauthApplicationRepository,
	oauthAccessTokenRepo usecase.OauthAccessTokenRepository,
	oauthAccessGrantRepo usecase.OauthAccessGrantRepository,
	oauthRefreshTokenRepo usecase.OauthRefreshTokenRepository,
	formatter usecase.Formatter,
	config entity.OauthPlugin,
	sqldb db.DB,
	apiError module.ApiError,
) {
	authorizationUsecase := usecase.NewAuthorization(oauthApplicationRepo, oauthAccessTokenRepo, oauthAccessGrantRepo, oauthRefreshTokenRepo, formatter, config, sqldb, apiError)

	appModule.Controller().InjectHTTP(
		http.NewGrant(authorizationUsecase, apiError),
		http.NewToken(authorizationUsecase, apiError),
		http.NewRevoke(authorizationUsecase, apiError),
	)

	// appModule.InjectDownStreamPlugin(downstream.NewOauth(oauthAccessTokenRepo, sqldb))
}
