package model

import (
	"context"
	"database/sql"
	"time"

	"github.com/codefluence-x/altair/core"
	"github.com/codefluence-x/altair/entity"
	"github.com/codefluence-x/altair/query"
)

type oauthAccessGrant struct {
	db *sql.DB
}

func OauthAccessGrant(db *sql.DB) core.OauthAccessGrantModel {
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
			&oauthAccessGrant.Code,
			&oauthAccessGrant.RedirectURI,
			&oauthAccessGrant.Scopes,
			&oauthAccessGrant.ExpiresIn,
			&oauthAccessGrant.CreatedAt,
			&oauthAccessGrant.RevokedAT,
		)
	})

	return oauthAccessGrant, err
}
