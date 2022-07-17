package downstream_test

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"testing"
	"time"

	coreEntity "github.com/kodefluence/altair/entity"
	coreMock "github.com/kodefluence/altair/mock"
	mockdb "github.com/kodefluence/monorepo/db/mock"
	"github.com/kodefluence/monorepo/exception"

	"github.com/kodefluence/altair/provider/plugin/oauth/downstream"
	"github.com/kodefluence/altair/provider/plugin/oauth/entity"
	"github.com/kodefluence/altair/provider/plugin/oauth/mock"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/kodefluence/aurelia"
	"github.com/stretchr/testify/assert"
)

func TestOauth(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	sqldb := mockdb.NewMockDB(mockCtrl)

	t.Run("Name", func(t *testing.T) {
		t.Run("Return oauth-plugin", func(t *testing.T) {
			oauthAccessTokenModel := mock.NewMockOauthAccessTokenModel(mockCtrl)
			oauthPlugin := downstream.NewOauth(oauthAccessTokenModel, sqldb)
			assert.Equal(t, "oauth-plugin", oauthPlugin.Name())
		})
	})

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

					r, _ := http.NewRequest("GET", "https://github.com/kodefluence/altair", nil)

					routePath := coreEntity.RouterPath{Auth: "oauth"}

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
					oauthAccessTokenModel.EXPECT().OneByToken(gomock.Any(), token, sqldb).Return(entityAccessToken, nil)

					oauthPlugin := downstream.NewOauth(oauthAccessTokenModel, sqldb)

					err := oauthPlugin.Intervene(c, r, routePath)

					assert.Nil(t, err)
					assert.Equal(t, strconv.Itoa(entityAccessToken.ResourceOwnerID), r.Header.Get("Resource-Owner-ID"))
					assert.Equal(t, strconv.Itoa(entityAccessToken.OauthApplicationID), r.Header.Get("Oauth-Application-ID"))
				})
			})

			t.Run("Normal scenario with valid route scope", func(t *testing.T) {
				t.Run("Return nil", func(t *testing.T) {
					token := "token"

					c := &gin.Context{}
					c.Request = &http.Request{
						Header: http.Header{},
					}
					c.Request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

					r, _ := http.NewRequest("GET", "https://github.com/kodefluence/altair", nil)

					routePath := coreEntity.RouterPath{Auth: "oauth", Scope: "public"}

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
					oauthAccessTokenModel.EXPECT().OneByToken(gomock.Any(), token, sqldb).Return(entityAccessToken, nil)

					oauthPlugin := downstream.NewOauth(oauthAccessTokenModel, sqldb)

					err := oauthPlugin.Intervene(c, r, routePath)

					assert.Nil(t, err)
					assert.Equal(t, strconv.Itoa(entityAccessToken.ResourceOwnerID), r.Header.Get("Resource-Owner-ID"))
					assert.Equal(t, strconv.Itoa(entityAccessToken.OauthApplicationID), r.Header.Get("Oauth-Application-ID"))
				})
			})

			t.Run("Auth is not oauth type", func(t *testing.T) {
				t.Run("Return nil", func(t *testing.T) {
					token := "token"

					c := &gin.Context{}
					c.Request = &http.Request{
						Header: http.Header{},
					}
					c.Request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

					r, _ := http.NewRequest("GET", "https://github.com/kodefluence/altair", nil)

					routePath := coreEntity.RouterPath{Auth: "none"}

					oauthAccessTokenModel := mock.NewMockOauthAccessTokenModel(mockCtrl)
					oauthAccessTokenModel.EXPECT().OneByToken(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

					oauthPlugin := downstream.NewOauth(oauthAccessTokenModel, sqldb)

					err := oauthPlugin.Intervene(c, r, routePath)

					assert.Nil(t, err)
				})
			})

			t.Run("Parse token error", func(t *testing.T) {
				t.Run("Token is not provided", func(t *testing.T) {
					t.Run("Return error", func(t *testing.T) {
						token := "invalid-token"

						c := &gin.Context{}
						c.Request = &http.Request{
							Header: http.Header{},
						}
						c.Request.Header.Add("Authorization", fmt.Sprintf("%s", token))

						responseWritterMock := coreMock.NewMockResponseWriter(mockCtrl)
						responseWritterMock.EXPECT().WriteHeaderNow().AnyTimes()
						responseWritterMock.EXPECT().WriteHeader(gomock.Any()).AnyTimes()
						responseWritterMock.EXPECT().Status().Return(http.StatusUnauthorized).AnyTimes()

						c.Writer = responseWritterMock

						r, _ := http.NewRequest("GET", "https://github.com/kodefluence/altair", nil)

						routePath := coreEntity.RouterPath{Auth: "oauth"}

						oauthAccessTokenModel := mock.NewMockOauthAccessTokenModel(mockCtrl)
						oauthAccessTokenModel.EXPECT().OneByToken(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

						oauthPlugin := downstream.NewOauth(oauthAccessTokenModel, sqldb)

						err := oauthPlugin.Intervene(c, r, routePath)

						assert.NotNil(t, err)
					})
				})

				t.Run("Token format is invalid", func(t *testing.T) {
					t.Run("Return error", func(t *testing.T) {
						token := "invalid token"

						c := &gin.Context{}
						c.Request = &http.Request{
							Header: http.Header{},
						}
						c.Request.Header.Add("Authorization", fmt.Sprintf("NotBearer %s", token))

						responseWritterMock := coreMock.NewMockResponseWriter(mockCtrl)
						responseWritterMock.EXPECT().WriteHeaderNow().AnyTimes()
						responseWritterMock.EXPECT().WriteHeader(gomock.Any()).AnyTimes()
						responseWritterMock.EXPECT().Status().Return(http.StatusUnauthorized).AnyTimes()

						c.Writer = responseWritterMock

						r, _ := http.NewRequest("GET", "https://github.com/kodefluence/altair", nil)

						routePath := coreEntity.RouterPath{Auth: "oauth"}

						oauthAccessTokenModel := mock.NewMockOauthAccessTokenModel(mockCtrl)
						oauthAccessTokenModel.EXPECT().OneByToken(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

						oauthPlugin := downstream.NewOauth(oauthAccessTokenModel, sqldb)

						err := oauthPlugin.Intervene(c, r, routePath)

						assert.NotNil(t, err)
					})
				})
			})

			t.Run("Oauth token not found", func(t *testing.T) {
				t.Run("Return nil with status 401", func(t *testing.T) {
					token := "token"

					c := &gin.Context{}
					c.Request = &http.Request{
						Header: http.Header{},
					}
					c.Request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

					responseWritterMock := coreMock.NewMockResponseWriter(mockCtrl)
					responseWritterMock.EXPECT().WriteHeaderNow().AnyTimes()
					responseWritterMock.EXPECT().WriteHeader(gomock.Any()).AnyTimes()
					responseWritterMock.EXPECT().Status().Return(http.StatusUnauthorized).AnyTimes()

					c.Writer = responseWritterMock

					r, _ := http.NewRequest("GET", "https://github.com/kodefluence/altair", nil)

					routePath := coreEntity.RouterPath{Auth: "oauth"}

					oauthAccessTokenModel := mock.NewMockOauthAccessTokenModel(mockCtrl)
					oauthAccessTokenModel.EXPECT().OneByToken(gomock.Any(), token, sqldb).Return(entity.OauthAccessToken{}, exception.Throw(sql.ErrNoRows, exception.WithType(exception.NotFound)))

					oauthPlugin := downstream.NewOauth(oauthAccessTokenModel, sqldb)

					err := oauthPlugin.Intervene(c, r, routePath)

					assert.NotNil(t, err)
					assert.Equal(t, http.StatusUnauthorized, c.Writer.Status())
				})
			})

			t.Run("Oauth token model error", func(t *testing.T) {
				t.Run("Return nil with status 503", func(t *testing.T) {
					token := "token"

					c := &gin.Context{}
					c.Request = &http.Request{
						Header: http.Header{},
					}
					c.Request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

					responseWritterMock := coreMock.NewMockResponseWriter(mockCtrl)
					responseWritterMock.EXPECT().WriteHeaderNow().AnyTimes()
					responseWritterMock.EXPECT().WriteHeader(gomock.Any()).AnyTimes()
					responseWritterMock.EXPECT().Status().Return(http.StatusServiceUnavailable).AnyTimes()

					c.Writer = responseWritterMock

					r, _ := http.NewRequest("GET", "https://github.com/kodefluence/altair", nil)

					routePath := coreEntity.RouterPath{Auth: "oauth"}

					oauthAccessTokenModel := mock.NewMockOauthAccessTokenModel(mockCtrl)
					oauthAccessTokenModel.EXPECT().OneByToken(gomock.Any(), token, sqldb).Return(entity.OauthAccessToken{}, exception.Throw(errors.New("unexpected error")))

					oauthPlugin := downstream.NewOauth(oauthAccessTokenModel, sqldb)

					err := oauthPlugin.Intervene(c, r, routePath)

					assert.NotNil(t, err)
					assert.Equal(t, http.StatusServiceUnavailable, c.Writer.Status())
				})
			})

			t.Run("Route scope is invalid", func(t *testing.T) {
				t.Run("Return error", func(t *testing.T) {
					token := "token"

					c := &gin.Context{}
					c.Request = &http.Request{
						Header: http.Header{},
					}
					c.Request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

					responseWritterMock := coreMock.NewMockResponseWriter(mockCtrl)
					responseWritterMock.EXPECT().WriteHeaderNow().AnyTimes()
					responseWritterMock.EXPECT().WriteHeader(gomock.Any()).AnyTimes()
					responseWritterMock.EXPECT().Status().Return(http.StatusForbidden).AnyTimes()

					c.Writer = responseWritterMock

					r, _ := http.NewRequest("GET", "https://github.com/kodefluence/altair", nil)

					routePath := coreEntity.RouterPath{Auth: "oauth", Scope: "user"}

					entityAccessToken := entity.OauthAccessToken{
						ID:                 1,
						OauthApplicationID: 1,
						ResourceOwnerID:    1,
						Token:              aurelia.Hash("x", "y"),
						Scopes: sql.NullString{
							String: "public",
							Valid:  true,
						},
						ExpiresIn: time.Now().Add(time.Hour * 4),
						CreatedAt: time.Now(),
					}

					oauthAccessTokenModel := mock.NewMockOauthAccessTokenModel(mockCtrl)
					oauthAccessTokenModel.EXPECT().OneByToken(gomock.Any(), token, sqldb).Return(entityAccessToken, nil)

					oauthPlugin := downstream.NewOauth(oauthAccessTokenModel, sqldb)

					err := oauthPlugin.Intervene(c, r, routePath)

					assert.NotNil(t, err)
					assert.Equal(t, http.StatusForbidden, c.Writer.Status())
				})
			})

			t.Run("Access token is already expired", func(t *testing.T) {
				t.Run("Return error", func(t *testing.T) {
					token := "token"

					c := &gin.Context{}
					c.Request = &http.Request{
						Header: http.Header{},
					}
					c.Request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

					responseWritterMock := coreMock.NewMockResponseWriter(mockCtrl)
					responseWritterMock.EXPECT().WriteHeaderNow().AnyTimes()
					responseWritterMock.EXPECT().WriteHeader(gomock.Any()).AnyTimes()
					responseWritterMock.EXPECT().Status().Return(http.StatusUnauthorized).AnyTimes()

					c.Writer = responseWritterMock

					r, _ := http.NewRequest("GET", "https://github.com/kodefluence/altair", nil)

					routePath := coreEntity.RouterPath{Auth: "oauth", Scope: "user"}

					entityAccessToken := entity.OauthAccessToken{
						ID:                 1,
						OauthApplicationID: 1,
						ResourceOwnerID:    1,
						Token:              aurelia.Hash("x", "y"),
						Scopes: sql.NullString{
							String: "user",
							Valid:  true,
						},
						ExpiresIn: time.Now().Add(-time.Hour * 4),
						CreatedAt: time.Now(),
					}

					oauthAccessTokenModel := mock.NewMockOauthAccessTokenModel(mockCtrl)
					oauthAccessTokenModel.EXPECT().OneByToken(gomock.Any(), token, sqldb).Return(entityAccessToken, nil)

					oauthPlugin := downstream.NewOauth(oauthAccessTokenModel, sqldb)

					err := oauthPlugin.Intervene(c, r, routePath)

					assert.NotNil(t, err)
					assert.Equal(t, http.StatusUnauthorized, c.Writer.Status())
				})
			})
		})
	})
}
