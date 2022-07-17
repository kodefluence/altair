package service_test

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/kodefluence/altair/provider/plugin/oauth/entity"
	"github.com/kodefluence/altair/provider/plugin/oauth/eobject"
	"github.com/kodefluence/altair/provider/plugin/oauth/formatter"
	"github.com/kodefluence/altair/provider/plugin/oauth/mock"
	"github.com/kodefluence/altair/provider/plugin/oauth/service"
	"github.com/kodefluence/altair/util"
	"github.com/kodefluence/aurelia"
	"github.com/kodefluence/monorepo/db"
	mockdb "github.com/kodefluence/monorepo/db/mock"
	"github.com/kodefluence/monorepo/exception"
	"github.com/kodefluence/monorepo/kontext"
	"github.com/stretchr/testify/assert"
)

func TestAuthorizationRefreshToken(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	sqldb := mockdb.NewMockDB(mockCtrl)
	mockTx := mockdb.NewMockTX(mockCtrl)

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

					mockCall := []*gomock.Call{}
					mockCall = append(mockCall,
						oauthApplicationModel.EXPECT().
							OneByUIDandSecret(gomock.Any(), *accessTokenRequest.ClientUID, *accessTokenRequest.ClientSecret, sqldb).
							Return(oauthApplication, nil),
						oauthValidator.EXPECT().ValidateTokenGrant(ctx, accessTokenRequest).Return(nil),
					)

					insideTransactionCall := []*gomock.Call{}

					sqlTransactionCall := sqldb.EXPECT().Transaction(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(ctx kontext.Context, transactionKey string, f func(tx db.TX) exception.Exception) exception.Exception {
						assert.Equal(t, "authorization-grant-token-from-refresh-token", transactionKey)

						insideTransactionCall = append(insideTransactionCall,
							oauthRefreshTokenModel.EXPECT().OneByToken(gomock.Any(), *accessTokenRequest.RefreshToken, mockTx).Return(oauthRefreshToken, nil),
							oauthValidator.EXPECT().ValidateTokenRefreshToken(gomock.Any(), oauthRefreshToken).Return(nil),
							oauthAccessTokenModel.EXPECT().One(gomock.Any(), oauthRefreshToken.OauthAccessTokenID, mockTx).Return(oldAccessToken, nil),
							oauthAccessTokenModel.EXPECT().Create(gomock.Any(), gomock.Any(), mockTx).DoAndReturn(func(ktx kontext.Context, insertable entity.OauthAccessTokenInsertable, tx db.TX) (int, error) {
								assert.Equal(t, oauthAccessTokenInsertable.ResourceOwnerID, insertable.ResourceOwnerID)
								assert.Equal(t, oauthAccessTokenInsertable.OauthApplicationID, insertable.OauthApplicationID)
								return 1000, nil
							}),
							oauthAccessTokenModel.EXPECT().One(gomock.Any(), 1000, mockTx).Return(oauthAccessToken, nil),
							oauthRefreshTokenModel.EXPECT().Revoke(gomock.Any(), oauthRefreshToken.Token, mockTx).Return(nil),
							oauthRefreshTokenModel.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(ktx kontext.Context, insertable entity.OauthRefreshTokenInsertable, tx db.TX) (int, error) {
								assert.Equal(t, oauthRefreshTokenInsertable.OauthAccessTokenID, insertable.OauthAccessTokenID)
								return 2, nil
							}),
							oauthRefreshTokenModel.EXPECT().One(gomock.Any(), 2, mockTx).Return(newOauthRefreshToken, nil),
						)
						return f(mockTx)
					})

					mockCall = append(mockCall, sqlTransactionCall)
					mockCall = append(mockCall, insideTransactionCall...)

					gomock.InOrder(mockCall...)

					authorizationService := service.NewAuthorization(oauthApplicationModel, oauthAccessTokenModel, oauthAccessGrantModel, oauthRefreshTokenModel, modelFormatter, oauthValidator, oauthFormatter, true, sqldb)
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

					mockCall := []*gomock.Call{}
					mockCall = append(mockCall,
						oauthApplicationModel.EXPECT().
							OneByUIDandSecret(gomock.Any(), *accessTokenRequest.ClientUID, *accessTokenRequest.ClientSecret, sqldb).
							Return(oauthApplication, nil),
						oauthValidator.EXPECT().ValidateTokenGrant(ctx, accessTokenRequest).Return(nil),
					)

					insideTransactionCall := []*gomock.Call{}

					sqlTransactionCall := sqldb.EXPECT().Transaction(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(ctx kontext.Context, transactionKey string, f func(tx db.TX) exception.Exception) exception.Exception {
						assert.Equal(t, "authorization-grant-token-from-refresh-token", transactionKey)

						insideTransactionCall = append(insideTransactionCall,
							oauthRefreshTokenModel.EXPECT().OneByToken(gomock.Any(), gomock.Any(), gomock.Any()).Times(0),
							oauthValidator.EXPECT().ValidateTokenRefreshToken(gomock.Any(), gomock.Any()).Times(0),
							oauthAccessTokenModel.EXPECT().One(gomock.Any(), gomock.Any(), gomock.Any()).Times(0),
							oauthAccessTokenModel.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).Times(0),
							oauthAccessTokenModel.EXPECT().One(gomock.Any(), gomock.Any(), gomock.Any()).Times(0),
							oauthRefreshTokenModel.EXPECT().Revoke(gomock.Any(), gomock.Any(), gomock.Any()).Times(0),
							oauthRefreshTokenModel.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).Times(0),
						)
						return f(mockTx)
					}).Times(0)

					mockCall = append(mockCall, sqlTransactionCall)
					mockCall = append(mockCall, insideTransactionCall...)
					gomock.InOrder(mockCall...)

					authorizationService := service.NewAuthorization(oauthApplicationModel, oauthAccessTokenModel, oauthAccessGrantModel, oauthRefreshTokenModel, modelFormatter, oauthValidator, oauthFormatter, false, sqldb)
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

					mockCall := []*gomock.Call{}
					mockCall = append(mockCall,
						oauthApplicationModel.EXPECT().
							OneByUIDandSecret(gomock.Any(), *accessTokenRequest.ClientUID, *accessTokenRequest.ClientSecret, sqldb).
							Return(oauthApplication, nil),
						oauthValidator.EXPECT().ValidateTokenGrant(ctx, accessTokenRequest).Return(nil),
					)

					insideTransactionCall := []*gomock.Call{}

					sqlTransactionCall := sqldb.EXPECT().Transaction(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(ctx kontext.Context, transactionKey string, f func(tx db.TX) exception.Exception) exception.Exception {
						assert.Equal(t, "authorization-grant-token-from-refresh-token", transactionKey)

						insideTransactionCall = append(insideTransactionCall,
							oauthRefreshTokenModel.EXPECT().OneByToken(gomock.Any(), *accessTokenRequest.RefreshToken, mockTx).Return(entity.OauthRefreshToken{}, exception.Throw(sql.ErrNoRows, exception.WithType(exception.NotFound))),
							oauthValidator.EXPECT().ValidateTokenRefreshToken(gomock.Any(), gomock.Any()).Times(0),
							oauthAccessTokenModel.EXPECT().One(gomock.Any(), gomock.Any(), gomock.Any()).Times(0),
							oauthAccessTokenModel.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).Times(0),
							oauthAccessTokenModel.EXPECT().One(gomock.Any(), gomock.Any(), gomock.Any()).Times(0),
							oauthRefreshTokenModel.EXPECT().Revoke(gomock.Any(), gomock.Any(), gomock.Any()).Times(0),
							oauthRefreshTokenModel.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).Times(0),
						)
						return f(mockTx)
					})

					mockCall = append(mockCall, sqlTransactionCall)
					mockCall = append(mockCall, insideTransactionCall...)

					gomock.InOrder(mockCall...)

					authorizationService := service.NewAuthorization(oauthApplicationModel, oauthAccessTokenModel, oauthAccessGrantModel, oauthRefreshTokenModel, modelFormatter, oauthValidator, oauthFormatter, true, sqldb)
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

					mockCall := []*gomock.Call{}
					mockCall = append(mockCall,
						oauthApplicationModel.EXPECT().
							OneByUIDandSecret(gomock.Any(), *accessTokenRequest.ClientUID, *accessTokenRequest.ClientSecret, sqldb).
							Return(oauthApplication, nil),
						oauthValidator.EXPECT().ValidateTokenGrant(ctx, accessTokenRequest).Return(nil),
					)

					insideTransactionCall := []*gomock.Call{}

					sqlTransactionCall := sqldb.EXPECT().Transaction(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(ctx kontext.Context, transactionKey string, f func(tx db.TX) exception.Exception) exception.Exception {
						assert.Equal(t, "authorization-grant-token-from-refresh-token", transactionKey)

						insideTransactionCall = append(insideTransactionCall,
							oauthRefreshTokenModel.EXPECT().OneByToken(gomock.Any(), *accessTokenRequest.RefreshToken, mockTx).Return(entity.OauthRefreshToken{}, exception.Throw(errors.New("unexpected error"))),
							oauthValidator.EXPECT().ValidateTokenRefreshToken(gomock.Any(), gomock.Any()).Times(0),
							oauthAccessTokenModel.EXPECT().One(gomock.Any(), gomock.Any(), gomock.Any()).Times(0),
							oauthAccessTokenModel.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).Times(0),
							oauthAccessTokenModel.EXPECT().One(gomock.Any(), gomock.Any(), gomock.Any()).Times(0),
							oauthRefreshTokenModel.EXPECT().Revoke(gomock.Any(), gomock.Any(), gomock.Any()).Times(0),
							oauthRefreshTokenModel.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).Times(0),
						)
						return f(mockTx)
					})

					mockCall = append(mockCall, sqlTransactionCall)
					mockCall = append(mockCall, insideTransactionCall...)

					gomock.InOrder(mockCall...)

					authorizationService := service.NewAuthorization(oauthApplicationModel, oauthAccessTokenModel, oauthAccessGrantModel, oauthRefreshTokenModel, modelFormatter, oauthValidator, oauthFormatter, true, sqldb)
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

					errorObject := eobject.ForbiddenError(ctx, "access_token", "refresh token already used")

					expectedError := &entity.Error{
						HttpStatus: http.StatusForbidden,
						Errors:     eobject.Wrap(errorObject),
					}

					mockCall := []*gomock.Call{}
					mockCall = append(mockCall,
						oauthApplicationModel.EXPECT().
							OneByUIDandSecret(gomock.Any(), *accessTokenRequest.ClientUID, *accessTokenRequest.ClientSecret, sqldb).
							Return(oauthApplication, nil),
						oauthValidator.EXPECT().ValidateTokenGrant(ctx, accessTokenRequest).Return(nil),
					)

					insideTransactionCall := []*gomock.Call{}

					sqlTransactionCall := sqldb.EXPECT().Transaction(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(ctx kontext.Context, transactionKey string, f func(tx db.TX) exception.Exception) exception.Exception {
						assert.Equal(t, "authorization-grant-token-from-refresh-token", transactionKey)

						insideTransactionCall = append(insideTransactionCall,
							oauthRefreshTokenModel.EXPECT().OneByToken(gomock.Any(), *accessTokenRequest.RefreshToken, mockTx).Return(oauthRefreshToken, nil),
							oauthValidator.EXPECT().ValidateTokenRefreshToken(gomock.Any(), oauthRefreshToken).Return(exception.Throw(errorObject, exception.WithTitle(errorObject.Code), exception.WithDetail(errorObject.Message), exception.WithType(exception.Forbidden))),
							oauthAccessTokenModel.EXPECT().One(gomock.Any(), gomock.Any(), gomock.Any()).Times(0),
							oauthAccessTokenModel.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).Times(0),
							oauthAccessTokenModel.EXPECT().One(gomock.Any(), gomock.Any(), gomock.Any()).Times(0),
							oauthRefreshTokenModel.EXPECT().Revoke(gomock.Any(), gomock.Any(), gomock.Any()).Times(0),
							oauthRefreshTokenModel.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).Times(0),
						)
						return f(mockTx)
					})

					mockCall = append(mockCall, sqlTransactionCall)
					mockCall = append(mockCall, insideTransactionCall...)

					gomock.InOrder(mockCall...)

					authorizationService := service.NewAuthorization(oauthApplicationModel, oauthAccessTokenModel, oauthAccessGrantModel, oauthRefreshTokenModel, modelFormatter, oauthValidator, oauthFormatter, true, sqldb)
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

					mockCall := []*gomock.Call{}
					mockCall = append(mockCall,
						oauthApplicationModel.EXPECT().
							OneByUIDandSecret(gomock.Any(), *accessTokenRequest.ClientUID, *accessTokenRequest.ClientSecret, sqldb).
							Return(oauthApplication, nil),
						oauthValidator.EXPECT().ValidateTokenGrant(ctx, accessTokenRequest).Return(nil),
					)

					insideTransactionCall := []*gomock.Call{}

					sqlTransactionCall := sqldb.EXPECT().Transaction(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(ctx kontext.Context, transactionKey string, f func(tx db.TX) exception.Exception) exception.Exception {
						assert.Equal(t, "authorization-grant-token-from-refresh-token", transactionKey)

						insideTransactionCall = append(insideTransactionCall,
							oauthRefreshTokenModel.EXPECT().OneByToken(gomock.Any(), *accessTokenRequest.RefreshToken, mockTx).Return(oauthRefreshToken, nil),
							oauthValidator.EXPECT().ValidateTokenRefreshToken(gomock.Any(), oauthRefreshToken).Return(nil),
							oauthAccessTokenModel.EXPECT().One(gomock.Any(), oauthRefreshToken.OauthAccessTokenID, mockTx).Return(entity.OauthAccessToken{}, exception.Throw(errors.New("Unexpected error"))),
							oauthAccessTokenModel.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).Times(0),
							oauthAccessTokenModel.EXPECT().One(gomock.Any(), gomock.Any(), gomock.Any()).Times(0),
							oauthRefreshTokenModel.EXPECT().Revoke(gomock.Any(), gomock.Any(), gomock.Any()).Times(0),
							oauthRefreshTokenModel.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).Times(0),
						)
						return f(mockTx)
					})

					mockCall = append(mockCall, sqlTransactionCall)
					mockCall = append(mockCall, insideTransactionCall...)

					gomock.InOrder(mockCall...)

					authorizationService := service.NewAuthorization(oauthApplicationModel, oauthAccessTokenModel, oauthAccessGrantModel, oauthRefreshTokenModel, modelFormatter, oauthValidator, oauthFormatter, true, sqldb)
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

					mockCall := []*gomock.Call{}
					mockCall = append(mockCall,
						oauthApplicationModel.EXPECT().
							OneByUIDandSecret(gomock.Any(), *accessTokenRequest.ClientUID, *accessTokenRequest.ClientSecret, sqldb).
							Return(oauthApplication, nil),
						oauthValidator.EXPECT().ValidateTokenGrant(ctx, accessTokenRequest).Return(nil),
					)

					insideTransactionCall := []*gomock.Call{}

					sqlTransactionCall := sqldb.EXPECT().Transaction(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(ctx kontext.Context, transactionKey string, f func(tx db.TX) exception.Exception) exception.Exception {
						assert.Equal(t, "authorization-grant-token-from-refresh-token", transactionKey)

						insideTransactionCall = append(insideTransactionCall,
							oauthRefreshTokenModel.EXPECT().OneByToken(gomock.Any(), *accessTokenRequest.RefreshToken, mockTx).Return(oauthRefreshToken, nil),
							oauthValidator.EXPECT().ValidateTokenRefreshToken(gomock.Any(), oauthRefreshToken).Return(nil),
							oauthAccessTokenModel.EXPECT().One(gomock.Any(), oauthRefreshToken.OauthAccessTokenID, mockTx).Return(entity.OauthAccessToken{}, exception.Throw(sql.ErrNoRows, exception.WithType(exception.NotFound))),
							oauthAccessTokenModel.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).Times(0),
							oauthAccessTokenModel.EXPECT().One(gomock.Any(), gomock.Any(), gomock.Any()).Times(0),
							oauthRefreshTokenModel.EXPECT().Revoke(gomock.Any(), gomock.Any(), gomock.Any()).Times(0),
							oauthRefreshTokenModel.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).Times(0),
						)
						return f(mockTx)
					})

					mockCall = append(mockCall, sqlTransactionCall)
					mockCall = append(mockCall, insideTransactionCall...)

					gomock.InOrder(mockCall...)

					authorizationService := service.NewAuthorization(oauthApplicationModel, oauthAccessTokenModel, oauthAccessGrantModel, oauthRefreshTokenModel, modelFormatter, oauthValidator, oauthFormatter, true, sqldb)
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

					mockCall := []*gomock.Call{}
					mockCall = append(mockCall,
						oauthApplicationModel.EXPECT().
							OneByUIDandSecret(gomock.Any(), *accessTokenRequest.ClientUID, *accessTokenRequest.ClientSecret, sqldb).
							Return(oauthApplication, nil),
						oauthValidator.EXPECT().ValidateTokenGrant(ctx, accessTokenRequest).Return(nil),
					)

					insideTransactionCall := []*gomock.Call{}

					sqlTransactionCall := sqldb.EXPECT().Transaction(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(ctx kontext.Context, transactionKey string, f func(tx db.TX) exception.Exception) exception.Exception {
						assert.Equal(t, "authorization-grant-token-from-refresh-token", transactionKey)

						insideTransactionCall = append(insideTransactionCall,
							oauthRefreshTokenModel.EXPECT().OneByToken(gomock.Any(), *accessTokenRequest.RefreshToken, mockTx).Return(oauthRefreshToken, nil),
							oauthValidator.EXPECT().ValidateTokenRefreshToken(gomock.Any(), oauthRefreshToken).Return(nil),
							oauthAccessTokenModel.EXPECT().One(gomock.Any(), oauthRefreshToken.OauthAccessTokenID, mockTx).Return(oldAccessToken, nil),
							oauthAccessTokenModel.EXPECT().Create(gomock.Any(), gomock.Any(), mockTx).Return(0, exception.Throw(errors.New("unexpected error"))),
							oauthAccessTokenModel.EXPECT().One(gomock.Any(), gomock.Any(), gomock.Any()).Times(0),
							oauthRefreshTokenModel.EXPECT().Revoke(gomock.Any(), gomock.Any(), gomock.Any()).Times(0),
							oauthRefreshTokenModel.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).Times(0),
						)
						return f(mockTx)
					})

					mockCall = append(mockCall, sqlTransactionCall)
					mockCall = append(mockCall, insideTransactionCall...)

					gomock.InOrder(mockCall...)

					authorizationService := service.NewAuthorization(oauthApplicationModel, oauthAccessTokenModel, oauthAccessGrantModel, oauthRefreshTokenModel, modelFormatter, oauthValidator, oauthFormatter, true, sqldb)
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

					mockCall := []*gomock.Call{}
					mockCall = append(mockCall,
						oauthApplicationModel.EXPECT().
							OneByUIDandSecret(gomock.Any(), *accessTokenRequest.ClientUID, *accessTokenRequest.ClientSecret, sqldb).
							Return(oauthApplication, nil),
						oauthValidator.EXPECT().ValidateTokenGrant(ctx, accessTokenRequest).Return(nil),
					)

					insideTransactionCall := []*gomock.Call{}

					sqlTransactionCall := sqldb.EXPECT().Transaction(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(ctx kontext.Context, transactionKey string, f func(tx db.TX) exception.Exception) exception.Exception {
						assert.Equal(t, "authorization-grant-token-from-refresh-token", transactionKey)

						insideTransactionCall = append(insideTransactionCall,
							oauthRefreshTokenModel.EXPECT().OneByToken(gomock.Any(), *accessTokenRequest.RefreshToken, mockTx).Return(oauthRefreshToken, nil),
							oauthValidator.EXPECT().ValidateTokenRefreshToken(gomock.Any(), oauthRefreshToken).Return(nil),
							oauthAccessTokenModel.EXPECT().One(gomock.Any(), oauthRefreshToken.OauthAccessTokenID, mockTx).Return(oldAccessToken, nil),
							oauthAccessTokenModel.EXPECT().Create(gomock.Any(), gomock.Any(), mockTx).DoAndReturn(func(ktx kontext.Context, insertable entity.OauthAccessTokenInsertable, tx db.TX) (int, error) {
								assert.Equal(t, oauthAccessTokenInsertable.ResourceOwnerID, insertable.ResourceOwnerID)
								assert.Equal(t, oauthAccessTokenInsertable.OauthApplicationID, insertable.OauthApplicationID)
								return 1000, nil
							}),
							oauthAccessTokenModel.EXPECT().One(gomock.Any(), 1000, mockTx).Return(entity.OauthAccessToken{}, exception.Throw(errors.New("unexpected error"))),
							oauthRefreshTokenModel.EXPECT().Revoke(gomock.Any(), gomock.Any(), gomock.Any()).Times(0),
							oauthRefreshTokenModel.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).Times(0),
						)
						return f(mockTx)
					})

					mockCall = append(mockCall, sqlTransactionCall)
					mockCall = append(mockCall, insideTransactionCall...)

					gomock.InOrder(mockCall...)

					authorizationService := service.NewAuthorization(oauthApplicationModel, oauthAccessTokenModel, oauthAccessGrantModel, oauthRefreshTokenModel, modelFormatter, oauthValidator, oauthFormatter, true, sqldb)
					_, err := authorizationService.Token(ctx, accessTokenRequest)

					expectedError := &entity.Error{
						HttpStatus: http.StatusInternalServerError,
						Errors:     eobject.Wrap(eobject.InternalServerError(ctx)),
					}

					assert.Equal(t, expectedError, err)
				})
			})

			t.Run("When refresh token request valid and there is error when revoking old refresh token", func(t *testing.T) {
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

					mockCall := []*gomock.Call{}
					mockCall = append(mockCall,
						oauthApplicationModel.EXPECT().
							OneByUIDandSecret(gomock.Any(), *accessTokenRequest.ClientUID, *accessTokenRequest.ClientSecret, sqldb).
							Return(oauthApplication, nil),
						oauthValidator.EXPECT().ValidateTokenGrant(ctx, accessTokenRequest).Return(nil),
					)

					insideTransactionCall := []*gomock.Call{}

					sqlTransactionCall := sqldb.EXPECT().Transaction(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(ctx kontext.Context, transactionKey string, f func(tx db.TX) exception.Exception) exception.Exception {
						assert.Equal(t, "authorization-grant-token-from-refresh-token", transactionKey)

						insideTransactionCall = append(insideTransactionCall,
							oauthRefreshTokenModel.EXPECT().OneByToken(gomock.Any(), *accessTokenRequest.RefreshToken, mockTx).Return(oauthRefreshToken, nil),
							oauthValidator.EXPECT().ValidateTokenRefreshToken(gomock.Any(), oauthRefreshToken).Return(nil),
							oauthAccessTokenModel.EXPECT().One(gomock.Any(), oauthRefreshToken.OauthAccessTokenID, mockTx).Return(oldAccessToken, nil),
							oauthAccessTokenModel.EXPECT().Create(gomock.Any(), gomock.Any(), mockTx).DoAndReturn(func(ktx kontext.Context, insertable entity.OauthAccessTokenInsertable, tx db.TX) (int, error) {
								assert.Equal(t, oauthAccessTokenInsertable.ResourceOwnerID, insertable.ResourceOwnerID)
								assert.Equal(t, oauthAccessTokenInsertable.OauthApplicationID, insertable.OauthApplicationID)
								return 1000, nil
							}),
							oauthAccessTokenModel.EXPECT().One(gomock.Any(), 1000, mockTx).Return(oauthAccessToken, nil),
							oauthRefreshTokenModel.EXPECT().Revoke(gomock.Any(), oauthRefreshToken.Token, mockTx).Return(exception.Throw(errors.New("unexpected error"))),
							oauthRefreshTokenModel.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).Times(0),
						)
						return f(mockTx)
					})

					mockCall = append(mockCall, sqlTransactionCall)
					mockCall = append(mockCall, insideTransactionCall...)

					gomock.InOrder(mockCall...)

					authorizationService := service.NewAuthorization(oauthApplicationModel, oauthAccessTokenModel, oauthAccessGrantModel, oauthRefreshTokenModel, modelFormatter, oauthValidator, oauthFormatter, true, sqldb)
					_, err := authorizationService.Token(ctx, accessTokenRequest)

					expectedError := &entity.Error{
						HttpStatus: http.StatusInternalServerError,
						Errors:     eobject.Wrap(eobject.InternalServerError(ctx)),
					}

					assert.Equal(t, expectedError, err)
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

					mockCall := []*gomock.Call{}
					mockCall = append(mockCall,
						oauthApplicationModel.EXPECT().
							OneByUIDandSecret(gomock.Any(), *accessTokenRequest.ClientUID, *accessTokenRequest.ClientSecret, sqldb).
							Return(oauthApplication, nil),
						oauthValidator.EXPECT().ValidateTokenGrant(ctx, accessTokenRequest).Return(nil),
					)

					insideTransactionCall := []*gomock.Call{}

					sqlTransactionCall := sqldb.EXPECT().Transaction(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(ctx kontext.Context, transactionKey string, f func(tx db.TX) exception.Exception) exception.Exception {
						assert.Equal(t, "authorization-grant-token-from-refresh-token", transactionKey)

						insideTransactionCall = append(insideTransactionCall,
							oauthRefreshTokenModel.EXPECT().OneByToken(gomock.Any(), *accessTokenRequest.RefreshToken, mockTx).Return(oauthRefreshToken, nil),
							oauthValidator.EXPECT().ValidateTokenRefreshToken(gomock.Any(), oauthRefreshToken).Return(nil),
							oauthAccessTokenModel.EXPECT().One(gomock.Any(), oauthRefreshToken.OauthAccessTokenID, mockTx).Return(oldAccessToken, nil),
							oauthAccessTokenModel.EXPECT().Create(gomock.Any(), gomock.Any(), mockTx).DoAndReturn(func(ktx kontext.Context, insertable entity.OauthAccessTokenInsertable, tx db.TX) (int, error) {
								assert.Equal(t, oauthAccessTokenInsertable.ResourceOwnerID, insertable.ResourceOwnerID)
								assert.Equal(t, oauthAccessTokenInsertable.OauthApplicationID, insertable.OauthApplicationID)
								return 1000, nil
							}),
							oauthAccessTokenModel.EXPECT().One(gomock.Any(), 1000, mockTx).Return(oauthAccessToken, nil),
							oauthRefreshTokenModel.EXPECT().Revoke(gomock.Any(), oauthRefreshToken.Token, mockTx).Return(nil),
							oauthRefreshTokenModel.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).Return(0, exception.Throw(errors.New("unexpected error"))),
						)
						return f(mockTx)
					})

					mockCall = append(mockCall, sqlTransactionCall)
					mockCall = append(mockCall, insideTransactionCall...)

					gomock.InOrder(mockCall...)

					authorizationService := service.NewAuthorization(oauthApplicationModel, oauthAccessTokenModel, oauthAccessGrantModel, oauthRefreshTokenModel, modelFormatter, oauthValidator, oauthFormatter, true, sqldb)
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

					mockCall := []*gomock.Call{}
					mockCall = append(mockCall,
						oauthApplicationModel.EXPECT().
							OneByUIDandSecret(gomock.Any(), *accessTokenRequest.ClientUID, *accessTokenRequest.ClientSecret, sqldb).
							Return(oauthApplication, nil),
						oauthValidator.EXPECT().ValidateTokenGrant(ctx, accessTokenRequest).Return(nil),
					)

					insideTransactionCall := []*gomock.Call{}

					sqlTransactionCall := sqldb.EXPECT().Transaction(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(ctx kontext.Context, transactionKey string, f func(tx db.TX) exception.Exception) exception.Exception {
						assert.Equal(t, "authorization-grant-token-from-refresh-token", transactionKey)

						insideTransactionCall = append(insideTransactionCall,
							oauthRefreshTokenModel.EXPECT().OneByToken(gomock.Any(), *accessTokenRequest.RefreshToken, mockTx).Return(oauthRefreshToken, nil),
							oauthValidator.EXPECT().ValidateTokenRefreshToken(gomock.Any(), oauthRefreshToken).Return(nil),
							oauthAccessTokenModel.EXPECT().One(gomock.Any(), oauthRefreshToken.OauthAccessTokenID, mockTx).Return(oldAccessToken, nil),
							oauthAccessTokenModel.EXPECT().Create(gomock.Any(), gomock.Any(), mockTx).DoAndReturn(func(ktx kontext.Context, insertable entity.OauthAccessTokenInsertable, tx db.TX) (int, error) {
								assert.Equal(t, oauthAccessTokenInsertable.ResourceOwnerID, insertable.ResourceOwnerID)
								assert.Equal(t, oauthAccessTokenInsertable.OauthApplicationID, insertable.OauthApplicationID)
								return 1000, nil
							}),
							oauthAccessTokenModel.EXPECT().One(gomock.Any(), 1000, mockTx).Return(oauthAccessToken, nil),
							oauthRefreshTokenModel.EXPECT().Revoke(gomock.Any(), oauthRefreshToken.Token, mockTx).Return(nil),
							oauthRefreshTokenModel.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(ktx kontext.Context, insertable entity.OauthRefreshTokenInsertable, tx db.TX) (int, error) {
								assert.Equal(t, oauthRefreshTokenInsertable.OauthAccessTokenID, insertable.OauthAccessTokenID)
								return 2, nil
							}),
							oauthRefreshTokenModel.EXPECT().One(gomock.Any(), 2, mockTx).Return(entity.OauthRefreshToken{}, exception.Throw(errors.New("unexpected error"))),
						)
						return f(mockTx)
					})

					mockCall = append(mockCall, sqlTransactionCall)
					mockCall = append(mockCall, insideTransactionCall...)

					gomock.InOrder(mockCall...)

					authorizationService := service.NewAuthorization(oauthApplicationModel, oauthAccessTokenModel, oauthAccessGrantModel, oauthRefreshTokenModel, modelFormatter, oauthValidator, oauthFormatter, true, sqldb)
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
