package oauth

import (
	auth "github.com/codefluence-x/altair/controller/oauth/authorization"
	"github.com/codefluence-x/altair/core"
)

type authorization struct{}

func Authorization() core.AuthorizationDispatcher {
	return authorization{}
}

func (a authorization) Grant(authService core.Authorization) core.Controller {
	return auth.Grant(authService)
}

func (a authorization) Revoke(authService core.Authorization) core.Controller {
	return auth.Revoke(authService)
}
