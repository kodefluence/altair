package http_test

import (
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
	"github.com/kodefluence/monorepo/jsonapi"
	"github.com/kodefluence/monorepo/kontext"
	"github.com/stretchr/testify/assert"
)

func TestOne(t *testing.T) {
	type responseOne struct {
		Data entity.OauthApplicationJSON `json:"data"`
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	apierror := apierror.Provide()

	t.Run("Method", func(t *testing.T) {
		applicationManager := mock.NewMockApplicationManager(mockCtrl)
		assert.Equal(t, "GET", applicationHttp.NewOne(applicationManager, apierror).Method())
	})

	t.Run("Path", func(t *testing.T) {
		applicationManager := mock.NewMockApplicationManager(mockCtrl)
		assert.Equal(t, "/oauth/applications/:id", applicationHttp.NewOne(applicationManager, apierror).Path())
	})

	t.Run("Control", func(t *testing.T) {
		t.Run("Given id params", func(t *testing.T) {
			t.Run("Return oauth application data", func(t *testing.T) {
				apiEngine := gin.New()

				oauthApplicationJSON := entity.OauthApplicationJSON{
					ID:           util.IntToPointer(1),
					OwnerID:      util.IntToPointer(1),
					Description:  util.StringToPointer("Application 1"),
					Scopes:       util.StringToPointer("public user"),
					ClientUID:    util.StringToPointer("clientuid01"),
					ClientSecret: util.StringToPointer("clientsecret01"),
				}

				applicationManager := mock.NewMockApplicationManager(mockCtrl)
				applicationManager.EXPECT().One(gomock.Any(), 1).Return(oauthApplicationJSON, nil)

				ctrl := applicationHttp.NewOne(applicationManager, apierror)
				apiEngine.Handle(ctrl.Method(), ctrl.Path(), ctrl.Control)

				var response responseOne
				w := testhelper.PerformRequest(apiEngine, ctrl.Method(), "/oauth/applications/1", nil)

				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.Nil(t, err)

				assert.Equal(t, http.StatusOK, w.Code)
				assert.Equal(t, oauthApplicationJSON, response.Data)
			})

			t.Run("Unexpected error in application manager", func(t *testing.T) {
				t.Run("Return error", func(t *testing.T) {
					apiEngine := gin.New()

					oauthApplicationJSON := entity.OauthApplicationJSON{}

					applicationManager := mock.NewMockApplicationManager(mockCtrl)
					applicationManager.EXPECT().One(gomock.Any(), 1).Return(oauthApplicationJSON, jsonapi.BuildResponse(apierror.InternalServerError(kontext.Fabricate())).Errors)

					ctrl := applicationHttp.NewOne(applicationManager, apierror)
					apiEngine.Handle(ctrl.Method(), ctrl.Path(), ctrl.Control)

					w := testhelper.PerformRequest(apiEngine, ctrl.Method(), "/oauth/applications/1", nil)
					assert.Equal(t, http.StatusInternalServerError, w.Code)
					assert.Equal(t, "{\"errors\":[{\"title\":\"Internal server error\",\"detail\":\"Something is not right, help us fix this problem. Contribute to https://github.com/kodefluence/altair. Tracing code: '\\u003cnil\\u003e'\",\"code\":\"ERR0500\",\"status\":500}]}", string(w.Body.Bytes()))
				})
			})
		})
	})

	t.Run("Given invalid params", func(t *testing.T) {
		t.Run("Return bad request", func(t *testing.T) {
			apiEngine := gin.New()

			applicationManager := mock.NewMockApplicationManager(mockCtrl)
			ctrl := applicationHttp.NewOne(applicationManager, apierror)
			apiEngine.Handle(ctrl.Method(), ctrl.Path(), ctrl.Control)

			w := testhelper.PerformRequest(apiEngine, ctrl.Method(), "/oauth/applications/x", nil)
			assert.Equal(t, http.StatusBadRequest, w.Code)
			assert.Equal(t, "{\"errors\":[{\"title\":\"Bad request error\",\"detail\":\"You've send malformed request in your `url parameters: id is not integer`\",\"code\":\"ERR0400\",\"status\":400}]}", string(w.Body.Bytes()))
		})
	})
}
