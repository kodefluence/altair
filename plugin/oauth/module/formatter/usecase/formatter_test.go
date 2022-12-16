package usecase_test

import (
	"database/sql"
	"testing"
	"time"

	"github.com/kodefluence/altair/plugin/oauth/entity"
	"github.com/kodefluence/altair/plugin/oauth/module/formatter/usecase"
	"github.com/kodefluence/altair/util"
	"github.com/kodefluence/aurelia"
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
						RevokedAt: sql.NullTime{
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
						RevokedAt: sql.NullTime{
							Time:  time.Now(),
							Valid: true,
						},
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					},
				}

				oauthApplicationJSON := newFormatter().ApplicationList(oauthApplications)
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

				insertable := newFormatter().OauthApplicationInsertable(oauthApplicationJSON)

				assert.Equal(t, *oauthApplicationJSON.OwnerID, insertable.OwnerID)
				assert.Equal(t, *oauthApplicationJSON.OwnerType, insertable.OwnerType)
				assert.Equal(t, *oauthApplicationJSON.Description, insertable.Description)
				assert.Equal(t, *oauthApplicationJSON.Scopes, insertable.Scopes)
				assert.NotEqual(t, "", insertable.ClientUID)
				assert.NotEqual(t, "", insertable.ClientSecret)
			})
		})
	})

	t.Run("AccessTokenFromAuthorizationRequestInsertable", func(t *testing.T) {
		t.Run("Given authorization request and oauth application", func(t *testing.T) {
			t.Run("Return oauth access token insertable", func(t *testing.T) {
				authorizationRequest := entity.AuthorizationRequestJSON{
					ResourceOwnerID: util.IntToPointer(1),
					Scopes:          util.StringToPointer("users public"),
				}

				application := entity.OauthApplication{
					ID: 1,
				}

				insertable := newFormatter().AccessTokenFromAuthorizationRequestInsertable(authorizationRequest, application)

				assert.Equal(t, application.ID, insertable.OauthApplicationID)
				assert.Equal(t, *authorizationRequest.ResourceOwnerID, insertable.ResourceOwnerID)
				assert.Equal(t, *authorizationRequest.Scopes, insertable.Scopes)
				assert.NotEqual(t, "", insertable.Token)
				assert.NotEqual(t, time.Time{}, insertable.ExpiresIn)
			})
		})
	})

	t.Run("AccessTokenClientCredentialInsertable", func(t *testing.T) {
		t.Run("Given oauth application and scopes", func(t *testing.T) {
			t.Run("Return oauth access token insertable", func(t *testing.T) {
				application := entity.OauthApplication{
					ID: 1,
				}

				scopes := util.StringToPointer("publc")

				insertable := newFormatter().AccessTokenClientCredentialInsertable(application, scopes)

				assert.Equal(t, application.ID, insertable.OauthApplicationID)
				assert.Equal(t, 0, insertable.ResourceOwnerID)
				assert.Equal(t, *scopes, insertable.Scopes)
			})
		})
	})

	t.Run("AccessGrantFromAuthorizationRequestInsertable", func(t *testing.T) {
		t.Run("Given authorization request and oauth application", func(t *testing.T) {
			t.Run("Return oauth access grant insertable", func(t *testing.T) {
				authorizationRequest := entity.AuthorizationRequestJSON{
					ResourceOwnerID: util.IntToPointer(1),
					Scopes:          util.StringToPointer("users public"),
					RedirectURI:     util.StringToPointer("https://github.com"),
				}

				application := entity.OauthApplication{
					ID: 1,
				}

				insertable := newFormatter().AccessGrantFromAuthorizationRequestInsertable(authorizationRequest, application)

				assert.Equal(t, application.ID, insertable.OauthApplicationID)
				assert.Equal(t, *authorizationRequest.ResourceOwnerID, insertable.ResourceOwnerID)
				assert.Equal(t, *authorizationRequest.Scopes, insertable.Scopes)
				assert.Equal(t, *authorizationRequest.RedirectURI, insertable.RedirectURI)
				assert.NotEqual(t, "", insertable.Code)
				assert.NotEqual(t, time.Time{}, insertable.ExpiresIn)
			})
		})
	})

	t.Run("AccessTokenFromOauthAccessGrantInsertable", func(t *testing.T) {
		t.Run("Given authorization request and oauth application", func(t *testing.T) {
			t.Run("Return oauth access grant insertable", func(t *testing.T) {
				oauthAccessGrant := entity.OauthAccessGrant{
					ID:                 1,
					Code:               "authorization_code",
					CreatedAt:          time.Now(),
					ExpiresIn:          time.Now().Add(time.Hour),
					OauthApplicationID: 1,
					RedirectURI: sql.NullString{
						String: "http://localhost:8000/redirect_uri",
						Valid:  true,
					},
					ResourceOwnerID: 1,
					RevokedAT: sql.NullTime{
						Valid: false,
					},
					Scopes: sql.NullString{
						String: "user store",
						Valid:  true,
					},
				}

				application := entity.OauthApplication{
					ID: 1,
				}

				insertable := newFormatter().AccessTokenFromOauthAccessGrantInsertable(oauthAccessGrant, application)

				assert.Equal(t, application.ID, insertable.OauthApplicationID)
				assert.Equal(t, oauthAccessGrant.ResourceOwnerID, insertable.ResourceOwnerID)
				assert.Equal(t, oauthAccessGrant.Scopes.String, insertable.Scopes)
				assert.NotEqual(t, time.Time{}, insertable.ExpiresIn)
			})
		})
	})

	t.Run("AccessTokenFromOauthRefreshTokenInsertable", func(t *testing.T) {
		t.Run("Given application and access token", func(t *testing.T) {
			t.Run("Return oauth access token insertable", func(t *testing.T) {
				oauthAccessToken := entity.OauthAccessToken{
					Scopes: sql.NullString{
						String: "public",
						Valid:  true,
					},
					ResourceOwnerID: 1,
				}

				application := entity.OauthApplication{
					ID: 1,
				}

				insertable := newFormatter().AccessTokenFromOauthRefreshTokenInsertable(application, oauthAccessToken)

				assert.Equal(t, application.ID, insertable.OauthApplicationID)
				assert.Equal(t, oauthAccessToken.ResourceOwnerID, insertable.ResourceOwnerID)
				assert.Equal(t, oauthAccessToken.Scopes.String, insertable.Scopes)
				assert.NotEqual(t, time.Time{}, insertable.ExpiresIn)
			})
		})
	})

	t.Run("AccessTokenFromOauthRefreshTokenInsertable", func(t *testing.T) {
		t.Run("Given application and access token", func(t *testing.T) {
			t.Run("Return oauth access token insertable", func(t *testing.T) {
				oauthAccessToken := entity.OauthAccessToken{
					ID: 1,
					Scopes: sql.NullString{
						String: "public",
						Valid:  true,
					},
					ResourceOwnerID: 1,
				}

				application := entity.OauthApplication{
					ID: 1,
				}

				insertable := newFormatter().RefreshTokenInsertable(application, oauthAccessToken)

				assert.Equal(t, oauthAccessToken.ID, insertable.OauthAccessTokenID)
			})
		})
	})

	////

	t.Run("AccessGrant", func(t *testing.T) {
		t.Run("Given authorization request and oauth access grant data", func(t *testing.T) {
			t.Run("Code not revoked", func(t *testing.T) {
				oauthAccessGrant := entity.OauthAccessGrant{
					ID:                 1,
					OauthApplicationID: 1,
					ResourceOwnerID:    1,
					Code:               util.SHA1(),
					RevokedAT: sql.NullTime{
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

				output := newFormatter().AccessGrant(oauthAccessGrant)

				assert.Equal(t, &oauthAccessGrant.ID, output.ID)
				assert.Equal(t, &oauthAccessGrant.OauthApplicationID, output.OauthApplicationID)
				assert.Equal(t, &oauthAccessGrant.ResourceOwnerID, output.ResourceOwnerID)
				assert.Equal(t, &oauthAccessGrant.Code, output.Code)
				assert.Equal(t, &oauthAccessGrant.RedirectURI.String, output.RedirectURI)
				assert.Equal(t, &oauthAccessGrant.CreatedAt, output.CreatedAt)
				assert.LessOrEqual(t, *output.ExpiresIn, int(time.Until(oauthAccessGrant.ExpiresIn).Seconds()))
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
					RevokedAT: sql.NullTime{
						Valid: true,
						Time:  time.Now(),
					},
				}

				output := newFormatter().AccessGrant(oauthAccessGrant)

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

					output := newFormatter().AccessToken(oauthAccessToken, *authorizationReq.RedirectURI, nil)

					assert.Equal(t, &oauthAccessToken.ID, output.ID)
					assert.Equal(t, &oauthAccessToken.OauthApplicationID, output.OauthApplicationID)
					assert.Equal(t, &oauthAccessToken.ResourceOwnerID, output.ResourceOwnerID)
					assert.Equal(t, &oauthAccessToken.Token, output.Token)
					assert.Equal(t, &oauthAccessToken.CreatedAt, output.CreatedAt)
					assert.Equal(t, &oauthAccessToken.Scopes.String, output.Scopes)
					assert.LessOrEqual(t, *output.ExpiresIn, int(time.Until(oauthAccessToken.ExpiresIn).Seconds()))
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

					output := newFormatter().AccessToken(oauthAccessToken, *authorizationReq.RedirectURI, nil)

					assert.Equal(t, &oauthAccessToken.ID, output.ID)
					assert.Equal(t, &oauthAccessToken.OauthApplicationID, output.OauthApplicationID)
					assert.Equal(t, &oauthAccessToken.ResourceOwnerID, output.ResourceOwnerID)
					assert.Equal(t, &oauthAccessToken.Token, output.Token)
					assert.Equal(t, &oauthAccessToken.CreatedAt, output.CreatedAt)
					assert.Equal(t, &oauthAccessToken.Scopes.String, output.Scopes)
					assert.LessOrEqual(t, *output.ExpiresIn, int(time.Until(oauthAccessToken.ExpiresIn).Seconds()))
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
					RevokedAT: sql.NullTime{
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
					RevokedAT: sql.NullTime{
						Valid: false,
					},
				}

				oauthRefreshTokenJSON := newFormatter().RefreshToken(oauthRefreshToken)

				output := newFormatter().AccessToken(oauthAccessToken, *authorizationReq.RedirectURI, &oauthRefreshTokenJSON)

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
						RevokedAT: sql.NullTime{
							Valid: false,
						},
					}

					output := newFormatter().RefreshToken(oauthRefreshToken)

					assert.Equal(t, &oauthRefreshToken.CreatedAt, output.CreatedAt)
					assert.Equal(t, &oauthRefreshToken.Token, output.Token)
					assert.LessOrEqual(t, *output.ExpiresIn, int(time.Until(oauthRefreshToken.ExpiresIn).Seconds()))
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
					RevokedAT: sql.NullTime{
						Time:  time.Now().Add(-time.Hour),
						Valid: true,
					},
				}

				output := newFormatter().RefreshToken(oauthRefreshToken)

				assert.Equal(t, &oauthRefreshToken.CreatedAt, output.CreatedAt)
				assert.Equal(t, &oauthRefreshToken.Token, output.Token)
				assert.Equal(t, 0, *output.ExpiresIn)
				assert.Equal(t, &oauthRefreshToken.RevokedAT.Time, output.RevokedAT)
			})
		})
	})
}

func newFormatter() *usecase.Formatter {
	return usecase.NewFormatter(
		time.Hour, time.Hour, time.Hour,
	)
}
