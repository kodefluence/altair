package http_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/kodefluence/altair/module/apierror"
	"github.com/kodefluence/altair/module/controller"
	"github.com/kodefluence/altair/plugin/oauth/entity"
	authorizationHttp "github.com/kodefluence/altair/plugin/oauth/module/authorization/controller/http"
	"github.com/kodefluence/altair/plugin/oauth/module/authorization/controller/http/mock"
	"github.com/kodefluence/altair/testhelper"
	"github.com/kodefluence/altair/util"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestRevoke(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	t.Run("Method", func(t *testing.T) {
		authorizationService := mock.NewMockAuthorization(mockCtrl)
		assert.Equal(t, "POST", authorizationHttp.NewRevoke(authorizationService, apierror.Provide()).Method())
	})

	t.Run("Path", func(t *testing.T) {
		authorizationService := mock.NewMockAuthorization(mockCtrl)
		assert.Equal(t, "/oauth/authorizations/revoke", authorizationHttp.NewRevoke(authorizationService, apierror.Provide()).Path())
	})

	t.Run("Control", func(t *testing.T) {
		t.Run("Given request with json body", func(t *testing.T) {
			t.Run("Return message with status 200", func(t *testing.T) {
				apiEngine := gin.Default()

				revokeTokenRequest := entity.RevokeAccessTokenRequestJSON{
					Token: util.StringToPointer("some-cool-token"),
				}
				encodedBytes, err := json.Marshal(revokeTokenRequest)
				assert.Nil(t, err)

				authorizationService := mock.NewMockAuthorization(mockCtrl)
				authorizationService.EXPECT().RevokeToken(gomock.Any(), revokeTokenRequest).Return(nil)

				ctrl := authorizationHttp.NewRevoke(authorizationService, apierror.Provide())
				controller.Provide(apiEngine.Handle, apierror.Provide(), &cobra.Command{}).InjectHTTP(ctrl)

				w := testhelper.PerformRequest(apiEngine, ctrl.Method(), ctrl.Path(), bytes.NewReader(encodedBytes))
				responseByte, err := io.ReadAll(w.Body)
				assert.Nil(t, err)
				assert.Equal(t, http.StatusOK, w.Code)
				assert.Equal(t, "{}", string(responseByte))

			})

			t.Run("Unexpected error in authorization services", func(t *testing.T) {
				t.Run("Return entity error status", func(t *testing.T) {
					apiEngine := gin.Default()

					revokeTokenRequest := entity.RevokeAccessTokenRequestJSON{
						Token: nil,
					}
					encodedBytes, err := json.Marshal(revokeTokenRequest)
					assert.Nil(t, err)

					authorizationService := mock.NewMockAuthorization(mockCtrl)
					authorizationService.EXPECT().RevokeToken(gomock.Any(), revokeTokenRequest).Return(testhelper.ErrInternalServer())

					ctrl := authorizationHttp.NewRevoke(authorizationService, apierror.Provide())
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

				ctrl := authorizationHttp.NewRevoke(authorizationService, apierror.Provide())
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

				ctrl := authorizationHttp.NewRevoke(authorizationService, apierror.Provide())
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
