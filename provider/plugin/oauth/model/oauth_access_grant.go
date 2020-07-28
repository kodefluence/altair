package model

import (
	"context"
	"database/sql"
	"time"

	"github.com/codefluence-x/altair/provider/plugin/oauth/entity"
	"github.com/codefluence-x/altair/provider/plugin/oauth/query"
)

// OauthAccessGrant an interface to access oauth_access_grants table
type OauthAccessGrant struct {
	db *sql.DB
}

// NewOauthAccessGrant create new OauthAccessGrant interface
func NewOauthAccessGrant(db *sql.DB) *OauthAccessGrant {
	return &OauthAccessGrant{
		db: db,
	}
}

// Name return model name
func (oag *OauthAccessGrant) Name() string {
	return "oauth-access-grant-model"
}

// One selecting one record from database using id
func (oag *OauthAccessGrant) One(ctx context.Context, ID int) (entity.OauthAccessGrant, error) {
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

// OneByCode selecting one record from database using code
func (oag *OauthAccessGrant) OneByCode(ctx context.Context, code string) (entity.OauthAccessGrant, error) {
	var oauthAccessGrant entity.OauthAccessGrant

	err := monitor(ctx, oag.Name(), query.SelectOneOauthAccessGrantByCode, func() error {
		ctxWithTimeout, cf := context.WithTimeout(ctx, time.Second*10)
		defer cf()

		row := oag.db.QueryRowContext(ctxWithTimeout, query.SelectOneOauthAccessGrantByCode, code)
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

// Create new record
func (oag *OauthAccessGrant) Create(ctx context.Context, data entity.OauthAccessGrantInsertable, txs ...*sql.Tx) (int, error) {
	var lastInsertedID int
	var dbExecutable DBExecutable

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
		lastInsertedID = int(id)

		return err
	})

	return lastInsertedID, err
}

// Revoke fill revoked_at of oauth_access_grants row
func (oag *OauthAccessGrant) Revoke(ctx context.Context, code string, txs ...*sql.Tx) error {
	return monitor(ctx, oag.Name(), query.RevokeAuthorizationCode, func() error {
		result, err := oag.db.Exec(query.RevokeAuthorizationCode, code)
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
