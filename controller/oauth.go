package controller

import "github.com/codefluence-x/altair/core"

type oauthDispatcher struct{}

func Oauth() core.OauthDispatcher {
	return oauthDispatcher{}
}
