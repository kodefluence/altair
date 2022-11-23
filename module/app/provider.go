package app

import (
	"github.com/kodefluence/altair/module"
	"github.com/kodefluence/altair/module/app/usecase"
)

func Provide(controller module.Controller) module.App {
	return usecase.NewApp(controller)
}
