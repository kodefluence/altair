package model_test

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/codefluence-x/altair/entity"
	"github.com/codefluence-x/altair/model"
	"github.com/stretchr/testify/assert"
)

var oauthAccessTokenModelRows = []string{
	"id",
	"oauth_application_id",
	"resource_owner_id",
	"token",
	"scopes",
	"expires_in",
	"created_at",
	"revoked_at",
}

func TestOauthAccessToken(t *testing.T) {

	t.Run("Name", func(t *testing.T) {
		t.Run("Return a model's name", func(t *testing.T) {
			db, _, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}

			assert.Equal(t, "oauth-access-token-model", model.OauthAccessToken(db).Name())
		})
	})

	t.Run("One", func(t *testing.T) {
		t.Run("Given oauth access token id", func(t *testing.T) {
			t.Run("Return entity oauth access token", func(t *testing.T) {
				db, mockdb, err := sqlmock.New()
				if err != nil {
					t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
				}

				data := entity.OauthAccessToken{
					ID: 1,
				}

				rows := sqlmock.NewRows(oauthAccessTokenModelRows).
					AddRow(
						data.ID,
						data.OauthApplicationID,
						data.ResourceOwnerID,
						data.Token,
						data.Scopes,
						data.ExpiresIn,
						data.CreatedAt,
						data.RevokedAT,
					)

				mockdb.ExpectQuery(`select \* from oauth_access_tokens where id = \? limit 1`).
					WithArgs(1).
					WillReturnRows(rows)

				oauthAccessTokenModel := model.OauthAccessToken(db)
				dataFromDB, err := oauthAccessTokenModel.One(context.Background(), 1)

				assert.Nil(t, err)
				assert.Nil(t, mockdb.ExpectationsWereMet())
				assert.Equal(t, data, dataFromDB)
			})

			// t.Run("Row not found in database", func(t *testing.T) {
			// 	t.Run("Return sql.ErrNoRows error", func(t *testing.T) {
			// 		db, mockdb, err := sqlmock.New()
			// 		if err != nil {
			// 			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			// 		}

			// 		mockdb.ExpectQuery(`select \* from oauth_applications where id = \?`).
			// 			WithArgs(1).
			// 			WillReturnError(sql.ErrNoRows)

			// 		oauthApplicationModel := model.OauthApplication(db)
			// 		dataFromDB, err := oauthApplicationModel.One(context.Background(), 1)

			// 		assert.NotNil(t, err)
			// 		assert.Nil(t, mockdb.ExpectationsWereMet())
			// 		assert.Equal(t, entity.OauthApplication{}, dataFromDB)
			// 	})
			// })
		})
	})
}
