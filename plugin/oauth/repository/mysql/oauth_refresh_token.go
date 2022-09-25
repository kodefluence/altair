package mysql

import (
	"context"
	"errors"
	"time"

	"github.com/kodefluence/altair/plugin/oauth/entity"
	"github.com/kodefluence/monorepo/db"
	"github.com/kodefluence/monorepo/exception"
	"github.com/kodefluence/monorepo/kontext"
)

// OauthRefreshToken is a connector to oauth_refresh_tokens table
type OauthRefreshToken struct{}

// NewOauthRefreshToken create new OauthRefreshToken struct
func NewOauthRefreshToken() *OauthRefreshToken {
	return &OauthRefreshToken{}
}

// OneByToken selecting one oauth refresh token data based on given token data
func (*OauthRefreshToken) OneByToken(ktx kontext.Context, token string, tx db.TX) (entity.OauthRefreshToken, exception.Exception) {
	var oauthRefreshToken entity.OauthRefreshToken

	ctxWithTimeout, cf := context.WithTimeout(ktx.Ctx(), time.Second*10)
	defer cf()

	row := tx.QueryRowContext(
		kontext.Fabricate(kontext.WithDefaultContext(ctxWithTimeout)),
		"oauth-refresh-token-one-by-token",
		"select id, oauth_access_token_id, token, expires_in, created_at, revoked_at from oauth_refresh_tokens where token = ? and revoked_at is null limit 1",
		token,
	)
	err := row.Scan(
		&oauthRefreshToken.ID,
		&oauthRefreshToken.OauthAccessTokenID,
		&oauthRefreshToken.Token,
		&oauthRefreshToken.ExpiresIn,
		&oauthRefreshToken.CreatedAt,
		&oauthRefreshToken.RevokedAT,
	)

	return oauthRefreshToken, err
}

// One selecting one oauth refresh token data
func (*OauthRefreshToken) One(ktx kontext.Context, ID int, tx db.TX) (entity.OauthRefreshToken, exception.Exception) {
	var oauthRefreshToken entity.OauthRefreshToken

	ctxWithTimeout, cf := context.WithTimeout(ktx.Ctx(), time.Second*10)
	defer cf()

	row := tx.QueryRowContext(
		kontext.Fabricate(kontext.WithDefaultContext(ctxWithTimeout)),
		"oauth-refresh-token-one",
		"select id, oauth_access_token_id, token, expires_in, created_at, revoked_at from oauth_refresh_tokens where id = ? limit 1",
		ID,
	)
	err := row.Scan(
		&oauthRefreshToken.ID,
		&oauthRefreshToken.OauthAccessTokenID,
		&oauthRefreshToken.Token,
		&oauthRefreshToken.ExpiresIn,
		&oauthRefreshToken.CreatedAt,
		&oauthRefreshToken.RevokedAT,
	)

	return oauthRefreshToken, err
}

// Create new oauth refresh token based on oauth refresh token insertable
func (*OauthRefreshToken) Create(ktx kontext.Context, data entity.OauthRefreshTokenInsertable, tx db.TX) (int, exception.Exception) {
	result, err := tx.ExecContext(
		ktx,
		"oauth-refresh-token-create",
		"insert into oauth_refresh_tokens (oauth_access_token_id, token, expires_in, created_at, revoked_at) values(?, ?, ?, now(), null)",
		data.OauthAccessTokenID, data.Token, data.ExpiresIn)
	if err != nil {
		return 0, err
	}

	ID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(ID), nil
}

// Revoke given token
func (*OauthRefreshToken) Revoke(ktx kontext.Context, token string, tx db.TX) exception.Exception {
	result, err := tx.ExecContext(ktx, "oauth-refresh-token-revoke", "update oauth_refresh_tokens set revoked_at = now() where token = ?", token)
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
