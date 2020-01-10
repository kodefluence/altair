package oauth

import (
	app "github.com/codefluence-x/altair/controller/oauth/application"
	"github.com/codefluence-x/altair/core"
)

type application struct{}

func Application() core.OauthApplicationDispatcher {
	return application{}
}

func (a application) List(applicationManager core.ApplicationManager) core.Controller {
	return app.List(applicationManager)
}

func (a application) One(applicationManager core.ApplicationManager) core.Controller {
	return app.One(applicationManager)
}

func (a application) Create(applicationManager core.ApplicationManager) core.Controller {
	return app.Create(applicationManager)
}
