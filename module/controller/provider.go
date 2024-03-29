package controller

import (
	"github.com/kodefluence/altair/module"
	"github.com/kodefluence/altair/module/controller/usecase"
	"github.com/spf13/cobra"
)

func Provide(httpInjector usecase.HttpInjector, apiError module.ApiError, rootCommand *cobra.Command) module.Controller {
	return usecase.NewController(httpInjector, apiError, rootCommand)
}
