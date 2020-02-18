package formatter_test

import (
	"testing"
	"time"

	"github.com/codefluence-x/altair/entity"
	"github.com/codefluence-x/altair/formatter"
	"github.com/codefluence-x/altair/util"
	"github.com/codefluence-x/aurelia"
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

func TestOauthFormatter(t *testing.T) {

	t.Run("AccessToken", func(t *testing.T) {
		t.Run("Given authorization request and oauth access token data", func(t *testing.T) {
			t.Run("Token not revoked", func(t *testing.T) {
				t.Run("Return oauth access token json", func(t *testing.T) {
					authorizationReq := entity.AuthorizationRequestJSON{
						ResponseType:    util.StringToPointer("token"),
						ResourceOwnerID: util.IntToPointer(1),
						ClientUID:       util.StringToPointer(aurelia.Hash("", "")),
						ClientSecret:    util.StringToPointer(aurelia.Hash("", "")),
						RedirectURI:     util.StringToPointer("http://github.com"),
						Scopes:          util.StringToPointer("public users"),
					}

					oauthAccessToken := entity.OauthAccessToken{
						ID:                 1,
						OauthApplicationID: 1,
						ResourceOwnerID:    1,
						Token:              aurelia.Hash("", ""),
						Scopes:             "public users stores",
						ExpiresIn:          time.Now().Add(time.Hour),
						CreatedAt:          time.Now(),
					}

					output := formatter.Oauth().AccessToken(authorizationReq, oauthAccessToken)

					assert.Equal(t, &oauthAccessToken.ID, output.ID)
					assert.Equal(t, &oauthAccessToken.OauthApplicationID, output.OauthApplicationID)
					assert.Equal(t, &oauthAccessToken.ResourceOwnerID, output.ResourceOwnerID)
					assert.Equal(t, &oauthAccessToken.Token, output.Token)
					assert.Equal(t, &oauthAccessToken.CreatedAt, output.CreatedAt)
					assert.LessOrEqual(t, *output.ExpiresIn, int(oauthAccessToken.ExpiresIn.Sub(time.Now()).Seconds()))
					assert.Greater(t, *output.ExpiresIn, 3500)
					assert.Nil(t, output.RevokedAT)

					assert.Equal(t, authorizationReq.RedirectURI, output.RedirectURI)
					assert.Equal(t, authorizationReq.Scopes, output.Scopes)

					assert.Equal(t, &oauthAccessToken.ID, output.ID)
				})
			})

			t.Run("Token already revoked", func(t *testing.T) {
				authorizationReq := entity.AuthorizationRequestJSON{
					ResponseType:    util.StringToPointer("token"),
					ResourceOwnerID: util.IntToPointer(1),
					ClientUID:       util.StringToPointer(aurelia.Hash("", "")),
					ClientSecret:    util.StringToPointer(aurelia.Hash("", "")),
					RedirectURI:     util.StringToPointer("http://github.com"),
					Scopes:          util.StringToPointer("public users"),
				}

				oauthAccessToken := entity.OauthAccessToken{
					ID:                 1,
					OauthApplicationID: 1,
					ResourceOwnerID:    1,
					Token:              aurelia.Hash("", ""),
					Scopes:             "public users stores",
					ExpiresIn:          time.Now().Add(-time.Hour),
					CreatedAt:          time.Now().Add(-time.Hour * 2),
					RevokedAT: mysql.NullTime{
						Valid: true,
						Time:  time.Now(),
					},
				}

				output := formatter.Oauth().AccessToken(authorizationReq, oauthAccessToken)

				assert.Equal(t, &oauthAccessToken.ID, output.ID)
				assert.Equal(t, &oauthAccessToken.OauthApplicationID, output.OauthApplicationID)
				assert.Equal(t, &oauthAccessToken.ResourceOwnerID, output.ResourceOwnerID)
				assert.Equal(t, &oauthAccessToken.Token, output.Token)
				assert.Equal(t, &oauthAccessToken.CreatedAt, output.CreatedAt)
				assert.Equal(t, 0, *output.ExpiresIn)
				assert.Equal(t, oauthAccessToken.RevokedAT.Time, *output.RevokedAT)

				assert.Equal(t, authorizationReq.RedirectURI, output.RedirectURI)
				assert.Equal(t, authorizationReq.Scopes, output.Scopes)

				assert.Equal(t, &oauthAccessToken.ID, output.ID)
			})
		})
	})
}
