package formatter_test

import (
	"database/sql"
	"testing"
	"time"

	"github.com/codefluence-x/altair/provider/plugin/oauth/entity"
	"github.com/codefluence-x/altair/provider/plugin/oauth/formatter"
	"github.com/go-sql-driver/mysql"

	"github.com/codefluence-x/altair/util"
	"github.com/stretchr/testify/assert"
)

func TestModel(t *testing.T) {
	tokenExpiresIn := time.Hour * 24
	codeExpiresIn := time.Hour * 24

	t.Run("AccessTokenFromAuthorizationRequest", func(t *testing.T) {
		t.Run("Given authorization request and oauth application", func(t *testing.T) {
			t.Run("Return oauth access token insertable", func(t *testing.T) {
				authorizationRequest := entity.AuthorizationRequestJSON{
					ResourceOwnerID: util.IntToPointer(1),
					Scopes:          util.StringToPointer("users public"),
				}

				application := entity.OauthApplication{
					ID: 1,
				}

				modelFormatter := formatter.NewModel(tokenExpiresIn, codeExpiresIn)
				insertable := modelFormatter.AccessTokenFromAuthorizationRequest(authorizationRequest, application)

				assert.Equal(t, application.ID, insertable.OauthApplicationID)
				assert.Equal(t, *authorizationRequest.ResourceOwnerID, insertable.ResourceOwnerID)
				assert.Equal(t, *authorizationRequest.Scopes, insertable.Scopes)
				assert.NotEqual(t, "", insertable.Token)
				assert.NotEqual(t, time.Time{}, insertable.ExpiresIn)
			})
		})
	})

	t.Run("AccessGrantFromAuthorizationRequest", func(t *testing.T) {
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

				modelFormatter := formatter.NewModel(tokenExpiresIn, codeExpiresIn)
				insertable := modelFormatter.AccessGrantFromAuthorizationRequest(authorizationRequest, application)

				assert.Equal(t, application.ID, insertable.OauthApplicationID)
				assert.Equal(t, *authorizationRequest.ResourceOwnerID, insertable.ResourceOwnerID)
				assert.Equal(t, *authorizationRequest.Scopes, insertable.Scopes)
				assert.Equal(t, *authorizationRequest.RedirectURI, insertable.RedirectURI)
				assert.NotEqual(t, "", insertable.Code)
				assert.NotEqual(t, time.Time{}, insertable.ExpiresIn)
			})
		})
	})

	t.Run("OauthApplication", func(t *testing.T) {
		t.Run("Given authorization request and oauth application", func(t *testing.T) {
			t.Run("Return oauth access grant insertable", func(t *testing.T) {

				oauthApplicationJSON := entity.OauthApplicationJSON{
					OwnerID:     util.IntToPointer(1),
					OwnerType:   util.StringToPointer("confidential"),
					Description: util.StringToPointer("Application 1"),
					Scopes:      util.StringToPointer("public user"),
				}

				modelFormatter := formatter.NewModel(tokenExpiresIn, codeExpiresIn)
				insertable := modelFormatter.OauthApplication(oauthApplicationJSON)

				assert.Equal(t, *oauthApplicationJSON.OwnerID, insertable.OwnerID)
				assert.Equal(t, *oauthApplicationJSON.OwnerType, insertable.OwnerType)
				assert.Equal(t, *oauthApplicationJSON.Description, insertable.Description)
				assert.Equal(t, *oauthApplicationJSON.Scopes, insertable.Scopes)
				assert.NotEqual(t, "", insertable.ClientUID)
				assert.NotEqual(t, "", insertable.ClientSecret)

				// assert.Equal(t, application.ID, insertable.OauthApplicationID)
			})
		})
	})

	t.Run("AccessTokenFromOauthAccessGrant", func(t *testing.T) {
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
					RevokedAT: mysql.NullTime{
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

				modelFormatter := formatter.NewModel(tokenExpiresIn, codeExpiresIn)
				insertable := modelFormatter.AccessTokenFromOauthAccessGrant(oauthAccessGrant, application)

				assert.Equal(t, application.ID, insertable.OauthApplicationID)
				assert.Equal(t, oauthAccessGrant.ResourceOwnerID, insertable.ResourceOwnerID)
				assert.Equal(t, oauthAccessGrant.Scopes.String, insertable.Scopes)
				assert.NotEqual(t, time.Time{}, insertable.ExpiresIn)
			})
		})
	})
}
