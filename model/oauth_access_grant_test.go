package model_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/codefluence-x/altair/entity"
	"github.com/codefluence-x/altair/model"
	"github.com/stretchr/testify/assert"
)

var oauthAccessGrantModelRows = []string{
	"id",
	"oauth_application_id",
	"resource_owner_id",
	"scopes",
	"code",
	"redirect_uri",
	"expires_in",
	"created_at",
	"revoked_at",
}

func TestOauthAccessGrant(t *testing.T) {

	t.Run("Name", func(t *testing.T) {
		t.Run("Return a model's name", func(t *testing.T) {
			db, _, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}

			assert.Equal(t, "oauth-access-grant-model", model.OauthAccessGrant(db).Name())
		})
	})

	t.Run("One", func(t *testing.T) {
		t.Run("Given context and oauth access grant ID", func(t *testing.T) {
			t.Run("Return entity oauth access token", func(t *testing.T) {
				db, mockdb, err := sqlmock.New()
				if err != nil {
					t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
				}

				data := entity.OauthAccessGrant{
					ID: 1,
				}

				rows := sqlmock.NewRows(oauthAccessGrantModelRows).
					AddRow(
						data.ID,
						data.OauthApplicationID,
						data.ResourceOwnerID,
						data.Code,
						data.RedirectURI,
						data.Scopes,
						data.ExpiresIn,
						data.CreatedAt,
						data.RevokedAT,
					)

				mockdb.ExpectQuery(`select \* from oauth_access_grants where id = \? limit 1`).
					WithArgs(1).
					WillReturnRows(rows)

				oauthAccessGrant := model.OauthAccessGrant(db)
				dataFromDB, err := oauthAccessGrant.One(context.Background(), 1)

				assert.Nil(t, err)
				assert.Nil(t, mockdb.ExpectationsWereMet())
				assert.Equal(t, data, dataFromDB)
			})

			t.Run("Row not found in database", func(t *testing.T) {
				db, mockdb, err := sqlmock.New()
				if err != nil {
					t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
				}

				mockdb.ExpectQuery(`select \* from oauth_access_grants where id = \? limit 1`).
					WithArgs(1).
					WillReturnError(sql.ErrNoRows)

				oauthAccessGrant := model.OauthAccessGrant(db)
				dataFromDB, err := oauthAccessGrant.One(context.Background(), 1)

				assert.NotNil(t, err)
				assert.Equal(t, sql.ErrNoRows, err)
				assert.Nil(t, mockdb.ExpectationsWereMet())
				assert.Equal(t, entity.OauthAccessGrant{}, dataFromDB)
			})
		})
	})
}
