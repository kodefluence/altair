package model_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/codefluence-x/altair/entity"
	"github.com/codefluence-x/altair/model"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

var oauthApplicationModelRows = []string{
	"id",
	"owner_id",
	"owner_type",
	"description",
	"scopes",
	"client_uid",
	"client_secret",
	"revoked_at",
	"created_at",
	"updated_at",
}

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

				rows := sqlmock.NewRows(oauthApplicationModelRows)

				for _, x := range expectedOauthApplications {
					rows.AddRow(
						x.ID, x.OwnerID, x.OwnerType, x.Description, x.Scopes, x.ClientUID,
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

					rows := sqlmock.NewRows(append(oauthApplicationModelRows, "some_new_column_maybe"))

					for _, x := range expectedOauthApplications {
						rows.AddRow(
							x.ID, x.OwnerID, x.OwnerType, x.Description, x.Scopes, x.ClientUID,
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

	t.Run("Count", func(t *testing.T) {
		t.Run("Return total data of oauth application", func(t *testing.T) {
			db, mockdb, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}

			rows := sqlmock.NewRows([]string{"total"})
			rows.AddRow(100)

			mockdb.ExpectQuery(`select count\(\*\) as total from oauth_applications where revoked_at is null`).
				WillReturnRows(rows)

			oauthApplicationModel := model.OauthApplication(db)
			total, err := oauthApplicationModel.Count(context.Background())

			assert.Equal(t, 100, total)
			assert.Nil(t, err)
			assert.Nil(t, mockdb.ExpectationsWereMet())
		})
	})

	t.Run("Create", func(t *testing.T) {
		t.Run("Given context and oauth application json entity", func(t *testing.T) {
			t.Run("Return last inserted id", func(t *testing.T) {
				db, mockdb, err := sqlmock.New()
				if err != nil {
					t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
				}

				oauthApplication := entity.OauthApplicationInsertable{
					OwnerID:      nil,
					OwnerType:    "confidential",
					Description:  "This is description",
					Scopes:       "This is scopes",
					ClientUID:    "This is client uid",
					ClientSecret: "This is client secret",
				}

				mockdb.ExpectExec(`insert into oauth_applications \(owner_id, owner_type, description, scopes, client_uid, client_secret, revoked_at, created_at, updated_at\) values\(\?, \?, \?, \?, \?, \?, null, now\(\), now\(\)\)`).
					WithArgs(oauthApplication.OwnerID, oauthApplication.OwnerType, oauthApplication.Description, oauthApplication.Scopes, sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))

				oauthApplicationModel := model.OauthApplication(db)
				lastInsertedID, err := oauthApplicationModel.Create(context.Background(), oauthApplication)

				assert.Nil(t, err)
				assert.Nil(t, mockdb.ExpectationsWereMet())
				assert.Equal(t, 1, lastInsertedID)
			})

			t.Run("Given sql transactions", func(t *testing.T) {
				t.Run("Return last inserted id", func(t *testing.T) {
					db, mockdb, err := sqlmock.New()
					if err != nil {
						t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
					}

					oauthApplication := entity.OauthApplicationInsertable{
						OwnerID:      nil,
						OwnerType:    "confidential",
						Description:  "This is description",
						Scopes:       "This is scopes",
						ClientUID:    "This is client uid",
						ClientSecret: "This is client secret",
					}

					mockdb.ExpectBegin()
					tx, _ := db.Begin()

					mockdb.ExpectExec(`insert into oauth_applications \(owner_id, owner_type, description, scopes, client_uid, client_secret, revoked_at, created_at, updated_at\) values\(\?, \?, \?, \?, \?, \?, null, now\(\), now\(\)\)`).
						WithArgs(oauthApplication.OwnerID, oauthApplication.OwnerType, oauthApplication.Description, oauthApplication.Scopes, sqlmock.AnyArg(), sqlmock.AnyArg()).
						WillReturnResult(sqlmock.NewResult(1, 1))

					oauthApplicationModel := model.OauthApplication(db)
					lastInsertedID, err := oauthApplicationModel.Create(context.Background(), oauthApplication, tx)

					assert.Nil(t, err)
					assert.Nil(t, mockdb.ExpectationsWereMet())
					assert.Equal(t, 1, lastInsertedID)
				})
			})

			t.Run("Return zero value of last inserted id and error", func(t *testing.T) {
				t.Run("Execution error", func(t *testing.T) {
					db, mockdb, err := sqlmock.New()
					if err != nil {
						t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
					}

					oauthApplication := entity.OauthApplicationInsertable{
						OwnerID:      nil,
						OwnerType:    "confidential",
						Description:  "This is description",
						Scopes:       "This is scopes",
						ClientUID:    "This is client uid",
						ClientSecret: "This is client secret",
					}

					mockdb.ExpectExec(`insert into oauth_applications \(owner_id, owner_type, description, scopes, client_uid, client_secret, revoked_at, created_at, updated_at\) values\(\?, \?, \?, \?, \?, \?, null, now\(\), now\(\)\)`).
						WithArgs(oauthApplication.OwnerID, oauthApplication.OwnerType, oauthApplication.Description, oauthApplication.Scopes, sqlmock.AnyArg(), sqlmock.AnyArg()).
						WillReturnError(errors.New("unexpected error"))

					oauthApplicationModel := model.OauthApplication(db)
					lastInsertedID, err := oauthApplicationModel.Create(context.Background(), oauthApplication)

					assert.NotNil(t, err)
					assert.Nil(t, mockdb.ExpectationsWereMet())
					assert.Equal(t, 0, lastInsertedID)
				})

				t.Run("Get last inserted id error", func(t *testing.T) {
					db, mockdb, err := sqlmock.New()
					if err != nil {
						t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
					}

					oauthApplication := entity.OauthApplicationInsertable{
						OwnerID:      nil,
						OwnerType:    "confidential",
						Description:  "This is description",
						Scopes:       "This is scopes",
						ClientUID:    "This is client uid",
						ClientSecret: "This is client secret",
					}

					mockdb.ExpectExec(`insert into oauth_applications \(owner_id, owner_type, description, scopes, client_uid, client_secret, revoked_at, created_at, updated_at\) values\(\?, \?, \?, \?, \?, \?, null, now\(\), now\(\)\)`).
						WithArgs(oauthApplication.OwnerID, oauthApplication.OwnerType, oauthApplication.Description, oauthApplication.Scopes, sqlmock.AnyArg(), sqlmock.AnyArg()).
						WillReturnResult(sqlmock.NewErrorResult(errors.New("unexpected error")))

					oauthApplicationModel := model.OauthApplication(db)
					lastInsertedID, err := oauthApplicationModel.Create(context.Background(), oauthApplication)

					assert.NotNil(t, err)
					assert.Nil(t, mockdb.ExpectationsWereMet())
					assert.Equal(t, 0, lastInsertedID)
				})
			})
		})
	})

	t.Run("One", func(t *testing.T) {
		t.Run("Given oauth application id", func(t *testing.T) {
			t.Run("Return entity oauth application", func(t *testing.T) {
				db, mockdb, err := sqlmock.New()
				if err != nil {
					t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
				}

				data := entity.OauthApplication{
					ID: 1,
				}

				rows := sqlmock.NewRows(oauthApplicationModelRows).
					AddRow(
						data.ID, data.OwnerID, data.OwnerType, data.Description,
						data.Scopes, data.ClientUID, data.ClientSecret,
						data.RevokedAt, data.CreatedAt, data.UpdatedAt,
					)

				mockdb.ExpectQuery(`select \* from oauth_applications where id = \?`).
					WithArgs(1).
					WillReturnRows(rows)

				oauthApplicationModel := model.OauthApplication(db)
				dataFromDB, err := oauthApplicationModel.One(context.Background(), 1)

				assert.Nil(t, err)
				assert.Nil(t, mockdb.ExpectationsWereMet())
				assert.Equal(t, data, dataFromDB)
			})

			t.Run("Row not found in database", func(t *testing.T) {
				t.Run("Return sql.ErrNoRows error", func(t *testing.T) {
					db, mockdb, err := sqlmock.New()
					if err != nil {
						t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
					}

					mockdb.ExpectQuery(`select \* from oauth_applications where id = \?`).
						WithArgs(1).
						WillReturnError(sql.ErrNoRows)

					oauthApplicationModel := model.OauthApplication(db)
					dataFromDB, err := oauthApplicationModel.One(context.Background(), 1)

					assert.NotNil(t, err)
					assert.Nil(t, mockdb.ExpectationsWereMet())
					assert.Equal(t, entity.OauthApplication{}, dataFromDB)
				})
			})
		})
	})

	t.Run("OneByUIDandSecret", func(t *testing.T) {
		t.Run("Given oauth application id", func(t *testing.T) {
			t.Run("Return entity oauth application", func(t *testing.T) {
				db, mockdb, err := sqlmock.New()
				if err != nil {
					t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
				}

				data := entity.OauthApplication{
					ID: 1,
				}

				rows := sqlmock.NewRows(oauthApplicationModelRows).
					AddRow(
						data.ID, data.OwnerID, data.OwnerType, data.Description,
						data.Scopes, data.ClientUID, data.ClientSecret,
						data.RevokedAt, data.CreatedAt, data.UpdatedAt,
					)

				mockdb.ExpectQuery(`select \* from oauth_applications where client_uid = \? and client_secret = \? limit 1`).
					WithArgs("sample_client_uid", "sample_client_secret").
					WillReturnRows(rows)

				oauthApplicationModel := model.OauthApplication(db)
				dataFromDB, err := oauthApplicationModel.OneByUIDandSecret(context.Background(), "sample_client_uid", "sample_client_secret")

				assert.Nil(t, err)
				assert.Nil(t, mockdb.ExpectationsWereMet())
				assert.Equal(t, data, dataFromDB)
			})

			t.Run("Row not found in database", func(t *testing.T) {
				t.Run("Return sql.ErrNoRows error", func(t *testing.T) {
					db, mockdb, err := sqlmock.New()
					if err != nil {
						t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
					}

					mockdb.ExpectQuery(`select \* from oauth_applications where client_uid = \? and client_secret = \? limit 1`).
						WithArgs("sample_client_uid", "sample_client_secret").
						WillReturnError(sql.ErrNoRows)

					oauthApplicationModel := model.OauthApplication(db)
					dataFromDB, err := oauthApplicationModel.OneByUIDandSecret(context.Background(), "sample_client_uid", "sample_client_secret")

					assert.NotNil(t, err)
					assert.Nil(t, mockdb.ExpectationsWereMet())
					assert.Equal(t, entity.OauthApplication{}, dataFromDB)
				})
			})
		})
	})
}
