package downstream_test

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/codefluence-x/altair/entity"
	"github.com/codefluence-x/altair/mock"
	"github.com/codefluence-x/altair/plugin"
	"github.com/codefluence-x/aurelia"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestOauth(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	t.Run("Intervene", func(t *testing.T) {
		t.Run("Given gin.Context and http.Request", func(t *testing.T) {
			t.Run("Normal scenario", func(t *testing.T) {
				t.Run("Return nil", func(t *testing.T) {
					token := "token"

					c := &gin.Context{}
					c.Request = &http.Request{
						Header: http.Header{},
					}
					c.Request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

					r, _ := http.NewRequest("GET", "https://github.com/codefluence-x/altair", nil)

					entityAccessToken := entity.OauthAccessToken{
						ID:                 1,
						OauthApplicationID: 1,
						ResourceOwnerID:    1,
						Token:              aurelia.Hash("x", "y"),
						Scopes: sql.NullString{
							String: "public user",
							Valid:  true,
						},
						ExpiresIn: time.Now().Add(time.Hour * 4),
						CreatedAt: time.Now(),
					}

					oauthAccessTokenModel := mock.NewMockOauthAccessTokenModel(mockCtrl)
					oauthAccessTokenModel.EXPECT().OneByToken(c, token).Return(entityAccessToken, nil)

					oauthPlugin := plugin.DownStream().Oauth(oauthAccessTokenModel)

					err := oauthPlugin.Intervene(c, r)

					assert.Nil(t, err)
					assert.Equal(t, strconv.Itoa(entityAccessToken.ResourceOwnerID), r.Header.Get("Resource-Owner-ID"))
					assert.Equal(t, strconv.Itoa(entityAccessToken.OauthApplicationID), r.Header.Get("Oauth-Application-ID"))
				})
			})

			t.Run("Parse token error", func(t *testing.T) {
				t.Run("Token is not provided", func(t *testing.T) {
					t.Run("Return error", func(t *testing.T) {

					})
				})

				t.Run("Token format is invalid", func(t *testing.T) {
					t.Run("Return error", func(t *testing.T) {

					})
				})
			})

			t.Run("Oauth token not found", func(t *testing.T) {
				t.Run("Return nil with status 401", func(t *testing.T) {

				})
			})

			t.Run("Oauth token model error", func(t *testing.T) {
				t.Run("Return nil with status 503", func(t *testing.T) {

				})
			})
		})
	})
}
