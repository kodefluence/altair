package core

type OauthDispatcher interface {
	Application() OauthApplicationDispatcher
	Authorization() AuthorizationDispatcher
}

type OauthApplicationDispatcher interface {
	List(applicationManager ApplicationManager) Controller
	One(applicationManager ApplicationManager) Controller
	Create(applicationManager ApplicationManager) Controller
}

type AuthorizationDispatcher interface {
	Grant(authorization Authorization) Controller
}

type RouteDispatcher interface {
	Compiler() RouteCompiler
	Generator() RouteGenerator
}

type DownstreamPluginDispatcher interface {
	Oauth(oauthAccessTokenModel OauthAccessTokenModel) DownStreamPlugin
}
