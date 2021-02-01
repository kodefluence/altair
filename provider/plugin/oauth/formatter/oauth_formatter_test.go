package formatter_test

import (
	"database/sql"
	"testing"
	"time"

	"github.com/codefluence-x/altair/provider/plugin/oauth/entity"
	"github.com/codefluence-x/altair/provider/plugin/oauth/formatter"

	"github.com/codefluence-x/altair/util"
	"github.com/codefluence-x/aurelia"
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

func TestOauthFormatter(t *testing.T) {

	t.Run("AccessGrant", func(t *testing.T) {
		t.Run("Given authorization request and oauth access grant data", func(t *testing.T) {
			t.Run("Code not revoked", func(t *testing.T) {
				oauthAccessGrant := entity.OauthAccessGrant{
					ID:                 1,
					OauthApplicationID: 1,
					ResourceOwnerID:    1,
					Code:               util.SHA1(),
					RevokedAT: mysql.NullTime{
						Time:  time.Time{},
						Valid: false,
					},
					RedirectURI: sql.NullString{
						String: "https://github.com",
						Valid:  true,
					},
					Scopes: sql.NullString{
						String: "public users stores",
						Valid:  true,
					},
					ExpiresIn: time.Now().Add(time.Hour),
					CreatedAt: time.Now(),
				}

				output := formatter.Oauth().AccessGrant(oauthAccessGrant)

				assert.Equal(t, &oauthAccessGrant.ID, output.ID)
				assert.Equal(t, &oauthAccessGrant.OauthApplicationID, output.OauthApplicationID)
				assert.Equal(t, &oauthAccessGrant.ResourceOwnerID, output.ResourceOwnerID)
				assert.Equal(t, &oauthAccessGrant.Code, output.Code)
				assert.Equal(t, &oauthAccessGrant.RedirectURI.String, output.RedirectURI)
				assert.Equal(t, &oauthAccessGrant.CreatedAt, output.CreatedAt)
				assert.LessOrEqual(t, *output.ExpiresIn, int(oauthAccessGrant.ExpiresIn.Sub(time.Now()).Seconds()))
				assert.Greater(t, *output.ExpiresIn, 3500)
				assert.Nil(t, output.RevokedAT)

				assert.Equal(t, &oauthAccessGrant.RedirectURI.String, output.RedirectURI)
				assert.Equal(t, &oauthAccessGrant.Scopes.String, output.Scopes)

				assert.Equal(t, &oauthAccessGrant.ID, output.ID)
			})

			t.Run("Code already revoked", func(t *testing.T) {
				oauthAccessGrant := entity.OauthAccessGrant{
					ID:                 1,
					OauthApplicationID: 1,
					ResourceOwnerID:    1,
					Code:               util.SHA1(),
					RedirectURI: sql.NullString{
						String: "https://github.com",
						Valid:  true,
					},
					Scopes: sql.NullString{
						String: "public users stores",
						Valid:  true,
					},
					ExpiresIn: time.Now().Add(-time.Hour),
					CreatedAt: time.Now().Add(-time.Hour * 2),
					RevokedAT: mysql.NullTime{
						Valid: true,
						Time:  time.Now(),
					},
				}

				output := formatter.Oauth().AccessGrant(oauthAccessGrant)

				assert.Equal(t, &oauthAccessGrant.ID, output.ID)
				assert.Equal(t, &oauthAccessGrant.OauthApplicationID, output.OauthApplicationID)
				assert.Equal(t, &oauthAccessGrant.ResourceOwnerID, output.ResourceOwnerID)
				assert.Equal(t, &oauthAccessGrant.Code, output.Code)
				assert.Equal(t, &oauthAccessGrant.RedirectURI.String, output.RedirectURI)
				assert.Equal(t, &oauthAccessGrant.CreatedAt, output.CreatedAt)
				assert.Equal(t, 0, *output.ExpiresIn)
				assert.Equal(t, &oauthAccessGrant.RevokedAT.Time, output.RevokedAT)

				assert.Equal(t, &oauthAccessGrant.RedirectURI.String, output.RedirectURI)
				assert.Equal(t, &oauthAccessGrant.Scopes.String, output.Scopes)

				assert.Equal(t, &oauthAccessGrant.ID, output.ID)
			})
		})
	})

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
						Scopes: sql.NullString{
							String: "public users stores",
							Valid:  true,
						},
						ExpiresIn: time.Now().Add(time.Hour),
						CreatedAt: time.Now(),
					}

					output := formatter.Oauth().AccessToken(oauthAccessToken, *authorizationReq.RedirectURI, nil)

					assert.Equal(t, &oauthAccessToken.ID, output.ID)
					assert.Equal(t, &oauthAccessToken.OauthApplicationID, output.OauthApplicationID)
					assert.Equal(t, &oauthAccessToken.ResourceOwnerID, output.ResourceOwnerID)
					assert.Equal(t, &oauthAccessToken.Token, output.Token)
					assert.Equal(t, &oauthAccessToken.CreatedAt, output.CreatedAt)
					assert.Equal(t, &oauthAccessToken.Scopes.String, output.Scopes)
					assert.LessOrEqual(t, *output.ExpiresIn, int(oauthAccessToken.ExpiresIn.Sub(time.Now()).Seconds()))
					assert.Greater(t, *output.ExpiresIn, 3500)
					assert.Nil(t, output.RevokedAT)

					assert.Equal(t, authorizationReq.RedirectURI, output.RedirectURI)

					assert.Equal(t, &oauthAccessToken.ID, output.ID)
				})
			})

			t.Run("Token not revoked with refresh token", func(t *testing.T) {
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
						Scopes: sql.NullString{
							String: "public users stores",
							Valid:  true,
						},
						ExpiresIn: time.Now().Add(time.Hour),
						CreatedAt: time.Now(),
					}

					output := formatter.Oauth().AccessToken(oauthAccessToken, *authorizationReq.RedirectURI, nil)

					assert.Equal(t, &oauthAccessToken.ID, output.ID)
					assert.Equal(t, &oauthAccessToken.OauthApplicationID, output.OauthApplicationID)
					assert.Equal(t, &oauthAccessToken.ResourceOwnerID, output.ResourceOwnerID)
					assert.Equal(t, &oauthAccessToken.Token, output.Token)
					assert.Equal(t, &oauthAccessToken.CreatedAt, output.CreatedAt)
					assert.Equal(t, &oauthAccessToken.Scopes.String, output.Scopes)
					assert.LessOrEqual(t, *output.ExpiresIn, int(oauthAccessToken.ExpiresIn.Sub(time.Now()).Seconds()))
					assert.Greater(t, *output.ExpiresIn, 3500)
					assert.Nil(t, output.RevokedAT)

					assert.Equal(t, authorizationReq.RedirectURI, output.RedirectURI)

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
					Scopes: sql.NullString{
						String: "public users stores",
						Valid:  true,
					},
					ExpiresIn: time.Now().Add(-time.Hour),
					CreatedAt: time.Now().Add(-time.Hour * 2),
					RevokedAT: mysql.NullTime{
						Valid: true,
						Time:  time.Now(),
					},
				}

				oauthRefreshToken := entity.OauthRefreshToken{
					ID:                 1,
					OauthAccessTokenID: 2,
					Token:              "token",
					ExpiresIn:          time.Now().Add(time.Hour),
					CreatedAt:          time.Now(),
					RevokedAT: mysql.NullTime{
						Valid: false,
					},
				}

				oauthRefreshTokenJSON := formatter.Oauth().RefreshToken(oauthRefreshToken)

				output := formatter.Oauth().AccessToken(oauthAccessToken, *authorizationReq.RedirectURI, &oauthRefreshTokenJSON)

				assert.Equal(t, &oauthAccessToken.ID, output.ID)
				assert.Equal(t, &oauthAccessToken.OauthApplicationID, output.OauthApplicationID)
				assert.Equal(t, &oauthAccessToken.ResourceOwnerID, output.ResourceOwnerID)
				assert.Equal(t, &oauthAccessToken.Token, output.Token)
				assert.Equal(t, &oauthAccessToken.CreatedAt, output.CreatedAt)
				assert.Equal(t, &oauthAccessToken.Scopes.String, output.Scopes)
				assert.Equal(t, 0, *output.ExpiresIn)
				assert.Equal(t, oauthAccessToken.RevokedAT.Time, *output.RevokedAT)

				assert.Equal(t, authorizationReq.RedirectURI, output.RedirectURI)

				assert.Equal(t, &oauthAccessToken.ID, output.ID)
			})
		})
	})

	t.Run("RefreshToken", func(t *testing.T) {
		t.Run("Given oauth refresh token data", func(t *testing.T) {
			t.Run("Token not revoked", func(t *testing.T) {
				t.Run("Return oauth refresh token json", func(t *testing.T) {
					oauthRefreshToken := entity.OauthRefreshToken{
						ID:                 1,
						OauthAccessTokenID: 2,
						Token:              "token",
						ExpiresIn:          time.Now().Add(time.Hour),
						CreatedAt:          time.Now(),
						RevokedAT: mysql.NullTime{
							Valid: false,
						},
					}

					output := formatter.Oauth().RefreshToken(oauthRefreshToken)

					assert.Equal(t, &oauthRefreshToken.CreatedAt, output.CreatedAt)
					assert.Equal(t, &oauthRefreshToken.Token, output.Token)
					assert.LessOrEqual(t, *output.ExpiresIn, int(oauthRefreshToken.ExpiresIn.Sub(time.Now()).Seconds()))
					assert.Nil(t, output.RevokedAT)
				})
			})

			t.Run("Token already revoked", func(t *testing.T) {
				oauthRefreshToken := entity.OauthRefreshToken{
					ID:                 1,
					OauthAccessTokenID: 2,
					Token:              "token",
					ExpiresIn:          time.Now().Add(-time.Hour),
					CreatedAt:          time.Now().Add(-time.Hour * 2),
					RevokedAT: mysql.NullTime{
						Time:  time.Now().Add(-time.Hour),
						Valid: true,
					},
				}

				output := formatter.Oauth().RefreshToken(oauthRefreshToken)

				assert.Equal(t, &oauthRefreshToken.CreatedAt, output.CreatedAt)
				assert.Equal(t, &oauthRefreshToken.Token, output.Token)
				assert.Equal(t, 0, *output.ExpiresIn)
				assert.Equal(t, &oauthRefreshToken.RevokedAT.Time, output.RevokedAT)
			})
		})
	})
}
