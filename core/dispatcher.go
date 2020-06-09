package core

type RouteDispatcher interface {
	Compiler() RouteCompiler
	Generator() RouteGenerator
}
