package model

import (
	"context"
	"database/sql"
	"time"

	"github.com/codefluence-x/altair/entity"
	"github.com/codefluence-x/altair/provider/plugin/oauth/query"
)

type oauthApplication struct {
	db *sql.DB
}

func OauthApplication(db *sql.DB) *oauthApplication {
	return &oauthApplication{
		db: db,
	}
}

func (oa *oauthApplication) Name() string {
	return "oauth-application-model"
}

func (oa *oauthApplication) Paginate(ctx context.Context, offset, limit int) ([]entity.OauthApplication, error) {
	var oauthApplications []entity.OauthApplication

	err := monitor(ctx, oa.Name(), query.PaginateOauthApplication, func() error {

		ctxWithTimeout, cf := context.WithTimeout(ctx, time.Second*10)
		defer cf()

		rows, err := oa.db.QueryContext(ctxWithTimeout, query.PaginateOauthApplication, offset, limit)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var oauthApplication entity.OauthApplication

			err := rows.Scan(
				&oauthApplication.ID, &oauthApplication.OwnerID, &oauthApplication.OwnerType, &oauthApplication.Description,
				&oauthApplication.Scopes, &oauthApplication.ClientUID, &oauthApplication.ClientSecret,
				&oauthApplication.RevokedAt, &oauthApplication.CreatedAt, &oauthApplication.UpdatedAt,
			)
			if err != nil {
				return err
			}

			oauthApplications = append(oauthApplications, oauthApplication)
		}

		return rows.Err()
	})

	return oauthApplications, err
}

func (oa *oauthApplication) Count(ctx context.Context) (int, error) {
	var total int

	err := monitor(ctx, oa.Name(), query.CountOauthApplication, func() error {

		ctxWithTimeout, cf := context.WithTimeout(ctx, time.Second*10)
		defer cf()

		row := oa.db.QueryRowContext(ctxWithTimeout, query.CountOauthApplication)
		return row.Scan(&total)
	})

	return total, err
}

func (oa *oauthApplication) One(ctx context.Context, ID int) (entity.OauthApplication, error) {
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

func (oa *oauthApplication) OneByUIDandSecret(ctx context.Context, clientUID, clientSecret string) (entity.OauthApplication, error) {
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

func (oa *oauthApplication) Create(ctx context.Context, data entity.OauthApplicationInsertable, txs ...*sql.Tx) (int, error) {
	var lastInsertedId int
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

		lastInsertedId = int(id)
		return nil
	})

	return lastInsertedId, err
}
