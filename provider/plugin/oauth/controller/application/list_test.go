package application_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/codefluence-x/altair/entity"
	"github.com/codefluence-x/altair/eobject"
	"github.com/codefluence-x/altair/provider/plugin/oauth/controller"
	"github.com/codefluence-x/altair/provider/plugin/oauth/mock"
	"github.com/codefluence-x/altair/testhelper"
	"github.com/codefluence-x/altair/util"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type responseList struct {
	Data []entity.OauthApplicationJSON `json:"data"`
	Meta gin.H                         `json:"meta"`
}

func TestList(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	t.Run("Method", func(t *testing.T) {
		applicationManager := mock.NewMockApplicationManager(mockCtrl)
		assert.Equal(t, "GET", controller.Application().List(applicationManager).Method())
	})

	t.Run("Path", func(t *testing.T) {
		applicationManager := mock.NewMockApplicationManager(mockCtrl)
		assert.Equal(t, "/oauth/applications", controller.Application().List(applicationManager).Path())
	})

	t.Run("Control", func(t *testing.T) {
		t.Run("Given request with offset and limit", func(t *testing.T) {
			t.Run("Return list of oauth application", func(t *testing.T) {
				apiEngine := gin.Default()

				oauthApplicationJSONs := []entity.OauthApplicationJSON{
					entity.OauthApplicationJSON{
						ID:           util.IntToPointer(1),
						OwnerID:      util.IntToPointer(1),
						Description:  util.StringToPointer("Application 1"),
						Scopes:       util.StringToPointer("public user"),
						ClientUID:    util.StringToPointer("clientuid01"),
						ClientSecret: util.StringToPointer("clientsecret01"),
					},
					entity.OauthApplicationJSON{
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

				ctrl := controller.Application().List(applicationManager)
				apiEngine.Handle(ctrl.Method(), ctrl.Path(), ctrl.Control)

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

					expectedError := &entity.Error{
						HttpStatus: http.StatusInternalServerError,
						Errors:     eobject.Wrap(eobject.InternalServerError(context.Background())),
					}
					applicationManager := mock.NewMockApplicationManager(mockCtrl)
					applicationManager.EXPECT().List(gomock.Any(), 0, 10).Return([]entity.OauthApplicationJSON(nil), 0, expectedError)

					ctrl := controller.Application().List(applicationManager)
					apiEngine.Handle(ctrl.Method(), ctrl.Path(), ctrl.Control)

					var response testhelper.ErrorResponse
					w := testhelper.PerformRequest(apiEngine, ctrl.Method(), ctrl.Path(), nil)

					err := json.Unmarshal(w.Body.Bytes(), &response)
					assert.Nil(t, err)

					assert.Equal(t, expectedError.HttpStatus, w.Code)
					assert.Equal(t, expectedError.Errors, response.Errors)
				})
			})
		})

		t.Run("Given request with invalid offset", func(t *testing.T) {
			t.Run("Return bad request error", func(t *testing.T) {
				apiEngine := gin.Default()

				applicationManager := mock.NewMockApplicationManager(mockCtrl)
				applicationManager.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

				ctrl := controller.Application().List(applicationManager)
				apiEngine.Handle(ctrl.Method(), ctrl.Path(), ctrl.Control)

				expectedError := &entity.Error{
					HttpStatus: http.StatusBadRequest,
					Errors:     eobject.Wrap(eobject.BadRequestError("query parameters: offset")),
				}

				var response testhelper.ErrorResponse
				w := testhelper.PerformRequest(apiEngine, ctrl.Method(), ctrl.Path()+"?offset=invalid", nil)

				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.Nil(t, err)

				assert.Equal(t, expectedError.HttpStatus, w.Code)
				assert.Equal(t, expectedError.Errors, response.Errors)
			})
		})

		t.Run("Given request with invalid limit", func(t *testing.T) {
			t.Run("Return bad request error", func(t *testing.T) {
				apiEngine := gin.Default()

				applicationManager := mock.NewMockApplicationManager(mockCtrl)
				applicationManager.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

				ctrl := controller.Application().List(applicationManager)
				apiEngine.Handle(ctrl.Method(), ctrl.Path(), ctrl.Control)

				expectedError := &entity.Error{
					HttpStatus: http.StatusBadRequest,
					Errors:     eobject.Wrap(eobject.BadRequestError("query parameters: limit")),
				}

				var response testhelper.ErrorResponse
				w := testhelper.PerformRequest(apiEngine, ctrl.Method(), ctrl.Path()+"?limit=invalid", nil)

				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.Nil(t, err)

				assert.Equal(t, expectedError.HttpStatus, w.Code)
				assert.Equal(t, expectedError.Errors, response.Errors)
			})
		})
	})
}
