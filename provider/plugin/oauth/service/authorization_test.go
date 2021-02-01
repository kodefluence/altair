package service_test

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"testing"

	"github.com/codefluence-x/altair/provider/plugin/oauth/entity"
	"github.com/codefluence-x/altair/provider/plugin/oauth/eobject"
	"github.com/codefluence-x/altair/provider/plugin/oauth/mock"
	"github.com/codefluence-x/altair/provider/plugin/oauth/service"
	"github.com/codefluence-x/altair/util"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAuthorization(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	t.Run("RevokeToken", func(t *testing.T) {
		t.Run("Given context and revoke access token request", func(t *testing.T) {
			t.Run("Run gracefully", func(t *testing.T) {
				t.Run("Return nil", func(t *testing.T) {
					oauthApplicationModel := mock.NewMockOauthApplicationModel(mockCtrl)
					oauthAccessTokenModel := mock.NewMockOauthAccessTokenModel(mockCtrl)
					oauthAccessGrantModel := mock.NewMockOauthAccessGrantModel(mockCtrl)
					oauthRefreshTokenModel := mock.NewMockOauthRefreshTokenModel(mockCtrl)
					oauthValidator := mock.NewMockOauthValidator(mockCtrl)
					modelFormatterMock := mock.NewMockModelFormater(mockCtrl)
					oauthFormatterMock := mock.NewMockOauthFormatter(mockCtrl)

					ctx := context.WithValue(context.Background(), "request_id", uuid.New().String())

					revokeRequest := entity.RevokeAccessTokenRequestJSON{
						Token: util.StringToPointer("some-cool-token"),
					}

					oauthAccessTokenModel.EXPECT().Revoke(ctx, *revokeRequest.Token).Return(nil)

					authorizationService := service.NewAuthorization(oauthApplicationModel, oauthAccessTokenModel, oauthAccessGrantModel, oauthRefreshTokenModel, modelFormatterMock, oauthValidator, oauthFormatterMock, false)
					err := authorizationService.RevokeToken(ctx, revokeRequest)
					assert.Nil(t, err)
				})
			})

			t.Run("Revoke error", func(t *testing.T) {
				t.Run("Token not found", func(t *testing.T) {
					t.Run("Return 404 error", func(t *testing.T) {
						oauthApplicationModel := mock.NewMockOauthApplicationModel(mockCtrl)
						oauthAccessTokenModel := mock.NewMockOauthAccessTokenModel(mockCtrl)
						oauthAccessGrantModel := mock.NewMockOauthAccessGrantModel(mockCtrl)
						oauthRefreshTokenModel := mock.NewMockOauthRefreshTokenModel(mockCtrl)
						oauthValidator := mock.NewMockOauthValidator(mockCtrl)
						modelFormatterMock := mock.NewMockModelFormater(mockCtrl)
						oauthFormatterMock := mock.NewMockOauthFormatter(mockCtrl)

						ctx := context.WithValue(context.Background(), "request_id", uuid.New().String())

						revokeRequest := entity.RevokeAccessTokenRequestJSON{
							Token: util.StringToPointer("some-cool-token"),
						}

						oauthAccessTokenModel.EXPECT().Revoke(ctx, *revokeRequest.Token).Return(sql.ErrNoRows)

						expectedError := &entity.Error{
							HttpStatus: http.StatusNotFound,
							Errors:     eobject.Wrap(eobject.NotFoundError(ctx, "token")),
						}

						authorizationService := service.NewAuthorization(oauthApplicationModel, oauthAccessTokenModel, oauthAccessGrantModel, oauthRefreshTokenModel, modelFormatterMock, oauthValidator, oauthFormatterMock, false)
						err := authorizationService.RevokeToken(ctx, revokeRequest)
						assert.Equal(t, expectedError, err)
					})
				})

				t.Run("Other error", func(t *testing.T) {
					t.Run("Return 500 error", func(t *testing.T) {
						oauthApplicationModel := mock.NewMockOauthApplicationModel(mockCtrl)
						oauthAccessTokenModel := mock.NewMockOauthAccessTokenModel(mockCtrl)
						oauthAccessGrantModel := mock.NewMockOauthAccessGrantModel(mockCtrl)
						oauthRefreshTokenModel := mock.NewMockOauthRefreshTokenModel(mockCtrl)
						oauthValidator := mock.NewMockOauthValidator(mockCtrl)
						modelFormatterMock := mock.NewMockModelFormater(mockCtrl)
						oauthFormatterMock := mock.NewMockOauthFormatter(mockCtrl)

						ctx := context.WithValue(context.Background(), "request_id", uuid.New().String())

						revokeRequest := entity.RevokeAccessTokenRequestJSON{
							Token: util.StringToPointer("some-cool-token"),
						}

						oauthAccessTokenModel.EXPECT().Revoke(ctx, *revokeRequest.Token).Return(errors.New("unexpected error"))

						expectedError := &entity.Error{
							HttpStatus: http.StatusInternalServerError,
							Errors:     eobject.Wrap(eobject.InternalServerError(ctx)),
						}

						authorizationService := service.NewAuthorization(oauthApplicationModel, oauthAccessTokenModel, oauthAccessGrantModel, oauthRefreshTokenModel, modelFormatterMock, oauthValidator, oauthFormatterMock, false)
						err := authorizationService.RevokeToken(ctx, revokeRequest)
						assert.Equal(t, expectedError, err)
					})
				})
			})
		})

		t.Run("Given context and revoke access token request with nil token", func(t *testing.T) {
			t.Run("Return 422 error", func(t *testing.T) {
				oauthApplicationModel := mock.NewMockOauthApplicationModel(mockCtrl)
				oauthAccessTokenModel := mock.NewMockOauthAccessTokenModel(mockCtrl)
				oauthAccessGrantModel := mock.NewMockOauthAccessGrantModel(mockCtrl)
				oauthRefreshTokenModel := mock.NewMockOauthRefreshTokenModel(mockCtrl)
				oauthValidator := mock.NewMockOauthValidator(mockCtrl)
				modelFormatterMock := mock.NewMockModelFormater(mockCtrl)
				oauthFormatterMock := mock.NewMockOauthFormatter(mockCtrl)

				ctx := context.WithValue(context.Background(), "request_id", uuid.New().String())

				revokeRequest := entity.RevokeAccessTokenRequestJSON{
					Token: nil,
				}

				oauthAccessTokenModel.EXPECT().Revoke(gomock.Any(), gomock.Any()).Times(0)

				expectedError := &entity.Error{
					HttpStatus: http.StatusUnprocessableEntity,
					Errors:     eobject.Wrap(eobject.ValidationError("token is empty")),
				}

				authorizationService := service.NewAuthorization(oauthApplicationModel, oauthAccessTokenModel, oauthAccessGrantModel, oauthRefreshTokenModel, modelFormatterMock, oauthValidator, oauthFormatterMock, false)
				err := authorizationService.RevokeToken(ctx, revokeRequest)
				assert.Equal(t, expectedError, err)
			})
		})
	})
}
