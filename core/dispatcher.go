package core

type OauthDispatcher interface {
	Application() OauthApplicationDispatcher
}

type OauthApplicationDispatcher interface {
	List(applicationManager ApplicationManager) Controller
	One(applicationManager ApplicationManager) Controller
	Create(applicationManager ApplicationManager) Controller
}
