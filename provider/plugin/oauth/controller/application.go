package controller

import (
	"github.com/codefluence-x/altair/core"
	app "github.com/codefluence-x/altair/provider/plugin/oauth/controller/application"
	"github.com/codefluence-x/altair/provider/plugin/oauth/interfaces"
)

type application struct{}

func Application() interfaces.OauthApplicationDispatcher {
	return application{}
}

func (a application) List(applicationManager interfaces.ApplicationManager) core.Controller {
	return app.List(applicationManager)
}

func (a application) One(applicationManager interfaces.ApplicationManager) core.Controller {
	return app.One(applicationManager)
}

func (a application) Create(applicationManager interfaces.ApplicationManager) core.Controller {
	return app.Create(applicationManager)
}

func (a application) Update(applicationManager interfaces.ApplicationManager) core.Controller {
	return app.Update(applicationManager)
}
