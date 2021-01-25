package model

import (
	"context"
	"database/sql"
	"time"

	"github.com/codefluence-x/altair/provider/plugin/oauth/entity"
	"github.com/codefluence-x/altair/provider/plugin/oauth/query"
)

// OauthRefreshToken is a connector to oauth_refresh_tokens table
type OauthRefreshToken struct {
	db *sql.DB
}

// NewOauthRefreshToken create new OauthRefreshToken struct
func NewOauthRefreshToken(db *sql.DB) *OauthRefreshToken {
	return &OauthRefreshToken{
		db: db,
	}
}

// Name of the model
func (ort *OauthRefreshToken) Name() string {
	return "oauth-refresh-token-model"
}

// OneByToken selecting one oauth refresh token data based on given token data
func (ort *OauthRefreshToken) OneByToken(ctx context.Context, token string) (entity.OauthRefreshToken, error) {
	var OauthRefreshToken entity.OauthRefreshToken

	err := monitor(ctx, ort.Name(), query.SelectOneOauthRefreshTokenByToken, func() error {
		ctxWithTimeout, cf := context.WithTimeout(ctx, time.Second*10)
		defer cf()

		row := ort.db.QueryRowContext(ctxWithTimeout, query.SelectOneOauthRefreshTokenByToken, token)
		return row.Scan(
			&OauthRefreshToken.ID,
			&OauthRefreshToken.OauthAccessTokenID,
			&OauthRefreshToken.Token,
			&OauthRefreshToken.ExpiresIn,
			&OauthRefreshToken.CreatedAt,
			&OauthRefreshToken.RevokedAT,
		)
	})

	return OauthRefreshToken, err
}

// One selecting one oauth refresh token data
func (ort *OauthRefreshToken) One(ctx context.Context, ID int) (entity.OauthRefreshToken, error) {
	var OauthRefreshToken entity.OauthRefreshToken

	err := monitor(ctx, ort.Name(), query.SelectOneOauthRefreshToken, func() error {
		ctxWithTimeout, cf := context.WithTimeout(ctx, time.Second*10)
		defer cf()

		row := ort.db.QueryRowContext(ctxWithTimeout, query.SelectOneOauthRefreshToken, ID)
		return row.Scan(
			&OauthRefreshToken.ID,
			&OauthRefreshToken.OauthAccessTokenID,
			&OauthRefreshToken.Token,
			&OauthRefreshToken.ExpiresIn,
			&OauthRefreshToken.CreatedAt,
			&OauthRefreshToken.RevokedAT,
		)
	})

	return OauthRefreshToken, err
}

// Create new oauth refresh token based on oauth refresh token insertable
func (ort *OauthRefreshToken) Create(ctx context.Context, data entity.OauthRefreshTokenInsertable, txs ...*sql.Tx) (int, error) {
	var lastInsertedID int
	var dbExecutable DBExecutable

	dbExecutable = ort.db
	if len(txs) > 0 {
		dbExecutable = txs[0]
	}

	err := monitor(ctx, ort.Name(), query.InsertOauthRefreshToken, func() error {
		result, err := dbExecutable.Exec(query.InsertOauthRefreshToken, data.OauthAccessTokenID, data.Token, data.ExpiresIn)
		if err != nil {
			return err
		}

		id, err := result.LastInsertId()
		lastInsertedID = int(id)

		return err
	})

	return lastInsertedID, err
}

// Revoke given token
func (ort *OauthRefreshToken) Revoke(ctx context.Context, token string) error {
	return monitor(ctx, ort.Name(), query.RevokeRefreshToken, func() error {
		result, err := ort.db.Exec(query.RevokeRefreshToken, token)
		if err != nil {
			return err
		}

		affectedRows, err := result.RowsAffected()
		if err != nil {
			return err
		}

		if affectedRows == 0 {
			return sql.ErrNoRows
		}

		return nil
	})
}
