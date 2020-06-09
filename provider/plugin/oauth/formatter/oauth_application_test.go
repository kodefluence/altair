package formatter_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/codefluence-x/altair/provider/plugin/oauth/entity"
	"github.com/codefluence-x/altair/provider/plugin/oauth/formatter"

	"github.com/go-sql-driver/mysql"
	"gotest.tools/assert"
)

func TestOauthApplication(t *testing.T) {

	t.Run("Application list", func(t *testing.T) {
		t.Run("Given context and array of entity.OauthApplication", func(t *testing.T) {
			t.Run("Return array of entity.OauthApplicationJSON", func(t *testing.T) {
				oauthApplications := []entity.OauthApplication{
					entity.OauthApplication{
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
					entity.OauthApplication{
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

				oauthApplicationJSON := formatter.OauthApplication().ApplicationList(context.Background(), oauthApplications)
				assert.Equal(t, len(oauthApplications), len(oauthApplicationJSON))
			})
		})
	})
}
