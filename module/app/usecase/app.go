package usecase

import "github.com/kodefluence/altair/module"

type App struct {
	controller module.Controller
}

func NewApp(controller module.Controller) *App {
	return &App{
		controller: controller,
	}
}

func (a *App) Config() module.Config {
	return nil
}

func (a *App) Controller() module.Controller {
	return a.controller
}
