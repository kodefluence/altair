package entity

import "time"

//
//
// JSON Struct
//
//

// OauthAccessGrantJSON is a json response from OauthAccessGrant
type OauthAccessGrantJSON struct {
	ID                 *int       `json:"id"`
	OauthApplicationID *int       `json:"oauth_application_id"`
	ResourceOwnerID    *int       `json:"resource_owner_id"`
	Code               *string    `json:"code"`
	RedirectURI        *string    `json:"redirect_uri"`
	Scopes             *string    `json:"scopes"`
	ExpiresIn          *int       `json:"expires_in"`
	CreatedAt          *time.Time `json:"created_at"`
	RevokedAT          *time.Time `json:"revoked_at"`
}

// OauthAccessTokenJSON is a json response from OauthAccessToken
type OauthAccessTokenJSON struct {
	ID                 *int                   `json:"id"`
	OauthApplicationID *int                   `json:"oauth_application_id"`
	ResourceOwnerID    *int                   `json:"resource_owner_id"`
	Token              *string                `json:"token"`
	Scopes             *string                `json:"scopes"`
	ExpiresIn          *int                   `json:"expires_in"`
	RedirectURI        *string                `json:"redirect_uri,omitempty"`
	CreatedAt          *time.Time             `json:"created_at"`
	RevokedAT          *time.Time             `json:"revoked_at"`
	RefreshToken       *OauthRefreshTokenJSON `json:"refresh_token,omitempty"`
}

type OauthRefreshTokenJSON struct {
	Token     *string    `json:"token"`
	ExpiresIn *int       `json:"expires_in"`
	CreatedAt *time.Time `json:"created_at"`
	RevokedAT *time.Time `json:"revoked_at"`
}

type OauthApplicationJSON struct {
	ID           *int       `json:"id"`
	OwnerID      *int       `json:"owner_id"`
	OwnerType    *string    `json:"owner_type"`
	Description  *string    `json:"description"`
	Scopes       *string    `json:"scopes"`
	ClientUID    *string    `json:"client_uid"`
	ClientSecret *string    `json:"client_secret"`
	RevokedAt    *time.Time `json:"revoked_at"`
	CreatedAt    *time.Time `json:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at"`
}

type OauthApplicationUpdateJSON struct {
	Description *string `json:"description"`
	Scopes      *string `json:"scopes"`
}

type AuthorizationRequestJSON struct {
	ResponseType *string `json:"response_type"`

	ResourceOwnerID *int `json:"resource_owner_id"`

	ClientUID    *string `json:"client_uid"`
	ClientSecret *string `json:"client_secret"`

	RedirectURI *string `json:"redirect_uri"`
	Scopes      *string `json:"scopes"`
}

type RevokeAccessTokenRequestJSON struct {
	Token *string `json:"token"`
}

type AccessTokenRequestJSON struct {
	GrantType *string `json:"grant_type"`

	ClientUID    *string `json:"client_uid"`
	ClientSecret *string `json:"client_secret"`

	RefreshToken *string `json:"refresh_token"`

	Code        *string `json:"code"`
	RedirectURI *string `json:"redirect_uri"`

	Scope *string `json:"scope"`
}
