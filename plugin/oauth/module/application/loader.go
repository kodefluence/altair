package application

import (
	"github.com/kodefluence/altair/module"
	"github.com/kodefluence/altair/plugin/oauth/module/application/controller/command"
	"github.com/kodefluence/altair/plugin/oauth/module/application/controller/downstream"
	"github.com/kodefluence/altair/plugin/oauth/module/application/controller/http"
	"github.com/kodefluence/altair/plugin/oauth/module/application/usecase"
	"github.com/kodefluence/monorepo/db"
)

func Load(
	appModule module.App,
	sqldb db.DB,
	oauthApplicationRepo usecase.OauthApplicationRepository,
	formatter usecase.Formatter,
	apiError module.ApiError,
) {
	applicationManager := usecase.NewApplicationManager(sqldb, oauthApplicationRepo, apiError, formatter)
	appModule.Controller().InjectHTTP(
		http.NewCreate(applicationManager, apiError),
		http.NewOne(applicationManager, apiError),
		http.NewList(applicationManager, apiError),
		http.NewUpdate(applicationManager, apiError),
	)
	appModule.Controller().InjectDownstream(downstream.NewApplicationValidation(oauthApplicationRepo, sqldb))
	appModule.Controller().InjectCommand(command.NewCreateOauthApplication(applicationManager))
}
