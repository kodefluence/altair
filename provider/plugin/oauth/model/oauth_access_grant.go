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

// OauthAccessGrant an interface to access oauth_access_grants table
type OauthAccessGrant struct{}

// NewOauthAccessGrant create new OauthAccessGrant interface
func NewOauthAccessGrant() *OauthAccessGrant {
	return &OauthAccessGrant{}
}

// One selecting one record from database using id
func (*OauthAccessGrant) One(ktx kontext.Context, ID int, tx db.TX) (entity.OauthAccessGrant, exception.Exception) {
	var oauthAccessGrant entity.OauthAccessGrant

	ctxWithTimeout, cf := context.WithTimeout(ktx.Ctx(), time.Second*10)
	defer cf()

	row := tx.QueryRowContext(kontext.Fabricate(kontext.WithDefaultContext(ctxWithTimeout)), "oauth-access-grant-one", query.SelectOneOauthAccessGrant, ID)
	err := row.Scan(
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

	return oauthAccessGrant, err
}

// OneByCode selecting one record from database using code
func (*OauthAccessGrant) OneByCode(ktx kontext.Context, code string, tx db.TX) (entity.OauthAccessGrant, exception.Exception) {
	var oauthAccessGrant entity.OauthAccessGrant

	ctxWithTimeout, cf := context.WithTimeout(ktx.Ctx(), time.Second*10)
	defer cf()

	row := tx.QueryRowContext(kontext.Fabricate(kontext.WithDefaultContext(ctxWithTimeout)), "oauth-access-grant-one-by-code", query.SelectOneOauthAccessGrantByCode, code)
	err := row.Scan(
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

	return oauthAccessGrant, err
}

// Create new record
func (*OauthAccessGrant) Create(ktx kontext.Context, data entity.OauthAccessGrantInsertable, tx db.TX) (int, exception.Exception) {
	result, err := tx.ExecContext(ktx, "oauth-access-grant-create", query.InsertOauthAccessGrant, data.OauthApplicationID, data.ResourceOwnerID, data.Scopes, data.Code, data.RedirectURI, data.ExpiresIn)
	if err != nil {
		return 0, err
	}

	ID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(ID), nil
}

// Revoke fill revoked_at of oauth_access_grants row
func (*OauthAccessGrant) Revoke(ktx kontext.Context, code string, tx db.TX) exception.Exception {
	result, err := tx.ExecContext(ktx, "oauth-access-grant-revoke", query.RevokeAuthorizationCode, code)
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
