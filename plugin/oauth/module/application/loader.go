package application

import (
	"github.com/kodefluence/altair/core"
	"github.com/kodefluence/altair/module"
	"github.com/kodefluence/altair/plugin/oauth/module/application/controller/http"
	"github.com/kodefluence/altair/plugin/oauth/module/application/usecase"
	"github.com/kodefluence/monorepo/db"
)

func Load(
	appBearer core.AppBearer,
	sqldb db.DB,
	oauthApplicationRepo usecase.OauthApplicationRepository,
	formatter usecase.Formatter,
	apiError module.ApiError,
) {
	applicationManager := usecase.NewApplicationManager(sqldb, oauthApplicationRepo, apiError, formatter)
	appBearer.InjectController(http.NewCreate(applicationManager, apiError))
	appBearer.InjectController(http.NewOne(applicationManager, apiError))
	appBearer.InjectController(http.NewList(applicationManager, apiError))
	appBearer.InjectController(http.NewUpdate(applicationManager, apiError))
}
