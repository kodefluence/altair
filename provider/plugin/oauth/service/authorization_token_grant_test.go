package service_test

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/go-sql-driver/mysql"
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

func TestAuthorizationToken(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	sqldb := mockdb.NewMockDB(mockCtrl)
	mockTx := mockdb.NewMockTX(mockCtrl)

	t.Run("Token", func(t *testing.T) {
		t.Run("Given context and access token request", func(t *testing.T) {
			t.Run("When access token request valid and there is no error in database side", func(t *testing.T) {
				t.Run("Then it will return access token response", func(t *testing.T) {
					oauthApplicationModel := mock.NewMockOauthApplicationModel(mockCtrl)
					oauthAccessTokenModel := mock.NewMockOauthAccessTokenModel(mockCtrl)
					oauthAccessGrantModel := mock.NewMockOauthAccessGrantModel(mockCtrl)
					oauthRefreshTokenModel := mock.NewMockOauthRefreshTokenModel(mockCtrl)
					oauthValidator := mock.NewMockOauthValidator(mockCtrl)
					modelFormatter := formatter.NewModel(time.Hour*4, time.Hour*2, time.Hour*2)
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
					oauthAccessTokenJSON := oauthFormatter.AccessToken(oauthAccessToken, *accessTokenRequest.RedirectURI, nil)

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
							oauthAccessGrantModel.EXPECT().OneByCode(gomock.Any(), *accessTokenRequest.Code, mockTx).Return(oauthAccessGrant, nil),
							oauthValidator.EXPECT().ValidateTokenAuthorizationCode(gomock.Any(), accessTokenRequest, oauthAccessGrant).Return(nil),
							modelFormatterMock.EXPECT().AccessTokenFromOauthAccessGrant(oauthAccessGrant, oauthApplication).Return(oauthAccessTokenInsertable),
							oauthAccessTokenModel.EXPECT().Create(gomock.Any(), oauthAccessTokenInsertable, mockTx).Return(1, nil),
							oauthAccessTokenModel.EXPECT().One(gomock.Any(), 1, mockTx).Return(oauthAccessToken, nil),
							oauthAccessGrantModel.EXPECT().Revoke(gomock.Any(), *accessTokenRequest.Code, mockTx).Return(nil),
							oauthFormatterMock.EXPECT().AccessToken(oauthAccessToken, oauthAccessGrant.RedirectURI.String, nil).Return(oauthAccessTokenJSON),
						)
						return f(mockTx)
					})

					mockCall = append(mockCall, sqlTransactionCall)
					mockCall = append(mockCall, insideTransactionCall...)

					gomock.InOrder(mockCall...)

					authorizationService := service.NewAuthorization(oauthApplicationModel, oauthAccessTokenModel, oauthAccessGrantModel, oauthRefreshTokenModel, modelFormatterMock, oauthValidator, oauthFormatterMock, false, sqldb)
					oauthAccessTokenOutput, err := authorizationService.Token(ctx, accessTokenRequest)

					assert.Nil(t, err)
					assert.Equal(t, oauthAccessTokenJSON, oauthAccessTokenOutput)
				})
			})

			t.Run("When revoking oauth access grant error", func(t *testing.T) {
				t.Run("Then it will return the error", func(t *testing.T) {
					oauthApplicationModel := mock.NewMockOauthApplicationModel(mockCtrl)
					oauthAccessTokenModel := mock.NewMockOauthAccessTokenModel(mockCtrl)
					oauthAccessGrantModel := mock.NewMockOauthAccessGrantModel(mockCtrl)
					oauthRefreshTokenModel := mock.NewMockOauthRefreshTokenModel(mockCtrl)
					oauthValidator := mock.NewMockOauthValidator(mockCtrl)
					modelFormatter := formatter.NewModel(time.Hour*4, time.Hour*2, time.Hour*2)
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

					expectedError := &entity.Error{
						HttpStatus: http.StatusInternalServerError,
						Errors:     eobject.Wrap(eobject.InternalServerError(ctx)),
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
							oauthAccessGrantModel.EXPECT().OneByCode(gomock.Any(), *accessTokenRequest.Code, mockTx).Return(oauthAccessGrant, nil),
							oauthValidator.EXPECT().ValidateTokenAuthorizationCode(gomock.Any(), accessTokenRequest, oauthAccessGrant).Return(nil),
							modelFormatterMock.EXPECT().AccessTokenFromOauthAccessGrant(oauthAccessGrant, oauthApplication).Return(oauthAccessTokenInsertable),
							oauthAccessTokenModel.EXPECT().Create(gomock.Any(), oauthAccessTokenInsertable, mockTx).Return(1, nil),
							oauthAccessTokenModel.EXPECT().One(gomock.Any(), 1, mockTx).Return(oauthAccessToken, nil),
							oauthAccessGrantModel.EXPECT().Revoke(gomock.Any(), *accessTokenRequest.Code, mockTx).Return(exception.Throw(errors.New("unexpected error"))),
							oauthFormatterMock.EXPECT().AccessToken(gomock.Any(), gomock.Any(), gomock.Any()).Times(0),
						)
						return f(mockTx)
					})

					mockCall = append(mockCall, sqlTransactionCall)
					mockCall = append(mockCall, insideTransactionCall...)

					gomock.InOrder(mockCall...)

					authorizationService := service.NewAuthorization(oauthApplicationModel, oauthAccessTokenModel, oauthAccessGrantModel, oauthRefreshTokenModel, modelFormatterMock, oauthValidator, oauthFormatterMock, false, sqldb)
					oauthAccessTokenOutput, err := authorizationService.Token(ctx, accessTokenRequest)

					assert.NotNil(t, err)
					assert.Equal(t, expectedError, err)
					assert.Equal(t, entity.OauthAccessTokenJSON{}, oauthAccessTokenOutput)
				})
			})

			t.Run("When application client and secret is not valid", func(t *testing.T) {
				t.Run("Then it will return unprocessable entity error", func(t *testing.T) {
					oauthApplicationModel := mock.NewMockOauthApplicationModel(mockCtrl)
					oauthAccessTokenModel := mock.NewMockOauthAccessTokenModel(mockCtrl)
					oauthAccessGrantModel := mock.NewMockOauthAccessGrantModel(mockCtrl)
					oauthRefreshTokenModel := mock.NewMockOauthRefreshTokenModel(mockCtrl)
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

					mockCall := []*gomock.Call{}
					mockCall = append(mockCall,
						oauthApplicationModel.EXPECT().
							OneByUIDandSecret(gomock.Any(), *accessTokenRequest.ClientUID, *accessTokenRequest.ClientSecret, sqldb).
							Return(entity.OauthApplication{}, exception.Throw(exception.Throw(errors.New("unexpected error")))),
						oauthValidator.EXPECT().ValidateTokenGrant(gomock.Any(), gomock.Any()).Times(0),
					)

					insideTransactionCall := []*gomock.Call{}

					sqlTransactionCall := sqldb.EXPECT().Transaction(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(ctx kontext.Context, transactionKey string, f func(tx db.TX) exception.Exception) exception.Exception {
						assert.Equal(t, "authorization-grant-token-from-refresh-token", transactionKey)

						insideTransactionCall = append(insideTransactionCall,
							oauthAccessGrantModel.EXPECT().OneByCode(gomock.Any(), gomock.Any(), sqldb).Times(0),
							oauthValidator.EXPECT().ValidateTokenAuthorizationCode(gomock.Any(), gomock.Any(), gomock.Any()).Times(0),
							modelFormatterMock.EXPECT().AccessTokenFromOauthAccessGrant(gomock.Any(), gomock.Any()).Times(0),
							oauthAccessTokenModel.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).Times(0),
							oauthAccessTokenModel.EXPECT().One(gomock.Any(), gomock.Any(), gomock.Any()).Times(0),
							oauthAccessGrantModel.EXPECT().Revoke(gomock.Any(), gomock.Any(), gomock.Any()).Times(0),
							oauthFormatterMock.EXPECT().AccessToken(gomock.Any(), gomock.Any(), gomock.Any()).Times(0),
						)
						return f(mockTx)
					}).Times(0)

					mockCall = append(mockCall, sqlTransactionCall)
					mockCall = append(mockCall, insideTransactionCall...)

					gomock.InOrder(mockCall...)

					authorizationService := service.NewAuthorization(oauthApplicationModel, oauthAccessTokenModel, oauthAccessGrantModel, oauthRefreshTokenModel, modelFormatterMock, oauthValidator, oauthFormatterMock, false, sqldb)
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
					oauthRefreshTokenModel := mock.NewMockOauthRefreshTokenModel(mockCtrl)
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
							OneByUIDandSecret(gomock.Any(), *accessTokenRequest.ClientUID, *accessTokenRequest.ClientSecret, sqldb).
							Return(oauthApplication, nil),
						oauthValidator.EXPECT().ValidateTokenGrant(ctx, accessTokenRequest).Return(expectedError),
						oauthAccessGrantModel.EXPECT().OneByCode(gomock.Any(), gomock.Any(), sqldb).Times(0),
						oauthValidator.EXPECT().ValidateTokenAuthorizationCode(gomock.Any(), gomock.Any(), gomock.Any()).Times(0),
						modelFormatterMock.EXPECT().AccessTokenFromOauthAccessGrant(gomock.Any(), gomock.Any()).Times(0),
						oauthAccessTokenModel.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).Times(0),
						oauthAccessTokenModel.EXPECT().One(gomock.Any(), gomock.Any(), gomock.Any()).Times(0),
						oauthAccessGrantModel.EXPECT().Revoke(gomock.Any(), gomock.Any(), gomock.Any()).Times(0),
						oauthFormatterMock.EXPECT().AccessToken(gomock.Any(), gomock.Any(), gomock.Any()).Times(0),
					)

					authorizationService := service.NewAuthorization(oauthApplicationModel, oauthAccessTokenModel, oauthAccessGrantModel, oauthRefreshTokenModel, modelFormatterMock, oauthValidator, oauthFormatterMock, false, sqldb)
					oauthAccessTokenOutput, err := authorizationService.Token(ctx, accessTokenRequest)

					assert.NotNil(t, err)
					assert.Equal(t, expectedError, err)
					assert.Equal(t, entity.OauthAccessTokenJSON{}, oauthAccessTokenOutput)
				})
			})

			t.Run("When oauth access grants is not valid", func(t *testing.T) {
				t.Run("Then it will return unprocessable entity error", func(t *testing.T) {
					oauthApplicationModel := mock.NewMockOauthApplicationModel(mockCtrl)
					oauthAccessTokenModel := mock.NewMockOauthAccessTokenModel(mockCtrl)
					oauthAccessGrantModel := mock.NewMockOauthAccessGrantModel(mockCtrl)
					oauthRefreshTokenModel := mock.NewMockOauthRefreshTokenModel(mockCtrl)
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

					errorObject := eobject.ForbiddenError(ctx, "access_token", "authorization code already used")

					expectedError := &entity.Error{
						HttpStatus: http.StatusForbidden,
						Errors: eobject.Wrap(
							errorObject,
						),
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
							oauthAccessGrantModel.EXPECT().OneByCode(gomock.Any(), *accessTokenRequest.Code, mockTx).Return(oauthAccessGrant, nil),
							oauthValidator.EXPECT().ValidateTokenAuthorizationCode(gomock.Any(), accessTokenRequest, oauthAccessGrant).Return(exception.Throw(errorObject, exception.WithTitle(errorObject.Code), exception.WithDetail(errorObject.Message), exception.WithType(exception.Forbidden))),
							modelFormatterMock.EXPECT().AccessTokenFromOauthAccessGrant(gomock.Any(), gomock.Any()).Times(0),
							oauthAccessTokenModel.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).Times(0),
							oauthAccessTokenModel.EXPECT().One(gomock.Any(), gomock.Any(), gomock.Any()).Times(0),
							oauthAccessGrantModel.EXPECT().Revoke(gomock.Any(), gomock.Any(), gomock.Any()).Times(0),
							oauthFormatterMock.EXPECT().AccessToken(gomock.Any(), gomock.Any(), gomock.Any()).Times(0),
						)
						return f(mockTx)
					})

					mockCall = append(mockCall, sqlTransactionCall)
					mockCall = append(mockCall, insideTransactionCall...)

					gomock.InOrder(mockCall...)

					authorizationService := service.NewAuthorization(oauthApplicationModel, oauthAccessTokenModel, oauthAccessGrantModel, oauthRefreshTokenModel, modelFormatterMock, oauthValidator, oauthFormatterMock, false, sqldb)
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
					oauthRefreshTokenModel := mock.NewMockOauthRefreshTokenModel(mockCtrl)
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
							oauthAccessGrantModel.EXPECT().OneByCode(gomock.Any(), *accessTokenRequest.Code, mockTx).Return(entity.OauthAccessGrant{}, exception.Throw(sql.ErrNoRows, exception.WithType(exception.NotFound))),
							oauthValidator.EXPECT().ValidateTokenAuthorizationCode(gomock.Any(), gomock.Any(), gomock.Any()).Times(0),
							modelFormatterMock.EXPECT().AccessTokenFromOauthAccessGrant(gomock.Any(), gomock.Any()).Times(0),
							oauthAccessTokenModel.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).Times(0),
							oauthAccessTokenModel.EXPECT().One(gomock.Any(), gomock.Any(), gomock.Any()).Times(0),
							oauthAccessGrantModel.EXPECT().Revoke(gomock.Any(), gomock.Any(), gomock.Any()).Times(0),
							oauthFormatterMock.EXPECT().AccessToken(gomock.Any(), gomock.Any(), gomock.Any()).Times(0),
						)
						return f(mockTx)
					})

					mockCall = append(mockCall, sqlTransactionCall)
					mockCall = append(mockCall, insideTransactionCall...)

					gomock.InOrder(mockCall...)

					authorizationService := service.NewAuthorization(oauthApplicationModel, oauthAccessTokenModel, oauthAccessGrantModel, oauthRefreshTokenModel, modelFormatterMock, oauthValidator, oauthFormatterMock, false, sqldb)
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
					oauthRefreshTokenModel := mock.NewMockOauthRefreshTokenModel(mockCtrl)
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
							oauthAccessGrantModel.EXPECT().OneByCode(gomock.Any(), *accessTokenRequest.Code, mockTx).Return(entity.OauthAccessGrant{}, exception.Throw(errors.New("unexpected error"))),
							oauthValidator.EXPECT().ValidateTokenAuthorizationCode(gomock.Any(), gomock.Any(), gomock.Any()).Times(0),
							modelFormatterMock.EXPECT().AccessTokenFromOauthAccessGrant(gomock.Any(), gomock.Any()).Times(0),
							oauthAccessTokenModel.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).Times(0),
							oauthAccessTokenModel.EXPECT().One(gomock.Any(), gomock.Any(), gomock.Any()).Times(0),
							oauthAccessGrantModel.EXPECT().Revoke(gomock.Any(), gomock.Any(), gomock.Any()).Times(0),
							oauthFormatterMock.EXPECT().AccessToken(gomock.Any(), gomock.Any(), gomock.Any()).Times(0),
						)
						return f(mockTx)
					})

					mockCall = append(mockCall, sqlTransactionCall)
					mockCall = append(mockCall, insideTransactionCall...)

					gomock.InOrder(mockCall...)

					authorizationService := service.NewAuthorization(oauthApplicationModel, oauthAccessTokenModel, oauthAccessGrantModel, oauthRefreshTokenModel, modelFormatterMock, oauthValidator, oauthFormatterMock, false, sqldb)
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
					oauthRefreshTokenModel := mock.NewMockOauthRefreshTokenModel(mockCtrl)
					oauthValidator := mock.NewMockOauthValidator(mockCtrl)
					modelFormatter := formatter.NewModel(time.Hour*4, time.Hour*2, time.Hour*2)
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
							oauthAccessGrantModel.EXPECT().OneByCode(gomock.Any(), *accessTokenRequest.Code, mockTx).Return(oauthAccessGrant, nil),
							oauthValidator.EXPECT().ValidateTokenAuthorizationCode(gomock.Any(), accessTokenRequest, oauthAccessGrant).Return(nil),
							modelFormatterMock.EXPECT().AccessTokenFromOauthAccessGrant(oauthAccessGrant, oauthApplication).Return(oauthAccessTokenInsertable),
							oauthAccessTokenModel.EXPECT().Create(gomock.Any(), oauthAccessTokenInsertable, mockTx).Return(0, exception.Throw(errors.New("unexpected error"))),
							oauthAccessTokenModel.EXPECT().One(gomock.Any(), gomock.Any(), gomock.Any()).Times(0),
							oauthAccessGrantModel.EXPECT().Revoke(gomock.Any(), gomock.Any(), gomock.Any()).Times(0),
							oauthFormatterMock.EXPECT().AccessToken(gomock.Any(), gomock.Any(), gomock.Any()).Times(0),
						)
						return f(mockTx)
					})

					mockCall = append(mockCall, sqlTransactionCall)
					mockCall = append(mockCall, insideTransactionCall...)

					gomock.InOrder(mockCall...)

					authorizationService := service.NewAuthorization(oauthApplicationModel, oauthAccessTokenModel, oauthAccessGrantModel, oauthRefreshTokenModel, modelFormatterMock, oauthValidator, oauthFormatterMock, false, sqldb)
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
					oauthRefreshTokenModel := mock.NewMockOauthRefreshTokenModel(mockCtrl)
					oauthValidator := mock.NewMockOauthValidator(mockCtrl)
					modelFormatter := formatter.NewModel(time.Hour*4, time.Hour*2, time.Hour*2)
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
							oauthAccessGrantModel.EXPECT().OneByCode(gomock.Any(), *accessTokenRequest.Code, mockTx).Return(oauthAccessGrant, nil),
							oauthValidator.EXPECT().ValidateTokenAuthorizationCode(gomock.Any(), accessTokenRequest, oauthAccessGrant).Return(nil),
							modelFormatterMock.EXPECT().AccessTokenFromOauthAccessGrant(oauthAccessGrant, oauthApplication).Return(oauthAccessTokenInsertable),
							oauthAccessTokenModel.EXPECT().Create(gomock.Any(), oauthAccessTokenInsertable, mockTx).Return(1, nil),
							oauthAccessTokenModel.EXPECT().One(gomock.Any(), 1, mockTx).Return(entity.OauthAccessToken{}, exception.Throw(sql.ErrNoRows, exception.WithType(exception.NotFound))),
							oauthAccessGrantModel.EXPECT().Revoke(gomock.Any(), gomock.Any(), gomock.Any()).Times(0),
							oauthFormatterMock.EXPECT().AccessToken(gomock.Any(), gomock.Any(), gomock.Any()).Times(0),
						)
						return f(mockTx)
					})

					mockCall = append(mockCall, sqlTransactionCall)
					mockCall = append(mockCall, insideTransactionCall...)

					gomock.InOrder(mockCall...)

					authorizationService := service.NewAuthorization(oauthApplicationModel, oauthAccessTokenModel, oauthAccessGrantModel, oauthRefreshTokenModel, modelFormatterMock, oauthValidator, oauthFormatterMock, false, sqldb)
					oauthAccessTokenOutput, err := authorizationService.Token(ctx, accessTokenRequest)

					assert.NotNil(t, err)
					assert.Equal(t, expectedError, err)
					assert.Equal(t, entity.OauthAccessTokenJSON{}, oauthAccessTokenOutput)
				})
			})
		})
	})

}
