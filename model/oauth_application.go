package model

import (
	"context"
	"database/sql"
	"time"

	"github.com/codefluence-x/altair/core"
	"github.com/codefluence-x/altair/entity"
	"github.com/codefluence-x/altair/query"
)

type oauthApplication struct {
	db *sql.DB
}

func OauthApplication(db *sql.DB) core.OauthApplicationModel {
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
				&oauthApplication.ID, &oauthApplication.OwnerID, &oauthApplication.Description,
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
