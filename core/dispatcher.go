package core

type OauthDispatcher interface {
	Application() OauthApplicationDispatcher
}

type OauthApplicationDispatcher interface{}
