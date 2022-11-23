package http_test

import (
	"bytes"
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
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestUpdate(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	apierror := apierror.Provide()

	t.Run("Method", func(t *testing.T) {
		applicationManager := mock.NewMockApplicationManager(mockCtrl)
		assert.Equal(t, "PUT", applicationHttp.NewUpdate(applicationManager, apierror).Method())
	})

	t.Run("Path", func(t *testing.T) {
		applicationManager := mock.NewMockApplicationManager(mockCtrl)
		assert.Equal(t, "/oauth/applications/:id", applicationHttp.NewUpdate(applicationManager, apierror).Path())
	})

	t.Run("Control", func(t *testing.T) {
		t.Run("Given request with json body", func(t *testing.T) {
			t.Run("Return oauth application data with status 200", func(t *testing.T) {
				apiEngine := gin.Default()

				oauthApplicationUpdateJSON := entity.OauthApplicationUpdateJSON{
					Description: util.StringToPointer("Application 1"),
					Scopes:      util.StringToPointer("public user"),
				}
				encodedBytes, err := json.Marshal(oauthApplicationUpdateJSON)
				assert.Nil(t, err)

				oauthApplicationJSON := entity.OauthApplicationJSON{
					OwnerID:      util.IntToPointer(1),
					Description:  util.StringToPointer("Application 1"),
					Scopes:       util.StringToPointer("public user"),
					ClientUID:    util.StringToPointer("clientuid01"),
					ClientSecret: util.StringToPointer("clientsecret01"),
				}

				applicationManager := mock.NewMockApplicationManager(mockCtrl)
				applicationManager.EXPECT().Update(gomock.Any(), 1, oauthApplicationUpdateJSON).Return(oauthApplicationJSON, nil)

				ctrl := applicationHttp.NewUpdate(applicationManager, apierror)
				controller.Provide(apiEngine.Handle, apierror, &cobra.Command{}).InjectHTTP(ctrl)

				var response responseOne
				w := testhelper.PerformRequest(apiEngine, ctrl.Method(), "/oauth/applications/1", bytes.NewReader(encodedBytes))

				err = json.Unmarshal(w.Body.Bytes(), &response)
				assert.Nil(t, err)

				assert.Equal(t, http.StatusOK, w.Code)
				assert.Equal(t, oauthApplicationJSON, response.Data)
			})

			t.Run("Unexpected error in application manager", func(t *testing.T) {
				apiEngine := gin.Default()

				oauthApplicationUpdateJSON := entity.OauthApplicationUpdateJSON{
					Description: util.StringToPointer("Application 1"),
					Scopes:      util.StringToPointer("public user"),
				}
				encodedBytes, err := json.Marshal(oauthApplicationUpdateJSON)
				assert.Nil(t, err)

				oauthApplicationJSON := entity.OauthApplicationJSON{}

				applicationManager := mock.NewMockApplicationManager(mockCtrl)
				applicationManager.EXPECT().Update(gomock.Any(), 1, oauthApplicationUpdateJSON).Return(oauthApplicationJSON, testhelper.ErrInternalServer())

				ctrl := applicationHttp.NewUpdate(applicationManager, apierror)
				controller.Provide(apiEngine.Handle, apierror, &cobra.Command{}).InjectHTTP(ctrl)

				var response responseOne
				w := testhelper.PerformRequest(apiEngine, ctrl.Method(), "/oauth/applications/1", bytes.NewReader(encodedBytes))

				err = json.Unmarshal(w.Body.Bytes(), &response)
				assert.Nil(t, err)

				assert.Equal(t, http.StatusInternalServerError, w.Code)
			})
		})

		t.Run("Given invalid request body", func(t *testing.T) {
			t.Run("Return bad request", func(t *testing.T) {
				apiEngine := gin.Default()

				applicationManager := mock.NewMockApplicationManager(mockCtrl)
				applicationManager.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

				ctrl := applicationHttp.NewUpdate(applicationManager, apierror)
				controller.Provide(apiEngine.Handle, apierror, &cobra.Command{}).InjectHTTP(ctrl)

				w := testhelper.PerformRequest(apiEngine, ctrl.Method(), "/oauth/applications/1", testhelper.MockErrorIoReader{})

				assert.Equal(t, http.StatusBadRequest, w.Code)
				assert.Equal(t, "{\"errors\":[{\"title\":\"Bad request error\",\"detail\":\"You've send malformed request in your `request body`\",\"code\":\"ERR0400\",\"status\":400}]}", string(w.Body.Bytes()))
			})

			t.Run("Given invalid url params", func(t *testing.T) {
				t.Run("Return bad request", func(t *testing.T) {

					oauthApplicationUpdateJSON := entity.OauthApplicationUpdateJSON{
						Description: util.StringToPointer("Application 1"),
						Scopes:      util.StringToPointer("public user"),
					}
					encodedBytes, err := json.Marshal(oauthApplicationUpdateJSON)
					assert.Nil(t, err)

					apiEngine := gin.Default()

					applicationManager := mock.NewMockApplicationManager(mockCtrl)
					applicationManager.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

					ctrl := applicationHttp.NewUpdate(applicationManager, apierror)
					controller.Provide(apiEngine.Handle, apierror, &cobra.Command{}).InjectHTTP(ctrl)

					w := testhelper.PerformRequest(apiEngine, ctrl.Method(), "/oauth/applications/s", bytes.NewReader(encodedBytes))

					assert.Equal(t, http.StatusBadRequest, w.Code)
					assert.Equal(t, "{\"errors\":[{\"title\":\"Bad request error\",\"detail\":\"You've send malformed request in your `url parameters: id is not integer`\",\"code\":\"ERR0400\",\"status\":400}]}", string(w.Body.Bytes()))
				})
			})

			t.Run("Given request body but not json", func(t *testing.T) {
				t.Run("Return bad request", func(t *testing.T) {
					apiEngine := gin.Default()

					applicationManager := mock.NewMockApplicationManager(mockCtrl)
					applicationManager.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

					ctrl := applicationHttp.NewUpdate(applicationManager, apierror)
					controller.Provide(apiEngine.Handle, apierror, &cobra.Command{}).InjectHTTP(ctrl)

					w := testhelper.PerformRequest(apiEngine, ctrl.Method(), "/oauth/applications/1", bytes.NewReader([]byte(`this is gonna be error`)))

					assert.Equal(t, http.StatusBadRequest, w.Code)
					assert.Equal(t, "{\"errors\":[{\"title\":\"Bad request error\",\"detail\":\"You've send malformed request in your `invalid json format`\",\"code\":\"ERR0400\",\"status\":400}]}", string(w.Body.Bytes()))
				})
			})
		})
	})
}
