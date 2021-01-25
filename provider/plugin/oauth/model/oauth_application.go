package model

import (
	"context"
	"database/sql"
	"time"

	"github.com/codefluence-x/altair/provider/plugin/oauth/entity"
	"github.com/codefluence-x/altair/provider/plugin/oauth/query"
)

type OauthApplication struct {
	db *sql.DB
}

func NewOauthApplication(db *sql.DB) *OauthApplication {
	return &OauthApplication{
		db: db,
	}
}

func (oa *OauthApplication) Name() string {
	return "oauth-application-model"
}

func (oa *OauthApplication) Paginate(ctx context.Context, offset, limit int) ([]entity.OauthApplication, error) {
	var OauthApplications []entity.OauthApplication

	err := monitor(ctx, oa.Name(), query.PaginateOauthApplication, func() error {

		ctxWithTimeout, cf := context.WithTimeout(ctx, time.Second*10)
		defer cf()

		rows, err := oa.db.QueryContext(ctxWithTimeout, query.PaginateOauthApplication, offset, limit)
		if err != nil {
			return err
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
				return err
			}

			OauthApplications = append(OauthApplications, OauthApplication)
		}

		return rows.Err()
	})

	return OauthApplications, err
}

func (oa *OauthApplication) Count(ctx context.Context) (int, error) {
	var total int

	err := monitor(ctx, oa.Name(), query.CountOauthApplication, func() error {

		ctxWithTimeout, cf := context.WithTimeout(ctx, time.Second*10)
		defer cf()

		row := oa.db.QueryRowContext(ctxWithTimeout, query.CountOauthApplication)
		return row.Scan(&total)
	})

	return total, err
}

func (oa *OauthApplication) One(ctx context.Context, ID int) (entity.OauthApplication, error) {
	var data entity.OauthApplication

	err := monitor(ctx, oa.Name(), query.SelectOneOauthApplication, func() error {

		ctxWithTimeout, cf := context.WithTimeout(ctx, time.Second*10)
		defer cf()

		row := oa.db.QueryRowContext(ctxWithTimeout, query.SelectOneOauthApplication, ID)
		return row.Scan(
			&data.ID, &data.OwnerID, &data.OwnerType, &data.Description,
			&data.Scopes, &data.ClientUID, &data.ClientSecret,
			&data.RevokedAt, &data.CreatedAt, &data.UpdatedAt,
		)
	})

	return data, err
}

func (oa *OauthApplication) OneByUIDandSecret(ctx context.Context, clientUID, clientSecret string) (entity.OauthApplication, error) {
	var data entity.OauthApplication

	err := monitor(ctx, oa.Name(), query.SelectOneByUIDandSecret, func() error {

		ctxWithTimeout, cf := context.WithTimeout(ctx, time.Second*10)
		defer cf()

		row := oa.db.QueryRowContext(ctxWithTimeout, query.SelectOneByUIDandSecret, clientUID, clientSecret)
		return row.Scan(
			&data.ID, &data.OwnerID, &data.OwnerType, &data.Description,
			&data.Scopes, &data.ClientUID, &data.ClientSecret,
			&data.RevokedAt, &data.CreatedAt, &data.UpdatedAt,
		)
	})

	return data, err
}

func (oa *OauthApplication) Create(ctx context.Context, data entity.OauthApplicationInsertable, txs ...*sql.Tx) (int, error) {
	var lastInsertedID int
	var dbExecutable DBExecutable

	dbExecutable = oa.db
	if len(txs) > 0 {
		dbExecutable = txs[0]
	}

	err := monitor(ctx, oa.Name(), query.InsertOauthApplication, func() error {

		res, err := dbExecutable.Exec(query.InsertOauthApplication, data.OwnerID, data.OwnerType, data.Description, data.Scopes, data.ClientUID, data.ClientSecret)
		if err != nil {
			return err
		}

		id, err := res.LastInsertId()
		if err != nil {
			return err
		}

		lastInsertedID = int(id)
		return nil
	})

	return lastInsertedID, err
}

func (oa *OauthApplication) Update(ctx context.Context, ID int, data entity.OauthApplicationUpdateable, txs ...*sql.Tx) error {
	var dbExecutable DBExecutable

	dbExecutable = oa.db
	if len(txs) > 0 {
		dbExecutable = txs[0]
	}

	err := monitor(ctx, oa.Name(), query.UpdateOauthApplication, func() error {

		res, err := dbExecutable.Exec(query.UpdateOauthApplication, data.Description, data.Scopes, ID)
		if err != nil {
			return err
		}

		affectedRows, err := res.RowsAffected()
		if err != nil {
			return err
		} else if affectedRows == 0 {
			return sql.ErrNoRows
		}

		return nil
	})

	return err
}
