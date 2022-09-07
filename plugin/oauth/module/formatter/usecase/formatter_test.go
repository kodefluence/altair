package usecase_test

import (
	"database/sql"
	"testing"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/kodefluence/altair/plugin/oauth/entity"
	"github.com/kodefluence/altair/plugin/oauth/module/formatter/usecase"
	"github.com/kodefluence/altair/util"
	"github.com/stretchr/testify/assert"
)

func TestOauthApplication(t *testing.T) {
	t.Run("Application list", func(t *testing.T) {
		t.Run("Given context and array of entity.OauthApplication", func(t *testing.T) {
			t.Run("Return array of entity.OauthApplicationJSON", func(t *testing.T) {
				oauthApplications := []entity.OauthApplication{
					{
						ID: 1,
						OwnerID: sql.NullInt64{
							Int64: 1,
							Valid: true,
						},
						Description: sql.NullString{
							String: "Application 01",
							Valid:  true,
						},
						Scopes: sql.NullString{
							String: "public users",
							Valid:  true,
						},
						ClientUID:    "clientuid01",
						ClientSecret: "clientsecret01",
						RevokedAt: mysql.NullTime{
							Time:  time.Now(),
							Valid: true,
						},
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					},
					{
						ID: 2,
						OwnerID: sql.NullInt64{
							Int64: 2,
							Valid: true,
						},
						Description: sql.NullString{
							String: "Application 02",
							Valid:  true,
						},
						Scopes: sql.NullString{
							String: "public users",
							Valid:  true,
						},
						ClientUID:    "clientuid02",
						ClientSecret: "clientsecret02",
						RevokedAt: mysql.NullTime{
							Time:  time.Now(),
							Valid: true,
						},
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					},
				}

				oauthApplicationJSON := usecase.NewFormatter().ApplicationList(oauthApplications)
				assert.Equal(t, len(oauthApplications), len(oauthApplicationJSON))
			})
		})
	})

	t.Run("OauthApplicationInsertable", func(t *testing.T) {
		t.Run("Given authorization request and oauth application", func(t *testing.T) {
			t.Run("Return oauth access grant insertable", func(t *testing.T) {

				oauthApplicationJSON := entity.OauthApplicationJSON{
					OwnerID:     util.IntToPointer(1),
					OwnerType:   util.StringToPointer("confidential"),
					Description: util.StringToPointer("Application 1"),
					Scopes:      util.StringToPointer("public user"),
				}

				insertable := usecase.NewFormatter().OauthApplicationInsertable(oauthApplicationJSON)

				assert.Equal(t, *oauthApplicationJSON.OwnerID, insertable.OwnerID)
				assert.Equal(t, *oauthApplicationJSON.OwnerType, insertable.OwnerType)
				assert.Equal(t, *oauthApplicationJSON.Description, insertable.Description)
				assert.Equal(t, *oauthApplicationJSON.Scopes, insertable.Scopes)
				assert.NotEqual(t, "", insertable.ClientUID)
				assert.NotEqual(t, "", insertable.ClientSecret)
			})
		})
	})
}
