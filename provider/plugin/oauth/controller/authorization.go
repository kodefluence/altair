package controller

import (
	"github.com/codefluence-x/altair/core"
	auth "github.com/codefluence-x/altair/provider/plugin/oauth/controller/authorization"
	"github.com/codefluence-x/altair/provider/plugin/oauth/interfaces"
)

type authorization struct{}

func Authorization() interfaces.AuthorizationDispatcher {
	return authorization{}
}

func (a authorization) Grant(authService interfaces.Authorization) core.Controller {
	return auth.Grant(authService)
}

func (a authorization) Revoke(authService interfaces.Authorization) core.Controller {
	return auth.Revoke(authService)
}

func (a authorization) Token(authService interfaces.Authorization) core.Controller {
	return auth.Token(authService)
}
