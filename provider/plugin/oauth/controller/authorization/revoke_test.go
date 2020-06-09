package authorization_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/codefluence-x/altair/provider/plugin/oauth/entity"
	"github.com/codefluence-x/altair/provider/plugin/oauth/eobject"

	"github.com/codefluence-x/altair/testhelper"

	"github.com/codefluence-x/altair/provider/plugin/oauth/controller"
	"github.com/codefluence-x/altair/provider/plugin/oauth/mock"

	"github.com/codefluence-x/altair/util"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type ErrorResponse struct {
	Errors []entity.ErrorObject `json:"errors"`
}

func TestRevoke(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	t.Run("Method", func(t *testing.T) {
		authorizationService := mock.NewMockAuthorization(mockCtrl)
		assert.Equal(t, "POST", controller.Authorization().Revoke(authorizationService).Method())
	})

	t.Run("Path", func(t *testing.T) {
		authorizationService := mock.NewMockAuthorization(mockCtrl)
		assert.Equal(t, "/oauth/authorizations/revoke", controller.Authorization().Revoke(authorizationService).Path())
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

				ctrl := controller.Authorization().Revoke(authorizationService)
				apiEngine.Handle(ctrl.Method(), ctrl.Path(), ctrl.Control)

				var response gin.H
				w := testhelper.PerformRequest(apiEngine, ctrl.Method(), ctrl.Path(), bytes.NewReader(encodedBytes))

				err = json.Unmarshal(w.Body.Bytes(), &response)
				assert.Nil(t, err)

				err = json.Unmarshal(w.Body.Bytes(), &response)
				assert.Nil(t, err)

				assert.Equal(t, http.StatusOK, w.Code)
				assert.Equal(t, "Access token has been successfully revoked.", response["message"])
			})

			t.Run("Unexpected error in authorization services", func(t *testing.T) {
				t.Run("Return entity error status", func(t *testing.T) {
					apiEngine := gin.Default()

					revokeTokenRequest := entity.RevokeAccessTokenRequestJSON{
						Token: nil,
					}
					encodedBytes, err := json.Marshal(revokeTokenRequest)
					assert.Nil(t, err)

					expectedError := &entity.Error{
						HttpStatus: http.StatusUnprocessableEntity,
						Errors:     eobject.Wrap(eobject.ValidationError("token is empty")),
					}

					authorizationService := mock.NewMockAuthorization(mockCtrl)
					authorizationService.EXPECT().RevokeToken(gomock.Any(), revokeTokenRequest).Return(expectedError)

					ctrl := controller.Authorization().Revoke(authorizationService)
					apiEngine.Handle(ctrl.Method(), ctrl.Path(), ctrl.Control)

					var response ErrorResponse
					w := testhelper.PerformRequest(apiEngine, ctrl.Method(), ctrl.Path(), bytes.NewReader(encodedBytes))

					err = json.Unmarshal(w.Body.Bytes(), &response)
					assert.Nil(t, err)

					assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
					assert.Equal(t, expectedError.Errors, response.Errors)
				})
			})
		})

		t.Run("Given invalid request body", func(t *testing.T) {
			t.Run("Return bad request", func(t *testing.T) {
				apiEngine := gin.Default()

				authorizationService := mock.NewMockAuthorization(mockCtrl)
				authorizationService.EXPECT().Grantor(gomock.Any(), gomock.Any()).Times(0)

				ctrl := controller.Authorization().Revoke(authorizationService)
				apiEngine.Handle(ctrl.Method(), ctrl.Path(), ctrl.Control)

				expectedError := &entity.Error{
					HttpStatus: http.StatusBadRequest,
					Errors:     eobject.Wrap(eobject.BadRequestError("request body")),
				}

				var response ErrorResponse
				w := testhelper.PerformRequest(apiEngine, ctrl.Method(), ctrl.Path(), testhelper.MockErrorIoReader{})

				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.Nil(t, err)

				assert.Equal(t, expectedError.HttpStatus, w.Code)
				assert.Equal(t, expectedError.Errors, response.Errors)
			})
		})

		t.Run("Given request body but not json", func(t *testing.T) {
			t.Run("Return bad request", func(t *testing.T) {
				apiEngine := gin.Default()

				authorizationService := mock.NewMockAuthorization(mockCtrl)
				authorizationService.EXPECT().Grantor(gomock.Any(), gomock.Any()).Times(0)

				ctrl := controller.Authorization().Revoke(authorizationService)
				apiEngine.Handle(ctrl.Method(), ctrl.Path(), ctrl.Control)

				expectedError := &entity.Error{
					HttpStatus: http.StatusBadRequest,
					Errors:     eobject.Wrap(eobject.BadRequestError("request body")),
				}

				var response ErrorResponse
				w := testhelper.PerformRequest(apiEngine, ctrl.Method(), ctrl.Path(), bytes.NewReader([]byte(`this is gonna be error`)))

				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.Nil(t, err)

				assert.Equal(t, expectedError.HttpStatus, w.Code)
				assert.Equal(t, expectedError.Errors, response.Errors)
			})
		})
	})
}
