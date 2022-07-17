package controller

import (
	auth "github.com/kodefluence/altair/provider/plugin/oauth/controller/authorization"
	"github.com/kodefluence/altair/provider/plugin/oauth/interfaces"
)

// Authorization dispatch authorization related controller
type Authorization struct{}

// NewAuthorization return struct of Authorization
func NewAuthorization() *Authorization {
	return &Authorization{}
}

// Grant return handler of POST /oauth/authorizations
func (a *Authorization) Grant(authService interfaces.Authorization) *auth.GrantController {
	return auth.NewGrant(authService)
}

// Revoke return handler of POST /oauth/authorizations/revoke
func (a *Authorization) Revoke(authService interfaces.Authorization) *auth.RevokeController {
	return auth.NewRevoke(authService)
}

// Token return handler of POST /oauth/authorizations/token
func (a *Authorization) Token(authService interfaces.Authorization) *auth.TokenController {
	return auth.NewToken(authService)
}
