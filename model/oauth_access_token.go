package model

import (
	"context"
	"database/sql"
	"time"

	"github.com/codefluence-x/altair/core"
	"github.com/codefluence-x/altair/entity"
	"github.com/codefluence-x/altair/query"
)

type oauthAccessToken struct {
	db *sql.DB
}

func OauthAccessToken(db *sql.DB) core.OauthAccessTokenModel {
	return &oauthAccessToken{
		db: db,
	}
}

func (oat *oauthAccessToken) Name() string {
	return "oauth-access-token-model"
}

func (oat *oauthAccessToken) OneByToken(ctx context.Context, token string) (entity.OauthAccessToken, error) {
	var oauthAccessToken entity.OauthAccessToken

	err := monitor(ctx, oat.Name(), query.SelectOneOauthAccessTokenByToken, func() error {
		ctxWithTimeout, cf := context.WithTimeout(ctx, time.Second*10)
		defer cf()

		row := oat.db.QueryRowContext(ctxWithTimeout, query.SelectOneOauthAccessTokenByToken, token)
		return row.Scan(
			&oauthAccessToken.ID,
			&oauthAccessToken.OauthApplicationID,
			&oauthAccessToken.ResourceOwnerID,
			&oauthAccessToken.Token,
			&oauthAccessToken.Scopes,
			&oauthAccessToken.ExpiresIn,
			&oauthAccessToken.CreatedAt,
			&oauthAccessToken.RevokedAT,
		)
	})

	return oauthAccessToken, err
}

func (oat *oauthAccessToken) One(ctx context.Context, ID int) (entity.OauthAccessToken, error) {
	var oauthAccessToken entity.OauthAccessToken

	err := monitor(ctx, oat.Name(), query.SelectOneOauthAccessToken, func() error {
		ctxWithTimeout, cf := context.WithTimeout(ctx, time.Second*10)
		defer cf()

		row := oat.db.QueryRowContext(ctxWithTimeout, query.SelectOneOauthAccessToken, ID)
		return row.Scan(
			&oauthAccessToken.ID,
			&oauthAccessToken.OauthApplicationID,
			&oauthAccessToken.ResourceOwnerID,
			&oauthAccessToken.Token,
			&oauthAccessToken.Scopes,
			&oauthAccessToken.ExpiresIn,
			&oauthAccessToken.CreatedAt,
			&oauthAccessToken.RevokedAT,
		)
	})

	return oauthAccessToken, err
}

func (oat *oauthAccessToken) Create(ctx context.Context, data entity.OauthAccessTokenInsertable, txs ...*sql.Tx) (int, error) {
	var lastInsertedId int
	var dbExecutable core.DBExecutable

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
