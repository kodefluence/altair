package formatter_test

import (
	"testing"
	"time"

	"github.com/codefluence-x/altair/entity"
	"github.com/codefluence-x/altair/formatter"
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

				modelFormatter := formatter.Model(tokenExpiresIn, codeExpiresIn)
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

				modelFormatter := formatter.Model(tokenExpiresIn, codeExpiresIn)
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
}
