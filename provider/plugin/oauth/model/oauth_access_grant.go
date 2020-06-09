package model

import (
	"context"
	"database/sql"
	"time"

	"github.com/codefluence-x/altair/entity"
	"github.com/codefluence-x/altair/provider/plugin/oauth"
	"github.com/codefluence-x/altair/query"
)

type oauthAccessGrant struct {
	db *sql.DB
}

func OauthAccessGrant(db *sql.DB) oauth.OauthAccessGrantModel {
	return &oauthAccessGrant{
		db: db,
	}
}

func (oag *oauthAccessGrant) Name() string {
	return "oauth-access-grant-model"
}

func (oag *oauthAccessGrant) One(ctx context.Context, ID int) (entity.OauthAccessGrant, error) {
	var oauthAccessGrant entity.OauthAccessGrant

	err := monitor(ctx, oag.Name(), query.SelectOneOauthAccessGrant, func() error {
		ctxWithTimeout, cf := context.WithTimeout(ctx, time.Second*10)
		defer cf()

		row := oag.db.QueryRowContext(ctxWithTimeout, query.SelectOneOauthAccessGrant, ID)
		return row.Scan(
			&oauthAccessGrant.ID,
			&oauthAccessGrant.OauthApplicationID,
			&oauthAccessGrant.ResourceOwnerID,
			&oauthAccessGrant.Scopes,
			&oauthAccessGrant.Code,
			&oauthAccessGrant.RedirectURI,
			&oauthAccessGrant.ExpiresIn,
			&oauthAccessGrant.CreatedAt,
			&oauthAccessGrant.RevokedAT,
		)
	})

	return oauthAccessGrant, err
}

func (oag *oauthAccessGrant) Create(ctx context.Context, data entity.OauthAccessGrantInsertable, txs ...*sql.Tx) (int, error) {
	var lastInsertedId int
	var dbExecutable oauth.DBExecutable

	dbExecutable = oag.db
	if len(txs) > 0 {
		dbExecutable = txs[0]
	}

	err := monitor(ctx, oag.Name(), query.InsertOauthAccessGrant, func() error {
		result, err := dbExecutable.Exec(query.InsertOauthAccessGrant, data.OauthApplicationID, data.ResourceOwnerID, data.Scopes, data.Code, data.RedirectURI, data.ExpiresIn)
		if err != nil {
			return err
		}

		id, err := result.LastInsertId()
		lastInsertedId = int(id)

		return err
	})

	return lastInsertedId, err
}
