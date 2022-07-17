package application_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/kodefluence/altair/provider/plugin/oauth/controller"
	"github.com/kodefluence/altair/provider/plugin/oauth/entity"
	"github.com/kodefluence/altair/provider/plugin/oauth/eobject"
	"github.com/kodefluence/altair/provider/plugin/oauth/mock"
	"github.com/kodefluence/altair/testhelper"
	"github.com/kodefluence/altair/util"
	"github.com/stretchr/testify/assert"
)

func TestUpdate(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	t.Run("Method", func(t *testing.T) {
		applicationManager := mock.NewMockApplicationManager(mockCtrl)
		assert.Equal(t, "PUT", controller.NewApplication().Update(applicationManager).Method())
	})

	t.Run("Path", func(t *testing.T) {
		applicationManager := mock.NewMockApplicationManager(mockCtrl)
		assert.Equal(t, "/oauth/applications/:id", controller.NewApplication().Update(applicationManager).Path())
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

				ctrl := controller.NewApplication().Update(applicationManager)
				apiEngine.Handle(ctrl.Method(), ctrl.Path(), ctrl.Control)

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

				expectedError := &entity.Error{
					HttpStatus: http.StatusInternalServerError,
					Errors:     eobject.Wrap(eobject.InternalServerError(context.Background())),
				}

				applicationManager := mock.NewMockApplicationManager(mockCtrl)
				applicationManager.EXPECT().Update(gomock.Any(), 1, oauthApplicationUpdateJSON).Return(entity.OauthApplicationJSON{}, expectedError)

				ctrl := controller.NewApplication().Update(applicationManager)
				apiEngine.Handle(ctrl.Method(), ctrl.Path(), ctrl.Control)

				var response ErrorResponse
				w := testhelper.PerformRequest(apiEngine, ctrl.Method(), "/oauth/applications/1", bytes.NewReader(encodedBytes))

				err = json.Unmarshal(w.Body.Bytes(), &response)
				assert.Nil(t, err)

				assert.Equal(t, expectedError.HttpStatus, w.Code)
				assert.Equal(t, expectedError.Errors, response.Errors)
			})
		})

		t.Run("Given invalid request body", func(t *testing.T) {
			t.Run("Return bad request", func(t *testing.T) {
				apiEngine := gin.Default()

				applicationManager := mock.NewMockApplicationManager(mockCtrl)
				applicationManager.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

				ctrl := controller.NewApplication().Update(applicationManager)
				apiEngine.Handle(ctrl.Method(), ctrl.Path(), ctrl.Control)

				expectedError := &entity.Error{
					HttpStatus: http.StatusBadRequest,
					Errors:     eobject.Wrap(eobject.BadRequestError("request body")),
				}

				var response ErrorResponse
				w := testhelper.PerformRequest(apiEngine, ctrl.Method(), "/oauth/applications/1", testhelper.MockErrorIoReader{})

				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.Nil(t, err)

				assert.Equal(t, expectedError.HttpStatus, w.Code)
				assert.Equal(t, expectedError.Errors, response.Errors)
			})
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

				ctrl := controller.NewApplication().Update(applicationManager)
				apiEngine.Handle(ctrl.Method(), ctrl.Path(), ctrl.Control)

				expectedError := &entity.Error{
					HttpStatus: http.StatusBadRequest,
					Errors:     eobject.Wrap(eobject.BadRequestError("url parameters: id is not integer")),
				}

				var response ErrorResponse
				w := testhelper.PerformRequest(apiEngine, ctrl.Method(), "/oauth/applications/s", bytes.NewReader(encodedBytes))

				err = json.Unmarshal(w.Body.Bytes(), &response)
				assert.Nil(t, err)

				assert.Equal(t, expectedError.HttpStatus, w.Code)
				assert.Equal(t, expectedError.Errors, response.Errors)
			})
		})

		t.Run("Given request body but not json", func(t *testing.T) {
			t.Run("Return bad request", func(t *testing.T) {
				apiEngine := gin.Default()

				applicationManager := mock.NewMockApplicationManager(mockCtrl)
				applicationManager.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

				ctrl := controller.NewApplication().Update(applicationManager)
				apiEngine.Handle(ctrl.Method(), ctrl.Path(), ctrl.Control)

				expectedError := &entity.Error{
					HttpStatus: http.StatusBadRequest,
					Errors:     eobject.Wrap(eobject.BadRequestError("request body")),
				}

				var response ErrorResponse
				w := testhelper.PerformRequest(apiEngine, ctrl.Method(), "/oauth/applications/1", bytes.NewReader([]byte(`this is gonna be error`)))

				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.Nil(t, err)

				assert.Equal(t, expectedError.HttpStatus, w.Code)
				assert.Equal(t, expectedError.Errors, response.Errors)
			})
		})
	})
}
