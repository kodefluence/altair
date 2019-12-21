package model_test

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/codefluence-x/altair/entity"
	"github.com/codefluence-x/altair/model"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestOauthApplication(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	t.Run("Name", func(t *testing.T) {
		t.Run("Return a model's name", func(t *testing.T) {
			db, _, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}

			assert.Equal(t, "oauth-application-model", model.OauthApplication(db).Name())
		})
	})

	t.Run("Paginate", func(t *testing.T) {
		t.Run("Given offset and limit", func(t *testing.T) {
			t.Run("Return array of oauth applications data", func(t *testing.T) {
				db, mockdb, err := sqlmock.New()
				if err != nil {
					t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
				}

				expectedOauthApplications := []entity.OauthApplication{
					entity.OauthApplication{
						ID: 1,
					},
					entity.OauthApplication{
						ID: 2,
					},
				}

				rows := sqlmock.NewRows([]string{
					"id",
					"owner_id",
					"description",
					"scopes",
					"client_uid",
					"client_secret",
					"revoked_at",
					"created_at",
					"updated_at",
				})

				for _, x := range expectedOauthApplications {
					rows.AddRow(
						x.ID, x.OwnerID, x.Description, x.Scopes, x.ClientUID,
						x.ClientSecret, x.RevokedAt, x.CreatedAt, x.UpdatedAt,
					)
				}

				mockdb.ExpectQuery(`select \* from oauth_applications limit \?, \?`).WillReturnRows(rows)

				oauthApplicationModel := model.OauthApplication(db)
				oauthApplications, err := oauthApplicationModel.Paginate(context.Background(), 0, 10)

				assert.Nil(t, err)
				assert.Equal(t, expectedOauthApplications, oauthApplications)
				assert.Nil(t, mockdb.ExpectationsWereMet())
			})

			t.Run("Return an error", func(t *testing.T) {
				t.Run("Query error", func(t *testing.T) {
					db, mockdb, err := sqlmock.New()
					if err != nil {
						t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
					}

					mockdb.ExpectQuery(`select \* from oauth_applications limit \?, \?`).WillReturnError(errors.New("Unexpected error"))

					oauthApplicationModel := model.OauthApplication(db)
					oauthApplications, err := oauthApplicationModel.Paginate(context.Background(), 0, 10)

					assert.NotNil(t, err)
					assert.Equal(t, []entity.OauthApplication(nil), oauthApplications)
					assert.Nil(t, mockdb.ExpectationsWereMet())
				})

				t.Run("Row scan error", func(t *testing.T) {
					db, mockdb, err := sqlmock.New()
					if err != nil {
						t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
					}

					expectedOauthApplications := []entity.OauthApplication{
						entity.OauthApplication{
							ID: 1,
						},
						entity.OauthApplication{
							ID: 2,
						},
					}

					rows := sqlmock.NewRows([]string{
						"id",
						"owner_id",
						"description",
						"scopes",
						"client_uid",
						"client_secret",
						"revoked_at",
						"created_at",
						"updated_at",
						"some_new_column_maybe",
					})

					for _, x := range expectedOauthApplications {
						rows.AddRow(
							x.ID, x.OwnerID, x.Description, x.Scopes, x.ClientUID,
							x.ClientSecret, x.RevokedAt, x.CreatedAt, x.UpdatedAt, x.UpdatedAt,
						)
					}

					mockdb.ExpectQuery(`select \* from oauth_applications limit \?, \?`).WillReturnRows(rows)

					oauthApplicationModel := model.OauthApplication(db)
					oauthApplications, err := oauthApplicationModel.Paginate(context.Background(), 0, 10)

					assert.NotNil(t, err)
					assert.Equal(t, []entity.OauthApplication(nil), oauthApplications)
					assert.Nil(t, mockdb.ExpectationsWereMet())
				})
			})
		})
	})
}
