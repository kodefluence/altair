package controller

import (
	app "github.com/codefluence-x/altair/provider/plugin/oauth/controller/application"
	"github.com/codefluence-x/altair/provider/plugin/oauth/interfaces"
)

// Application dispatch oauth application related controller
type Application struct{}

// NewApplication return struct of Application
func NewApplication() *Application {
	return &Application{}
}

// List return handler of GET /oauth/applications
func (a Application) List(applicationManager interfaces.ApplicationManager) *app.ListController {
	return app.NewList(applicationManager)
}

// One return handler of GET /oauth/applications/:id
func (a Application) One(applicationManager interfaces.ApplicationManager) *app.OneController {
	return app.NewOne(applicationManager)
}

// Create return handler of POST /oauth/applications
func (a Application) Create(applicationManager interfaces.ApplicationManager) *app.CreateController {
	return app.NewCreate(applicationManager)
}

// Update return handler of PUT /oauth/applications
func (a Application) Update(applicationManager interfaces.ApplicationManager) *app.UpdateController {
	return app.NewUpdate(applicationManager)
}
