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
	"github.com/google/uuid"
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

func TestAuthorizationGrantor(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	sqldb := mockdb.NewMockDB(mockCtrl)
	mockTx := mockdb.NewMockTX(mockCtrl)

	t.Run("Grantor", func(t *testing.T) {
		t.Run("Given context and authorization request with a response type of code", func(t *testing.T) {
			t.Run("Return entity.OauthAccessGrantJSON and nil", func(t *testing.T) {
				oauthApplicationModel := mock.NewMockOauthApplicationModel(mockCtrl)
				oauthAccessTokenModel := mock.NewMockOauthAccessTokenModel(mockCtrl)
				oauthAccessGrantModel := mock.NewMockOauthAccessGrantModel(mockCtrl)
				oauthRefreshTokenModel := mock.NewMockOauthRefreshTokenModel(mockCtrl)
				oauthValidator := mock.NewMockOauthValidator(mockCtrl)
				modelFormatter := formatter.NewModel(time.Hour*4, time.Hour*2, time.Hour*2)
				modelFormatterMock := mock.NewMockModelFormater(mockCtrl)
				oauthFormatter := formatter.Oauth()
				oauthFormatterMock := mock.NewMockOauthFormatter(mockCtrl)

				ctx := context.WithValue(context.Background(), "request_id", uuid.New().String())

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

				mockCall := []*gomock.Call{}
				mockCall = append(mockCall,
					oauthApplicationModel.EXPECT().
						OneByUIDandSecret(gomock.Any(), *authorizationRequest.ClientUID, *authorizationRequest.ClientSecret, sqldb).
						Return(oauthApplication, nil),
					oauthValidator.EXPECT().ValidateAuthorizationGrant(ctx, authorizationRequest, oauthApplication).Return(nil),
				)

				insideTransactionCall := []*gomock.Call{}

				sqlTransactionCall := sqldb.EXPECT().Transaction(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(ctx kontext.Context, transactionKey string, f func(tx db.TX) exception.Exception) exception.Exception {
					assert.Equal(t, "authorization-grant-authorization-code", transactionKey)

					insideTransactionCall = append(insideTransactionCall,
						modelFormatterMock.EXPECT().AccessGrantFromAuthorizationRequest(authorizationRequest, oauthApplication).Return(oauthAccessGrantInsertable),
						oauthAccessGrantModel.EXPECT().Create(gomock.Any(), oauthAccessGrantInsertable, mockTx).Return(oauthAccessGrant.ID, nil),
						oauthAccessGrantModel.EXPECT().One(gomock.Any(), oauthAccessGrant.ID, mockTx).Return(oauthAccessGrant, nil),
						oauthFormatterMock.EXPECT().AccessGrant(oauthAccessGrant).Return(oauthAccessGrantJSON),
					)
					return f(mockTx)
				})

				mockCall = append(mockCall, sqlTransactionCall)
				mockCall = append(mockCall, insideTransactionCall...)

				gomock.InOrder(mockCall...)

				authorizationService := service.NewAuthorization(oauthApplicationModel, oauthAccessTokenModel, oauthAccessGrantModel, oauthRefreshTokenModel, modelFormatterMock, oauthValidator, oauthFormatterMock, false, sqldb)
				results, err := authorizationService.Grantor(ctx, authorizationRequest)
				assert.Nil(t, err)
				assert.Equal(t, oauthAccessGrantJSON, results)
			})

			t.Run("Oauth authorization grant validation failed", func(t *testing.T) {
				t.Run("Return entity.OauthAccessGrantJSON and nil", func(t *testing.T) {
					oauthApplicationModel := mock.NewMockOauthApplicationModel(mockCtrl)
					oauthAccessTokenModel := mock.NewMockOauthAccessTokenModel(mockCtrl)
					oauthAccessGrantModel := mock.NewMockOauthAccessGrantModel(mockCtrl)
					oauthRefreshTokenModel := mock.NewMockOauthRefreshTokenModel(mockCtrl)
					oauthValidator := mock.NewMockOauthValidator(mockCtrl)
					modelFormatterMock := mock.NewMockModelFormater(mockCtrl)
					oauthFormatterMock := mock.NewMockOauthFormatter(mockCtrl)

					ctx := context.WithValue(context.Background(), "request_id", uuid.New().String())

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

					mockCall := []*gomock.Call{}
					mockCall = append(mockCall,
						oauthApplicationModel.EXPECT().
							OneByUIDandSecret(gomock.Any(), *authorizationRequest.ClientUID, *authorizationRequest.ClientSecret, sqldb).
							Return(oauthApplication, nil),
						oauthValidator.EXPECT().ValidateAuthorizationGrant(ctx, authorizationRequest, oauthApplication).Return(expectedError),
					)

					insideTransactionCall := []*gomock.Call{}

					sqlTransactionCall := sqldb.EXPECT().Transaction(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(ctx kontext.Context, transactionKey string, f func(tx db.TX) exception.Exception) exception.Exception {
						assert.Equal(t, "authorization-grant-authorization-code", transactionKey)

						insideTransactionCall = append(insideTransactionCall,
							modelFormatterMock.EXPECT().AccessGrantFromAuthorizationRequest(gomock.Any(), gomock.Any()).Times(0),
							oauthAccessGrantModel.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).Times(0),
							oauthAccessGrantModel.EXPECT().One(gomock.Any(), gomock.Any(), gomock.Any()).Times(0),
							oauthFormatterMock.EXPECT().AccessGrant(gomock.Any()).Times(0),
						)
						return f(mockTx)
					}).Times(0)

					mockCall = append(mockCall, sqlTransactionCall)
					mockCall = append(mockCall, insideTransactionCall...)

					gomock.InOrder(mockCall...)

					authorizationService := service.NewAuthorization(oauthApplicationModel, oauthAccessTokenModel, oauthAccessGrantModel, oauthRefreshTokenModel, modelFormatterMock, oauthValidator, oauthFormatterMock, false, sqldb)
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
						oauthRefreshTokenModel := mock.NewMockOauthRefreshTokenModel(mockCtrl)
						oauthValidator := mock.NewMockOauthValidator(mockCtrl)
						modelFormatterMock := mock.NewMockModelFormater(mockCtrl)
						oauthFormatterMock := mock.NewMockOauthFormatter(mockCtrl)

						ctx := context.WithValue(context.Background(), "request_id", uuid.New().String())

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

						mockCall := []*gomock.Call{}
						mockCall = append(mockCall,
							oauthApplicationModel.EXPECT().OneByUIDandSecret(gomock.Any(), *authorizationRequest.ClientUID, *authorizationRequest.ClientSecret, sqldb).Return(entity.OauthApplication{}, exception.Throw(sql.ErrNoRows, exception.WithType(exception.NotFound))),
							oauthValidator.EXPECT().ValidateAuthorizationGrant(gomock.Any(), gomock.Any(), gomock.Any()).Times(0),
						)

						insideTransactionCall := []*gomock.Call{}

						sqlTransactionCall := sqldb.EXPECT().Transaction(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(ctx kontext.Context, transactionKey string, f func(tx db.TX) exception.Exception) exception.Exception {
							assert.Equal(t, "authorization-grant-authorization-code", transactionKey)

							insideTransactionCall = append(insideTransactionCall,
								modelFormatterMock.EXPECT().AccessGrantFromAuthorizationRequest(gomock.Any(), gomock.Any()).Times(0),
								oauthAccessGrantModel.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).Times(0),
								oauthAccessGrantModel.EXPECT().One(gomock.Any(), gomock.Any(), gomock.Any()).Times(0),
								oauthFormatterMock.EXPECT().AccessGrant(gomock.Any()).Times(0),
							)
							return f(mockTx)
						}).Times(0)

						mockCall = append(mockCall, sqlTransactionCall)
						mockCall = append(mockCall, insideTransactionCall...)

						gomock.InOrder(mockCall...)

						authorizationService := service.NewAuthorization(oauthApplicationModel, oauthAccessTokenModel, oauthAccessGrantModel, oauthRefreshTokenModel, modelFormatterMock, oauthValidator, oauthFormatterMock, false, sqldb)
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
					oauthRefreshTokenModel := mock.NewMockOauthRefreshTokenModel(mockCtrl)
					oauthValidator := mock.NewMockOauthValidator(mockCtrl)
					modelFormatter := formatter.NewModel(time.Hour*4, time.Hour*2, time.Hour*2)
					modelFormatterMock := mock.NewMockModelFormater(mockCtrl)
					oauthFormatterMock := mock.NewMockOauthFormatter(mockCtrl)

					ctx := context.WithValue(context.Background(), "request_id", uuid.New().String())

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

					mockCall := []*gomock.Call{}
					mockCall = append(mockCall,
						oauthApplicationModel.EXPECT().
							OneByUIDandSecret(gomock.Any(), *authorizationRequest.ClientUID, *authorizationRequest.ClientSecret, sqldb).
							Return(oauthApplication, nil),
						oauthValidator.EXPECT().ValidateAuthorizationGrant(ctx, authorizationRequest, oauthApplication).Return(nil),
					)

					insideTransactionCall := []*gomock.Call{}

					sqlTransactionCall := sqldb.EXPECT().Transaction(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(ctx kontext.Context, transactionKey string, f func(tx db.TX) exception.Exception) exception.Exception {
						assert.Equal(t, "authorization-grant-authorization-code", transactionKey)

						insideTransactionCall = append(insideTransactionCall,
							modelFormatterMock.EXPECT().AccessGrantFromAuthorizationRequest(authorizationRequest, oauthApplication).Return(oauthAccessGrantInsertable),
							oauthAccessGrantModel.EXPECT().Create(gomock.Any(), oauthAccessGrantInsertable, mockTx).Return(0, exception.Throw(exception.Throw(errors.New("unexpected error")))),
							oauthAccessGrantModel.EXPECT().One(gomock.Any(), gomock.Any(), gomock.Any()).Times(0),
							oauthFormatterMock.EXPECT().AccessGrant(gomock.Any()).Times(0),
						)
						return f(mockTx)
					})

					mockCall = append(mockCall, sqlTransactionCall)
					mockCall = append(mockCall, insideTransactionCall...)

					expectedErr := &entity.Error{
						HttpStatus: http.StatusInternalServerError,
						Errors:     eobject.Wrap(eobject.InternalServerError(ctx)),
					}

					authorizationService := service.NewAuthorization(oauthApplicationModel, oauthAccessTokenModel, oauthAccessGrantModel, oauthRefreshTokenModel, modelFormatterMock, oauthValidator, oauthFormatterMock, false, sqldb)
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
					oauthRefreshTokenModel := mock.NewMockOauthRefreshTokenModel(mockCtrl)
					oauthValidator := mock.NewMockOauthValidator(mockCtrl)
					modelFormatter := formatter.NewModel(time.Hour*4, time.Hour*2, time.Hour*2)
					modelFormatterMock := mock.NewMockModelFormater(mockCtrl)
					oauthFormatterMock := mock.NewMockOauthFormatter(mockCtrl)

					ctx := context.WithValue(context.Background(), "request_id", uuid.New().String())

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

					mockCall := []*gomock.Call{}
					mockCall = append(mockCall,
						oauthApplicationModel.EXPECT().
							OneByUIDandSecret(gomock.Any(), *authorizationRequest.ClientUID, *authorizationRequest.ClientSecret, sqldb).
							Return(oauthApplication, nil),
						oauthValidator.EXPECT().ValidateAuthorizationGrant(ctx, authorizationRequest, oauthApplication).Return(nil),
					)

					insideTransactionCall := []*gomock.Call{}

					sqlTransactionCall := sqldb.EXPECT().Transaction(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(ctx kontext.Context, transactionKey string, f func(tx db.TX) exception.Exception) exception.Exception {
						assert.Equal(t, "authorization-grant-authorization-code", transactionKey)

						insideTransactionCall = append(insideTransactionCall,
							modelFormatterMock.EXPECT().AccessGrantFromAuthorizationRequest(authorizationRequest, oauthApplication).Return(oauthAccessGrantInsertable),
							oauthAccessGrantModel.EXPECT().Create(gomock.Any(), oauthAccessGrantInsertable, mockTx).Return(oauthAccessGrant.ID, nil),
							oauthAccessGrantModel.EXPECT().One(gomock.Any(), oauthAccessGrant.ID, mockTx).Return(entity.OauthAccessGrant{}, exception.Throw(exception.Throw(errors.New("unexpected error")))),
							oauthFormatterMock.EXPECT().AccessGrant(gomock.Any()).Times(0),
						)
						return f(mockTx)
					})

					mockCall = append(mockCall, sqlTransactionCall)
					mockCall = append(mockCall, insideTransactionCall...)

					expectedErr := &entity.Error{
						HttpStatus: http.StatusInternalServerError,
						Errors:     eobject.Wrap(eobject.InternalServerError(ctx)),
					}

					authorizationService := service.NewAuthorization(oauthApplicationModel, oauthAccessTokenModel, oauthAccessGrantModel, oauthRefreshTokenModel, modelFormatterMock, oauthValidator, oauthFormatterMock, false, sqldb)
					results, err := authorizationService.Grantor(ctx, authorizationRequest)
					assert.NotNil(t, err)
					assert.Equal(t, expectedErr, err)
					assert.Equal(t, entity.OauthAccessGrantJSON{}, results)
				})
			})
		})
	})

}
