package application_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/codefluence-x/altair/entity"
	"github.com/codefluence-x/altair/eobject"
	"github.com/codefluence-x/altair/mock"
	"github.com/codefluence-x/altair/provider/plugin/oauth/controller"
	"github.com/codefluence-x/altair/util"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type responseOne struct {
	Data entity.OauthApplicationJSON `json:"data"`
}

func TestOne(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	t.Run("Method", func(t *testing.T) {
		applicationManager := mock.NewMockApplicationManager(mockCtrl)
		assert.Equal(t, "GET", controller.Application().One(applicationManager).Method())
	})

	t.Run("Path", func(t *testing.T) {
		applicationManager := mock.NewMockApplicationManager(mockCtrl)
		assert.Equal(t, "/oauth/applications/:id", controller.Application().One(applicationManager).Path())
	})

	t.Run("Control", func(t *testing.T) {
		t.Run("Given id params", func(t *testing.T) {
			t.Run("Return oauth application data", func(t *testing.T) {
				apiEngine := gin.Default()

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

				ctrl := controller.Application().One(applicationManager)
				apiEngine.Handle(ctrl.Method(), ctrl.Path(), ctrl.Control)

				var response responseOne
				w := mock.PerformRequest(apiEngine, ctrl.Method(), "/oauth/applications/1", nil)

				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.Nil(t, err)

				assert.Equal(t, http.StatusOK, w.Code)
				assert.Equal(t, oauthApplicationJSON, response.Data)
			})

			t.Run("Unexpected error in application manager", func(t *testing.T) {
				t.Run("Return error", func(t *testing.T) {
					apiEngine := gin.Default()

					oauthApplicationJSON := entity.OauthApplicationJSON{}

					expectedError := &entity.Error{
						HttpStatus: http.StatusInternalServerError,
						Errors:     eobject.Wrap(eobject.InternalServerError(context.Background())),
					}

					applicationManager := mock.NewMockApplicationManager(mockCtrl)
					applicationManager.EXPECT().One(gomock.Any(), 1).Return(oauthApplicationJSON, expectedError)

					ctrl := controller.Application().One(applicationManager)
					apiEngine.Handle(ctrl.Method(), ctrl.Path(), ctrl.Control)

					var response mock.ErrorResponse
					w := mock.PerformRequest(apiEngine, ctrl.Method(), "/oauth/applications/1", nil)

					err := json.Unmarshal(w.Body.Bytes(), &response)
					assert.Nil(t, err)

					assert.Equal(t, expectedError.HttpStatus, w.Code)
					assert.Equal(t, expectedError.Errors, response.Errors)
				})
			})
		})

		t.Run("Given invalid params", func(t *testing.T) {
			t.Run("Return bad request", func(t *testing.T) {
				apiEngine := gin.Default()

				applicationManager := mock.NewMockApplicationManager(mockCtrl)
				applicationManager.EXPECT().One(gomock.Any(), gomock.Any()).Times(0)

				ctrl := controller.Application().One(applicationManager)
				apiEngine.Handle(ctrl.Method(), ctrl.Path(), ctrl.Control)

				expectedError := &entity.Error{
					HttpStatus: http.StatusBadRequest,
					Errors:     eobject.Wrap(eobject.BadRequestError("url parameters: id is not integer")),
				}

				var response mock.ErrorResponse
				w := mock.PerformRequest(apiEngine, ctrl.Method(), "/oauth/applications/x", nil)

				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.Nil(t, err)

				assert.Equal(t, expectedError.HttpStatus, w.Code)
				assert.Equal(t, expectedError.Errors, response.Errors)
			})
		})
	})
}
