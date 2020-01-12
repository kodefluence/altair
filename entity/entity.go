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
	Description  string
	Scopes       string
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
	RedirectURI        string
	Scopes             string
	ExpiresIn          time.Time
	CreatedAt          time.Time
	RevokedAT          mysql.NullTime
}

type OauthAccessToken struct {
	ID                 int
	OauthApplicationID int
	ResourceOwnerID    int
	Token              string
	Scopes             string
	ExpiresIn          time.Time
	CreatedAt          time.Time
	RevokedAT          mysql.NullTime
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

type AuthorizationRequestJSON struct {
	ResponseType *string `json:"response_type"`

	ClientUID    *string `json:"client_uid"`
	ClientSecret *string `json:"client_secret"`

	RedirectURI *string `json:"redirect_uri"`
	Scopes      *string `json:"scopes"`
}

type AccessTokenRequestJSON struct {
	GrantType *string `json:"grant_type"`

	ClientUID    *string `json:"client_uid"`
	ClientSecret *string `json:"client_secret"`

	Code        *string `json:"code"`
	RedirectURI *string `json:"redirect_uri"`
}
