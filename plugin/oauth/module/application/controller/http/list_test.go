package http_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/kodefluence/altair/module/apierror"
	"github.com/kodefluence/altair/module/controller"
	"github.com/kodefluence/altair/plugin/oauth/entity"
	applicationHttp "github.com/kodefluence/altair/plugin/oauth/module/application/controller/http"
	"github.com/kodefluence/altair/plugin/oauth/module/application/controller/http/mock"
	"github.com/kodefluence/altair/testhelper"
	"github.com/kodefluence/altair/util"
	"github.com/kodefluence/monorepo/jsonapi"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

type responseList struct {
	Data []entity.OauthApplicationJSON `json:"data"`
	Meta jsonapi.Meta                  `json:"meta"`
}

func TestList(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	apierror := apierror.Provide()

	t.Run("Method", func(t *testing.T) {
		applicationManager := mock.NewMockApplicationManager(mockCtrl)
		assert.Equal(t, "GET", applicationHttp.NewList(applicationManager, apierror).Method())
	})

	t.Run("Path", func(t *testing.T) {
		applicationManager := mock.NewMockApplicationManager(mockCtrl)
		assert.Equal(t, "/oauth/applications", applicationHttp.NewList(applicationManager, apierror).Path())
	})

	t.Run("Control", func(t *testing.T) {
		t.Run("Given request with offset and limit", func(t *testing.T) {
			t.Run("Return list of oauth application", func(t *testing.T) {
				apiEngine := gin.Default()

				oauthApplicationJSONs := []entity.OauthApplicationJSON{
					{
						ID:           util.IntToPointer(1),
						OwnerID:      util.IntToPointer(1),
						Description:  util.StringToPointer("Application 1"),
						Scopes:       util.StringToPointer("public user"),
						ClientUID:    util.StringToPointer("clientuid01"),
						ClientSecret: util.StringToPointer("clientsecret01"),
					},
					{
						ID:           util.IntToPointer(2),
						OwnerID:      util.IntToPointer(2),
						Description:  util.StringToPointer("Application 2"),
						Scopes:       util.StringToPointer("public user"),
						ClientUID:    util.StringToPointer("clientuid02"),
						ClientSecret: util.StringToPointer("clientsecret02"),
					},
				}

				applicationManager := mock.NewMockApplicationManager(mockCtrl)
				applicationManager.EXPECT().List(gomock.Any(), 0, 10).Return(oauthApplicationJSONs, 10, nil)

				ctrl := applicationHttp.NewList(applicationManager, apierror)
				controller.Provide(apiEngine.Handle, apierror, &cobra.Command{}).InjectHTTP(ctrl)

				var response responseList
				w := testhelper.PerformRequest(apiEngine, ctrl.Method(), ctrl.Path(), nil)

				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.Nil(t, err)

				assert.Equal(t, http.StatusOK, w.Code)
				assert.Equal(t, oauthApplicationJSONs, response.Data)
			})

			t.Run("Application manager error", func(t *testing.T) {
				t.Run("Return internal server error", func(t *testing.T) {
					apiEngine := gin.Default()

					applicationManager := mock.NewMockApplicationManager(mockCtrl)
					applicationManager.EXPECT().List(gomock.Any(), 0, 10).Return([]entity.OauthApplicationJSON(nil), 0, testhelper.ErrInternalServer())

					ctrl := applicationHttp.NewList(applicationManager, apierror)
					controller.Provide(apiEngine.Handle, apierror, &cobra.Command{}).InjectHTTP(ctrl)

					w := testhelper.PerformRequest(apiEngine, ctrl.Method(), ctrl.Path(), nil)

					assert.Equal(t, http.StatusInternalServerError, w.Code)
				})
			})
		})

		t.Run("Given request with invalid offset", func(t *testing.T) {
			t.Run("Return bad request error", func(t *testing.T) {
				apiEngine := gin.Default()

				applicationManager := mock.NewMockApplicationManager(mockCtrl)

				ctrl := applicationHttp.NewList(applicationManager, apierror)
				controller.Provide(apiEngine.Handle, apierror, &cobra.Command{}).InjectHTTP(ctrl)

				w := testhelper.PerformRequest(apiEngine, ctrl.Method(), ctrl.Path()+"?offset=invalid", nil)

				assert.Equal(t, http.StatusBadRequest, w.Code)
				assert.Equal(t, "{\"errors\":[{\"title\":\"Bad request error\",\"detail\":\"You've send malformed request in your `query parameters: offset`\",\"code\":\"ERR0400\",\"status\":400}]}", string(w.Body.Bytes()))

			})
		})

		t.Run("Given request with invalid limit", func(t *testing.T) {
			t.Run("Return bad request error", func(t *testing.T) {
				apiEngine := gin.Default()

				applicationManager := mock.NewMockApplicationManager(mockCtrl)

				ctrl := applicationHttp.NewList(applicationManager, apierror)
				controller.Provide(apiEngine.Handle, apierror, &cobra.Command{}).InjectHTTP(ctrl)

				w := testhelper.PerformRequest(apiEngine, ctrl.Method(), ctrl.Path()+"?limit=invalid", nil)

				assert.Equal(t, http.StatusBadRequest, w.Code)
				assert.Equal(t, "{\"errors\":[{\"title\":\"Bad request error\",\"detail\":\"You've send malformed request in your `query parameters: limit`\",\"code\":\"ERR0400\",\"status\":400}]}", string(w.Body.Bytes()))
			})
		})
	})
}
