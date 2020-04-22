package plugin

import (
	"github.com/codefluence-x/altair/core"
	"github.com/codefluence-x/altair/plugin/downstream"
)

type downstreamDispatcher struct{}

func DownStream() core.DownstreamPluginDispatcher {
	return &downstreamDispatcher{}
}

func (d *downstreamDispatcher) Oauth(oauthAccessTokenModel core.OauthAccessTokenModel) core.DownStreamPlugin {
	return downstream.Oauth(oauthAccessTokenModel)
}
