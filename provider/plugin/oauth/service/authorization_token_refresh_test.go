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
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestAuthorizationRefreshToken(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	t.Run("Token", func(t *testing.T) {
		t.Run("Given context and refresh token request", func(t *testing.T) {
			t.Run("When refresh token request valid and there is no error in database side", func(t *testing.T) {
				t.Run("Then it will return access token response", func(t *testing.T) {
					oauthApplicationModel := mock.NewMockOauthApplicationModel(mockCtrl)
					oauthAccessTokenModel := mock.NewMockOauthAccessTokenModel(mockCtrl)
					oauthAccessGrantModel := mock.NewMockOauthAccessGrantModel(mockCtrl)
					oauthRefreshTokenModel := mock.NewMockOauthRefreshTokenModel(mockCtrl)
					oauthValidator := mock.NewMockOauthValidator(mockCtrl)
					modelFormatter := formatter.NewModel(time.Hour*4, time.Hour*2, time.Hour*2)
					oauthFormatter := formatter.Oauth()

					ctx := context.Background()

					accessTokenRequest := entity.AccessTokenRequestJSON{
						ClientSecret: util.StringToPointer("client_secret"),
						ClientUID:    util.StringToPointer("client_uid"),
						RefreshToken: util.StringToPointer("abcdef_123456"),
						GrantType:    util.StringToPointer("refresh_token"),
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

					oldAccessToken := entity.OauthAccessToken{
						ID: 999,
					}

					oauthRefreshToken := entity.OauthRefreshToken{
						ID:                 1,
						OauthAccessTokenID: oldAccessToken.ID,
						Token:              *accessTokenRequest.RefreshToken,
					}

					oauthAccessToken := entity.OauthAccessToken{
						ID:                 1000,
						OauthApplicationID: oauthApplication.ID,
						ResourceOwnerID:    oldAccessToken.ResourceOwnerID,
						Token:              aurelia.Hash("x", "y"),
						Scopes: sql.NullString{
							String: oldAccessToken.Scopes.String,
							Valid:  true,
						},
						ExpiresIn: time.Now().Add(time.Hour * 4),
						CreatedAt: time.Now(),
					}

					newOauthRefreshToken := entity.OauthRefreshToken{
						ID:                 2,
						OauthAccessTokenID: oauthAccessToken.ID,
						Token:              "newly created token",
					}

					oauthRefreshTokenJSON := oauthFormatter.RefreshToken(newOauthRefreshToken)

					oauthAccessTokenInsertable := modelFormatter.AccessTokenFromOauthRefreshToken(oauthApplication, oldAccessToken)
					oauthRefreshTokenInsertable := modelFormatter.RefreshToken(oauthApplication, oauthAccessToken)
					oauthAccessTokenJSON := oauthFormatter.AccessToken(oauthAccessToken, *accessTokenRequest.RedirectURI, &oauthRefreshTokenJSON)

					gomock.InOrder(
						oauthApplicationModel.EXPECT().
							OneByUIDandSecret(ctx, *accessTokenRequest.ClientUID, *accessTokenRequest.ClientSecret).
							Return(oauthApplication, nil),
						oauthValidator.EXPECT().ValidateTokenGrant(ctx, accessTokenRequest).Return(nil),
						oauthRefreshTokenModel.EXPECT().OneByToken(ctx, *accessTokenRequest.RefreshToken).Return(oauthRefreshToken, nil),
						oauthValidator.EXPECT().ValidateTokenRefreshToken(ctx, oauthRefreshToken).Return(nil),
						oauthAccessTokenModel.EXPECT().One(ctx, oauthRefreshToken.OauthAccessTokenID).Return(oldAccessToken, nil),
						oauthAccessTokenModel.EXPECT().Create(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, insertable entity.OauthAccessTokenInsertable) (int, error) {
							assert.Equal(t, oauthAccessTokenInsertable.ResourceOwnerID, insertable.ResourceOwnerID)
							assert.Equal(t, oauthAccessTokenInsertable.OauthApplicationID, insertable.OauthApplicationID)
							return 1000, nil
						}),
						oauthAccessTokenModel.EXPECT().One(ctx, 1000).Return(oauthAccessToken, nil),
						oauthRefreshTokenModel.EXPECT().Revoke(ctx, oauthRefreshToken.Token).Return(nil),
						oauthRefreshTokenModel.EXPECT().Create(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, insertable entity.OauthRefreshTokenInsertable) (int, error) {
							assert.Equal(t, oauthRefreshTokenInsertable.OauthAccessTokenID, insertable.OauthAccessTokenID)
							return 2, nil
						}),
						oauthRefreshTokenModel.EXPECT().One(ctx, 2).Return(newOauthRefreshToken, nil),
					)

					authorizationService := service.NewAuthorization(oauthApplicationModel, oauthAccessTokenModel, oauthAccessGrantModel, oauthRefreshTokenModel, modelFormatter, oauthValidator, oauthFormatter, true)
					oauthAccessTokenOutput, err := authorizationService.Token(ctx, accessTokenRequest)

					assert.Nil(t, err)
					assert.Equal(t, oauthAccessTokenJSON, oauthAccessTokenOutput)
				})
			})

			t.Run("When refresh token toggle is off", func(t *testing.T) {
				t.Run("Then it will return unprocessable entity error", func(t *testing.T) {
					oauthApplicationModel := mock.NewMockOauthApplicationModel(mockCtrl)
					oauthAccessTokenModel := mock.NewMockOauthAccessTokenModel(mockCtrl)
					oauthAccessGrantModel := mock.NewMockOauthAccessGrantModel(mockCtrl)
					oauthRefreshTokenModel := mock.NewMockOauthRefreshTokenModel(mockCtrl)
					oauthValidator := mock.NewMockOauthValidator(mockCtrl)
					modelFormatter := formatter.NewModel(time.Hour*4, time.Hour*2, time.Hour*2)
					oauthFormatter := formatter.Oauth()

					ctx := context.Background()

					accessTokenRequest := entity.AccessTokenRequestJSON{
						ClientSecret: util.StringToPointer("client_secret"),
						ClientUID:    util.StringToPointer("client_uid"),
						RefreshToken: util.StringToPointer("abcdef_123456"),
						GrantType:    util.StringToPointer("refresh_token"),
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

					gomock.InOrder(
						oauthApplicationModel.EXPECT().
							OneByUIDandSecret(ctx, *accessTokenRequest.ClientUID, *accessTokenRequest.ClientSecret).
							Return(oauthApplication, nil),
						oauthValidator.EXPECT().ValidateTokenGrant(ctx, accessTokenRequest).Return(nil),
						oauthRefreshTokenModel.EXPECT().OneByToken(gomock.Any(), gomock.Any()).Times(0),
						oauthValidator.EXPECT().ValidateTokenRefreshToken(gomock.Any(), gomock.Any()).Times(0),
						oauthAccessTokenModel.EXPECT().One(gomock.Any(), gomock.Any()).Times(0),
						oauthAccessTokenModel.EXPECT().Create(gomock.Any(), gomock.Any()).Times(0),
						oauthAccessTokenModel.EXPECT().One(gomock.Any(), gomock.Any()).Times(0),
						oauthRefreshTokenModel.EXPECT().Revoke(gomock.Any(), gomock.Any()).Times(0),
						oauthRefreshTokenModel.EXPECT().Create(gomock.Any(), gomock.Any()).Times(0),
					)

					authorizationService := service.NewAuthorization(oauthApplicationModel, oauthAccessTokenModel, oauthAccessGrantModel, oauthRefreshTokenModel, modelFormatter, oauthValidator, oauthFormatter, false)
					_, err := authorizationService.Token(ctx, accessTokenRequest)

					expectedError := &entity.Error{
						HttpStatus: http.StatusUnprocessableEntity,
						Errors:     eobject.Wrap(eobject.ValidationError(`grant_type can't be empty`)),
					}

					assert.Equal(t, expectedError, err)
				})
			})

			t.Run("When refresh token is not found", func(t *testing.T) {
				t.Run("Then it will return 404", func(t *testing.T) {
					oauthApplicationModel := mock.NewMockOauthApplicationModel(mockCtrl)
					oauthAccessTokenModel := mock.NewMockOauthAccessTokenModel(mockCtrl)
					oauthAccessGrantModel := mock.NewMockOauthAccessGrantModel(mockCtrl)
					oauthRefreshTokenModel := mock.NewMockOauthRefreshTokenModel(mockCtrl)
					oauthValidator := mock.NewMockOauthValidator(mockCtrl)
					modelFormatter := formatter.NewModel(time.Hour*4, time.Hour*2, time.Hour*2)
					oauthFormatter := formatter.Oauth()

					ctx := context.Background()

					accessTokenRequest := entity.AccessTokenRequestJSON{
						ClientSecret: util.StringToPointer("client_secret"),
						ClientUID:    util.StringToPointer("client_uid"),
						RefreshToken: util.StringToPointer("abcdef_123456"),
						GrantType:    util.StringToPointer("refresh_token"),
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

					gomock.InOrder(
						oauthApplicationModel.EXPECT().
							OneByUIDandSecret(ctx, *accessTokenRequest.ClientUID, *accessTokenRequest.ClientSecret).
							Return(oauthApplication, nil),
						oauthValidator.EXPECT().ValidateTokenGrant(ctx, accessTokenRequest).Return(nil),
						oauthRefreshTokenModel.EXPECT().OneByToken(ctx, *accessTokenRequest.RefreshToken).Return(entity.OauthRefreshToken{}, sql.ErrNoRows),
						oauthValidator.EXPECT().ValidateTokenRefreshToken(gomock.Any(), gomock.Any()).Times(0),
						oauthAccessTokenModel.EXPECT().One(gomock.Any(), gomock.Any()).Times(0),
						oauthAccessTokenModel.EXPECT().Create(gomock.Any(), gomock.Any()).Times(0),
						oauthAccessTokenModel.EXPECT().One(gomock.Any(), gomock.Any()).Times(0),
						oauthRefreshTokenModel.EXPECT().Revoke(gomock.Any(), gomock.Any()).Times(0),
						oauthRefreshTokenModel.EXPECT().Create(gomock.Any(), gomock.Any()).Times(0),
					)

					authorizationService := service.NewAuthorization(oauthApplicationModel, oauthAccessTokenModel, oauthAccessGrantModel, oauthRefreshTokenModel, modelFormatter, oauthValidator, oauthFormatter, true)
					_, err := authorizationService.Token(ctx, accessTokenRequest)

					expectedError := &entity.Error{
						HttpStatus: http.StatusNotFound,
						Errors:     eobject.Wrap(eobject.NotFoundError(ctx, "refresh_token")),
					}

					assert.Equal(t, expectedError, err)
				})
			})

			t.Run("When get refresh token got unexpected error", func(t *testing.T) {
				t.Run("Then it will return 500", func(t *testing.T) {
					oauthApplicationModel := mock.NewMockOauthApplicationModel(mockCtrl)
					oauthAccessTokenModel := mock.NewMockOauthAccessTokenModel(mockCtrl)
					oauthAccessGrantModel := mock.NewMockOauthAccessGrantModel(mockCtrl)
					oauthRefreshTokenModel := mock.NewMockOauthRefreshTokenModel(mockCtrl)
					oauthValidator := mock.NewMockOauthValidator(mockCtrl)
					modelFormatter := formatter.NewModel(time.Hour*4, time.Hour*2, time.Hour*2)
					oauthFormatter := formatter.Oauth()

					ctx := context.Background()

					accessTokenRequest := entity.AccessTokenRequestJSON{
						ClientSecret: util.StringToPointer("client_secret"),
						ClientUID:    util.StringToPointer("client_uid"),
						RefreshToken: util.StringToPointer("abcdef_123456"),
						GrantType:    util.StringToPointer("refresh_token"),
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

					gomock.InOrder(
						oauthApplicationModel.EXPECT().
							OneByUIDandSecret(ctx, *accessTokenRequest.ClientUID, *accessTokenRequest.ClientSecret).
							Return(oauthApplication, nil),
						oauthValidator.EXPECT().ValidateTokenGrant(ctx, accessTokenRequest).Return(nil),
						oauthRefreshTokenModel.EXPECT().OneByToken(ctx, *accessTokenRequest.RefreshToken).Return(entity.OauthRefreshToken{}, errors.New("unexpected error")),
						oauthValidator.EXPECT().ValidateTokenRefreshToken(gomock.Any(), gomock.Any()).Times(0),
						oauthAccessTokenModel.EXPECT().One(gomock.Any(), gomock.Any()).Times(0),
						oauthAccessTokenModel.EXPECT().Create(gomock.Any(), gomock.Any()).Times(0),
						oauthAccessTokenModel.EXPECT().One(gomock.Any(), gomock.Any()).Times(0),
						oauthRefreshTokenModel.EXPECT().Revoke(gomock.Any(), gomock.Any()).Times(0),
						oauthRefreshTokenModel.EXPECT().Create(gomock.Any(), gomock.Any()).Times(0),
					)

					authorizationService := service.NewAuthorization(oauthApplicationModel, oauthAccessTokenModel, oauthAccessGrantModel, oauthRefreshTokenModel, modelFormatter, oauthValidator, oauthFormatter, true)
					_, err := authorizationService.Token(ctx, accessTokenRequest)

					expectedError := &entity.Error{
						HttpStatus: http.StatusInternalServerError,
						Errors:     eobject.Wrap(eobject.InternalServerError(ctx)),
					}

					assert.Equal(t, expectedError, err)
				})
			})

			t.Run("When validate oauth refresh token is invalid", func(t *testing.T) {
				t.Run("Then it will return error", func(t *testing.T) {
					oauthApplicationModel := mock.NewMockOauthApplicationModel(mockCtrl)
					oauthAccessTokenModel := mock.NewMockOauthAccessTokenModel(mockCtrl)
					oauthAccessGrantModel := mock.NewMockOauthAccessGrantModel(mockCtrl)
					oauthRefreshTokenModel := mock.NewMockOauthRefreshTokenModel(mockCtrl)
					oauthValidator := mock.NewMockOauthValidator(mockCtrl)
					modelFormatter := formatter.NewModel(time.Hour*4, time.Hour*2, time.Hour*2)
					oauthFormatter := formatter.Oauth()

					ctx := context.Background()

					accessTokenRequest := entity.AccessTokenRequestJSON{
						ClientSecret: util.StringToPointer("client_secret"),
						ClientUID:    util.StringToPointer("client_uid"),
						RefreshToken: util.StringToPointer("abcdef_123456"),
						GrantType:    util.StringToPointer("refresh_token"),
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

					oldAccessToken := entity.OauthAccessToken{
						ID: 999,
					}

					oauthRefreshToken := entity.OauthRefreshToken{
						ID:                 1,
						OauthAccessTokenID: oldAccessToken.ID,
						Token:              *accessTokenRequest.RefreshToken,
					}

					expectedError := &entity.Error{
						HttpStatus: http.StatusUnprocessableEntity,
						Errors:     eobject.Wrap(eobject.ValidationError(`refresh token can't be empty`)),
					}

					gomock.InOrder(
						oauthApplicationModel.EXPECT().
							OneByUIDandSecret(ctx, *accessTokenRequest.ClientUID, *accessTokenRequest.ClientSecret).
							Return(oauthApplication, nil),
						oauthValidator.EXPECT().ValidateTokenGrant(ctx, accessTokenRequest).Return(nil),
						oauthRefreshTokenModel.EXPECT().OneByToken(ctx, *accessTokenRequest.RefreshToken).Return(oauthRefreshToken, nil),
						oauthValidator.EXPECT().ValidateTokenRefreshToken(ctx, oauthRefreshToken).Return(expectedError),
						oauthAccessTokenModel.EXPECT().One(gomock.Any(), gomock.Any()).Times(0),
						oauthAccessTokenModel.EXPECT().Create(gomock.Any(), gomock.Any()).Times(0),
						oauthAccessTokenModel.EXPECT().One(gomock.Any(), gomock.Any()).Times(0),
						oauthRefreshTokenModel.EXPECT().Revoke(gomock.Any(), gomock.Any()).Times(0),
						oauthRefreshTokenModel.EXPECT().Create(gomock.Any(), gomock.Any()).Times(0),
					)

					authorizationService := service.NewAuthorization(oauthApplicationModel, oauthAccessTokenModel, oauthAccessGrantModel, oauthRefreshTokenModel, modelFormatter, oauthValidator, oauthFormatter, true)
					_, err := authorizationService.Token(ctx, accessTokenRequest)

					assert.Equal(t, expectedError, err)
				})
			})

			t.Run("When there is an error when selecting old access token", func(t *testing.T) {
				t.Run("Then it will return internal server error", func(t *testing.T) {
					oauthApplicationModel := mock.NewMockOauthApplicationModel(mockCtrl)
					oauthAccessTokenModel := mock.NewMockOauthAccessTokenModel(mockCtrl)
					oauthAccessGrantModel := mock.NewMockOauthAccessGrantModel(mockCtrl)
					oauthRefreshTokenModel := mock.NewMockOauthRefreshTokenModel(mockCtrl)
					oauthValidator := mock.NewMockOauthValidator(mockCtrl)
					modelFormatter := formatter.NewModel(time.Hour*4, time.Hour*2, time.Hour*2)
					oauthFormatter := formatter.Oauth()

					ctx := context.Background()

					accessTokenRequest := entity.AccessTokenRequestJSON{
						ClientSecret: util.StringToPointer("client_secret"),
						ClientUID:    util.StringToPointer("client_uid"),
						RefreshToken: util.StringToPointer("abcdef_123456"),
						GrantType:    util.StringToPointer("refresh_token"),
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

					oldAccessToken := entity.OauthAccessToken{
						ID: 999,
					}

					oauthRefreshToken := entity.OauthRefreshToken{
						ID:                 1,
						OauthAccessTokenID: oldAccessToken.ID,
						Token:              *accessTokenRequest.RefreshToken,
					}

					gomock.InOrder(
						oauthApplicationModel.EXPECT().
							OneByUIDandSecret(ctx, *accessTokenRequest.ClientUID, *accessTokenRequest.ClientSecret).
							Return(oauthApplication, nil),
						oauthValidator.EXPECT().ValidateTokenGrant(ctx, accessTokenRequest).Return(nil),
						oauthRefreshTokenModel.EXPECT().OneByToken(ctx, *accessTokenRequest.RefreshToken).Return(oauthRefreshToken, nil),
						oauthValidator.EXPECT().ValidateTokenRefreshToken(ctx, oauthRefreshToken).Return(nil),
						oauthAccessTokenModel.EXPECT().One(ctx, oauthRefreshToken.OauthAccessTokenID).Return(entity.OauthAccessToken{}, errors.New("Unexpected error")),
						oauthAccessTokenModel.EXPECT().Create(gomock.Any(), gomock.Any()).Times(0),
						oauthAccessTokenModel.EXPECT().One(gomock.Any(), gomock.Any()).Times(0),
						oauthRefreshTokenModel.EXPECT().Revoke(gomock.Any(), gomock.Any()).Times(0),
						oauthRefreshTokenModel.EXPECT().Create(gomock.Any(), gomock.Any()).Times(0),
					)

					authorizationService := service.NewAuthorization(oauthApplicationModel, oauthAccessTokenModel, oauthAccessGrantModel, oauthRefreshTokenModel, modelFormatter, oauthValidator, oauthFormatter, true)
					_, err := authorizationService.Token(ctx, accessTokenRequest)

					expectedError := &entity.Error{
						HttpStatus: http.StatusInternalServerError,
						Errors:     eobject.Wrap(eobject.InternalServerError(ctx)),
					}

					assert.Equal(t, expectedError, err)
				})
			})

			t.Run("When old access token is already revoked", func(t *testing.T) {
				t.Run("Then it will return unauthorzed error", func(t *testing.T) {
					oauthApplicationModel := mock.NewMockOauthApplicationModel(mockCtrl)
					oauthAccessTokenModel := mock.NewMockOauthAccessTokenModel(mockCtrl)
					oauthAccessGrantModel := mock.NewMockOauthAccessGrantModel(mockCtrl)
					oauthRefreshTokenModel := mock.NewMockOauthRefreshTokenModel(mockCtrl)
					oauthValidator := mock.NewMockOauthValidator(mockCtrl)
					modelFormatter := formatter.NewModel(time.Hour*4, time.Hour*2, time.Hour*2)
					oauthFormatter := formatter.Oauth()

					ctx := context.Background()

					accessTokenRequest := entity.AccessTokenRequestJSON{
						ClientSecret: util.StringToPointer("client_secret"),
						ClientUID:    util.StringToPointer("client_uid"),
						RefreshToken: util.StringToPointer("abcdef_123456"),
						GrantType:    util.StringToPointer("refresh_token"),
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

					oldAccessToken := entity.OauthAccessToken{
						ID: 999,
					}

					oauthRefreshToken := entity.OauthRefreshToken{
						ID:                 1,
						OauthAccessTokenID: oldAccessToken.ID,
						Token:              *accessTokenRequest.RefreshToken,
					}

					gomock.InOrder(
						oauthApplicationModel.EXPECT().
							OneByUIDandSecret(ctx, *accessTokenRequest.ClientUID, *accessTokenRequest.ClientSecret).
							Return(oauthApplication, nil),
						oauthValidator.EXPECT().ValidateTokenGrant(ctx, accessTokenRequest).Return(nil),
						oauthRefreshTokenModel.EXPECT().OneByToken(ctx, *accessTokenRequest.RefreshToken).Return(oauthRefreshToken, nil),
						oauthValidator.EXPECT().ValidateTokenRefreshToken(ctx, oauthRefreshToken).Return(nil),
						oauthAccessTokenModel.EXPECT().One(ctx, oauthRefreshToken.OauthAccessTokenID).Return(entity.OauthAccessToken{}, sql.ErrNoRows),
						oauthAccessTokenModel.EXPECT().Create(gomock.Any(), gomock.Any()).Times(0),
						oauthAccessTokenModel.EXPECT().One(gomock.Any(), gomock.Any()).Times(0),
						oauthRefreshTokenModel.EXPECT().Revoke(gomock.Any(), gomock.Any()).Times(0),
						oauthRefreshTokenModel.EXPECT().Create(gomock.Any(), gomock.Any()).Times(0),
					)

					authorizationService := service.NewAuthorization(oauthApplicationModel, oauthAccessTokenModel, oauthAccessGrantModel, oauthRefreshTokenModel, modelFormatter, oauthValidator, oauthFormatter, true)
					_, err := authorizationService.Token(ctx, accessTokenRequest)

					expectedError := &entity.Error{
						HttpStatus: http.StatusUnauthorized,
						Errors:     eobject.Wrap(eobject.UnauthorizedError()),
					}

					assert.Equal(t, expectedError, err)
				})
			})

			t.Run("When there is an error when creating access token", func(t *testing.T) {
				t.Run("Then it will return internal server error", func(t *testing.T) {
					oauthApplicationModel := mock.NewMockOauthApplicationModel(mockCtrl)
					oauthAccessTokenModel := mock.NewMockOauthAccessTokenModel(mockCtrl)
					oauthAccessGrantModel := mock.NewMockOauthAccessGrantModel(mockCtrl)
					oauthRefreshTokenModel := mock.NewMockOauthRefreshTokenModel(mockCtrl)
					oauthValidator := mock.NewMockOauthValidator(mockCtrl)
					modelFormatter := formatter.NewModel(time.Hour*4, time.Hour*2, time.Hour*2)
					oauthFormatter := formatter.Oauth()

					ctx := context.Background()

					accessTokenRequest := entity.AccessTokenRequestJSON{
						ClientSecret: util.StringToPointer("client_secret"),
						ClientUID:    util.StringToPointer("client_uid"),
						RefreshToken: util.StringToPointer("abcdef_123456"),
						GrantType:    util.StringToPointer("refresh_token"),
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

					oldAccessToken := entity.OauthAccessToken{
						ID: 999,
					}

					oauthRefreshToken := entity.OauthRefreshToken{
						ID:                 1,
						OauthAccessTokenID: oldAccessToken.ID,
						Token:              *accessTokenRequest.RefreshToken,
					}

					gomock.InOrder(
						oauthApplicationModel.EXPECT().
							OneByUIDandSecret(ctx, *accessTokenRequest.ClientUID, *accessTokenRequest.ClientSecret).
							Return(oauthApplication, nil),
						oauthValidator.EXPECT().ValidateTokenGrant(ctx, accessTokenRequest).Return(nil),
						oauthRefreshTokenModel.EXPECT().OneByToken(ctx, *accessTokenRequest.RefreshToken).Return(oauthRefreshToken, nil),
						oauthValidator.EXPECT().ValidateTokenRefreshToken(ctx, oauthRefreshToken).Return(nil),
						oauthAccessTokenModel.EXPECT().One(ctx, oauthRefreshToken.OauthAccessTokenID).Return(oldAccessToken, nil),
						oauthAccessTokenModel.EXPECT().Create(ctx, gomock.Any()).Return(0, errors.New("unexpected error")),
						oauthAccessTokenModel.EXPECT().One(gomock.Any(), gomock.Any()).Times(0),
						oauthRefreshTokenModel.EXPECT().Revoke(gomock.Any(), gomock.Any()).Times(0),
						oauthRefreshTokenModel.EXPECT().Create(gomock.Any(), gomock.Any()).Times(0),
					)

					authorizationService := service.NewAuthorization(oauthApplicationModel, oauthAccessTokenModel, oauthAccessGrantModel, oauthRefreshTokenModel, modelFormatter, oauthValidator, oauthFormatter, true)
					_, err := authorizationService.Token(ctx, accessTokenRequest)

					expectedError := &entity.Error{
						HttpStatus: http.StatusInternalServerError,
						Errors:     eobject.Wrap(eobject.InternalServerError(ctx)),
					}

					assert.Equal(t, expectedError, err)
				})
			})

			t.Run("When there is an error when selecting newly created access token", func(t *testing.T) {
				t.Run("Then it will return internal server error", func(t *testing.T) {
					oauthApplicationModel := mock.NewMockOauthApplicationModel(mockCtrl)
					oauthAccessTokenModel := mock.NewMockOauthAccessTokenModel(mockCtrl)
					oauthAccessGrantModel := mock.NewMockOauthAccessGrantModel(mockCtrl)
					oauthRefreshTokenModel := mock.NewMockOauthRefreshTokenModel(mockCtrl)
					oauthValidator := mock.NewMockOauthValidator(mockCtrl)
					modelFormatter := formatter.NewModel(time.Hour*4, time.Hour*2, time.Hour*2)
					oauthFormatter := formatter.Oauth()

					ctx := context.Background()

					accessTokenRequest := entity.AccessTokenRequestJSON{
						ClientSecret: util.StringToPointer("client_secret"),
						ClientUID:    util.StringToPointer("client_uid"),
						RefreshToken: util.StringToPointer("abcdef_123456"),
						GrantType:    util.StringToPointer("refresh_token"),
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

					oldAccessToken := entity.OauthAccessToken{
						ID: 999,
					}

					oauthRefreshToken := entity.OauthRefreshToken{
						ID:                 1,
						OauthAccessTokenID: oldAccessToken.ID,
						Token:              *accessTokenRequest.RefreshToken,
					}

					oauthAccessTokenInsertable := modelFormatter.AccessTokenFromOauthRefreshToken(oauthApplication, oldAccessToken)

					gomock.InOrder(
						oauthApplicationModel.EXPECT().
							OneByUIDandSecret(ctx, *accessTokenRequest.ClientUID, *accessTokenRequest.ClientSecret).
							Return(oauthApplication, nil),
						oauthValidator.EXPECT().ValidateTokenGrant(ctx, accessTokenRequest).Return(nil),
						oauthRefreshTokenModel.EXPECT().OneByToken(ctx, *accessTokenRequest.RefreshToken).Return(oauthRefreshToken, nil),
						oauthValidator.EXPECT().ValidateTokenRefreshToken(ctx, oauthRefreshToken).Return(nil),
						oauthAccessTokenModel.EXPECT().One(ctx, oauthRefreshToken.OauthAccessTokenID).Return(oldAccessToken, nil),
						oauthAccessTokenModel.EXPECT().Create(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, insertable entity.OauthAccessTokenInsertable) (int, error) {
							assert.Equal(t, oauthAccessTokenInsertable.ResourceOwnerID, insertable.ResourceOwnerID)
							assert.Equal(t, oauthAccessTokenInsertable.OauthApplicationID, insertable.OauthApplicationID)
							return 1000, nil
						}),
						oauthAccessTokenModel.EXPECT().One(ctx, 1000).Return(entity.OauthAccessToken{}, errors.New("unexpected error")),
						oauthRefreshTokenModel.EXPECT().Revoke(gomock.Any(), gomock.Any()).Times(0),
						oauthRefreshTokenModel.EXPECT().Create(gomock.Any(), gomock.Any()).Times(0),
					)

					authorizationService := service.NewAuthorization(oauthApplicationModel, oauthAccessTokenModel, oauthAccessGrantModel, oauthRefreshTokenModel, modelFormatter, oauthValidator, oauthFormatter, true)
					_, err := authorizationService.Token(ctx, accessTokenRequest)

					expectedError := &entity.Error{
						HttpStatus: http.StatusInternalServerError,
						Errors:     eobject.Wrap(eobject.InternalServerError(ctx)),
					}

					assert.Equal(t, expectedError, err)
				})
			})

			t.Run("When refresh token request valid and there is error when revoking old refresh token", func(t *testing.T) {
				t.Run("Then it will return access token response", func(t *testing.T) {
					oauthApplicationModel := mock.NewMockOauthApplicationModel(mockCtrl)
					oauthAccessTokenModel := mock.NewMockOauthAccessTokenModel(mockCtrl)
					oauthAccessGrantModel := mock.NewMockOauthAccessGrantModel(mockCtrl)
					oauthRefreshTokenModel := mock.NewMockOauthRefreshTokenModel(mockCtrl)
					oauthValidator := mock.NewMockOauthValidator(mockCtrl)
					modelFormatter := formatter.NewModel(time.Hour*4, time.Hour*2, time.Hour*2)
					oauthFormatter := formatter.Oauth()

					ctx := context.Background()

					accessTokenRequest := entity.AccessTokenRequestJSON{
						ClientSecret: util.StringToPointer("client_secret"),
						ClientUID:    util.StringToPointer("client_uid"),
						RefreshToken: util.StringToPointer("abcdef_123456"),
						GrantType:    util.StringToPointer("refresh_token"),
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

					oldAccessToken := entity.OauthAccessToken{
						ID: 999,
					}

					oauthRefreshToken := entity.OauthRefreshToken{
						ID:                 1,
						OauthAccessTokenID: oldAccessToken.ID,
						Token:              *accessTokenRequest.RefreshToken,
					}

					oauthAccessToken := entity.OauthAccessToken{
						ID:                 1000,
						OauthApplicationID: oauthApplication.ID,
						ResourceOwnerID:    oldAccessToken.ResourceOwnerID,
						Token:              aurelia.Hash("x", "y"),
						Scopes: sql.NullString{
							String: oldAccessToken.Scopes.String,
							Valid:  true,
						},
						ExpiresIn: time.Now().Add(time.Hour * 4),
						CreatedAt: time.Now(),
					}

					newOauthRefreshToken := entity.OauthRefreshToken{
						ID:                 2,
						OauthAccessTokenID: oauthAccessToken.ID,
						Token:              "newly created token",
					}

					oauthRefreshTokenJSON := oauthFormatter.RefreshToken(newOauthRefreshToken)

					oauthAccessTokenInsertable := modelFormatter.AccessTokenFromOauthRefreshToken(oauthApplication, oldAccessToken)
					oauthRefreshTokenInsertable := modelFormatter.RefreshToken(oauthApplication, oauthAccessToken)
					oauthAccessTokenJSON := oauthFormatter.AccessToken(oauthAccessToken, *accessTokenRequest.RedirectURI, &oauthRefreshTokenJSON)

					gomock.InOrder(
						oauthApplicationModel.EXPECT().
							OneByUIDandSecret(ctx, *accessTokenRequest.ClientUID, *accessTokenRequest.ClientSecret).
							Return(oauthApplication, nil),
						oauthValidator.EXPECT().ValidateTokenGrant(ctx, accessTokenRequest).Return(nil),
						oauthRefreshTokenModel.EXPECT().OneByToken(ctx, *accessTokenRequest.RefreshToken).Return(oauthRefreshToken, nil),
						oauthValidator.EXPECT().ValidateTokenRefreshToken(ctx, oauthRefreshToken).Return(nil),
						oauthAccessTokenModel.EXPECT().One(ctx, oauthRefreshToken.OauthAccessTokenID).Return(oldAccessToken, nil),
						oauthAccessTokenModel.EXPECT().Create(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, insertable entity.OauthAccessTokenInsertable) (int, error) {
							assert.Equal(t, oauthAccessTokenInsertable.ResourceOwnerID, insertable.ResourceOwnerID)
							assert.Equal(t, oauthAccessTokenInsertable.OauthApplicationID, insertable.OauthApplicationID)
							return 1000, nil
						}),
						oauthAccessTokenModel.EXPECT().One(ctx, 1000).Return(oauthAccessToken, nil),
						oauthRefreshTokenModel.EXPECT().Revoke(ctx, oauthRefreshToken.Token).Return(errors.New("unexpected error")),
						oauthRefreshTokenModel.EXPECT().Create(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, insertable entity.OauthRefreshTokenInsertable) (int, error) {
							assert.Equal(t, oauthRefreshTokenInsertable.OauthAccessTokenID, insertable.OauthAccessTokenID)
							return 2, nil
						}),
						oauthRefreshTokenModel.EXPECT().One(ctx, 2).Return(newOauthRefreshToken, nil),
					)

					authorizationService := service.NewAuthorization(oauthApplicationModel, oauthAccessTokenModel, oauthAccessGrantModel, oauthRefreshTokenModel, modelFormatter, oauthValidator, oauthFormatter, true)
					oauthAccessTokenOutput, err := authorizationService.Token(ctx, accessTokenRequest)

					assert.Nil(t, err)
					assert.Equal(t, oauthAccessTokenJSON, oauthAccessTokenOutput)
				})
			})

			t.Run("When creating new refresh token failed", func(t *testing.T) {
				t.Run("Then it will return internal server error", func(t *testing.T) {
					oauthApplicationModel := mock.NewMockOauthApplicationModel(mockCtrl)
					oauthAccessTokenModel := mock.NewMockOauthAccessTokenModel(mockCtrl)
					oauthAccessGrantModel := mock.NewMockOauthAccessGrantModel(mockCtrl)
					oauthRefreshTokenModel := mock.NewMockOauthRefreshTokenModel(mockCtrl)
					oauthValidator := mock.NewMockOauthValidator(mockCtrl)
					modelFormatter := formatter.NewModel(time.Hour*4, time.Hour*2, time.Hour*2)
					oauthFormatter := formatter.Oauth()

					ctx := context.Background()

					accessTokenRequest := entity.AccessTokenRequestJSON{
						ClientSecret: util.StringToPointer("client_secret"),
						ClientUID:    util.StringToPointer("client_uid"),
						RefreshToken: util.StringToPointer("abcdef_123456"),
						GrantType:    util.StringToPointer("refresh_token"),
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

					oldAccessToken := entity.OauthAccessToken{
						ID: 999,
					}

					oauthRefreshToken := entity.OauthRefreshToken{
						ID:                 1,
						OauthAccessTokenID: oldAccessToken.ID,
						Token:              *accessTokenRequest.RefreshToken,
					}

					oauthAccessToken := entity.OauthAccessToken{
						ID:                 1000,
						OauthApplicationID: oauthApplication.ID,
						ResourceOwnerID:    oldAccessToken.ResourceOwnerID,
						Token:              aurelia.Hash("x", "y"),
						Scopes: sql.NullString{
							String: oldAccessToken.Scopes.String,
							Valid:  true,
						},
						ExpiresIn: time.Now().Add(time.Hour * 4),
						CreatedAt: time.Now(),
					}

					oauthAccessTokenInsertable := modelFormatter.AccessTokenFromOauthRefreshToken(oauthApplication, oldAccessToken)

					gomock.InOrder(
						oauthApplicationModel.EXPECT().
							OneByUIDandSecret(ctx, *accessTokenRequest.ClientUID, *accessTokenRequest.ClientSecret).
							Return(oauthApplication, nil),
						oauthValidator.EXPECT().ValidateTokenGrant(ctx, accessTokenRequest).Return(nil),
						oauthRefreshTokenModel.EXPECT().OneByToken(ctx, *accessTokenRequest.RefreshToken).Return(oauthRefreshToken, nil),
						oauthValidator.EXPECT().ValidateTokenRefreshToken(ctx, oauthRefreshToken).Return(nil),
						oauthAccessTokenModel.EXPECT().One(ctx, oauthRefreshToken.OauthAccessTokenID).Return(oldAccessToken, nil),
						oauthAccessTokenModel.EXPECT().Create(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, insertable entity.OauthAccessTokenInsertable) (int, error) {
							assert.Equal(t, oauthAccessTokenInsertable.ResourceOwnerID, insertable.ResourceOwnerID)
							assert.Equal(t, oauthAccessTokenInsertable.OauthApplicationID, insertable.OauthApplicationID)
							return 1000, nil
						}),
						oauthAccessTokenModel.EXPECT().One(ctx, 1000).Return(oauthAccessToken, nil),
						oauthRefreshTokenModel.EXPECT().Revoke(ctx, oauthRefreshToken.Token).Return(nil),
						oauthRefreshTokenModel.EXPECT().Create(ctx, gomock.Any()).Return(0, errors.New("unexpected error")),
					)

					authorizationService := service.NewAuthorization(oauthApplicationModel, oauthAccessTokenModel, oauthAccessGrantModel, oauthRefreshTokenModel, modelFormatter, oauthValidator, oauthFormatter, true)
					_, err := authorizationService.Token(ctx, accessTokenRequest)

					expectedError := &entity.Error{
						HttpStatus: http.StatusInternalServerError,
						Errors:     eobject.Wrap(eobject.InternalServerError(ctx)),
					}

					assert.Equal(t, expectedError, err)
				})
			})

			t.Run("When refresh token request valid and there is error when selecting newly created refresh token", func(t *testing.T) {
				t.Run("Then it will return internal server error", func(t *testing.T) {
					oauthApplicationModel := mock.NewMockOauthApplicationModel(mockCtrl)
					oauthAccessTokenModel := mock.NewMockOauthAccessTokenModel(mockCtrl)
					oauthAccessGrantModel := mock.NewMockOauthAccessGrantModel(mockCtrl)
					oauthRefreshTokenModel := mock.NewMockOauthRefreshTokenModel(mockCtrl)
					oauthValidator := mock.NewMockOauthValidator(mockCtrl)
					modelFormatter := formatter.NewModel(time.Hour*4, time.Hour*2, time.Hour*2)
					oauthFormatter := formatter.Oauth()

					ctx := context.Background()

					accessTokenRequest := entity.AccessTokenRequestJSON{
						ClientSecret: util.StringToPointer("client_secret"),
						ClientUID:    util.StringToPointer("client_uid"),
						RefreshToken: util.StringToPointer("abcdef_123456"),
						GrantType:    util.StringToPointer("refresh_token"),
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

					oldAccessToken := entity.OauthAccessToken{
						ID: 999,
					}

					oauthRefreshToken := entity.OauthRefreshToken{
						ID:                 1,
						OauthAccessTokenID: oldAccessToken.ID,
						Token:              *accessTokenRequest.RefreshToken,
					}

					oauthAccessToken := entity.OauthAccessToken{
						ID:                 1000,
						OauthApplicationID: oauthApplication.ID,
						ResourceOwnerID:    oldAccessToken.ResourceOwnerID,
						Token:              aurelia.Hash("x", "y"),
						Scopes: sql.NullString{
							String: oldAccessToken.Scopes.String,
							Valid:  true,
						},
						ExpiresIn: time.Now().Add(time.Hour * 4),
						CreatedAt: time.Now(),
					}

					oauthAccessTokenInsertable := modelFormatter.AccessTokenFromOauthRefreshToken(oauthApplication, oldAccessToken)
					oauthRefreshTokenInsertable := modelFormatter.RefreshToken(oauthApplication, oauthAccessToken)

					gomock.InOrder(
						oauthApplicationModel.EXPECT().
							OneByUIDandSecret(ctx, *accessTokenRequest.ClientUID, *accessTokenRequest.ClientSecret).
							Return(oauthApplication, nil),
						oauthValidator.EXPECT().ValidateTokenGrant(ctx, accessTokenRequest).Return(nil),
						oauthRefreshTokenModel.EXPECT().OneByToken(ctx, *accessTokenRequest.RefreshToken).Return(oauthRefreshToken, nil),
						oauthValidator.EXPECT().ValidateTokenRefreshToken(ctx, oauthRefreshToken).Return(nil),
						oauthAccessTokenModel.EXPECT().One(ctx, oauthRefreshToken.OauthAccessTokenID).Return(oldAccessToken, nil),
						oauthAccessTokenModel.EXPECT().Create(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, insertable entity.OauthAccessTokenInsertable) (int, error) {
							assert.Equal(t, oauthAccessTokenInsertable.ResourceOwnerID, insertable.ResourceOwnerID)
							assert.Equal(t, oauthAccessTokenInsertable.OauthApplicationID, insertable.OauthApplicationID)
							return 1000, nil
						}),
						oauthAccessTokenModel.EXPECT().One(ctx, 1000).Return(oauthAccessToken, nil),
						oauthRefreshTokenModel.EXPECT().Revoke(ctx, oauthRefreshToken.Token).Return(nil),
						oauthRefreshTokenModel.EXPECT().Create(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, insertable entity.OauthRefreshTokenInsertable) (int, error) {
							assert.Equal(t, oauthRefreshTokenInsertable.OauthAccessTokenID, insertable.OauthAccessTokenID)
							return 2, nil
						}),
						oauthRefreshTokenModel.EXPECT().One(ctx, 2).Return(entity.OauthRefreshToken{}, errors.New("unexpected error")),
					)

					authorizationService := service.NewAuthorization(oauthApplicationModel, oauthAccessTokenModel, oauthAccessGrantModel, oauthRefreshTokenModel, modelFormatter, oauthValidator, oauthFormatter, true)
					_, err := authorizationService.Token(ctx, accessTokenRequest)

					expectedError := &entity.Error{
						HttpStatus: http.StatusInternalServerError,
						Errors:     eobject.Wrap(eobject.InternalServerError(ctx)),
					}

					assert.Equal(t, expectedError, err)
				})
			})
		})
	})
}
