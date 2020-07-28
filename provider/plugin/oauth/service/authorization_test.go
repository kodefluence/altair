package service_test

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/codefluence-x/altair/provider/plugin/oauth/entity"
	"github.com/codefluence-x/altair/provider/plugin/oauth/eobject"

	"github.com/codefluence-x/altair/provider/plugin/oauth/formatter"
	"github.com/codefluence-x/altair/provider/plugin/oauth/mock"
	"github.com/codefluence-x/altair/provider/plugin/oauth/service"

	"github.com/codefluence-x/altair/util"
	"github.com/codefluence-x/aurelia"
	"github.com/go-sql-driver/mysql"
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
				modelFormatter := formatter.NewModel(time.Hour*4, time.Hour*2)
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
					OwnerType: "confidential",
					Description: sql.NullString{
						String: "Application 01",
						Valid:  true,
					},
					Scopes: sql.NullString{
						String: "public users",
						Valid:  true,
					},
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
					Scopes: sql.NullString{
						String: *authorizationRequest.Scopes,
						Valid:  true,
					},
					ExpiresIn: time.Now().Add(time.Hour * 4),
					CreatedAt: time.Now(),
				}

				oauthAccessTokenInsertable := modelFormatter.AccessTokenFromAuthorizationRequest(authorizationRequest, oauthApplication)
				oauthAccessTokenJSON := oauthFormatter.AccessToken(oauthAccessToken, *authorizationRequest.RedirectURI)

				gomock.InOrder(
					oauthApplicationModel.EXPECT().
						OneByUIDandSecret(ctx, *authorizationRequest.ClientUID, *authorizationRequest.ClientSecret).
						Return(oauthApplication, nil),
					oauthValidator.EXPECT().ValidateAuthorizationGrant(ctx, authorizationRequest, oauthApplication).Return(nil),
					modelFormatterMock.EXPECT().AccessTokenFromAuthorizationRequest(authorizationRequest, oauthApplication).Return(oauthAccessTokenInsertable),
					oauthAccessTokenModel.EXPECT().Create(ctx, oauthAccessTokenInsertable).Return(oauthAccessToken.ID, nil),
					oauthAccessTokenModel.EXPECT().One(ctx, oauthAccessToken.ID).Return(oauthAccessToken, nil),
					oauthFormatterMock.EXPECT().AccessToken(oauthAccessToken, *authorizationRequest.RedirectURI).Return(oauthAccessTokenJSON),
				)

				authorizationService := service.NewAuthorization(oauthApplicationModel, oauthAccessTokenModel, oauthAccessGrantModel, modelFormatterMock, oauthValidator, oauthFormatterMock)
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

						authorizationService := service.NewAuthorization(oauthApplicationModel, oauthAccessTokenModel, oauthAccessGrantModel, modelFormatterMock, oauthValidator, oauthFormatterMock)
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

						authorizationService := service.NewAuthorization(oauthApplicationModel, oauthAccessTokenModel, oauthAccessGrantModel, modelFormatterMock, oauthValidator, oauthFormatterMock)
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
						OwnerType: "confidential",
						Description: sql.NullString{
							String: "Application 01",
							Valid:  true,
						},
						Scopes: sql.NullString{
							String: "public users",
							Valid:  true,
						},
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

					authorizationService := service.NewAuthorization(oauthApplicationModel, oauthAccessTokenModel, oauthAccessGrantModel, modelFormatterMock, oauthValidator, oauthFormatterMock)
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
					modelFormatter := formatter.NewModel(time.Hour*4, time.Hour*2)
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
						OwnerType: "confidential",
						Description: sql.NullString{
							String: "Application 01",
							Valid:  true,
						},
						Scopes: sql.NullString{
							String: "public users",
							Valid:  true,
						},
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

					authorizationService := service.NewAuthorization(oauthApplicationModel, oauthAccessTokenModel, oauthAccessGrantModel, modelFormatterMock, oauthValidator, oauthFormatterMock)
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
					modelFormatter := formatter.NewModel(time.Hour*4, time.Hour*2)
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
						OwnerType: "confidential",
						Description: sql.NullString{
							String: "Application 01",
							Valid:  true,
						},
						Scopes: sql.NullString{
							String: "public users",
							Valid:  true,
						},
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
						Scopes: sql.NullString{
							String: *authorizationRequest.Scopes,
							Valid:  true,
						},
						ExpiresIn: time.Now().Add(time.Hour * 4),
						CreatedAt: time.Now(),
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
						oauthFormatterMock.EXPECT().AccessToken(oauthAccessToken, *authorizationRequest.RedirectURI).Times(0),
					)

					expectedError := &entity.Error{
						HttpStatus: http.StatusInternalServerError,
						Errors:     eobject.Wrap(eobject.InternalServerError(ctx)),
					}

					authorizationService := service.NewAuthorization(oauthApplicationModel, oauthAccessTokenModel, oauthAccessGrantModel, modelFormatterMock, oauthValidator, oauthFormatterMock)
					results, err := authorizationService.Grantor(ctx, authorizationRequest)
					assert.NotNil(t, err)
					assert.Equal(t, expectedError, err)
					assert.Equal(t, entity.OauthAccessTokenJSON{}, results)
				})
			})
		})

		t.Run("Given context and authorization request with a response type of code", func(t *testing.T) {
			t.Run("Return entity.OauthAccessGrantJSON and nil", func(t *testing.T) {
				oauthApplicationModel := mock.NewMockOauthApplicationModel(mockCtrl)
				oauthAccessTokenModel := mock.NewMockOauthAccessTokenModel(mockCtrl)
				oauthAccessGrantModel := mock.NewMockOauthAccessGrantModel(mockCtrl)
				oauthValidator := mock.NewMockOauthValidator(mockCtrl)
				modelFormatter := formatter.NewModel(time.Hour*4, time.Hour*2)
				modelFormatterMock := mock.NewMockModelFormater(mockCtrl)
				oauthFormatter := formatter.Oauth()
				oauthFormatterMock := mock.NewMockOauthFormatter(mockCtrl)

				ctx := context.WithValue(context.Background(), "track_id", uuid.New().String())

				authorizationRequest := entity.AuthorizationRequestJSON{
					ResponseType:    util.StringToPointer("code"),
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
					OwnerType: "public",
					Description: sql.NullString{
						String: "Application 01",
						Valid:  true,
					},
					Scopes: sql.NullString{
						String: "public users",
						Valid:  true,
					},
					ClientUID:    *authorizationRequest.ClientUID,
					ClientSecret: *authorizationRequest.ClientSecret,
					CreatedAt:    time.Now().Add(-time.Hour * 4),
					UpdatedAt:    time.Now(),
				}

				oauthAccessGrant := entity.OauthAccessGrant{
					ID:                 1,
					OauthApplicationID: oauthApplication.ID,
					ResourceOwnerID:    *authorizationRequest.ResourceOwnerID,
					Code:               util.SHA1(),
					ExpiresIn:          time.Now().Add(time.Hour * 4),
					CreatedAt:          time.Now(),
					RedirectURI: sql.NullString{
						String: *authorizationRequest.RedirectURI,
						Valid:  true,
					},
					Scopes: sql.NullString{
						String: *authorizationRequest.Scopes,
						Valid:  true,
					},
					RevokedAT: mysql.NullTime{
						Time:  time.Now(),
						Valid: false,
					},
				}

				oauthAccessGrantInsertable := modelFormatter.AccessGrantFromAuthorizationRequest(authorizationRequest, oauthApplication)
				oauthAccessGrantJSON := oauthFormatter.AccessGrant(oauthAccessGrant)

				gomock.InOrder(
					oauthApplicationModel.EXPECT().
						OneByUIDandSecret(ctx, *authorizationRequest.ClientUID, *authorizationRequest.ClientSecret).
						Return(oauthApplication, nil),
					oauthValidator.EXPECT().ValidateAuthorizationGrant(ctx, authorizationRequest, oauthApplication).Return(nil),
					modelFormatterMock.EXPECT().AccessGrantFromAuthorizationRequest(authorizationRequest, oauthApplication).Return(oauthAccessGrantInsertable),
					oauthAccessGrantModel.EXPECT().Create(ctx, oauthAccessGrantInsertable).Return(oauthAccessGrant.ID, nil),
					oauthAccessGrantModel.EXPECT().One(ctx, oauthAccessGrant.ID).Return(oauthAccessGrant, nil),
					oauthFormatterMock.EXPECT().AccessGrant(oauthAccessGrant).Return(oauthAccessGrantJSON),
				)

				authorizationService := service.NewAuthorization(oauthApplicationModel, oauthAccessTokenModel, oauthAccessGrantModel, modelFormatterMock, oauthValidator, oauthFormatterMock)
				results, err := authorizationService.Grantor(ctx, authorizationRequest)
				assert.Nil(t, err)
				assert.Equal(t, oauthAccessGrantJSON, results)
			})

			t.Run("Oauth authorization grant validation failed", func(t *testing.T) {
				t.Run("Return entity.OauthAccessGrantJSON and nil", func(t *testing.T) {
					oauthApplicationModel := mock.NewMockOauthApplicationModel(mockCtrl)
					oauthAccessTokenModel := mock.NewMockOauthAccessTokenModel(mockCtrl)
					oauthAccessGrantModel := mock.NewMockOauthAccessGrantModel(mockCtrl)
					oauthValidator := mock.NewMockOauthValidator(mockCtrl)
					modelFormatterMock := mock.NewMockModelFormater(mockCtrl)
					oauthFormatterMock := mock.NewMockOauthFormatter(mockCtrl)

					ctx := context.WithValue(context.Background(), "track_id", uuid.New().String())

					authorizationRequest := entity.AuthorizationRequestJSON{
						ResponseType:    util.StringToPointer("code"),
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
						OwnerType: "public",
						Description: sql.NullString{
							String: "Application 01",
							Valid:  true,
						},
						Scopes: sql.NullString{
							String: "public users",
							Valid:  true,
						},
						ClientUID:    *authorizationRequest.ClientUID,
						ClientSecret: *authorizationRequest.ClientSecret,
						CreatedAt:    time.Now().Add(-time.Hour * 4),
						UpdatedAt:    time.Now(),
					}

					expectedError := &entity.Error{
						HttpStatus: http.StatusUnprocessableEntity,
						Errors:     eobject.Wrap(eobject.ValidationError("object `owner_type` must be either of `confidential` or `public`")),
					}

					gomock.InOrder(
						oauthApplicationModel.EXPECT().
							OneByUIDandSecret(ctx, *authorizationRequest.ClientUID, *authorizationRequest.ClientSecret).
							Return(oauthApplication, nil),
						oauthValidator.EXPECT().ValidateAuthorizationGrant(ctx, authorizationRequest, oauthApplication).Return(expectedError),
						modelFormatterMock.EXPECT().AccessGrantFromAuthorizationRequest(gomock.Any(), gomock.Any()).Times(0),
						oauthAccessGrantModel.EXPECT().Create(gomock.Any(), gomock.Any()).Times(0),
						oauthAccessGrantModel.EXPECT().One(gomock.Any(), gomock.Any()).Times(0),
						oauthFormatterMock.EXPECT().AccessGrant(gomock.Any()).Times(0),
					)

					authorizationService := service.NewAuthorization(oauthApplicationModel, oauthAccessTokenModel, oauthAccessGrantModel, modelFormatterMock, oauthValidator, oauthFormatterMock)
					results, err := authorizationService.Grantor(ctx, authorizationRequest)
					assert.NotNil(t, err)
					assert.Equal(t, expectedError, err)
					assert.Equal(t, entity.OauthAccessGrantJSON{}, results)
				})
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
							ResponseType:    util.StringToPointer("code"),
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
							modelFormatterMock.EXPECT().AccessGrantFromAuthorizationRequest(gomock.Any(), gomock.Any()).Times(0),
							oauthAccessGrantModel.EXPECT().Create(gomock.Any(), gomock.Any()).Times(0),
							oauthAccessGrantModel.EXPECT().One(gomock.Any(), gomock.Any()).Times(0),
							oauthFormatterMock.EXPECT().AccessGrant(gomock.Any()).Times(0),
						)

						authorizationService := service.NewAuthorization(oauthApplicationModel, oauthAccessTokenModel, oauthAccessGrantModel, modelFormatterMock, oauthValidator, oauthFormatterMock)
						results, err := authorizationService.Grantor(ctx, authorizationRequest)
						assert.NotNil(t, err)
						assert.Equal(t, expectedError, err)
						assert.Equal(t, entity.OauthAccessGrantJSON{}, results)
					})
				})
			})

			t.Run("Create oauth access grants failed", func(t *testing.T) {
				t.Run("Return error 500", func(t *testing.T) {
					oauthApplicationModel := mock.NewMockOauthApplicationModel(mockCtrl)
					oauthAccessTokenModel := mock.NewMockOauthAccessTokenModel(mockCtrl)
					oauthAccessGrantModel := mock.NewMockOauthAccessGrantModel(mockCtrl)
					oauthValidator := mock.NewMockOauthValidator(mockCtrl)
					modelFormatter := formatter.NewModel(time.Hour*4, time.Hour*2)
					modelFormatterMock := mock.NewMockModelFormater(mockCtrl)
					oauthFormatterMock := mock.NewMockOauthFormatter(mockCtrl)

					ctx := context.WithValue(context.Background(), "track_id", uuid.New().String())

					authorizationRequest := entity.AuthorizationRequestJSON{
						ResponseType:    util.StringToPointer("code"),
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
						OwnerType: "public",
						Description: sql.NullString{
							String: "Application 01",
							Valid:  true,
						},
						Scopes: sql.NullString{
							String: "public users",
							Valid:  true,
						},
						ClientUID:    *authorizationRequest.ClientUID,
						ClientSecret: *authorizationRequest.ClientSecret,
						CreatedAt:    time.Now().Add(-time.Hour * 4),
						UpdatedAt:    time.Now(),
					}

					oauthAccessGrantInsertable := modelFormatter.AccessGrantFromAuthorizationRequest(authorizationRequest, oauthApplication)

					gomock.InOrder(
						oauthApplicationModel.EXPECT().
							OneByUIDandSecret(ctx, *authorizationRequest.ClientUID, *authorizationRequest.ClientSecret).
							Return(oauthApplication, nil),
						oauthValidator.EXPECT().ValidateAuthorizationGrant(ctx, authorizationRequest, oauthApplication).Return(nil),
						modelFormatterMock.EXPECT().AccessGrantFromAuthorizationRequest(authorizationRequest, oauthApplication).Return(oauthAccessGrantInsertable),
						oauthAccessGrantModel.EXPECT().Create(ctx, oauthAccessGrantInsertable).Return(0, errors.New("unexpected error")),
						oauthAccessGrantModel.EXPECT().One(ctx, gomock.Any()).Times(0),
						oauthFormatterMock.EXPECT().AccessGrant(gomock.Any()).Times(0),
					)

					expectedErr := &entity.Error{
						HttpStatus: http.StatusInternalServerError,
						Errors:     eobject.Wrap(eobject.InternalServerError(ctx)),
					}

					authorizationService := service.NewAuthorization(oauthApplicationModel, oauthAccessTokenModel, oauthAccessGrantModel, modelFormatterMock, oauthValidator, oauthFormatterMock)
					results, err := authorizationService.Grantor(ctx, authorizationRequest)
					assert.NotNil(t, err)
					assert.Equal(t, expectedErr, err)
					assert.Equal(t, entity.OauthAccessGrantJSON{}, results)
				})
			})

			t.Run("Finding newly created oauth access grants failed", func(t *testing.T) {
				t.Run("Return error 500", func(t *testing.T) {
					oauthApplicationModel := mock.NewMockOauthApplicationModel(mockCtrl)
					oauthAccessTokenModel := mock.NewMockOauthAccessTokenModel(mockCtrl)
					oauthAccessGrantModel := mock.NewMockOauthAccessGrantModel(mockCtrl)
					oauthValidator := mock.NewMockOauthValidator(mockCtrl)
					modelFormatter := formatter.NewModel(time.Hour*4, time.Hour*2)
					modelFormatterMock := mock.NewMockModelFormater(mockCtrl)
					oauthFormatterMock := mock.NewMockOauthFormatter(mockCtrl)

					ctx := context.WithValue(context.Background(), "track_id", uuid.New().String())

					authorizationRequest := entity.AuthorizationRequestJSON{
						ResponseType:    util.StringToPointer("code"),
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
						OwnerType: "public",
						Description: sql.NullString{
							String: "Application 01",
							Valid:  true,
						},
						Scopes: sql.NullString{
							String: "public users",
							Valid:  true,
						},
						ClientUID:    *authorizationRequest.ClientUID,
						ClientSecret: *authorizationRequest.ClientSecret,
						CreatedAt:    time.Now().Add(-time.Hour * 4),
						UpdatedAt:    time.Now(),
					}

					oauthAccessGrant := entity.OauthAccessGrant{
						ID:                 1,
						OauthApplicationID: oauthApplication.ID,
						ResourceOwnerID:    *authorizationRequest.ResourceOwnerID,
						Code:               util.SHA1(),
						CreatedAt:          time.Now(),
						ExpiresIn:          time.Now().Add(time.Hour * 4),
						RedirectURI: sql.NullString{
							String: *authorizationRequest.RedirectURI,
							Valid:  true,
						},
						Scopes: sql.NullString{
							String: *authorizationRequest.Scopes,
							Valid:  true,
						},
						RevokedAT: mysql.NullTime{
							Time:  time.Now(),
							Valid: false,
						},
					}

					oauthAccessGrantInsertable := modelFormatter.AccessGrantFromAuthorizationRequest(authorizationRequest, oauthApplication)

					gomock.InOrder(
						oauthApplicationModel.EXPECT().
							OneByUIDandSecret(ctx, *authorizationRequest.ClientUID, *authorizationRequest.ClientSecret).
							Return(oauthApplication, nil),
						oauthValidator.EXPECT().ValidateAuthorizationGrant(ctx, authorizationRequest, oauthApplication).Return(nil),
						modelFormatterMock.EXPECT().AccessGrantFromAuthorizationRequest(authorizationRequest, oauthApplication).Return(oauthAccessGrantInsertable),
						oauthAccessGrantModel.EXPECT().Create(ctx, oauthAccessGrantInsertable).Return(oauthAccessGrant.ID, nil),
						oauthAccessGrantModel.EXPECT().One(ctx, oauthAccessGrant.ID).Return(entity.OauthAccessGrant{}, errors.New("unexpected error")),
						oauthFormatterMock.EXPECT().AccessGrant(gomock.Any()).Times(0),
					)

					expectedErr := &entity.Error{
						HttpStatus: http.StatusInternalServerError,
						Errors:     eobject.Wrap(eobject.InternalServerError(ctx)),
					}

					authorizationService := service.NewAuthorization(oauthApplicationModel, oauthAccessTokenModel, oauthAccessGrantModel, modelFormatterMock, oauthValidator, oauthFormatterMock)
					results, err := authorizationService.Grantor(ctx, authorizationRequest)
					assert.NotNil(t, err)
					assert.Equal(t, expectedErr, err)
					assert.Equal(t, entity.OauthAccessGrantJSON{}, results)
				})
			})
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

				authorizationService := service.NewAuthorization(oauthApplicationModel, oauthAccessTokenModel, oauthAccessGrantModel, modelFormatterMock, oauthValidator, oauthFormatterMock)
				results, err := authorizationService.Grantor(ctx, authorizationRequest)
				assert.NotNil(t, err)
				assert.Equal(t, expectedError, err)
				assert.Equal(t, nil, results)
			})
		})

		t.Run("Given context and authorization request with a nil client uid", func(t *testing.T) {
			t.Run("Return error 422", func(t *testing.T) {
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
					ClientUID:       nil,
					ClientSecret:    util.StringToPointer(aurelia.Hash("z", "a")),
					RedirectURI:     util.StringToPointer("http://github.com"),
					Scopes:          util.StringToPointer("public users"),
				}

				expectedError := &entity.Error{
					HttpStatus: http.StatusUnprocessableEntity,
					Errors:     eobject.Wrap(eobject.ValidationError("client_uid cannot be empty")),
				}

				authorizationService := service.NewAuthorization(oauthApplicationModel, oauthAccessTokenModel, oauthAccessGrantModel, modelFormatterMock, oauthValidator, oauthFormatterMock)
				results, err := authorizationService.Grantor(ctx, authorizationRequest)
				assert.NotNil(t, err)
				assert.Equal(t, expectedError, err)
				assert.Equal(t, entity.OauthAccessTokenJSON{}, results)
			})
		})

		t.Run("Given context and authorization request with a nil client secret", func(t *testing.T) {
			t.Run("Return error 422", func(t *testing.T) {
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
					ClientUID:       util.StringToPointer(aurelia.Hash("z", "a")),
					ClientSecret:    nil,
					RedirectURI:     util.StringToPointer("http://github.com"),
					Scopes:          util.StringToPointer("public users"),
				}

				expectedError := &entity.Error{
					HttpStatus: http.StatusUnprocessableEntity,
					Errors:     eobject.Wrap(eobject.ValidationError("client_secret cannot be empty")),
				}

				authorizationService := service.NewAuthorization(oauthApplicationModel, oauthAccessTokenModel, oauthAccessGrantModel, modelFormatterMock, oauthValidator, oauthFormatterMock)
				results, err := authorizationService.Grantor(ctx, authorizationRequest)
				assert.NotNil(t, err)
				assert.Equal(t, expectedError, err)
				assert.Equal(t, entity.OauthAccessTokenJSON{}, results)
			})
		})

		t.Run("Given context and authorization request with nil response type", func(t *testing.T) {
			t.Run("Return error 422", func(t *testing.T) {
				oauthApplicationModel := mock.NewMockOauthApplicationModel(mockCtrl)
				oauthAccessTokenModel := mock.NewMockOauthAccessTokenModel(mockCtrl)
				oauthAccessGrantModel := mock.NewMockOauthAccessGrantModel(mockCtrl)
				oauthValidator := mock.NewMockOauthValidator(mockCtrl)
				modelFormatterMock := mock.NewMockModelFormater(mockCtrl)
				oauthFormatterMock := mock.NewMockOauthFormatter(mockCtrl)

				ctx := context.WithValue(context.Background(), "track_id", uuid.New().String())

				authorizationRequest := entity.AuthorizationRequestJSON{
					ResponseType:    nil,
					ResourceOwnerID: util.IntToPointer(1),
					ClientUID:       util.StringToPointer(aurelia.Hash("x", "y")),
					ClientSecret:    util.StringToPointer(aurelia.Hash("z", "a")),
					RedirectURI:     util.StringToPointer("http://github.com"),
					Scopes:          util.StringToPointer("public users"),
				}

				expectedError := &entity.Error{
					HttpStatus: http.StatusUnprocessableEntity,
					Errors:     eobject.Wrap(eobject.ValidationError("response_type cannot be empty")),
				}

				authorizationService := service.NewAuthorization(oauthApplicationModel, oauthAccessTokenModel, oauthAccessGrantModel, modelFormatterMock, oauthValidator, oauthFormatterMock)
				results, err := authorizationService.Grantor(ctx, authorizationRequest)
				assert.NotNil(t, err)
				assert.Equal(t, expectedError, err)
				assert.Equal(t, nil, results)
			})
		})
	})

	t.Run("Token", func(t *testing.T) {
		t.Run("Given context and access token request", func(t *testing.T) {
			t.Run("When access token request valid and there is no error in database side", func(t *testing.T) {
				t.Run("Then it will return access token response", func(t *testing.T) {
					oauthApplicationModel := mock.NewMockOauthApplicationModel(mockCtrl)
					oauthAccessTokenModel := mock.NewMockOauthAccessTokenModel(mockCtrl)
					oauthAccessGrantModel := mock.NewMockOauthAccessGrantModel(mockCtrl)
					oauthValidator := mock.NewMockOauthValidator(mockCtrl)
					modelFormatter := formatter.NewModel(time.Hour*4, time.Hour*2)
					modelFormatterMock := mock.NewMockModelFormater(mockCtrl)
					oauthFormatter := formatter.Oauth()
					oauthFormatterMock := mock.NewMockOauthFormatter(mockCtrl)

					ctx := context.Background()

					accessTokenRequest := entity.AccessTokenRequestJSON{
						ClientSecret: util.StringToPointer("client_secret"),
						ClientUID:    util.StringToPointer("client_uid"),
						Code:         util.StringToPointer("abcdef_123456"),
						GrantType:    util.StringToPointer("authorization_code"),
						RedirectURI:  util.StringToPointer("http://localhost:8000/oauth_redirect"),
					}

					oauthApplication := entity.OauthApplication{
						ID: 1,
						OwnerID: sql.NullInt64{
							Int64: 1,
							Valid: true,
						},
						OwnerType: "confidential",
						Description: sql.NullString{
							String: "Application 01",
							Valid:  true,
						},
						Scopes: sql.NullString{
							String: "public users",
							Valid:  true,
						},
						ClientUID:    *accessTokenRequest.ClientUID,
						ClientSecret: *accessTokenRequest.ClientSecret,
						CreatedAt:    time.Now().Add(-time.Hour * 4),
						UpdatedAt:    time.Now(),
					}

					oauthAccessGrant := entity.OauthAccessGrant{
						ID:                 1,
						Code:               *accessTokenRequest.Code,
						CreatedAt:          time.Now().Add(-time.Hour),
						ExpiresIn:          time.Now().Add(time.Hour),
						OauthApplicationID: oauthApplication.ID,
						RedirectURI: sql.NullString{
							String: *accessTokenRequest.RedirectURI,
							Valid:  true,
						},
						ResourceOwnerID: 1,
						RevokedAT: mysql.NullTime{
							Valid: false,
						},
						Scopes: sql.NullString{
							String: "user store",
							Valid:  true,
						},
					}

					oauthAccessToken := entity.OauthAccessToken{
						ID:                 1,
						OauthApplicationID: oauthApplication.ID,
						ResourceOwnerID:    oauthAccessGrant.ResourceOwnerID,
						Token:              aurelia.Hash("x", "y"),
						Scopes: sql.NullString{
							String: oauthAccessGrant.Scopes.String,
							Valid:  true,
						},
						ExpiresIn: time.Now().Add(time.Hour * 4),
						CreatedAt: time.Now(),
					}

					oauthAccessTokenInsertable := modelFormatter.AccessTokenFromOauthAccessGrant(oauthAccessGrant, oauthApplication)
					oauthAccessTokenJSON := oauthFormatter.AccessToken(oauthAccessToken, *accessTokenRequest.RedirectURI)

					gomock.InOrder(
						oauthApplicationModel.EXPECT().
							OneByUIDandSecret(ctx, *accessTokenRequest.ClientUID, *accessTokenRequest.ClientSecret).
							Return(oauthApplication, nil),
						oauthValidator.EXPECT().ValidateTokenGrant(ctx, accessTokenRequest).Return(nil),
						oauthAccessGrantModel.EXPECT().OneByCode(ctx, *accessTokenRequest.Code).Return(oauthAccessGrant, nil),
						modelFormatterMock.EXPECT().AccessTokenFromOauthAccessGrant(oauthAccessGrant, oauthApplication).Return(oauthAccessTokenInsertable),
						oauthAccessTokenModel.EXPECT().Create(ctx, oauthAccessTokenInsertable).Return(1, nil),
						oauthAccessTokenModel.EXPECT().One(ctx, 1).Return(oauthAccessToken, nil),
						oauthFormatterMock.EXPECT().AccessToken(oauthAccessToken, oauthAccessGrant.RedirectURI.String).Return(oauthAccessTokenJSON),
					)

					authorizationService := service.NewAuthorization(oauthApplicationModel, oauthAccessTokenModel, oauthAccessGrantModel, modelFormatterMock, oauthValidator, oauthFormatterMock)
					oauthAccessTokenOutput, err := authorizationService.Token(ctx, accessTokenRequest)

					assert.Nil(t, err)
					assert.Equal(t, oauthAccessTokenJSON, oauthAccessTokenOutput)
				})
			})
		})

		t.Run("When application client and secret is not valid", func(t *testing.T) {
			t.Run("Then it will return unprocessable entity error", func(t *testing.T) {
				oauthApplicationModel := mock.NewMockOauthApplicationModel(mockCtrl)
				oauthAccessTokenModel := mock.NewMockOauthAccessTokenModel(mockCtrl)
				oauthAccessGrantModel := mock.NewMockOauthAccessGrantModel(mockCtrl)
				oauthValidator := mock.NewMockOauthValidator(mockCtrl)
				modelFormatterMock := mock.NewMockModelFormater(mockCtrl)
				oauthFormatterMock := mock.NewMockOauthFormatter(mockCtrl)

				ctx := context.Background()

				accessTokenRequest := entity.AccessTokenRequestJSON{
					ClientSecret: util.StringToPointer("client_secret"),
					ClientUID:    util.StringToPointer("client_uid"),
					Code:         util.StringToPointer("abcdef_123456"),
					GrantType:    util.StringToPointer("authorization_code"),
					RedirectURI:  util.StringToPointer("http://localhost:8000/oauth_redirect"),
				}

				expectedError := &entity.Error{
					HttpStatus: http.StatusInternalServerError,
					Errors:     eobject.Wrap(eobject.InternalServerError(ctx)),
				}

				gomock.InOrder(
					oauthApplicationModel.EXPECT().
						OneByUIDandSecret(ctx, *accessTokenRequest.ClientUID, *accessTokenRequest.ClientSecret).
						Return(entity.OauthApplication{}, errors.New("unexpected error")),
					oauthValidator.EXPECT().ValidateTokenGrant(gomock.Any(), gomock.Any()).Times(0),
					oauthAccessGrantModel.EXPECT().OneByCode(gomock.Any(), gomock.Any()).Times(0),
					modelFormatterMock.EXPECT().AccessTokenFromOauthAccessGrant(gomock.Any(), gomock.Any()).Times(0),
					oauthAccessTokenModel.EXPECT().Create(gomock.Any(), gomock.Any()).Times(0),
					oauthAccessTokenModel.EXPECT().One(gomock.Any(), gomock.Any()).Times(0),
					oauthFormatterMock.EXPECT().AccessToken(gomock.Any(), gomock.Any()).Times(0),
				)

				authorizationService := service.NewAuthorization(oauthApplicationModel, oauthAccessTokenModel, oauthAccessGrantModel, modelFormatterMock, oauthValidator, oauthFormatterMock)
				oauthAccessTokenOutput, err := authorizationService.Token(ctx, accessTokenRequest)

				assert.NotNil(t, err)
				assert.Equal(t, expectedError, err)
				assert.Equal(t, entity.OauthAccessTokenJSON{}, oauthAccessTokenOutput)
			})
		})

		t.Run("When access token request is not valid", func(t *testing.T) {
			t.Run("Then it will return unprocessable entity error", func(t *testing.T) {
				oauthApplicationModel := mock.NewMockOauthApplicationModel(mockCtrl)
				oauthAccessTokenModel := mock.NewMockOauthAccessTokenModel(mockCtrl)
				oauthAccessGrantModel := mock.NewMockOauthAccessGrantModel(mockCtrl)
				oauthValidator := mock.NewMockOauthValidator(mockCtrl)
				modelFormatterMock := mock.NewMockModelFormater(mockCtrl)
				oauthFormatterMock := mock.NewMockOauthFormatter(mockCtrl)

				ctx := context.Background()

				accessTokenRequest := entity.AccessTokenRequestJSON{
					ClientSecret: util.StringToPointer("client_secret"),
					ClientUID:    util.StringToPointer("client_uid"),
					Code:         util.StringToPointer("abcdef_123456"),
					GrantType:    util.StringToPointer("authorization_code"),
					RedirectURI:  util.StringToPointer("http://localhost:8000/oauth_redirect"),
				}

				oauthApplication := entity.OauthApplication{
					ID: 1,
					OwnerID: sql.NullInt64{
						Int64: 1,
						Valid: true,
					},
					OwnerType: "confidential",
					Description: sql.NullString{
						String: "Application 01",
						Valid:  true,
					},
					Scopes: sql.NullString{
						String: "public users",
						Valid:  true,
					},
					ClientUID:    *accessTokenRequest.ClientUID,
					ClientSecret: *accessTokenRequest.ClientSecret,
					CreatedAt:    time.Now().Add(-time.Hour * 4),
					UpdatedAt:    time.Now(),
				}

				expectedError := &entity.Error{
					HttpStatus: http.StatusUnprocessableEntity,
					Errors: eobject.Wrap(
						eobject.ValidationError(`grant_type can't be empty`),
					),
				}

				gomock.InOrder(
					oauthApplicationModel.EXPECT().
						OneByUIDandSecret(ctx, *accessTokenRequest.ClientUID, *accessTokenRequest.ClientSecret).
						Return(oauthApplication, nil),
					oauthValidator.EXPECT().ValidateTokenGrant(ctx, accessTokenRequest).Return(expectedError),
					oauthAccessGrantModel.EXPECT().OneByCode(gomock.Any(), gomock.Any()).Times(0),
					modelFormatterMock.EXPECT().AccessTokenFromOauthAccessGrant(gomock.Any(), gomock.Any()).Times(0),
					oauthAccessTokenModel.EXPECT().Create(gomock.Any(), gomock.Any()).Times(0),
					oauthAccessTokenModel.EXPECT().One(gomock.Any(), gomock.Any()).Times(0),
					oauthFormatterMock.EXPECT().AccessToken(gomock.Any(), gomock.Any()).Times(0),
				)

				authorizationService := service.NewAuthorization(oauthApplicationModel, oauthAccessTokenModel, oauthAccessGrantModel, modelFormatterMock, oauthValidator, oauthFormatterMock)
				oauthAccessTokenOutput, err := authorizationService.Token(ctx, accessTokenRequest)

				assert.NotNil(t, err)
				assert.Equal(t, expectedError, err)
				assert.Equal(t, entity.OauthAccessTokenJSON{}, oauthAccessTokenOutput)
			})
		})

		t.Run("When oauth access grant is not found based on authorization_code", func(t *testing.T) {
			t.Run("Then it will return not found error", func(t *testing.T) {
				oauthApplicationModel := mock.NewMockOauthApplicationModel(mockCtrl)
				oauthAccessTokenModel := mock.NewMockOauthAccessTokenModel(mockCtrl)
				oauthAccessGrantModel := mock.NewMockOauthAccessGrantModel(mockCtrl)
				oauthValidator := mock.NewMockOauthValidator(mockCtrl)
				modelFormatterMock := mock.NewMockModelFormater(mockCtrl)
				oauthFormatterMock := mock.NewMockOauthFormatter(mockCtrl)

				ctx := context.Background()

				accessTokenRequest := entity.AccessTokenRequestJSON{
					ClientSecret: util.StringToPointer("client_secret"),
					ClientUID:    util.StringToPointer("client_uid"),
					Code:         util.StringToPointer("abcdef_123456"),
					GrantType:    util.StringToPointer("authorization_code"),
					RedirectURI:  util.StringToPointer("http://localhost:8000/oauth_redirect"),
				}

				oauthApplication := entity.OauthApplication{
					ID: 1,
					OwnerID: sql.NullInt64{
						Int64: 1,
						Valid: true,
					},
					OwnerType: "confidential",
					Description: sql.NullString{
						String: "Application 01",
						Valid:  true,
					},
					Scopes: sql.NullString{
						String: "public users",
						Valid:  true,
					},
					ClientUID:    *accessTokenRequest.ClientUID,
					ClientSecret: *accessTokenRequest.ClientSecret,
					CreatedAt:    time.Now().Add(-time.Hour * 4),
					UpdatedAt:    time.Now(),
				}

				expectedError := &entity.Error{
					HttpStatus: http.StatusNotFound,
					Errors:     eobject.Wrap(eobject.NotFoundError(ctx, "authorization_code")),
				}

				gomock.InOrder(
					oauthApplicationModel.EXPECT().
						OneByUIDandSecret(ctx, *accessTokenRequest.ClientUID, *accessTokenRequest.ClientSecret).
						Return(oauthApplication, nil),
					oauthValidator.EXPECT().ValidateTokenGrant(ctx, accessTokenRequest).Return(nil),
					oauthAccessGrantModel.EXPECT().OneByCode(ctx, *accessTokenRequest.Code).Return(entity.OauthAccessGrant{}, sql.ErrNoRows),
					modelFormatterMock.EXPECT().AccessTokenFromOauthAccessGrant(gomock.Any(), gomock.Any()).Times(0),
					oauthAccessTokenModel.EXPECT().Create(gomock.Any(), gomock.Any()).Times(0),
					oauthAccessTokenModel.EXPECT().One(gomock.Any(), gomock.Any()).Times(0),
					oauthFormatterMock.EXPECT().AccessToken(gomock.Any(), gomock.Any()).Times(0),
				)

				authorizationService := service.NewAuthorization(oauthApplicationModel, oauthAccessTokenModel, oauthAccessGrantModel, modelFormatterMock, oauthValidator, oauthFormatterMock)
				oauthAccessTokenOutput, err := authorizationService.Token(ctx, accessTokenRequest)

				assert.NotNil(t, err)
				assert.Equal(t, expectedError, err)
				assert.Equal(t, entity.OauthAccessTokenJSON{}, oauthAccessTokenOutput)
			})
		})

		t.Run("When there is unexpected error when getting oauth access grant by code", func(t *testing.T) {
			t.Run("Then it will return internal server error", func(t *testing.T) {
				oauthApplicationModel := mock.NewMockOauthApplicationModel(mockCtrl)
				oauthAccessTokenModel := mock.NewMockOauthAccessTokenModel(mockCtrl)
				oauthAccessGrantModel := mock.NewMockOauthAccessGrantModel(mockCtrl)
				oauthValidator := mock.NewMockOauthValidator(mockCtrl)
				modelFormatterMock := mock.NewMockModelFormater(mockCtrl)
				oauthFormatterMock := mock.NewMockOauthFormatter(mockCtrl)

				ctx := context.Background()

				accessTokenRequest := entity.AccessTokenRequestJSON{
					ClientSecret: util.StringToPointer("client_secret"),
					ClientUID:    util.StringToPointer("client_uid"),
					Code:         util.StringToPointer("abcdef_123456"),
					GrantType:    util.StringToPointer("authorization_code"),
					RedirectURI:  util.StringToPointer("http://localhost:8000/oauth_redirect"),
				}

				oauthApplication := entity.OauthApplication{
					ID: 1,
					OwnerID: sql.NullInt64{
						Int64: 1,
						Valid: true,
					},
					OwnerType: "confidential",
					Description: sql.NullString{
						String: "Application 01",
						Valid:  true,
					},
					Scopes: sql.NullString{
						String: "public users",
						Valid:  true,
					},
					ClientUID:    *accessTokenRequest.ClientUID,
					ClientSecret: *accessTokenRequest.ClientSecret,
					CreatedAt:    time.Now().Add(-time.Hour * 4),
					UpdatedAt:    time.Now(),
				}

				expectedError := &entity.Error{
					HttpStatus: http.StatusInternalServerError,
					Errors:     eobject.Wrap(eobject.InternalServerError(ctx)),
				}

				gomock.InOrder(
					oauthApplicationModel.EXPECT().
						OneByUIDandSecret(ctx, *accessTokenRequest.ClientUID, *accessTokenRequest.ClientSecret).
						Return(oauthApplication, nil),
					oauthValidator.EXPECT().ValidateTokenGrant(ctx, accessTokenRequest).Return(nil),
					oauthAccessGrantModel.EXPECT().OneByCode(ctx, *accessTokenRequest.Code).Return(entity.OauthAccessGrant{}, errors.New("unexpected error")),
					modelFormatterMock.EXPECT().AccessTokenFromOauthAccessGrant(gomock.Any(), gomock.Any()).Times(0),
					oauthAccessTokenModel.EXPECT().Create(gomock.Any(), gomock.Any()).Times(0),
					oauthAccessTokenModel.EXPECT().One(gomock.Any(), gomock.Any()).Times(0),
					oauthFormatterMock.EXPECT().AccessToken(gomock.Any(), gomock.Any()).Times(0),
				)

				authorizationService := service.NewAuthorization(oauthApplicationModel, oauthAccessTokenModel, oauthAccessGrantModel, modelFormatterMock, oauthValidator, oauthFormatterMock)
				oauthAccessTokenOutput, err := authorizationService.Token(ctx, accessTokenRequest)

				assert.NotNil(t, err)
				assert.Equal(t, expectedError, err)
				assert.Equal(t, entity.OauthAccessTokenJSON{}, oauthAccessTokenOutput)
			})
		})

		t.Run("When redirect uri from access token request is not identic with the database one", func(t *testing.T) {
			t.Run("Then it will return forbidden error", func(t *testing.T) {
				oauthApplicationModel := mock.NewMockOauthApplicationModel(mockCtrl)
				oauthAccessTokenModel := mock.NewMockOauthAccessTokenModel(mockCtrl)
				oauthAccessGrantModel := mock.NewMockOauthAccessGrantModel(mockCtrl)
				oauthValidator := mock.NewMockOauthValidator(mockCtrl)
				modelFormatterMock := mock.NewMockModelFormater(mockCtrl)
				oauthFormatterMock := mock.NewMockOauthFormatter(mockCtrl)

				ctx := context.Background()

				accessTokenRequest := entity.AccessTokenRequestJSON{
					ClientSecret: util.StringToPointer("client_secret"),
					ClientUID:    util.StringToPointer("client_uid"),
					Code:         util.StringToPointer("abcdef_123456"),
					GrantType:    util.StringToPointer("authorization_code"),
					RedirectURI:  util.StringToPointer("http://localhost:8000/oauth_redirect"),
				}

				oauthApplication := entity.OauthApplication{
					ID: 1,
					OwnerID: sql.NullInt64{
						Int64: 1,
						Valid: true,
					},
					OwnerType: "confidential",
					Description: sql.NullString{
						String: "Application 01",
						Valid:  true,
					},
					Scopes: sql.NullString{
						String: "public users",
						Valid:  true,
					},
					ClientUID:    *accessTokenRequest.ClientUID,
					ClientSecret: *accessTokenRequest.ClientSecret,
					CreatedAt:    time.Now().Add(-time.Hour * 4),
					UpdatedAt:    time.Now(),
				}

				oauthAccessGrant := entity.OauthAccessGrant{
					ID:                 1,
					Code:               *accessTokenRequest.Code,
					CreatedAt:          time.Now().Add(-time.Hour),
					ExpiresIn:          time.Now().Add(time.Hour),
					OauthApplicationID: oauthApplication.ID,
					RedirectURI: sql.NullString{
						String: "http://localhost:8000/different_redirect_uri",
						Valid:  true,
					},
					ResourceOwnerID: 1,
					RevokedAT: mysql.NullTime{
						Valid: false,
					},
					Scopes: sql.NullString{
						String: "user store",
						Valid:  true,
					},
				}

				expectedError := &entity.Error{
					HttpStatus: http.StatusForbidden,
					Errors:     eobject.Wrap(eobject.ForbiddenError(ctx, "redirect_uri", "redirect uri is different from one that generated before")),
				}

				gomock.InOrder(
					oauthApplicationModel.EXPECT().
						OneByUIDandSecret(ctx, *accessTokenRequest.ClientUID, *accessTokenRequest.ClientSecret).
						Return(oauthApplication, nil),
					oauthValidator.EXPECT().ValidateTokenGrant(ctx, accessTokenRequest).Return(nil),
					oauthAccessGrantModel.EXPECT().OneByCode(ctx, *accessTokenRequest.Code).Return(oauthAccessGrant, nil),
					modelFormatterMock.EXPECT().AccessTokenFromOauthAccessGrant(gomock.Any(), gomock.Any()).Times(0),
					oauthAccessTokenModel.EXPECT().Create(gomock.Any(), gomock.Any()).Times(0),
					oauthAccessTokenModel.EXPECT().One(gomock.Any(), gomock.Any()).Times(0),
					oauthFormatterMock.EXPECT().AccessToken(gomock.Any(), gomock.Any()).Times(0),
				)

				authorizationService := service.NewAuthorization(oauthApplicationModel, oauthAccessTokenModel, oauthAccessGrantModel, modelFormatterMock, oauthValidator, oauthFormatterMock)
				oauthAccessTokenOutput, err := authorizationService.Token(ctx, accessTokenRequest)

				assert.NotNil(t, err)
				assert.Equal(t, expectedError, err)
				assert.Equal(t, entity.OauthAccessTokenJSON{}, oauthAccessTokenOutput)
			})
		})

		t.Run("When authorization code already expired", func(t *testing.T) {
			t.Run("Then it will return forbidden error", func(t *testing.T) {
				oauthApplicationModel := mock.NewMockOauthApplicationModel(mockCtrl)
				oauthAccessTokenModel := mock.NewMockOauthAccessTokenModel(mockCtrl)
				oauthAccessGrantModel := mock.NewMockOauthAccessGrantModel(mockCtrl)
				oauthValidator := mock.NewMockOauthValidator(mockCtrl)
				modelFormatterMock := mock.NewMockModelFormater(mockCtrl)
				oauthFormatterMock := mock.NewMockOauthFormatter(mockCtrl)

				ctx := context.Background()

				accessTokenRequest := entity.AccessTokenRequestJSON{
					ClientSecret: util.StringToPointer("client_secret"),
					ClientUID:    util.StringToPointer("client_uid"),
					Code:         util.StringToPointer("abcdef_123456"),
					GrantType:    util.StringToPointer("authorization_code"),
					RedirectURI:  util.StringToPointer("http://localhost:8000/oauth_redirect"),
				}

				oauthApplication := entity.OauthApplication{
					ID: 1,
					OwnerID: sql.NullInt64{
						Int64: 1,
						Valid: true,
					},
					OwnerType: "confidential",
					Description: sql.NullString{
						String: "Application 01",
						Valid:  true,
					},
					Scopes: sql.NullString{
						String: "public users",
						Valid:  true,
					},
					ClientUID:    *accessTokenRequest.ClientUID,
					ClientSecret: *accessTokenRequest.ClientSecret,
					CreatedAt:    time.Now().Add(-time.Hour * 4),
					UpdatedAt:    time.Now(),
				}

				oauthAccessGrant := entity.OauthAccessGrant{
					ID:                 1,
					Code:               *accessTokenRequest.Code,
					CreatedAt:          time.Now().Add(-time.Hour * 2),
					ExpiresIn:          time.Now().Add(-time.Hour),
					OauthApplicationID: oauthApplication.ID,
					RedirectURI: sql.NullString{
						String: *accessTokenRequest.RedirectURI,
						Valid:  true,
					},
					ResourceOwnerID: 1,
					RevokedAT: mysql.NullTime{
						Valid: false,
					},
					Scopes: sql.NullString{
						String: "user store",
						Valid:  true,
					},
				}

				expectedError := &entity.Error{
					HttpStatus: http.StatusForbidden,
					Errors:     eobject.Wrap(eobject.ForbiddenError(ctx, "access_token", "authorization code already expired")),
				}

				gomock.InOrder(
					oauthApplicationModel.EXPECT().
						OneByUIDandSecret(ctx, *accessTokenRequest.ClientUID, *accessTokenRequest.ClientSecret).
						Return(oauthApplication, nil),
					oauthValidator.EXPECT().ValidateTokenGrant(ctx, accessTokenRequest).Return(nil),
					oauthAccessGrantModel.EXPECT().OneByCode(ctx, *accessTokenRequest.Code).Return(oauthAccessGrant, nil),
					modelFormatterMock.EXPECT().AccessTokenFromOauthAccessGrant(gomock.Any(), gomock.Any()).Times(0),
					oauthAccessTokenModel.EXPECT().Create(gomock.Any(), gomock.Any()).Times(0),
					oauthAccessTokenModel.EXPECT().One(gomock.Any(), gomock.Any()).Times(0),
					oauthFormatterMock.EXPECT().AccessToken(gomock.Any(), gomock.Any()).Times(0),
				)

				authorizationService := service.NewAuthorization(oauthApplicationModel, oauthAccessTokenModel, oauthAccessGrantModel, modelFormatterMock, oauthValidator, oauthFormatterMock)
				oauthAccessTokenOutput, err := authorizationService.Token(ctx, accessTokenRequest)

				assert.NotNil(t, err)
				assert.Equal(t, expectedError, err)
				assert.Equal(t, entity.OauthAccessTokenJSON{}, oauthAccessTokenOutput)
			})
		})

		t.Run("When there is unexpected error when creating access token", func(t *testing.T) {
			t.Run("Then it will return internal server error", func(t *testing.T) {
				oauthApplicationModel := mock.NewMockOauthApplicationModel(mockCtrl)
				oauthAccessTokenModel := mock.NewMockOauthAccessTokenModel(mockCtrl)
				oauthAccessGrantModel := mock.NewMockOauthAccessGrantModel(mockCtrl)
				oauthValidator := mock.NewMockOauthValidator(mockCtrl)
				modelFormatter := formatter.NewModel(time.Hour*4, time.Hour*2)
				modelFormatterMock := mock.NewMockModelFormater(mockCtrl)
				oauthFormatterMock := mock.NewMockOauthFormatter(mockCtrl)

				ctx := context.Background()

				accessTokenRequest := entity.AccessTokenRequestJSON{
					ClientSecret: util.StringToPointer("client_secret"),
					ClientUID:    util.StringToPointer("client_uid"),
					Code:         util.StringToPointer("abcdef_123456"),
					GrantType:    util.StringToPointer("authorization_code"),
					RedirectURI:  util.StringToPointer("http://localhost:8000/oauth_redirect"),
				}

				oauthApplication := entity.OauthApplication{
					ID: 1,
					OwnerID: sql.NullInt64{
						Int64: 1,
						Valid: true,
					},
					OwnerType: "confidential",
					Description: sql.NullString{
						String: "Application 01",
						Valid:  true,
					},
					Scopes: sql.NullString{
						String: "public users",
						Valid:  true,
					},
					ClientUID:    *accessTokenRequest.ClientUID,
					ClientSecret: *accessTokenRequest.ClientSecret,
					CreatedAt:    time.Now().Add(-time.Hour * 4),
					UpdatedAt:    time.Now(),
				}

				oauthAccessGrant := entity.OauthAccessGrant{
					ID:                 1,
					Code:               *accessTokenRequest.Code,
					CreatedAt:          time.Now().Add(-time.Hour),
					ExpiresIn:          time.Now().Add(time.Hour),
					OauthApplicationID: oauthApplication.ID,
					RedirectURI: sql.NullString{
						String: *accessTokenRequest.RedirectURI,
						Valid:  true,
					},
					ResourceOwnerID: 1,
					RevokedAT: mysql.NullTime{
						Valid: false,
					},
					Scopes: sql.NullString{
						String: "user store",
						Valid:  true,
					},
				}

				oauthAccessTokenInsertable := modelFormatter.AccessTokenFromOauthAccessGrant(oauthAccessGrant, oauthApplication)

				expectedError := &entity.Error{
					HttpStatus: http.StatusInternalServerError,
					Errors:     eobject.Wrap(eobject.InternalServerError(ctx)),
				}

				gomock.InOrder(
					oauthApplicationModel.EXPECT().
						OneByUIDandSecret(ctx, *accessTokenRequest.ClientUID, *accessTokenRequest.ClientSecret).
						Return(oauthApplication, nil),
					oauthValidator.EXPECT().ValidateTokenGrant(ctx, accessTokenRequest).Return(nil),
					oauthAccessGrantModel.EXPECT().OneByCode(ctx, *accessTokenRequest.Code).Return(oauthAccessGrant, nil),
					modelFormatterMock.EXPECT().AccessTokenFromOauthAccessGrant(oauthAccessGrant, oauthApplication).Return(oauthAccessTokenInsertable),
					oauthAccessTokenModel.EXPECT().Create(ctx, oauthAccessTokenInsertable).Return(0, errors.New("unexpected error")),
					oauthAccessTokenModel.EXPECT().One(gomock.Any(), gomock.Any()).Times(0),
					oauthFormatterMock.EXPECT().AccessToken(gomock.Any(), gomock.Any()).Times(0),
				)

				authorizationService := service.NewAuthorization(oauthApplicationModel, oauthAccessTokenModel, oauthAccessGrantModel, modelFormatterMock, oauthValidator, oauthFormatterMock)
				oauthAccessTokenOutput, err := authorizationService.Token(ctx, accessTokenRequest)

				assert.NotNil(t, err)
				assert.Equal(t, expectedError, err)
				assert.Equal(t, entity.OauthAccessTokenJSON{}, oauthAccessTokenOutput)
			})
		})

		t.Run("When there is unexpected error when selecting access token", func(t *testing.T) {
			t.Run("Then it will return internal server error", func(t *testing.T) {
				oauthApplicationModel := mock.NewMockOauthApplicationModel(mockCtrl)
				oauthAccessTokenModel := mock.NewMockOauthAccessTokenModel(mockCtrl)
				oauthAccessGrantModel := mock.NewMockOauthAccessGrantModel(mockCtrl)
				oauthValidator := mock.NewMockOauthValidator(mockCtrl)
				modelFormatter := formatter.NewModel(time.Hour*4, time.Hour*2)
				modelFormatterMock := mock.NewMockModelFormater(mockCtrl)
				oauthFormatterMock := mock.NewMockOauthFormatter(mockCtrl)

				ctx := context.Background()

				accessTokenRequest := entity.AccessTokenRequestJSON{
					ClientSecret: util.StringToPointer("client_secret"),
					ClientUID:    util.StringToPointer("client_uid"),
					Code:         util.StringToPointer("abcdef_123456"),
					GrantType:    util.StringToPointer("authorization_code"),
					RedirectURI:  util.StringToPointer("http://localhost:8000/oauth_redirect"),
				}

				oauthApplication := entity.OauthApplication{
					ID: 1,
					OwnerID: sql.NullInt64{
						Int64: 1,
						Valid: true,
					},
					OwnerType: "confidential",
					Description: sql.NullString{
						String: "Application 01",
						Valid:  true,
					},
					Scopes: sql.NullString{
						String: "public users",
						Valid:  true,
					},
					ClientUID:    *accessTokenRequest.ClientUID,
					ClientSecret: *accessTokenRequest.ClientSecret,
					CreatedAt:    time.Now().Add(-time.Hour * 4),
					UpdatedAt:    time.Now(),
				}

				oauthAccessGrant := entity.OauthAccessGrant{
					ID:                 1,
					Code:               *accessTokenRequest.Code,
					CreatedAt:          time.Now().Add(-time.Hour),
					ExpiresIn:          time.Now().Add(time.Hour),
					OauthApplicationID: oauthApplication.ID,
					RedirectURI: sql.NullString{
						String: *accessTokenRequest.RedirectURI,
						Valid:  true,
					},
					ResourceOwnerID: 1,
					RevokedAT: mysql.NullTime{
						Valid: false,
					},
					Scopes: sql.NullString{
						String: "user store",
						Valid:  true,
					},
				}

				oauthAccessTokenInsertable := modelFormatter.AccessTokenFromOauthAccessGrant(oauthAccessGrant, oauthApplication)

				expectedError := &entity.Error{
					HttpStatus: http.StatusInternalServerError,
					Errors:     eobject.Wrap(eobject.InternalServerError(ctx)),
				}

				gomock.InOrder(
					oauthApplicationModel.EXPECT().
						OneByUIDandSecret(ctx, *accessTokenRequest.ClientUID, *accessTokenRequest.ClientSecret).
						Return(oauthApplication, nil),
					oauthValidator.EXPECT().ValidateTokenGrant(ctx, accessTokenRequest).Return(nil),
					oauthAccessGrantModel.EXPECT().OneByCode(ctx, *accessTokenRequest.Code).Return(oauthAccessGrant, nil),
					modelFormatterMock.EXPECT().AccessTokenFromOauthAccessGrant(oauthAccessGrant, oauthApplication).Return(oauthAccessTokenInsertable),
					oauthAccessTokenModel.EXPECT().Create(ctx, oauthAccessTokenInsertable).Return(1, nil),
					oauthAccessTokenModel.EXPECT().One(ctx, 1).Return(entity.OauthAccessToken{}, sql.ErrNoRows),
					oauthFormatterMock.EXPECT().AccessToken(gomock.Any(), gomock.Any()).Times(0),
				)

				authorizationService := service.NewAuthorization(oauthApplicationModel, oauthAccessTokenModel, oauthAccessGrantModel, modelFormatterMock, oauthValidator, oauthFormatterMock)
				oauthAccessTokenOutput, err := authorizationService.Token(ctx, accessTokenRequest)

				assert.NotNil(t, err)
				assert.Equal(t, expectedError, err)
				assert.Equal(t, entity.OauthAccessTokenJSON{}, oauthAccessTokenOutput)
			})
		})
	})

	t.Run("RevokeToken", func(t *testing.T) {
		t.Run("Given context and revoke access token request", func(t *testing.T) {
			t.Run("Run gracefully", func(t *testing.T) {
				t.Run("Return nil", func(t *testing.T) {
					oauthApplicationModel := mock.NewMockOauthApplicationModel(mockCtrl)
					oauthAccessTokenModel := mock.NewMockOauthAccessTokenModel(mockCtrl)
					oauthAccessGrantModel := mock.NewMockOauthAccessGrantModel(mockCtrl)
					oauthValidator := mock.NewMockOauthValidator(mockCtrl)
					modelFormatterMock := mock.NewMockModelFormater(mockCtrl)
					oauthFormatterMock := mock.NewMockOauthFormatter(mockCtrl)

					ctx := context.WithValue(context.Background(), "track_id", uuid.New().String())

					revokeRequest := entity.RevokeAccessTokenRequestJSON{
						Token: util.StringToPointer("some-cool-token"),
					}

					oauthAccessTokenModel.EXPECT().Revoke(ctx, *revokeRequest.Token).Return(nil)

					authorizationService := service.NewAuthorization(oauthApplicationModel, oauthAccessTokenModel, oauthAccessGrantModel, modelFormatterMock, oauthValidator, oauthFormatterMock)
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
						oauthValidator := mock.NewMockOauthValidator(mockCtrl)
						modelFormatterMock := mock.NewMockModelFormater(mockCtrl)
						oauthFormatterMock := mock.NewMockOauthFormatter(mockCtrl)

						ctx := context.WithValue(context.Background(), "track_id", uuid.New().String())

						revokeRequest := entity.RevokeAccessTokenRequestJSON{
							Token: util.StringToPointer("some-cool-token"),
						}

						oauthAccessTokenModel.EXPECT().Revoke(ctx, *revokeRequest.Token).Return(sql.ErrNoRows)

						expectedError := &entity.Error{
							HttpStatus: http.StatusNotFound,
							Errors:     eobject.Wrap(eobject.NotFoundError(ctx, "token")),
						}

						authorizationService := service.NewAuthorization(oauthApplicationModel, oauthAccessTokenModel, oauthAccessGrantModel, modelFormatterMock, oauthValidator, oauthFormatterMock)
						err := authorizationService.RevokeToken(ctx, revokeRequest)
						assert.Equal(t, expectedError, err)
					})
				})

				t.Run("Other error", func(t *testing.T) {
					t.Run("Return 500 error", func(t *testing.T) {
						oauthApplicationModel := mock.NewMockOauthApplicationModel(mockCtrl)
						oauthAccessTokenModel := mock.NewMockOauthAccessTokenModel(mockCtrl)
						oauthAccessGrantModel := mock.NewMockOauthAccessGrantModel(mockCtrl)
						oauthValidator := mock.NewMockOauthValidator(mockCtrl)
						modelFormatterMock := mock.NewMockModelFormater(mockCtrl)
						oauthFormatterMock := mock.NewMockOauthFormatter(mockCtrl)

						ctx := context.WithValue(context.Background(), "track_id", uuid.New().String())

						revokeRequest := entity.RevokeAccessTokenRequestJSON{
							Token: util.StringToPointer("some-cool-token"),
						}

						oauthAccessTokenModel.EXPECT().Revoke(ctx, *revokeRequest.Token).Return(errors.New("unexpected error"))

						expectedError := &entity.Error{
							HttpStatus: http.StatusInternalServerError,
							Errors:     eobject.Wrap(eobject.InternalServerError(ctx)),
						}

						authorizationService := service.NewAuthorization(oauthApplicationModel, oauthAccessTokenModel, oauthAccessGrantModel, modelFormatterMock, oauthValidator, oauthFormatterMock)
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
				oauthValidator := mock.NewMockOauthValidator(mockCtrl)
				modelFormatterMock := mock.NewMockModelFormater(mockCtrl)
				oauthFormatterMock := mock.NewMockOauthFormatter(mockCtrl)

				ctx := context.WithValue(context.Background(), "track_id", uuid.New().String())

				revokeRequest := entity.RevokeAccessTokenRequestJSON{
					Token: nil,
				}

				oauthAccessTokenModel.EXPECT().Revoke(gomock.Any(), gomock.Any()).Times(0)

				expectedError := &entity.Error{
					HttpStatus: http.StatusUnprocessableEntity,
					Errors:     eobject.Wrap(eobject.ValidationError("token is empty")),
				}

				authorizationService := service.NewAuthorization(oauthApplicationModel, oauthAccessTokenModel, oauthAccessGrantModel, modelFormatterMock, oauthValidator, oauthFormatterMock)
				err := authorizationService.RevokeToken(ctx, revokeRequest)
				assert.Equal(t, expectedError, err)
			})
		})
	})
}
