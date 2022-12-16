package entity

import (
	"database/sql"
	"time"
)

//
//
// Model Struct
//
//

// OauthApplication is a struct returned from interfaces.OauthApplicationModel
type OauthApplication struct {
	ID           int
	OwnerID      sql.NullInt64
	OwnerType    string
	Description  sql.NullString
	Scopes       sql.NullString
	ClientUID    string
	ClientSecret string
	RevokedAt    sql.NullTime
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// OauthAccessGrant is a struct returned from interfaces.OauthAccessGrantModel
type OauthAccessGrant struct {
	ID                 int
	OauthApplicationID int
	ResourceOwnerID    int
	Code               string
	RedirectURI        sql.NullString
	Scopes             sql.NullString
	ExpiresIn          time.Time
	CreatedAt          time.Time
	RevokedAT          sql.NullTime
}

// OauthAccessToken is a struct returned from interfaces.OauthAccessTokenModel
type OauthAccessToken struct {
	ID                 int
	OauthApplicationID int
	ResourceOwnerID    int
	Token              string
	Scopes             sql.NullString
	ExpiresIn          time.Time
	CreatedAt          time.Time
	RevokedAT          sql.NullTime
}

// OauthRefreshToken is a struct returned from interfaces.OauthRefreshTokenModel
type OauthRefreshToken struct {
	ID                 int
	OauthAccessTokenID int
	Token              string
	ExpiresIn          time.Time
	CreatedAt          time.Time
	RevokedAT          sql.NullTime
}
