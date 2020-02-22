package service_test

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/codefluence-x/altair/entity"
	"github.com/codefluence-x/altair/eobject"
	"github.com/codefluence-x/altair/formatter"
	"github.com/codefluence-x/altair/mock"
	"github.com/codefluence-x/altair/service"
	"github.com/codefluence-x/altair/util"
	"github.com/codefluence-x/aurelia"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAuthorization(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	t.Run("Grantor", func(t *testing.T) {
		t.Run("Given context and authorization request with a response type of token", func(t *testing.T) {
			t.Run("Return entity.OauthAccessTokenJSON and nil", func(t *testing.T) {
				oauthApplicationModel := mock.NewMockOauthApplicationModel(mockCtrl)
				oauthAccessTokenModel := mock.NewMockOauthAccessTokenModel(mockCtrl)
				oauthAccessGrantModel := mock.NewMockOauthAccessGrantModel(mockCtrl)
				oauthValidator := mock.NewMockOauthValidator(mockCtrl)
				modelFormatter := formatter.Model(time.Hour*4, time.Hour*2)
				modelFormatterMock := mock.NewMockModelFormater(mockCtrl)
				oauthFormatter := formatter.Oauth()
				oauthFormatterMock := mock.NewMockOauthFormatter(mockCtrl)

				ctx := context.WithValue(context.Background(), "track_id", uuid.New().String())

				authorizationRequest := entity.AuthorizationRequestJSON{
					ResponseType:    util.StringToPointer("token"),
					ResourceOwnerID: util.IntToPointer(1),
					ClientUID:       util.StringToPointer(aurelia.Hash("x", "y")),
					ClientSecret:    util.StringToPointer(aurelia.Hash("z", "a")),
					RedirectURI:     util.StringToPointer("http://github.com"),
					Scopes:          util.StringToPointer("public users"),
				}

				oauthApplication := entity.OauthApplication{
					ID: 1,
					OwnerID: sql.NullInt64{
						Int64: 1,
						Valid: true,
					},
					OwnerType:    "confidential",
					Description:  "Blablabla",
					Scopes:       "public users stores",
					ClientUID:    *authorizationRequest.ClientUID,
					ClientSecret: *authorizationRequest.ClientSecret,
					CreatedAt:    time.Now().Add(-time.Hour * 4),
					UpdatedAt:    time.Now(),
				}

				oauthAccessToken := entity.OauthAccessToken{
					ID:                 1,
					OauthApplicationID: oauthApplication.ID,
					ResourceOwnerID:    *authorizationRequest.ResourceOwnerID,
					Token:              aurelia.Hash("x", "y"),
					Scopes:             *authorizationRequest.Scopes,
					ExpiresIn:          time.Now().Add(time.Hour * 4),
					CreatedAt:          time.Now(),
				}

				oauthAccessTokenInsertable := modelFormatter.AccessTokenFromAuthorizationRequest(authorizationRequest, oauthApplication)
				oauthAccessTokenJSON := oauthFormatter.AccessToken(authorizationRequest, oauthAccessToken)

				gomock.InOrder(
					oauthApplicationModel.EXPECT().
						OneByUIDandSecret(ctx, *authorizationRequest.ClientUID, *authorizationRequest.ClientSecret).
						Return(oauthApplication, nil),
					oauthValidator.EXPECT().ValidateAuthorizationGrant(ctx, authorizationRequest, oauthApplication).Return(nil),
					modelFormatterMock.EXPECT().AccessTokenFromAuthorizationRequest(authorizationRequest, oauthApplication).Return(oauthAccessTokenInsertable),
					oauthAccessTokenModel.EXPECT().Create(ctx, oauthAccessTokenInsertable).Return(oauthAccessToken.ID, nil),
					oauthAccessTokenModel.EXPECT().One(ctx, oauthAccessToken.ID).Return(oauthAccessToken, nil),
					oauthFormatterMock.EXPECT().AccessToken(authorizationRequest, oauthAccessToken).Return(oauthAccessTokenJSON),
				)

				authorizationService := service.Authorization(oauthApplicationModel, oauthAccessTokenModel, oauthAccessGrantModel, modelFormatterMock, oauthValidator, oauthFormatterMock)
				results, err := authorizationService.Grantor(ctx, authorizationRequest)
				assert.Nil(t, err)
				assert.Equal(t, oauthAccessTokenJSON, results)
			})

			t.Run("Oauth application model return error", func(t *testing.T) {

				t.Run("Not found error", func(t *testing.T) {
					t.Run("Return error 404", func(t *testing.T) {
						oauthApplicationModel := mock.NewMockOauthApplicationModel(mockCtrl)
						oauthAccessTokenModel := mock.NewMockOauthAccessTokenModel(mockCtrl)
						oauthAccessGrantModel := mock.NewMockOauthAccessGrantModel(mockCtrl)
						oauthValidator := mock.NewMockOauthValidator(mockCtrl)
						modelFormatterMock := mock.NewMockModelFormater(mockCtrl)
						oauthFormatterMock := mock.NewMockOauthFormatter(mockCtrl)

						ctx := context.WithValue(context.Background(), "track_id", uuid.New().String())

						authorizationRequest := entity.AuthorizationRequestJSON{
							ResponseType:    util.StringToPointer("token"),
							ResourceOwnerID: util.IntToPointer(1),
							ClientUID:       util.StringToPointer(aurelia.Hash("x", "y")),
							ClientSecret:    util.StringToPointer(aurelia.Hash("z", "a")),
							RedirectURI:     util.StringToPointer("http://github.com"),
							Scopes:          util.StringToPointer("public users"),
						}

						expectedError := &entity.Error{
							HttpStatus: http.StatusNotFound,
							Errors:     eobject.Wrap(eobject.NotFoundError(ctx, "client_uid & client_secret")),
						}

						gomock.InOrder(
							oauthApplicationModel.EXPECT().OneByUIDandSecret(ctx, *authorizationRequest.ClientUID, *authorizationRequest.ClientSecret).Return(entity.OauthApplication{}, sql.ErrNoRows),
							oauthValidator.EXPECT().ValidateAuthorizationGrant(gomock.Any(), gomock.Any(), gomock.Any()).Times(0),
							modelFormatterMock.EXPECT().AccessTokenFromAuthorizationRequest(gomock.Any(), gomock.Any()).Times(0),
							oauthAccessTokenModel.EXPECT().Create(gomock.Any(), gomock.Any()).Times(0),
							oauthAccessTokenModel.EXPECT().One(gomock.Any(), gomock.Any()).Times(0),
							oauthFormatterMock.EXPECT().AccessToken(gomock.Any(), gomock.Any()).Times(0),
						)

						authorizationService := service.Authorization(oauthApplicationModel, oauthAccessTokenModel, oauthAccessGrantModel, modelFormatterMock, oauthValidator, oauthFormatterMock)
						results, err := authorizationService.Grantor(ctx, authorizationRequest)
						assert.NotNil(t, err)
						assert.Equal(t, expectedError, err)
						assert.Equal(t, entity.OauthAccessTokenJSON{}, results)
					})
				})

				t.Run("Unexpected error", func(t *testing.T) {
					t.Run("Return error 500", func(t *testing.T) {
						oauthApplicationModel := mock.NewMockOauthApplicationModel(mockCtrl)
						oauthAccessTokenModel := mock.NewMockOauthAccessTokenModel(mockCtrl)
						oauthAccessGrantModel := mock.NewMockOauthAccessGrantModel(mockCtrl)
						oauthValidator := mock.NewMockOauthValidator(mockCtrl)
						modelFormatterMock := mock.NewMockModelFormater(mockCtrl)
						oauthFormatterMock := mock.NewMockOauthFormatter(mockCtrl)

						ctx := context.WithValue(context.Background(), "track_id", uuid.New().String())

						authorizationRequest := entity.AuthorizationRequestJSON{
							ResponseType:    util.StringToPointer("token"),
							ResourceOwnerID: util.IntToPointer(1),
							ClientUID:       util.StringToPointer(aurelia.Hash("x", "y")),
							ClientSecret:    util.StringToPointer(aurelia.Hash("z", "a")),
							RedirectURI:     util.StringToPointer("http://github.com"),
							Scopes:          util.StringToPointer("public users"),
						}

						expectedError := &entity.Error{
							HttpStatus: http.StatusInternalServerError,
							Errors:     eobject.Wrap(eobject.InternalServerError(ctx)),
						}

						gomock.InOrder(
							oauthApplicationModel.EXPECT().OneByUIDandSecret(ctx, *authorizationRequest.ClientUID, *authorizationRequest.ClientSecret).Return(entity.OauthApplication{}, errors.New("unexpected error")),
							oauthValidator.EXPECT().ValidateAuthorizationGrant(gomock.Any(), gomock.Any(), gomock.Any()).Times(0),
							modelFormatterMock.EXPECT().AccessTokenFromAuthorizationRequest(gomock.Any(), gomock.Any()).Times(0),
							oauthAccessTokenModel.EXPECT().Create(gomock.Any(), gomock.Any()).Times(0),
							oauthAccessTokenModel.EXPECT().One(gomock.Any(), gomock.Any()).Times(0),
							oauthFormatterMock.EXPECT().AccessToken(gomock.Any(), gomock.Any()).Times(0),
						)

						authorizationService := service.Authorization(oauthApplicationModel, oauthAccessTokenModel, oauthAccessGrantModel, modelFormatterMock, oauthValidator, oauthFormatterMock)
						results, err := authorizationService.Grantor(ctx, authorizationRequest)
						assert.NotNil(t, err)
						assert.Equal(t, expectedError, err)
						assert.Equal(t, entity.OauthAccessTokenJSON{}, results)
					})
				})
			})

			t.Run("Oauth authorization grant validation failed", func(t *testing.T) {
				t.Run("Return entity.Error", func(t *testing.T) {
					oauthApplicationModel := mock.NewMockOauthApplicationModel(mockCtrl)
					oauthAccessTokenModel := mock.NewMockOauthAccessTokenModel(mockCtrl)
					oauthAccessGrantModel := mock.NewMockOauthAccessGrantModel(mockCtrl)
					oauthValidator := mock.NewMockOauthValidator(mockCtrl)
					modelFormatterMock := mock.NewMockModelFormater(mockCtrl)
					oauthFormatterMock := mock.NewMockOauthFormatter(mockCtrl)

					ctx := context.WithValue(context.Background(), "track_id", uuid.New().String())

					authorizationRequest := entity.AuthorizationRequestJSON{
						ResponseType:    util.StringToPointer("token"),
						ResourceOwnerID: util.IntToPointer(1),
						ClientUID:       util.StringToPointer(aurelia.Hash("x", "y")),
						ClientSecret:    util.StringToPointer(aurelia.Hash("z", "a")),
						RedirectURI:     util.StringToPointer("http://github.com"),
						Scopes:          util.StringToPointer("public users"),
					}

					oauthApplication := entity.OauthApplication{
						ID: 1,
						OwnerID: sql.NullInt64{
							Int64: 1,
							Valid: true,
						},
						OwnerType:    "confidential",
						Description:  "Blablabla",
						Scopes:       "public users stores",
						ClientUID:    *authorizationRequest.ClientUID,
						ClientSecret: *authorizationRequest.ClientSecret,
						CreatedAt:    time.Now().Add(-time.Hour * 4),
						UpdatedAt:    time.Now(),
					}

					expectedError := &entity.Error{
						HttpStatus: http.StatusForbidden,
						Errors:     eobject.Wrap(eobject.ForbiddenError(ctx, "access_token", "your response type is not allowed in this application")),
					}

					gomock.InOrder(
						oauthApplicationModel.EXPECT().
							OneByUIDandSecret(ctx, *authorizationRequest.ClientUID, *authorizationRequest.ClientSecret).
							Return(oauthApplication, nil),
						oauthValidator.EXPECT().ValidateAuthorizationGrant(ctx, authorizationRequest, oauthApplication).Return(expectedError),
						modelFormatterMock.EXPECT().AccessTokenFromAuthorizationRequest(gomock.Any(), gomock.Any()).Times(0),
						oauthAccessTokenModel.EXPECT().Create(gomock.Any(), gomock.Any()).Times(0),
						oauthAccessTokenModel.EXPECT().One(gomock.Any(), gomock.Any()).Times(0),
						oauthFormatterMock.EXPECT().AccessToken(gomock.Any(), gomock.Any()).Times(0),
					)

					authorizationService := service.Authorization(oauthApplicationModel, oauthAccessTokenModel, oauthAccessGrantModel, modelFormatterMock, oauthValidator, oauthFormatterMock)
					results, err := authorizationService.Grantor(ctx, authorizationRequest)
					assert.NotNil(t, err)
					assert.Equal(t, expectedError, err)
					assert.Equal(t, entity.OauthAccessTokenJSON{}, results)
				})
			})

			t.Run("Failed created oauth access token", func(t *testing.T) {
				t.Run("Return error 500", func(t *testing.T) {
					oauthApplicationModel := mock.NewMockOauthApplicationModel(mockCtrl)
					oauthAccessTokenModel := mock.NewMockOauthAccessTokenModel(mockCtrl)
					oauthAccessGrantModel := mock.NewMockOauthAccessGrantModel(mockCtrl)
					oauthValidator := mock.NewMockOauthValidator(mockCtrl)
					modelFormatter := formatter.Model(time.Hour*4, time.Hour*2)
					modelFormatterMock := mock.NewMockModelFormater(mockCtrl)
					oauthFormatterMock := mock.NewMockOauthFormatter(mockCtrl)

					ctx := context.WithValue(context.Background(), "track_id", uuid.New().String())

					authorizationRequest := entity.AuthorizationRequestJSON{
						ResponseType:    util.StringToPointer("token"),
						ResourceOwnerID: util.IntToPointer(1),
						ClientUID:       util.StringToPointer(aurelia.Hash("x", "y")),
						ClientSecret:    util.StringToPointer(aurelia.Hash("z", "a")),
						RedirectURI:     util.StringToPointer("http://github.com"),
						Scopes:          util.StringToPointer("public users"),
					}

					oauthApplication := entity.OauthApplication{
						ID: 1,
						OwnerID: sql.NullInt64{
							Int64: 1,
							Valid: true,
						},
						OwnerType:    "confidential",
						Description:  "Blablabla",
						Scopes:       "public users stores",
						ClientUID:    *authorizationRequest.ClientUID,
						ClientSecret: *authorizationRequest.ClientSecret,
						CreatedAt:    time.Now().Add(-time.Hour * 4),
						UpdatedAt:    time.Now(),
					}

					oauthAccessTokenInsertable := modelFormatter.AccessTokenFromAuthorizationRequest(authorizationRequest, oauthApplication)

					gomock.InOrder(
						oauthApplicationModel.EXPECT().
							OneByUIDandSecret(ctx, *authorizationRequest.ClientUID, *authorizationRequest.ClientSecret).
							Return(oauthApplication, nil),
						oauthValidator.EXPECT().ValidateAuthorizationGrant(ctx, authorizationRequest, oauthApplication).Return(nil),
						modelFormatterMock.EXPECT().AccessTokenFromAuthorizationRequest(authorizationRequest, oauthApplication).Return(oauthAccessTokenInsertable),
						oauthAccessTokenModel.EXPECT().Create(ctx, oauthAccessTokenInsertable).Return(0, errors.New("unexpected error")),
						oauthAccessTokenModel.EXPECT().One(gomock.Any(), gomock.Any()).Times(0),
						oauthFormatterMock.EXPECT().AccessToken(gomock.Any(), gomock.Any()).Times(0),
					)

					expectedError := &entity.Error{
						HttpStatus: http.StatusInternalServerError,
						Errors:     eobject.Wrap(eobject.InternalServerError(ctx)),
					}

					authorizationService := service.Authorization(oauthApplicationModel, oauthAccessTokenModel, oauthAccessGrantModel, modelFormatterMock, oauthValidator, oauthFormatterMock)
					results, err := authorizationService.Grantor(ctx, authorizationRequest)
					assert.NotNil(t, err)
					assert.Equal(t, expectedError, err)
					assert.Equal(t, entity.OauthAccessTokenJSON{}, results)
				})
			})

			t.Run("Failed to fetch newly created access token", func(t *testing.T) {
				t.Run("Return error 500", func(t *testing.T) {
					oauthApplicationModel := mock.NewMockOauthApplicationModel(mockCtrl)
					oauthAccessTokenModel := mock.NewMockOauthAccessTokenModel(mockCtrl)
					oauthAccessGrantModel := mock.NewMockOauthAccessGrantModel(mockCtrl)
					oauthValidator := mock.NewMockOauthValidator(mockCtrl)
					modelFormatter := formatter.Model(time.Hour*4, time.Hour*2)
					modelFormatterMock := mock.NewMockModelFormater(mockCtrl)
					oauthFormatterMock := mock.NewMockOauthFormatter(mockCtrl)

					ctx := context.WithValue(context.Background(), "track_id", uuid.New().String())

					authorizationRequest := entity.AuthorizationRequestJSON{
						ResponseType:    util.StringToPointer("token"),
						ResourceOwnerID: util.IntToPointer(1),
						ClientUID:       util.StringToPointer(aurelia.Hash("x", "y")),
						ClientSecret:    util.StringToPointer(aurelia.Hash("z", "a")),
						RedirectURI:     util.StringToPointer("http://github.com"),
						Scopes:          util.StringToPointer("public users"),
					}

					oauthApplication := entity.OauthApplication{
						ID: 1,
						OwnerID: sql.NullInt64{
							Int64: 1,
							Valid: true,
						},
						OwnerType:    "confidential",
						Description:  "Blablabla",
						Scopes:       "public users stores",
						ClientUID:    *authorizationRequest.ClientUID,
						ClientSecret: *authorizationRequest.ClientSecret,
						CreatedAt:    time.Now().Add(-time.Hour * 4),
						UpdatedAt:    time.Now(),
					}

					oauthAccessToken := entity.OauthAccessToken{
						ID:                 1,
						OauthApplicationID: oauthApplication.ID,
						ResourceOwnerID:    *authorizationRequest.ResourceOwnerID,
						Token:              aurelia.Hash("x", "y"),
						Scopes:             *authorizationRequest.Scopes,
						ExpiresIn:          time.Now().Add(time.Hour * 4),
						CreatedAt:          time.Now(),
					}

					oauthAccessTokenInsertable := modelFormatter.AccessTokenFromAuthorizationRequest(authorizationRequest, oauthApplication)

					gomock.InOrder(
						oauthApplicationModel.EXPECT().
							OneByUIDandSecret(ctx, *authorizationRequest.ClientUID, *authorizationRequest.ClientSecret).
							Return(oauthApplication, nil),
						oauthValidator.EXPECT().ValidateAuthorizationGrant(ctx, authorizationRequest, oauthApplication).Return(nil),
						modelFormatterMock.EXPECT().AccessTokenFromAuthorizationRequest(authorizationRequest, oauthApplication).Return(oauthAccessTokenInsertable),
						oauthAccessTokenModel.EXPECT().Create(ctx, oauthAccessTokenInsertable).Return(oauthAccessToken.ID, nil),
						oauthAccessTokenModel.EXPECT().One(ctx, oauthAccessToken.ID).Return(entity.OauthAccessToken{}, errors.New("unexpected error")),
						oauthFormatterMock.EXPECT().AccessToken(authorizationRequest, oauthAccessToken).Times(0),
					)

					expectedError := &entity.Error{
						HttpStatus: http.StatusInternalServerError,
						Errors:     eobject.Wrap(eobject.InternalServerError(ctx)),
					}

					authorizationService := service.Authorization(oauthApplicationModel, oauthAccessTokenModel, oauthAccessGrantModel, modelFormatterMock, oauthValidator, oauthFormatterMock)
					results, err := authorizationService.Grantor(ctx, authorizationRequest)
					assert.NotNil(t, err)
					assert.Equal(t, expectedError, err)
					assert.Equal(t, entity.OauthAccessTokenJSON{}, results)
				})
			})
		})

		t.Run("Given context and authorization request with a response type of code", func(t *testing.T) {

		})

		t.Run("Given context and authorization request with a response type neither of code or token", func(t *testing.T) {
			t.Run("Return error 422", func(t *testing.T) {
				oauthApplicationModel := mock.NewMockOauthApplicationModel(mockCtrl)
				oauthAccessTokenModel := mock.NewMockOauthAccessTokenModel(mockCtrl)
				oauthAccessGrantModel := mock.NewMockOauthAccessGrantModel(mockCtrl)
				oauthValidator := mock.NewMockOauthValidator(mockCtrl)
				modelFormatterMock := mock.NewMockModelFormater(mockCtrl)
				oauthFormatterMock := mock.NewMockOauthFormatter(mockCtrl)

				ctx := context.WithValue(context.Background(), "track_id", uuid.New().String())

				authorizationRequest := entity.AuthorizationRequestJSON{
					ResponseType:    util.StringToPointer("others"),
					ResourceOwnerID: util.IntToPointer(1),
					ClientUID:       util.StringToPointer(aurelia.Hash("x", "y")),
					ClientSecret:    util.StringToPointer(aurelia.Hash("z", "a")),
					RedirectURI:     util.StringToPointer("http://github.com"),
					Scopes:          util.StringToPointer("public users"),
				}

				expectedError := &entity.Error{
					HttpStatus: http.StatusUnprocessableEntity,
					Errors:     eobject.Wrap(eobject.ValidationError("response_type is invalid. Should be either `token` or `code`.")),
				}

				authorizationService := service.Authorization(oauthApplicationModel, oauthAccessTokenModel, oauthAccessGrantModel, modelFormatterMock, oauthValidator, oauthFormatterMock)
				results, err := authorizationService.Grantor(ctx, authorizationRequest)
				assert.NotNil(t, err)
				assert.Equal(t, expectedError, err)
				assert.Equal(t, nil, results)
			})
		})
	})
}
