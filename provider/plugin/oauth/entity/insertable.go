package entity

import "time"

// OauthAccessTokenInsertable use for creating new access token data
type OauthAccessTokenInsertable struct {
	OauthApplicationID int
	ResourceOwnerID    int
	Token              string
	Scopes             interface{}
	ExpiresIn          time.Time
}

// OauthAccessGrantInsertable use for creating new access grant data
type OauthAccessGrantInsertable struct {
	OauthApplicationID int
	ResourceOwnerID    int
	Scopes             interface{}
	Code               string
	RedirectURI        interface{}
	ExpiresIn          time.Time
}

// OauthRefreshTokenInsertable use for creating new refresh token data
type OauthRefreshTokenInsertable struct {
	ExpiresIn          time.Time
	Token              string
	OauthAccessTokenID int
}
