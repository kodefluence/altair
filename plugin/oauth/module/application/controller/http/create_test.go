package http_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/kodefluence/altair/module/apierror"
	"github.com/kodefluence/altair/plugin/oauth/entity"
	applicationHttp "github.com/kodefluence/altair/plugin/oauth/module/application/controller/http"
	"github.com/kodefluence/altair/plugin/oauth/module/application/controller/http/mock"
	"github.com/kodefluence/altair/testhelper"
	"github.com/kodefluence/altair/util"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	apierror := apierror.Provide()

	t.Run("Method", func(t *testing.T) {
		applicationManager := mock.NewMockApplicationManager(mockCtrl)
		assert.Equal(t, "POST", applicationHttp.NewCreate(applicationManager, apierror).Method())
	})

	t.Run("Path", func(t *testing.T) {
		applicationManager := mock.NewMockApplicationManager(mockCtrl)
		assert.Equal(t, "/oauth/applications", applicationHttp.NewCreate(applicationManager, apierror).Path())
	})

	t.Run("Control", func(t *testing.T) {
		t.Run("Given request with json body", func(t *testing.T) {
			t.Run("Return oauth application data with status 202", func(t *testing.T) {
				apiEngine := gin.Default()

				oauthApplicationJSON := entity.OauthApplicationJSON{
					OwnerID:      util.IntToPointer(1),
					Description:  util.StringToPointer("Application 1"),
					Scopes:       util.StringToPointer("public user"),
					ClientUID:    util.StringToPointer("clientuid01"),
					ClientSecret: util.StringToPointer("clientsecret01"),
				}
				encodedBytes, err := json.Marshal(oauthApplicationJSON)
				assert.Nil(t, err)

				applicationManager := mock.NewMockApplicationManager(mockCtrl)
				applicationManager.EXPECT().Create(gomock.Any(), oauthApplicationJSON).Return(oauthApplicationJSON, nil)

				ctrl := applicationHttp.NewCreate(applicationManager, apierror)
				apiEngine.Handle(ctrl.Method(), ctrl.Path(), ctrl.Control)

				var response responseOne
				w := testhelper.PerformRequest(apiEngine, ctrl.Method(), ctrl.Path(), bytes.NewReader(encodedBytes))

				err = json.Unmarshal(w.Body.Bytes(), &response)
				assert.Nil(t, err)

				assert.Equal(t, http.StatusCreated, w.Code)
				assert.Equal(t, oauthApplicationJSON, response.Data)
			})

			t.Run("Unexpected error in application manager", func(t *testing.T) {
				apiEngine := gin.Default()

				oauthApplicationJSON := entity.OauthApplicationJSON{
					OwnerID:      util.IntToPointer(1),
					Description:  util.StringToPointer("Application 1"),
					Scopes:       util.StringToPointer("public user"),
					ClientUID:    util.StringToPointer("clientuid01"),
					ClientSecret: util.StringToPointer("clientsecret01"),
				}
				encodedBytes, err := json.Marshal(oauthApplicationJSON)
				assert.Nil(t, err)

				applicationManager := mock.NewMockApplicationManager(mockCtrl)
				applicationManager.EXPECT().Create(gomock.Any(), oauthApplicationJSON).Return(entity.OauthApplicationJSON{}, testhelper.ErrInternalServer())

				ctrl := applicationHttp.NewCreate(applicationManager, apierror)
				apiEngine.Handle(ctrl.Method(), ctrl.Path(), ctrl.Control)

				w := testhelper.PerformRequest(apiEngine, ctrl.Method(), ctrl.Path(), bytes.NewReader(encodedBytes))

				assert.Equal(t, http.StatusInternalServerError, w.Code)
				assert.Equal(t, "{\"errors\":[{\"title\":\"Internal server error\",\"detail\":\"Something is not right, help us fix this problem. Contribute to https://github.com/kodefluence/altair. Tracing code: '\\u003cnil\\u003e'\",\"code\":\"ERR0500\",\"status\":500}]}", string(w.Body.Bytes()))
			})
		})

		t.Run("Given invalid request body", func(t *testing.T) {
			t.Run("Return bad request", func(t *testing.T) {
				apiEngine := gin.Default()

				applicationManager := mock.NewMockApplicationManager(mockCtrl)
				applicationManager.EXPECT().Create(gomock.Any(), gomock.Any()).Times(0)

				ctrl := applicationHttp.NewCreate(applicationManager, apierror)
				apiEngine.Handle(ctrl.Method(), ctrl.Path(), ctrl.Control)

				w := testhelper.PerformRequest(apiEngine, ctrl.Method(), ctrl.Path(), testhelper.MockErrorIoReader{})

				assert.Equal(t, http.StatusBadRequest, w.Code)
				assert.Equal(t, "{\"errors\":[{\"title\":\"Bad request error\",\"detail\":\"You've send malformed request in your `request body`\",\"code\":\"ERR0400\",\"status\":400}]}", string(w.Body.Bytes()))
			})
		})

		t.Run("Given request body but not json", func(t *testing.T) {
			t.Run("Return bad request", func(t *testing.T) {
				apiEngine := gin.Default()

				applicationManager := mock.NewMockApplicationManager(mockCtrl)
				applicationManager.EXPECT().Create(gomock.Any(), gomock.Any()).Times(0)

				ctrl := applicationHttp.NewCreate(applicationManager, apierror)
				apiEngine.Handle(ctrl.Method(), ctrl.Path(), ctrl.Control)

				w := testhelper.PerformRequest(apiEngine, ctrl.Method(), ctrl.Path(), bytes.NewReader([]byte(`this is gonna be error`)))

				assert.Equal(t, http.StatusBadRequest, w.Code)
				assert.Equal(t, "{\"errors\":[{\"title\":\"Bad request error\",\"detail\":\"You've send malformed request in your `invalid json format`\",\"code\":\"ERR0400\",\"status\":400}]}", string(w.Body.Bytes()))
			})
		})
	})
}