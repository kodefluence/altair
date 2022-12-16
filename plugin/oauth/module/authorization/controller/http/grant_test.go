package http_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/kodefluence/altair/module/apierror"
	"github.com/kodefluence/altair/module/controller"
	"github.com/kodefluence/altair/plugin/oauth/entity"
	authorizationHttp "github.com/kodefluence/altair/plugin/oauth/module/authorization/controller/http"
	"github.com/kodefluence/altair/plugin/oauth/module/authorization/controller/http/mock"
	"github.com/kodefluence/altair/plugin/oauth/module/formatter"
	"github.com/kodefluence/altair/testhelper"
	"github.com/kodefluence/altair/util"
	"github.com/kodefluence/aurelia"
	"github.com/kodefluence/monorepo/jsonapi"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestGrant(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	t.Run("Method", func(t *testing.T) {
		authorizationService := mock.NewMockAuthorization(mockCtrl)
		assert.Equal(t, "POST", authorizationHttp.NewGrant(authorizationService, apierror.Provide()).Method())
	})

	t.Run("Path", func(t *testing.T) {
		authorizationService := mock.NewMockAuthorization(mockCtrl)
		assert.Equal(t, "/oauth/authorizations", authorizationHttp.NewGrant(authorizationService, apierror.Provide()).Path())
	})

	t.Run("Control", func(t *testing.T) {
		t.Run("Given request with json body", func(t *testing.T) {
			t.Run("Return oauth application data with status 202", func(t *testing.T) {
				apiEngine := gin.Default()

				authorizationRequest := entity.AuthorizationRequestJSON{
					ResponseType:    nil,
					ResourceOwnerID: util.ValueToPointer(1),
					ClientUID:       util.ValueToPointer(aurelia.Hash("x", "y")),
					ClientSecret:    util.ValueToPointer(aurelia.Hash("z", "a")),
					RedirectURI:     util.ValueToPointer("http://github.com"),
					Scopes:          util.ValueToPointer("public users"),
				}
				encodedBytes, err := json.Marshal(authorizationRequest)
				assert.Nil(t, err)

				oauthAccessToken := entity.OauthAccessToken{
					ID:                 1,
					OauthApplicationID: 1,
					ResourceOwnerID:    *authorizationRequest.ResourceOwnerID,
					Token:              aurelia.Hash("x", "y"),
					Scopes: sql.NullString{
						String: *authorizationRequest.Scopes,
						Valid:  true,
					},
					ExpiresIn: time.Now().Add(time.Hour * 4),
					CreatedAt: time.Now().Truncate(time.Second),
				}
				oauthAccessTokenJSON := formatter.Provide(time.Hour, time.Hour, time.Hour).AccessToken(oauthAccessToken, *authorizationRequest.RedirectURI, nil)
				expectedReponseByte, _ := json.Marshal(jsonapi.BuildResponse(jsonapi.WithData(oauthAccessTokenJSON)))

				authorizationService := mock.NewMockAuthorization(mockCtrl)
				authorizationService.EXPECT().GrantAuthorizationCode(gomock.Any(), authorizationRequest).Return(oauthAccessTokenJSON, nil)

				ctrl := authorizationHttp.NewGrant(authorizationService, apierror.Provide())
				controller.Provide(apiEngine.Handle, apierror.Provide(), &cobra.Command{}).InjectHTTP(ctrl)

				w := testhelper.PerformRequest(apiEngine, ctrl.Method(), ctrl.Path(), bytes.NewReader(encodedBytes))

				responseByte, err := io.ReadAll(w.Body)
				assert.Nil(t, err)
				assert.Equal(t, http.StatusOK, w.Code)
				assert.Equal(t, string(expectedReponseByte), string(responseByte))
			})

			t.Run("Unexpected error in authorization services", func(t *testing.T) {
				t.Run("Return entity error status", func(t *testing.T) {
					apiEngine := gin.Default()

					authorizationRequest := entity.AuthorizationRequestJSON{
						ResponseType:    nil,
						ResourceOwnerID: util.ValueToPointer(1),
						ClientUID:       util.ValueToPointer(aurelia.Hash("x", "y")),
						ClientSecret:    util.ValueToPointer(aurelia.Hash("z", "a")),
						RedirectURI:     util.ValueToPointer("http://github.com"),
						Scopes:          util.ValueToPointer("public users"),
					}
					encodedBytes, err := json.Marshal(authorizationRequest)
					assert.Nil(t, err)

					authorizationService := mock.NewMockAuthorization(mockCtrl)
					authorizationService.EXPECT().GrantAuthorizationCode(gomock.Any(), authorizationRequest).Return(entity.OauthAccessTokenJSON{}, testhelper.ErrInternalServer())

					ctrl := authorizationHttp.NewGrant(authorizationService, apierror.Provide())
					controller.Provide(apiEngine.Handle, apierror.Provide(), &cobra.Command{}).InjectHTTP(ctrl)

					w := testhelper.PerformRequest(apiEngine, ctrl.Method(), ctrl.Path(), bytes.NewReader(encodedBytes))
					assert.Equal(t, http.StatusInternalServerError, w.Code)
				})
			})
		})

		t.Run("Given invalid request body", func(t *testing.T) {
			t.Run("Return bad request", func(t *testing.T) {
				apiEngine := gin.Default()

				authorizationService := mock.NewMockAuthorization(mockCtrl)

				ctrl := authorizationHttp.NewGrant(authorizationService, apierror.Provide())
				controller.Provide(apiEngine.Handle, apierror.Provide(), &cobra.Command{}).InjectHTTP(ctrl)

				w := testhelper.PerformRequest(apiEngine, ctrl.Method(), ctrl.Path(), testhelper.MockErrorIoReader{})
				responseByte, err := io.ReadAll(w.Body)
				assert.Nil(t, err)
				assert.Equal(t, http.StatusBadRequest, w.Code)
				assert.Equal(t, "{\"errors\":[{\"title\":\"Bad request error\",\"detail\":\"You've send malformed request in your `request body`\",\"code\":\"ERR0400\",\"status\":400}]}", string(responseByte))
			})
		})

		t.Run("Given request body but not json", func(t *testing.T) {
			t.Run("Return bad request", func(t *testing.T) {
				apiEngine := gin.Default()

				authorizationService := mock.NewMockAuthorization(mockCtrl)

				ctrl := authorizationHttp.NewGrant(authorizationService, apierror.Provide())
				controller.Provide(apiEngine.Handle, apierror.Provide(), &cobra.Command{}).InjectHTTP(ctrl)

				w := testhelper.PerformRequest(apiEngine, ctrl.Method(), ctrl.Path(), bytes.NewReader([]byte(`this is gonna be error`)))
				responseByte, err := io.ReadAll(w.Body)
				assert.Nil(t, err)
				assert.Equal(t, http.StatusBadRequest, w.Code)
				assert.Equal(t, "{\"errors\":[{\"title\":\"Bad request error\",\"detail\":\"You've send malformed request in your `request body`\",\"code\":\"ERR0400\",\"status\":400}]}", string(responseByte))
			})
		})
	})
}
