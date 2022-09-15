package authorization

import (
	"github.com/kodefluence/altair/core"
	"github.com/kodefluence/altair/module"
	"github.com/kodefluence/altair/plugin/oauth/module/authorization/usecase"
	"github.com/kodefluence/monorepo/db"
)

func Load(
	appBearer core.AppBearer,
	sqldb db.DB,
	oauthApplicationRepo usecase.OauthApplicationRepository,
	formatter usecase.Formatter,
	apiError module.ApiError,
) {

}
