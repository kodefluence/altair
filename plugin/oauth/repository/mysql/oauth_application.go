package mysql

import (
	"context"
	"time"

	"github.com/kodefluence/altair/plugin/oauth/entity"
	"github.com/kodefluence/monorepo/db"
	"github.com/kodefluence/monorepo/exception"
	"github.com/kodefluence/monorepo/kontext"
)

// OauthApplication handle all database operation of `oauth_applications` table
type OauthApplication struct {
}

// NewOauthApplication create new OauthApplication struct with db.DB
func NewOauthApplication() *OauthApplication {
	return &OauthApplication{}
}

// Paginate oauth_applications data
func (*OauthApplication) Paginate(ktx kontext.Context, offset, limit int, tx db.TX) ([]entity.OauthApplication, exception.Exception) {
	var OauthApplications []entity.OauthApplication

	ctxWithTimeout, cf := context.WithTimeout(ktx.Ctx(), time.Second*10)
	defer cf()

	rows, err := tx.QueryContext(
		kontext.Fabricate(kontext.WithDefaultContext(ctxWithTimeout)),
		"oauth-application-paginate",
		"select id, owner_id, owner_type, description, scopes, client_uid, client_secret, revoked_at, created_at, updated_at from oauth_applications limit ?, ?",
		offset,
		limit,
	)
	if err != nil {
		return OauthApplications, err
	}
	defer rows.Close()

	for rows.Next() {
		var OauthApplication entity.OauthApplication

		err := rows.Scan(
			&OauthApplication.ID, &OauthApplication.OwnerID, &OauthApplication.OwnerType, &OauthApplication.Description,
			&OauthApplication.Scopes, &OauthApplication.ClientUID, &OauthApplication.ClientSecret,
			&OauthApplication.RevokedAt, &OauthApplication.CreatedAt, &OauthApplication.UpdatedAt,
		)
		if err != nil {
			return OauthApplications, err
		}

		OauthApplications = append(OauthApplications, OauthApplication)
	}

	return OauthApplications, rows.Err()
}

// Count all oauth_applications data
func (*OauthApplication) Count(ktx kontext.Context, tx db.TX) (int, exception.Exception) {
	var total int

	ctxWithTimeout, cf := context.WithTimeout(ktx.Ctx(), time.Second*10)
	defer cf()

	row := tx.QueryRowContext(
		kontext.Fabricate(kontext.WithDefaultContext(ctxWithTimeout)),
		"oauth-application-count",
		"select count(*) as total from oauth_applications where revoked_at is null",
	)
	if err := row.Scan(&total); err != nil {
		return total, err
	}

	return total, nil
}

// One get one oauth_applications by id
func (*OauthApplication) One(ktx kontext.Context, ID int, tx db.TX) (entity.OauthApplication, exception.Exception) {
	var data entity.OauthApplication

	ctxWithTimeout, cf := context.WithTimeout(ktx.Ctx(), time.Second*10)
	defer cf()

	row := tx.QueryRowContext(
		kontext.Fabricate(kontext.WithDefaultContext(ctxWithTimeout)),
		"oauth-application-one",
		"select id, owner_id, owner_type, description, scopes, client_uid, client_secret, revoked_at, created_at, updated_at from oauth_applications where id = ?",
		ID,
	)
	if err := row.Scan(
		&data.ID, &data.OwnerID, &data.OwnerType, &data.Description,
		&data.Scopes, &data.ClientUID, &data.ClientSecret,
		&data.RevokedAt, &data.CreatedAt, &data.UpdatedAt,
	); err != nil {
		return data, err
	}

	return data, nil
}

// OneByUIDandSecret get one oauth_applications by client uid and client secret
func (*OauthApplication) OneByUIDandSecret(ktx kontext.Context, clientUID, clientSecret string, tx db.TX) (entity.OauthApplication, exception.Exception) {
	var data entity.OauthApplication

	ctxWithTimeout, cf := context.WithTimeout(ktx.Ctx(), time.Second*10)
	defer cf()

	row := tx.QueryRowContext(
		kontext.Fabricate(kontext.WithDefaultContext(ctxWithTimeout)),
		"oauth-application-one-by-id-and-secret",
		"select id, owner_id, owner_type, description, scopes, client_uid, client_secret, revoked_at, created_at, updated_at from oauth_applications where client_uid = ? and client_secret = ? limit 1",
		clientUID,
		clientSecret,
	)
	if err := row.Scan(
		&data.ID, &data.OwnerID, &data.OwnerType, &data.Description,
		&data.Scopes, &data.ClientUID, &data.ClientSecret,
		&data.RevokedAt, &data.CreatedAt, &data.UpdatedAt,
	); err != nil {
		return data, err
	}

	return data, nil
}

// Create new oauth_applications data
func (*OauthApplication) Create(ktx kontext.Context, data entity.OauthApplicationInsertable, tx db.TX) (int, exception.Exception) {
	res, err := tx.ExecContext(
		ktx,
		"oauth-application-create",
		"insert into oauth_applications (owner_id, owner_type, description, scopes, client_uid, client_secret, revoked_at, created_at, updated_at) values(?, ?, ?, ?, ?, ?, null, now(), now())",
		data.OwnerID, data.OwnerType, data.Description, data.Scopes, data.ClientUID, data.ClientSecret)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// Update oauth_applications data
func (*OauthApplication) Update(ktx kontext.Context, ID int, data entity.OauthApplicationUpdateable, tx db.TX) exception.Exception {
	_, err := tx.ExecContext(
		kontext.Fabricate(),
		"oauth-application-update",
		"update oauth_applications set description = ?, scopes = ?, updated_at = now() where id = ?",
		data.Description, data.Scopes, ID)
	return err
}
