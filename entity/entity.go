package entity

import (
	"database/sql"
	"time"

	"github.com/go-sql-driver/mysql"
)

type OauthApplication struct {
	ID           int
	OwnerID      sql.NullInt64
	OwnerType    string
	Description  sql.NullString
	Scopes       sql.NullString
	ClientUID    string
	ClientSecret string
	RevokedAt    mysql.NullTime
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type OauthAccessGrant struct {
	ID                 int
	OauthApplicationID int
	ResourceOwnerID    int
	Code               string
	RedirectURI        sql.NullString
	Scopes             sql.NullString
	ExpiresIn          time.Time
	CreatedAt          time.Time
	RevokedAT          mysql.NullTime
}

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

type OauthAccessToken struct {
	ID                 int
	OauthApplicationID int
	ResourceOwnerID    int
	Token              string
	Scopes             sql.NullString
	ExpiresIn          time.Time
	CreatedAt          time.Time
	RevokedAT          mysql.NullTime
}

type OauthAccessTokenJSON struct {
	ID                 *int       `json:"id"`
	OauthApplicationID *int       `json:"oauth_application_id"`
	ResourceOwnerID    *int       `json:"resource_owner_id"`
	Token              *string    `json:"token"`
	Scopes             *string    `json:"scopes"`
	ExpiresIn          *int       `json:"expires_in"`
	RedirectURI        *string    `json:"redirect_uri"`
	CreatedAt          *time.Time `json:"created_at"`
	RevokedAT          *time.Time `json:"revoked_at"`
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

type OauthApplicationInsertable struct {
	OwnerID      interface{}
	OwnerType    string
	Description  interface{}
	Scopes       interface{}
	ClientUID    string
	ClientSecret string
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

	Code        *string `json:"code"`
	RedirectURI *string `json:"redirect_uri"`
}

type OauthAccessTokenInsertable struct {
	OauthApplicationID int
	ResourceOwnerID    int
	Token              string
	Scopes             interface{}
	ExpiresIn          time.Time
}

type OauthAccessGrantInsertable struct {
	OauthApplicationID int
	ResourceOwnerID    int
	Scopes             interface{}
	Code               string
	RedirectURI        interface{}
	ExpiresIn          time.Time
}

type RouteObject struct {
	Name   string                `yaml:"name"`
	Auth   string                `yaml:"auth"`
	Prefix string                `yaml:"prefix"`
	Host   string                `yaml:"host"`
	Path   map[string]RouterPath `yaml:"path"`
}

type RouterPath struct {
	Auth string `yaml:"auth"`
}
