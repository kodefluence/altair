package forwarder

import (
	"github.com/kodefluence/altair/core"
	"github.com/kodefluence/altair/forwarder/route"
)

type routeDispatcher struct{}

func Route() core.RouteDispatcher {
	return routeDispatcher{}
}

func (r routeDispatcher) Compiler() core.RouteCompiler {
	return route.Compiler()
}

func (r routeDispatcher) Generator() core.RouteGenerator {
	return route.Generator()
}
