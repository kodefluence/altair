package model_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/codefluence-x/altair/provider/plugin/oauth/entity"
	"github.com/codefluence-x/altair/provider/plugin/oauth/model"
	"github.com/stretchr/testify/assert"
)

var oauthRefreshTokenRows = []string{
	"id",
	"oauth_access_token_id",
	"token",
	"expires_in",
	"created_at",
	"revoked_at",
}

func TestOauthRefreshToken(t *testing.T) {

	t.Run("Name", func(t *testing.T) {
		t.Run("Return a model's name", func(t *testing.T) {
			db, _, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}

			assert.Equal(t, "oauth-refresh-token-model", model.NewOauthRefreshToken(db).Name())
		})
	})

	t.Run("One", func(t *testing.T) {
		t.Run("Given oauth refresh token id", func(t *testing.T) {
			t.Run("Return entity oauth refresh token", func(t *testing.T) {
				db, mockdb, err := sqlmock.New()
				if err != nil {
					t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
				}

				data := entity.OauthRefreshToken{
					ID: 1,
				}

				rows := sqlmock.NewRows(oauthRefreshTokenRows).
					AddRow(
						data.ID,
						data.OauthAccessTokenID,
						data.Token,
						data.ExpiresIn,
						data.CreatedAt,
						data.RevokedAT,
					)

				mockdb.ExpectQuery(`select id, oauth_access_token_id, token, expires_in, created_at, revoked_at from oauth_refresh_tokens where id = \? limit 1`).
					WithArgs(1).
					WillReturnRows(rows)

				oauthRefreshTokenModel := model.NewOauthRefreshToken(db)
				dataFromDB, err := oauthRefreshTokenModel.One(context.Background(), 1)

				assert.Nil(t, err)
				assert.Nil(t, mockdb.ExpectationsWereMet())
				assert.Equal(t, data, dataFromDB)
			})

			t.Run("Row not found in database", func(t *testing.T) {
				db, mockdb, err := sqlmock.New()
				if err != nil {
					t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
				}

				mockdb.ExpectQuery(`select id, oauth_access_token_id, token, expires_in, created_at, revoked_at from oauth_refresh_tokens where id = \? limit 1`).
					WithArgs(1).
					WillReturnError(sql.ErrNoRows)

				oauthRefreshTokenModel := model.NewOauthRefreshToken(db)
				dataFromDB, err := oauthRefreshTokenModel.One(context.Background(), 1)

				assert.NotNil(t, err)
				assert.Equal(t, sql.ErrNoRows, err)
				assert.Nil(t, mockdb.ExpectationsWereMet())
				assert.Equal(t, entity.OauthRefreshToken{}, dataFromDB)
			})
		})
	})

	t.Run("OneByToken", func(t *testing.T) {
		t.Run("Given oauth refresh token", func(t *testing.T) {
			t.Run("Return entity oauth refresh token", func(t *testing.T) {
				db, mockdb, err := sqlmock.New()
				if err != nil {
					t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
				}

				data := entity.OauthRefreshToken{
					ID: 1,
				}

				rows := sqlmock.NewRows(oauthRefreshTokenRows).
					AddRow(
						data.ID,
						data.OauthAccessTokenID,
						data.Token,
						data.ExpiresIn,
						data.CreatedAt,
						data.RevokedAT,
					)

				mockdb.ExpectQuery(`select id, oauth_access_token_id, token, expires_in, created_at, revoked_at from oauth_refresh_tokens where token = \? and revoked_at is null limit 1`).
					WithArgs("token").
					WillReturnRows(rows)

				oauthRefreshTokenModel := model.NewOauthRefreshToken(db)
				dataFromDB, err := oauthRefreshTokenModel.OneByToken(context.Background(), "token")

				assert.Nil(t, err)
				assert.Nil(t, mockdb.ExpectationsWereMet())
				assert.Equal(t, data, dataFromDB)
			})

			t.Run("Row not found in database", func(t *testing.T) {
				db, mockdb, err := sqlmock.New()
				if err != nil {
					t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
				}

				mockdb.ExpectQuery(`select id, oauth_access_token_id, token, expires_in, created_at, revoked_at from oauth_refresh_tokens where token = \? and revoked_at is null limit 1`).
					WithArgs("token").
					WillReturnError(sql.ErrNoRows)

				oauthRefreshTokenModel := model.NewOauthRefreshToken(db)
				dataFromDB, err := oauthRefreshTokenModel.OneByToken(context.Background(), "token")

				assert.NotNil(t, err)
				assert.Equal(t, sql.ErrNoRows, err)
				assert.Nil(t, mockdb.ExpectationsWereMet())
				assert.Equal(t, entity.OauthRefreshToken{}, dataFromDB)
			})
		})
	})

	t.Run("Create", func(t *testing.T) {
		t.Run("Given context and refresh token insertable", func(t *testing.T) {
			t.Run("Return last inserted id", func(t *testing.T) {
				db, mockdb, err := sqlmock.New()
				if err != nil {
					t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
				}

				insertable := entity.OauthRefreshTokenInsertable{
					OauthAccessTokenID: 1,
					Token:              "token",
					ExpiresIn:          time.Now().Add(time.Hour * 24),
				}

				mockdb.ExpectExec(`insert into oauth_refresh_tokens \(oauth_access_token_id, token, expires_in, created_at, revoked_at\) values\(\?, \?, \?, now\(\), null\)`).
					WithArgs(insertable.OauthAccessTokenID, insertable.Token, insertable.ExpiresIn).
					WillReturnResult(sqlmock.NewResult(1, 1))

				oauthRefreshTokenModel := model.NewOauthRefreshToken(db)
				lastInsertedID, err := oauthRefreshTokenModel.Create(context.Background(), insertable)

				assert.Nil(t, mockdb.ExpectationsWereMet())
				assert.Nil(t, err)
				assert.Equal(t, 1, lastInsertedID)
			})

			t.Run("Unexpected error", func(t *testing.T) {
				db, mockdb, err := sqlmock.New()
				if err != nil {
					t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
				}

				insertable := entity.OauthRefreshTokenInsertable{
					OauthAccessTokenID: 1,
					Token:              "token",
					ExpiresIn:          time.Now().Add(time.Hour * 24),
				}

				mockdb.ExpectExec(`insert into oauth_refresh_tokens \(oauth_access_token_id, token, expires_in, created_at, revoked_at\) values\(\?, \?, \?, now\(\), null\)`).
					WithArgs(insertable.OauthAccessTokenID, insertable.Token, insertable.ExpiresIn).
					WillReturnError(errors.New("unexpected error"))

				oauthRefreshTokenModel := model.NewOauthRefreshToken(db)
				lastInsertedID, err := oauthRefreshTokenModel.Create(context.Background(), insertable)

				assert.Nil(t, mockdb.ExpectationsWereMet())
				assert.NotNil(t, err)
				assert.Equal(t, 0, lastInsertedID)
			})

			t.Run("Get last inserted id error", func(t *testing.T) {
				db, mockdb, err := sqlmock.New()
				if err != nil {
					t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
				}

				insertable := entity.OauthRefreshTokenInsertable{
					OauthAccessTokenID: 1,
					Token:              "token",
					ExpiresIn:          time.Now().Add(time.Hour * 24),
				}

				mockdb.ExpectExec(`insert into oauth_refresh_tokens \(oauth_access_token_id, token, expires_in, created_at, revoked_at\) values\(\?, \?, \?, now\(\), null\)`).
					WithArgs(insertable.OauthAccessTokenID, insertable.Token, insertable.ExpiresIn).
					WillReturnResult(sqlmock.NewErrorResult(errors.New("unexpected error")))

				oauthRefreshTokenModel := model.NewOauthRefreshToken(db)
				lastInsertedID, err := oauthRefreshTokenModel.Create(context.Background(), insertable)

				assert.Nil(t, mockdb.ExpectationsWereMet())
				assert.NotNil(t, err)
				assert.Equal(t, 0, lastInsertedID)
			})
		})

		t.Run("Given context, refresh token insertable and database transaction", func(t *testing.T) {
			t.Run("Return last inserted id", func(t *testing.T) {
				db, mockdb, err := sqlmock.New()
				if err != nil {
					t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
				}

				insertable := entity.OauthRefreshTokenInsertable{
					OauthAccessTokenID: 1,
					Token:              "token",
					ExpiresIn:          time.Now().Add(time.Hour * 24),
				}

				mockdb.ExpectBegin()
				mockdb.ExpectExec(`insert into oauth_refresh_tokens \(oauth_access_token_id, token, expires_in, created_at, revoked_at\) values\(\?, \?, \?, now\(\), null\)`).
					WithArgs(insertable.OauthAccessTokenID, insertable.Token, insertable.ExpiresIn).
					WillReturnResult(sqlmock.NewResult(1, 1))

				tx, _ := db.Begin()

				oauthRefreshTokenModel := model.NewOauthRefreshToken(db)
				lastInsertedID, err := oauthRefreshTokenModel.Create(context.Background(), insertable, tx)

				assert.Nil(t, mockdb.ExpectationsWereMet())
				assert.Nil(t, err)
				assert.Equal(t, 1, lastInsertedID)
			})
		})
	})

	t.Run("Revoke", func(t *testing.T) {
		t.Run("Given context and token", func(t *testing.T) {
			t.Run("Run gracefully", func(t *testing.T) {
				t.Run("Return nil", func(t *testing.T) {
					db, mockdb, err := sqlmock.New()
					if err != nil {
						t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
					}

					mockdb.ExpectExec(`update oauth_refresh_tokens set revoked_at = now\(\) where token = \?`).
						WithArgs("token").WillReturnResult(sqlmock.NewResult(1, 1))

					oauthRefreshTokenModel := model.NewOauthRefreshToken(db)
					err = oauthRefreshTokenModel.Revoke(context.Background(), "token")
					assert.Nil(t, err)
				})
			})

			t.Run("Execution error", func(t *testing.T) {
				t.Run("Return error", func(t *testing.T) {
					db, mockdb, err := sqlmock.New()
					if err != nil {
						t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
					}

					mockdb.ExpectExec(`update oauth_refresh_tokens set revoked_at = now\(\) where token = \?`).
						WithArgs("token").WillReturnError(errors.New("unexpected error"))

					oauthRefreshTokenModel := model.NewOauthRefreshToken(db)
					err = oauthRefreshTokenModel.Revoke(context.Background(), "token")
					assert.NotNil(t, err)
				})
			})

			t.Run("Get rows affected error", func(t *testing.T) {
				t.Run("Return error", func(t *testing.T) {
					db, mockdb, err := sqlmock.New()
					if err != nil {
						t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
					}

					mockdb.ExpectExec(`update oauth_refresh_tokens set revoked_at = now\(\) where token = \?`).
						WithArgs("token").WillReturnResult(sqlmock.NewErrorResult(errors.New("unexpected error")))

					oauthRefreshTokenModel := model.NewOauthRefreshToken(db)
					err = oauthRefreshTokenModel.Revoke(context.Background(), "token")
					assert.NotNil(t, err)
				})
			})

			t.Run("No rows affected", func(t *testing.T) {
				t.Run("Return sql.ErrNoRows error", func(t *testing.T) {
					db, mockdb, err := sqlmock.New()
					if err != nil {
						t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
					}

					mockdb.ExpectExec(`update oauth_refresh_tokens set revoked_at = now\(\) where token = \?`).
						WithArgs("token").WillReturnResult(sqlmock.NewResult(1, 0))

					oauthRefreshTokenModel := model.NewOauthRefreshToken(db)
					err = oauthRefreshTokenModel.Revoke(context.Background(), "token")
					assert.NotNil(t, err)
					assert.Equal(t, sql.ErrNoRows, err)
				})
			})
		})
	})
}
