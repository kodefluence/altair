package model

import (
	"context"
	"errors"
	"time"

	"github.com/kodefluence/altair/provider/plugin/oauth/entity"
	"github.com/kodefluence/altair/provider/plugin/oauth/query"
	"github.com/kodefluence/monorepo/db"
	"github.com/kodefluence/monorepo/exception"
	"github.com/kodefluence/monorepo/kontext"
)

// OauthAccessToken handle all database operation to oauth_access_tokens
type OauthAccessToken struct{}

// NewOauthAccessToken create new OauthAccessTokens struct
func NewOauthAccessToken() *OauthAccessToken {
	return &OauthAccessToken{}
}

// OneByToken get oauth access token data by token string
func (*OauthAccessToken) OneByToken(ktx kontext.Context, token string, tx db.TX) (entity.OauthAccessToken, exception.Exception) {
	var oauthAccessToken entity.OauthAccessToken

	ctxWithTimeout, cf := context.WithTimeout(ktx.Ctx(), time.Second*10)
	defer cf()

	row := tx.QueryRowContext(kontext.Fabricate(kontext.WithDefaultContext(ctxWithTimeout)), "oauth-access-token-one-by-token", query.SelectOneOauthAccessTokenByToken, token)
	err := row.Scan(
		&oauthAccessToken.ID,
		&oauthAccessToken.OauthApplicationID,
		&oauthAccessToken.ResourceOwnerID,
		&oauthAccessToken.Token,
		&oauthAccessToken.Scopes,
		&oauthAccessToken.ExpiresIn,
		&oauthAccessToken.CreatedAt,
		&oauthAccessToken.RevokedAT,
	)

	return oauthAccessToken, err
}

// One get oauth access token data by id
func (*OauthAccessToken) One(ktx kontext.Context, ID int, tx db.TX) (entity.OauthAccessToken, exception.Exception) {
	var oauthAccessToken entity.OauthAccessToken

	ctxWithTimeout, cf := context.WithTimeout(ktx.Ctx(), time.Second*10)
	defer cf()

	row := tx.QueryRowContext(kontext.Fabricate(kontext.WithDefaultContext(ctxWithTimeout)), "oauth-access-token-one", query.SelectOneOauthAccessToken, ID)
	err := row.Scan(
		&oauthAccessToken.ID,
		&oauthAccessToken.OauthApplicationID,
		&oauthAccessToken.ResourceOwnerID,
		&oauthAccessToken.Token,
		&oauthAccessToken.Scopes,
		&oauthAccessToken.ExpiresIn,
		&oauthAccessToken.CreatedAt,
		&oauthAccessToken.RevokedAT,
	)

	return oauthAccessToken, err
}

// Create new oauth access token
func (*OauthAccessToken) Create(ktx kontext.Context, data entity.OauthAccessTokenInsertable, tx db.TX) (int, exception.Exception) {
	result, err := tx.ExecContext(ktx, "oauth-access-token-create", query.InsertOauthAccessToken, data.OauthApplicationID, data.ResourceOwnerID, data.Token, data.Scopes, data.ExpiresIn)
	if err != nil {
		return 0, err
	}

	ID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(ID), nil
}

// Revoke oauth access token
func (*OauthAccessToken) Revoke(ktx kontext.Context, token string, tx db.TX) exception.Exception {
	result, err := tx.ExecContext(ktx, "oauth-access-token-revoke", query.RevokeAccessToken, token)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return exception.Throw(errors.New("not found"), exception.WithType(exception.NotFound))
	}

	return nil
}
