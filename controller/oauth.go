package controller

import (
	"github.com/codefluence-x/altair/controller/oauth"
	"github.com/codefluence-x/altair/core"
)

type oauthDispatcher struct{}

func Oauth() core.OauthDispatcher {
	return oauthDispatcher{}
}

func (o oauthDispatcher) Application() core.OauthApplicationDispatcher {
	return oauth.Application()
}

func (o oauthDispatcher) Authorization() core.AuthorizationDispatcher {
	return oauth.Authorization()
}
