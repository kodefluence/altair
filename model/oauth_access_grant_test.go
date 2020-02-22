package model_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

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

	t.Run("Create", func(t *testing.T) {
		t.Run("Given context and access grant insertable", func(t *testing.T) {
			t.Run("Return last inserted id", func(t *testing.T) {
				db, mockdb, err := sqlmock.New()
				if err != nil {
					t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
				}

				insertable := entity.OauthAccessGrantInsertable{
					OauthApplicationID: 1,
					ResourceOwnerID:    1,
					Code:               "token",
					Scopes:             "public users stores",
					RedirectURI:        "https://github.com",
					ExpiresIn:          time.Now().Add(time.Hour * 24),
				}

				mockdb.ExpectExec(`insert into oauth_access_grants \(oauth_application_id, resource_owner_id, scopes, code, redirect_uri, expires_in, created_at, revoked_at\) values\(\?, \?, \?, \?, \?, \?, now\(\), null\)`).
					WithArgs(insertable.OauthApplicationID, insertable.ResourceOwnerID, insertable.Scopes, insertable.Code, insertable.RedirectURI, insertable.ExpiresIn).
					WillReturnResult(sqlmock.NewResult(1, 1))

				oauthAccessGrantModel := model.OauthAccessGrant(db)
				lastInsertedID, err := oauthAccessGrantModel.Create(context.Background(), insertable)

				assert.Nil(t, mockdb.ExpectationsWereMet())
				assert.Nil(t, err)
				assert.Equal(t, 1, lastInsertedID)
			})

			t.Run("Unexpected error", func(t *testing.T) {
				db, mockdb, err := sqlmock.New()
				if err != nil {
					t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
				}

				insertable := entity.OauthAccessGrantInsertable{
					OauthApplicationID: 1,
					ResourceOwnerID:    1,
					Code:               "token",
					Scopes:             "public users stores",
					RedirectURI:        "https://github.com",
					ExpiresIn:          time.Now().Add(time.Hour * 24),
				}

				mockdb.ExpectExec(`insert into oauth_access_grants \(oauth_application_id, resource_owner_id, scopes, code, redirect_uri, expires_in, created_at, revoked_at\) values\(\?, \?, \?, \?, \?, \?, now\(\), null\)`).
					WithArgs(insertable.OauthApplicationID, insertable.ResourceOwnerID, insertable.Scopes, insertable.Code, insertable.RedirectURI, insertable.ExpiresIn).
					WillReturnError(errors.New("unexpected error"))

				oauthAccessGrantModel := model.OauthAccessGrant(db)
				lastInsertedID, err := oauthAccessGrantModel.Create(context.Background(), insertable)

				assert.Nil(t, mockdb.ExpectationsWereMet())
				assert.NotNil(t, err)
				assert.Equal(t, 0, lastInsertedID)
			})

			t.Run("Get last inserted id error", func(t *testing.T) {
				db, mockdb, err := sqlmock.New()
				if err != nil {
					t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
				}

				insertable := entity.OauthAccessGrantInsertable{
					OauthApplicationID: 1,
					ResourceOwnerID:    1,
					Code:               "token",
					Scopes:             "public users stores",
					RedirectURI:        "https://github.com",
					ExpiresIn:          time.Now().Add(time.Hour * 24),
				}

				mockdb.ExpectExec(`insert into oauth_access_grants \(oauth_application_id, resource_owner_id, scopes, code, redirect_uri, expires_in, created_at, revoked_at\) values\(\?, \?, \?, \?, \?, \?, now\(\), null\)`).
					WithArgs(insertable.OauthApplicationID, insertable.ResourceOwnerID, insertable.Scopes, insertable.Code, insertable.RedirectURI, insertable.ExpiresIn).
					WillReturnResult(sqlmock.NewErrorResult(errors.New("unexpected error")))

				oauthAccessGrantModel := model.OauthAccessGrant(db)
				lastInsertedID, err := oauthAccessGrantModel.Create(context.Background(), insertable)

				assert.Nil(t, mockdb.ExpectationsWereMet())
				assert.NotNil(t, err)
				assert.Equal(t, 0, lastInsertedID)
			})
		})

		t.Run("Given context, access token insertable and database transaction", func(t *testing.T) {
			t.Run("Return last inserted id", func(t *testing.T) {
				db, mockdb, err := sqlmock.New()
				if err != nil {
					t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
				}

				insertable := entity.OauthAccessGrantInsertable{
					OauthApplicationID: 1,
					ResourceOwnerID:    1,
					Code:               "token",
					Scopes:             "public users stores",
					RedirectURI:        "https://github.com",
					ExpiresIn:          time.Now().Add(time.Hour * 24),
				}

				mockdb.ExpectBegin()
				tx, _ := db.Begin()

				mockdb.ExpectExec(`insert into oauth_access_grants \(oauth_application_id, resource_owner_id, scopes, code, redirect_uri, expires_in, created_at, revoked_at\) values\(\?, \?, \?, \?, \?, \?, now\(\), null\)`).
					WithArgs(insertable.OauthApplicationID, insertable.ResourceOwnerID, insertable.Scopes, insertable.Code, insertable.RedirectURI, insertable.ExpiresIn).
					WillReturnResult(sqlmock.NewResult(1, 1))

				oauthAccessGrantModel := model.OauthAccessGrant(db)
				lastInsertedID, err := oauthAccessGrantModel.Create(context.Background(), insertable, tx)

				assert.Nil(t, mockdb.ExpectationsWereMet())
				assert.Nil(t, err)
				assert.Equal(t, 1, lastInsertedID)
			})
		})
	})
}
