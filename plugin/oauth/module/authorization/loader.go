package authorization

import (
	"github.com/kodefluence/altair/core"
	"github.com/kodefluence/altair/module"
	"github.com/kodefluence/altair/plugin/oauth/entity"
	"github.com/kodefluence/altair/plugin/oauth/module/authorization/controller/http"
	"github.com/kodefluence/altair/plugin/oauth/module/authorization/usecase"
	"github.com/kodefluence/monorepo/db"
)

func Load(
	appBearer core.AppBearer,
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
	appBearer.InjectController(http.NewGrant(authorizationUsecase, apiError))
	appBearer.InjectController(http.NewToken(authorizationUsecase, apiError))
	appBearer.InjectController(http.NewRevoke(authorizationUsecase, apiError))
}
