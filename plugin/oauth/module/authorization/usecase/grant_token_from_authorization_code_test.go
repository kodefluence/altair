package usecase_test

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/kodefluence/altair/plugin/oauth/entity"
	"github.com/kodefluence/altair/plugin/oauth/module/authorization/usecase"
	"github.com/kodefluence/altair/util"
	"github.com/kodefluence/monorepo/db"
	"github.com/kodefluence/monorepo/exception"
	"github.com/kodefluence/monorepo/kontext"
	"github.com/stretchr/testify/suite"
)

type GrantTokenFromAuthorizationCodeTest struct {
	*AuthorizationBaseSuiteTest

	oauthApplication       entity.OauthApplication
	accessTokenRequestJSON entity.AccessTokenRequestJSON
	accessGrant            entity.OauthAccessGrant
	accessToken            entity.OauthAccessToken
	refreshToken           entity.OauthRefreshToken
	refreshTokenJSON       entity.OauthRefreshTokenJSON
}

func TestGrantTokenFromAuthorizationCode(t *testing.T) {
	suite.Run(t, &GrantTokenFromAuthorizationCodeTest{
		AuthorizationBaseSuiteTest: &AuthorizationBaseSuiteTest{},
	})
}

func (suite *GrantTokenFromAuthorizationCodeTest) SetupTest() {
	suite.oauthApplication = entity.OauthApplication{
		ID: 1,
		Scopes: sql.NullString{
			String: "public users",
			Valid:  true,
		},
		OwnerType: "confidential",
	}
	suite.accessTokenRequestJSON = entity.AccessTokenRequestJSON{
		GrantType:    util.StringToPointer("authorization_code"),
		ClientUID:    util.StringToPointer("client_uid"),
		ClientSecret: util.StringToPointer("client_secret"),
		RefreshToken: util.StringToPointer("some-refresh-token"),
		Code:         util.StringToPointer("some-code"),
		RedirectURI:  util.StringToPointer("https://github.com/kodefluence/altair"),
	}
	suite.accessGrant = entity.OauthAccessGrant{
		ID:                 1,
		OauthApplicationID: 1,
		ResourceOwnerID:    0,
		Code:               "some-authorization-code",
		RedirectURI: sql.NullString{
			String: *suite.accessTokenRequestJSON.RedirectURI,
			Valid:  true,
		},
		Scopes:    sql.NullString{},
		ExpiresIn: time.Now().Add(time.Hour),
		CreatedAt: time.Now().Add(-24 * time.Hour),
		RevokedAT: sql.NullTime{},
	}
	suite.accessToken = entity.OauthAccessToken{
		ID:                 1,
		OauthApplicationID: 1,
		ResourceOwnerID:    0,
		Token:              "some random token",
		Scopes:             sql.NullString{},
		ExpiresIn:          time.Time{},
		CreatedAt:          time.Time{},
		RevokedAT:          sql.NullTime{},
	}
	suite.refreshToken = entity.OauthRefreshToken{
		ID:                 1,
		OauthAccessTokenID: 1,
		Token:              "some token",
		ExpiresIn:          time.Time{},
		CreatedAt:          time.Time{},
		RevokedAT:          sql.NullTime{},
	}
}

func (suite *GrantTokenFromAuthorizationCodeTest) TestValidateTokenGrantSuiteTest() {
	suite.Run("Positive cases", func() {
		suite.Subtest("When all parameters is valid, then it would return nil", func() {
			suite.refreshTokenJSON = suite.formatter.RefreshToken(suite.refreshToken)
			gomock.InOrder(
				suite.oauthApplicationRepo.EXPECT().OneByUIDandSecret(suite.ktx, *suite.accessTokenRequestJSON.ClientUID, *suite.accessTokenRequestJSON.ClientSecret, suite.sqldb).Return(suite.oauthApplication, nil),
				suite.sqldb.EXPECT().Transaction(suite.ktx, "authorization-grant-token-from-refresh-token", gomock.Any()).DoAndReturn(func(ctx kontext.Context, transactionKey string, f func(tx db.TX) exception.Exception) exception.Exception {
					suite.oauthAccessGrantRepo.EXPECT().OneByCode(suite.ktx, *suite.accessTokenRequestJSON.Code, suite.sqldb).Return(suite.accessGrant, nil)
					suite.oauthAccessTokenRepo.EXPECT().Create(suite.ktx, gomock.Any(), suite.sqldb).DoAndReturn(func(ktx kontext.Context, data entity.OauthAccessTokenInsertable, tx db.TX) (int, exception.Exception) {
						suite.Assert().Equal(suite.accessGrant.Scopes.String, data.Scopes)
						suite.Assert().Equal(suite.oauthApplication.ID, data.OauthApplicationID)
						return 1, nil
					})
					suite.oauthAccessTokenRepo.EXPECT().One(suite.ktx, 1, suite.sqldb).Return(suite.accessToken, nil)
					suite.oauthAccessGrantRepo.EXPECT().Revoke(suite.ktx, *suite.accessTokenRequestJSON.Code, suite.sqldb).Return(nil)
					suite.oauthRefreshTokenRepo.EXPECT().Create(suite.ktx, gomock.Any(), suite.sqldb).DoAndReturn(func(ktx kontext.Context, data entity.OauthRefreshTokenInsertable, tx db.TX) (int, exception.Exception) {
						suite.Assert().Equal(suite.accessToken.ID, data.OauthAccessTokenID)
						return 1, nil
					})
					suite.oauthRefreshTokenRepo.EXPECT().One(suite.ktx, suite.refreshToken.ID, suite.sqldb).Return(suite.refreshToken, nil)
					return f(suite.sqldb)
				}),
			)

			accessTokenJSON, err := suite.authorization.GrantToken(suite.ktx, suite.accessTokenRequestJSON)
			byteAccessToken, _ := json.Marshal(accessTokenJSON)
			byteExpectedAccessToken, _ := json.Marshal(suite.formatter.AccessToken(suite.accessToken, suite.accessGrant.RedirectURI.String, &suite.refreshTokenJSON))
			suite.Assert().Nil(err)
			suite.Equal(string(byteExpectedAccessToken), string(byteAccessToken))
		})

		suite.Subtest("When all parameters is valid but refresh token is inactive, then it would return nil", func() {
			suite.config.Config.RefreshToken.Active = false
			suite.authorization = usecase.NewAuthorization(suite.oauthApplicationRepo, suite.oauthAccessTokenRepo, suite.oauthAccessGrantRepo, suite.oauthRefreshTokenRepo, suite.formatter, suite.config, suite.sqldb, suite.apiError)

			gomock.InOrder(
				suite.oauthApplicationRepo.EXPECT().OneByUIDandSecret(suite.ktx, *suite.accessTokenRequestJSON.ClientUID, *suite.accessTokenRequestJSON.ClientSecret, suite.sqldb).Return(suite.oauthApplication, nil),
				suite.sqldb.EXPECT().Transaction(suite.ktx, "authorization-grant-token-from-refresh-token", gomock.Any()).DoAndReturn(func(ctx kontext.Context, transactionKey string, f func(tx db.TX) exception.Exception) exception.Exception {
					suite.oauthAccessGrantRepo.EXPECT().OneByCode(suite.ktx, *suite.accessTokenRequestJSON.Code, suite.sqldb).Return(suite.accessGrant, nil)
					suite.oauthAccessTokenRepo.EXPECT().Create(suite.ktx, gomock.Any(), suite.sqldb).DoAndReturn(func(ktx kontext.Context, data entity.OauthAccessTokenInsertable, tx db.TX) (int, exception.Exception) {
						suite.Assert().Equal(suite.accessGrant.Scopes.String, data.Scopes)
						suite.Assert().Equal(suite.oauthApplication.ID, data.OauthApplicationID)
						return 1, nil
					})
					suite.oauthAccessTokenRepo.EXPECT().One(suite.ktx, 1, suite.sqldb).Return(suite.accessToken, nil)
					suite.oauthAccessGrantRepo.EXPECT().Revoke(suite.ktx, *suite.accessTokenRequestJSON.Code, suite.sqldb).Return(nil)
					return f(suite.sqldb)
				}),
			)

			accessTokenJSON, err := suite.authorization.GrantToken(suite.ktx, suite.accessTokenRequestJSON)
			byteAccessToken, _ := json.Marshal(accessTokenJSON)
			byteExpectedAccessToken, _ := json.Marshal(suite.formatter.AccessToken(suite.accessToken, suite.accessGrant.RedirectURI.String, nil))
			suite.Assert().Nil(err)
			suite.Equal(string(byteExpectedAccessToken), string(byteAccessToken))
		})
	})

	suite.Run("Negative cases", func() {
		suite.Subtest("When get grant token not found, then it would return error", func() {
			suite.refreshTokenJSON = suite.formatter.RefreshToken(suite.refreshToken)
			gomock.InOrder(
				suite.oauthApplicationRepo.EXPECT().OneByUIDandSecret(suite.ktx, *suite.accessTokenRequestJSON.ClientUID, *suite.accessTokenRequestJSON.ClientSecret, suite.sqldb).Return(suite.oauthApplication, nil),
				suite.sqldb.EXPECT().Transaction(suite.ktx, "authorization-grant-token-from-refresh-token", gomock.Any()).DoAndReturn(func(ctx kontext.Context, transactionKey string, f func(tx db.TX) exception.Exception) exception.Exception {
					suite.oauthAccessGrantRepo.EXPECT().OneByCode(suite.ktx, *suite.accessTokenRequestJSON.Code, suite.sqldb).Return(entity.OauthAccessGrant{}, exception.Throw(errors.New("not found"), exception.WithType(exception.NotFound)))
					return f(suite.sqldb)
				}),
			)

			_, err := suite.authorization.GrantToken(suite.ktx, suite.accessTokenRequestJSON)
			suite.Assert().NotNil(err)
			suite.Assert().Equal("JSONAPI Error:\n[Not found error] Detail: Resource of `authorization_code` is not found. Tracing code: `<nil>`, Code: ERR0404\n", err.Error())
			suite.Assert().Equal(http.StatusNotFound, err.HTTPStatus())
		})

		suite.Subtest("When get grant token failure, then it would return error", func() {
			suite.refreshTokenJSON = suite.formatter.RefreshToken(suite.refreshToken)
			gomock.InOrder(
				suite.oauthApplicationRepo.EXPECT().OneByUIDandSecret(suite.ktx, *suite.accessTokenRequestJSON.ClientUID, *suite.accessTokenRequestJSON.ClientSecret, suite.sqldb).Return(suite.oauthApplication, nil),
				suite.sqldb.EXPECT().Transaction(suite.ktx, "authorization-grant-token-from-refresh-token", gomock.Any()).DoAndReturn(func(ctx kontext.Context, transactionKey string, f func(tx db.TX) exception.Exception) exception.Exception {
					suite.oauthAccessGrantRepo.EXPECT().OneByCode(suite.ktx, *suite.accessTokenRequestJSON.Code, suite.sqldb).Return(entity.OauthAccessGrant{}, exception.Throw(errors.New("unexpected")))
					return f(suite.sqldb)
				}),
			)

			_, err := suite.authorization.GrantToken(suite.ktx, suite.accessTokenRequestJSON)
			suite.Assert().NotNil(err)
			suite.Assert().Equal("JSONAPI Error:\n[Internal server error] Detail: Something is not right, help us fix this problem. Contribute to https://github.com/kodefluence/altair. Tracing code: '<nil>', Code: ERR0500\n", err.Error())
			suite.Assert().Equal(http.StatusInternalServerError, err.HTTPStatus())
		})

		suite.Subtest("When access grant validation failure, then it would return error", func() {
			suite.accessGrant.RevokedAT = sql.NullTime{
				Time:  time.Now().Add(-1 * time.Hour),
				Valid: true,
			}
			suite.refreshTokenJSON = suite.formatter.RefreshToken(suite.refreshToken)
			gomock.InOrder(
				suite.oauthApplicationRepo.EXPECT().OneByUIDandSecret(suite.ktx, *suite.accessTokenRequestJSON.ClientUID, *suite.accessTokenRequestJSON.ClientSecret, suite.sqldb).Return(suite.oauthApplication, nil),
				suite.sqldb.EXPECT().Transaction(suite.ktx, "authorization-grant-token-from-refresh-token", gomock.Any()).DoAndReturn(func(ctx kontext.Context, transactionKey string, f func(tx db.TX) exception.Exception) exception.Exception {
					suite.oauthAccessGrantRepo.EXPECT().OneByCode(suite.ktx, *suite.accessTokenRequestJSON.Code, suite.sqldb).Return(suite.accessGrant, nil)
					return f(suite.sqldb)
				}),
			)

			_, err := suite.authorization.GrantToken(suite.ktx, suite.accessTokenRequestJSON)
			suite.Assert().NotNil(err)
			suite.Assert().Equal("JSONAPI Error:\n[Forbidden resource access] Detail: authorization code already used, Code: ERR0403\n", err.Error())
			suite.Assert().Equal(http.StatusForbidden, err.HTTPStatus())
		})

		suite.Subtest("When access token creation failure, then it would return error", func() {
			suite.refreshTokenJSON = suite.formatter.RefreshToken(suite.refreshToken)
			gomock.InOrder(
				suite.oauthApplicationRepo.EXPECT().OneByUIDandSecret(suite.ktx, *suite.accessTokenRequestJSON.ClientUID, *suite.accessTokenRequestJSON.ClientSecret, suite.sqldb).Return(suite.oauthApplication, nil),
				suite.sqldb.EXPECT().Transaction(suite.ktx, "authorization-grant-token-from-refresh-token", gomock.Any()).DoAndReturn(func(ctx kontext.Context, transactionKey string, f func(tx db.TX) exception.Exception) exception.Exception {
					suite.oauthAccessGrantRepo.EXPECT().OneByCode(suite.ktx, *suite.accessTokenRequestJSON.Code, suite.sqldb).Return(suite.accessGrant, nil)
					suite.oauthAccessTokenRepo.EXPECT().Create(suite.ktx, gomock.Any(), suite.sqldb).DoAndReturn(func(ktx kontext.Context, data entity.OauthAccessTokenInsertable, tx db.TX) (int, exception.Exception) {
						suite.Assert().Equal(suite.accessGrant.Scopes.String, data.Scopes)
						suite.Assert().Equal(suite.oauthApplication.ID, data.OauthApplicationID)
						return 1, exception.Throw(errors.New("unexpected"))
					})
					return f(suite.sqldb)
				}),
			)

			_, err := suite.authorization.GrantToken(suite.ktx, suite.accessTokenRequestJSON)
			suite.Assert().NotNil(err)
			suite.Assert().Equal("JSONAPI Error:\n[Internal server error] Detail: Something is not right, help us fix this problem. Contribute to https://github.com/kodefluence/altair. Tracing code: '<nil>', Code: ERR0500\n", err.Error())
			suite.Assert().Equal(http.StatusInternalServerError, err.HTTPStatus())
		})

		suite.Subtest("When find access token after creation failed, then it would return error", func() {
			suite.refreshTokenJSON = suite.formatter.RefreshToken(suite.refreshToken)
			gomock.InOrder(
				suite.oauthApplicationRepo.EXPECT().OneByUIDandSecret(suite.ktx, *suite.accessTokenRequestJSON.ClientUID, *suite.accessTokenRequestJSON.ClientSecret, suite.sqldb).Return(suite.oauthApplication, nil),
				suite.sqldb.EXPECT().Transaction(suite.ktx, "authorization-grant-token-from-refresh-token", gomock.Any()).DoAndReturn(func(ctx kontext.Context, transactionKey string, f func(tx db.TX) exception.Exception) exception.Exception {
					suite.oauthAccessGrantRepo.EXPECT().OneByCode(suite.ktx, *suite.accessTokenRequestJSON.Code, suite.sqldb).Return(suite.accessGrant, nil)
					suite.oauthAccessTokenRepo.EXPECT().Create(suite.ktx, gomock.Any(), suite.sqldb).DoAndReturn(func(ktx kontext.Context, data entity.OauthAccessTokenInsertable, tx db.TX) (int, exception.Exception) {
						suite.Assert().Equal(suite.accessGrant.Scopes.String, data.Scopes)
						suite.Assert().Equal(suite.oauthApplication.ID, data.OauthApplicationID)
						return 1, nil
					})
					suite.oauthAccessTokenRepo.EXPECT().One(suite.ktx, 1, suite.sqldb).Return(entity.OauthAccessToken{}, exception.Throw(errors.New("unexpected")))
					return f(suite.sqldb)
				}),
			)

			_, err := suite.authorization.GrantToken(suite.ktx, suite.accessTokenRequestJSON)
			suite.Assert().NotNil(err)
			suite.Assert().Equal("JSONAPI Error:\n[Internal server error] Detail: Something is not right, help us fix this problem. Contribute to https://github.com/kodefluence/altair. Tracing code: '<nil>', Code: ERR0500\n", err.Error())
			suite.Assert().Equal(http.StatusInternalServerError, err.HTTPStatus())
		})

		suite.Subtest("When revoke token error, then it would return error", func() {
			suite.refreshTokenJSON = suite.formatter.RefreshToken(suite.refreshToken)
			gomock.InOrder(
				suite.oauthApplicationRepo.EXPECT().OneByUIDandSecret(suite.ktx, *suite.accessTokenRequestJSON.ClientUID, *suite.accessTokenRequestJSON.ClientSecret, suite.sqldb).Return(suite.oauthApplication, nil),
				suite.sqldb.EXPECT().Transaction(suite.ktx, "authorization-grant-token-from-refresh-token", gomock.Any()).DoAndReturn(func(ctx kontext.Context, transactionKey string, f func(tx db.TX) exception.Exception) exception.Exception {
					suite.oauthAccessGrantRepo.EXPECT().OneByCode(suite.ktx, *suite.accessTokenRequestJSON.Code, suite.sqldb).Return(suite.accessGrant, nil)
					suite.oauthAccessTokenRepo.EXPECT().Create(suite.ktx, gomock.Any(), suite.sqldb).DoAndReturn(func(ktx kontext.Context, data entity.OauthAccessTokenInsertable, tx db.TX) (int, exception.Exception) {
						suite.Assert().Equal(suite.accessGrant.Scopes.String, data.Scopes)
						suite.Assert().Equal(suite.oauthApplication.ID, data.OauthApplicationID)
						return 1, nil
					})
					suite.oauthAccessTokenRepo.EXPECT().One(suite.ktx, 1, suite.sqldb).Return(suite.accessToken, nil)
					suite.oauthAccessGrantRepo.EXPECT().Revoke(suite.ktx, *suite.accessTokenRequestJSON.Code, suite.sqldb).Return(exception.Throw(errors.New("unexpected")))
					return f(suite.sqldb)
				}),
			)

			_, err := suite.authorization.GrantToken(suite.ktx, suite.accessTokenRequestJSON)
			suite.Assert().NotNil(err)
			suite.Assert().Equal("JSONAPI Error:\n[Internal server error] Detail: Something is not right, help us fix this problem. Contribute to https://github.com/kodefluence/altair. Tracing code: '<nil>', Code: ERR0500\n", err.Error())
			suite.Assert().Equal(http.StatusInternalServerError, err.HTTPStatus())
		})

		suite.Subtest("When refresh token creation failure, then it would return error", func() {
			suite.refreshTokenJSON = suite.formatter.RefreshToken(suite.refreshToken)
			gomock.InOrder(
				suite.oauthApplicationRepo.EXPECT().OneByUIDandSecret(suite.ktx, *suite.accessTokenRequestJSON.ClientUID, *suite.accessTokenRequestJSON.ClientSecret, suite.sqldb).Return(suite.oauthApplication, nil),
				suite.sqldb.EXPECT().Transaction(suite.ktx, "authorization-grant-token-from-refresh-token", gomock.Any()).DoAndReturn(func(ctx kontext.Context, transactionKey string, f func(tx db.TX) exception.Exception) exception.Exception {
					suite.oauthAccessGrantRepo.EXPECT().OneByCode(suite.ktx, *suite.accessTokenRequestJSON.Code, suite.sqldb).Return(suite.accessGrant, nil)
					suite.oauthAccessTokenRepo.EXPECT().Create(suite.ktx, gomock.Any(), suite.sqldb).DoAndReturn(func(ktx kontext.Context, data entity.OauthAccessTokenInsertable, tx db.TX) (int, exception.Exception) {
						suite.Assert().Equal(suite.accessGrant.Scopes.String, data.Scopes)
						suite.Assert().Equal(suite.oauthApplication.ID, data.OauthApplicationID)
						return 1, nil
					})
					suite.oauthAccessTokenRepo.EXPECT().One(suite.ktx, 1, suite.sqldb).Return(suite.accessToken, nil)
					suite.oauthAccessGrantRepo.EXPECT().Revoke(suite.ktx, *suite.accessTokenRequestJSON.Code, suite.sqldb).Return(nil)
					suite.oauthRefreshTokenRepo.EXPECT().Create(suite.ktx, gomock.Any(), suite.sqldb).DoAndReturn(func(ktx kontext.Context, data entity.OauthRefreshTokenInsertable, tx db.TX) (int, exception.Exception) {
						suite.Assert().Equal(suite.accessToken.ID, data.OauthAccessTokenID)
						return 1, nil
					})
					suite.oauthRefreshTokenRepo.EXPECT().One(suite.ktx, suite.refreshToken.ID, suite.sqldb).Return(entity.OauthRefreshToken{}, exception.Throw(errors.New("unexpected")))
					return f(suite.sqldb)
				}),
			)

			_, err := suite.authorization.GrantToken(suite.ktx, suite.accessTokenRequestJSON)
			suite.Assert().NotNil(err)
			suite.Assert().Equal("JSONAPI Error:\n[Internal server error] Detail: Something is not right, help us fix this problem. Contribute to https://github.com/kodefluence/altair. Tracing code: '<nil>', Code: ERR0500\n", err.Error())
			suite.Assert().Equal(http.StatusInternalServerError, err.HTTPStatus())
		})
	})
}

func (suite *GrantTokenFromAuthorizationCodeTest) Subtest(testcase string, subtest func()) {
	suite.SetupTest()
	suite.AuthorizationBaseSuiteTest.Subtest(testcase, subtest)
	suite.TearDownTest()
}
