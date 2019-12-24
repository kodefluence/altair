package core

type OauthDispatcher interface {
	Application() OauthApplicationDispatcher
}

type OauthApplicationDispatcher interface {
	List(applicationManager ApplicationManager) Controller
}
