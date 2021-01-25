package model

import (
	"context"
	"database/sql"
	"time"

	"github.com/codefluence-x/altair/provider/plugin/oauth/entity"
	"github.com/codefluence-x/altair/provider/plugin/oauth/query"
)

type OauthAccessToken struct {
	db *sql.DB
}

func NewOauthAccessToken(db *sql.DB) *OauthAccessToken {
	return &OauthAccessToken{
		db: db,
	}
}

func (oat *OauthAccessToken) Name() string {
	return "oauth-access-token-model"
}

func (oat *OauthAccessToken) OneByToken(ctx context.Context, token string) (entity.OauthAccessToken, error) {
	var OauthAccessToken entity.OauthAccessToken

	err := monitor(ctx, oat.Name(), query.SelectOneOauthAccessTokenByToken, func() error {
		ctxWithTimeout, cf := context.WithTimeout(ctx, time.Second*10)
		defer cf()

		row := oat.db.QueryRowContext(ctxWithTimeout, query.SelectOneOauthAccessTokenByToken, token)
		return row.Scan(
			&OauthAccessToken.ID,
			&OauthAccessToken.OauthApplicationID,
			&OauthAccessToken.ResourceOwnerID,
			&OauthAccessToken.Token,
			&OauthAccessToken.Scopes,
			&OauthAccessToken.ExpiresIn,
			&OauthAccessToken.CreatedAt,
			&OauthAccessToken.RevokedAT,
		)
	})

	return OauthAccessToken, err
}

func (oat *OauthAccessToken) One(ctx context.Context, ID int) (entity.OauthAccessToken, error) {
	var OauthAccessToken entity.OauthAccessToken

	err := monitor(ctx, oat.Name(), query.SelectOneOauthAccessToken, func() error {
		ctxWithTimeout, cf := context.WithTimeout(ctx, time.Second*10)
		defer cf()

		row := oat.db.QueryRowContext(ctxWithTimeout, query.SelectOneOauthAccessToken, ID)
		return row.Scan(
			&OauthAccessToken.ID,
			&OauthAccessToken.OauthApplicationID,
			&OauthAccessToken.ResourceOwnerID,
			&OauthAccessToken.Token,
			&OauthAccessToken.Scopes,
			&OauthAccessToken.ExpiresIn,
			&OauthAccessToken.CreatedAt,
			&OauthAccessToken.RevokedAT,
		)
	})

	return OauthAccessToken, err
}

func (oat *OauthAccessToken) Create(ctx context.Context, data entity.OauthAccessTokenInsertable, txs ...*sql.Tx) (int, error) {
	var lastInsertedId int
	var dbExecutable DBExecutable

	dbExecutable = oat.db
	if len(txs) > 0 {
		dbExecutable = txs[0]
	}

	err := monitor(ctx, oat.Name(), query.InsertOauthAccessToken, func() error {
		result, err := dbExecutable.Exec(query.InsertOauthAccessToken, data.OauthApplicationID, data.ResourceOwnerID, data.Token, data.Scopes, data.ExpiresIn)
		if err != nil {
			return err
		}

		id, err := result.LastInsertId()
		lastInsertedId = int(id)

		return err
	})

	return lastInsertedId, err
}

func (oat *OauthAccessToken) Revoke(ctx context.Context, token string) error {
	return monitor(ctx, oat.Name(), query.RevokeAccessToken, func() error {
		result, err := oat.db.Exec(query.RevokeAccessToken, token)
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
