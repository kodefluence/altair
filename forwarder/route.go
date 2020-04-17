package forwarder

import (
	"github.com/codefluence-x/altair/core"
	"github.com/codefluence-x/altair/forwarder/route"
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
