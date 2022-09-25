package usecase_test

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/golang/mock/gomock"
	"github.com/kodefluence/altair/plugin/oauth/entity"
	"github.com/kodefluence/altair/plugin/oauth/module/authorization/usecase"
	"github.com/kodefluence/altair/util"
	"github.com/kodefluence/monorepo/db"
	"github.com/kodefluence/monorepo/exception"
	"github.com/kodefluence/monorepo/kontext"
	"github.com/stretchr/testify/suite"
)

type ClientCredentialSuiteTest struct {
	*AuthorizationBaseSuiteTest

	oauthApplication       entity.OauthApplication
	accessTokenRequestJSON entity.AccessTokenRequestJSON
	accessToken            entity.OauthAccessToken
	refreshToken           entity.OauthRefreshToken
	refreshTokenJSON       entity.OauthRefreshTokenJSON
}

func TestClientCredential(t *testing.T) {
	suite.Run(t, &ClientCredentialSuiteTest{
		AuthorizationBaseSuiteTest: &AuthorizationBaseSuiteTest{},
	})
}

func (suite *ClientCredentialSuiteTest) SetupTest() {
	suite.oauthApplication = entity.OauthApplication{
		ID: 1,
		Scopes: sql.NullString{
			String: "public users",
			Valid:  true,
		},
		OwnerType: "confidential",
	}
	suite.accessTokenRequestJSON = entity.AccessTokenRequestJSON{
		GrantType:    util.StringToPointer("client_credentials"),
		ClientUID:    util.StringToPointer("client_uid"),
		ClientSecret: util.StringToPointer("client_secret"),
		Scope:        util.StringToPointer("public"),
	}
	suite.accessToken = entity.OauthAccessToken{
		ID:                 1,
		OauthApplicationID: 1,
		ResourceOwnerID:    0,
		Token:              "some random token",
		Scopes: sql.NullString{
			String: "public",
			Valid:  false,
		},
		ExpiresIn: time.Time{},
		CreatedAt: time.Time{},
		RevokedAT: mysql.NullTime{},
	}
	suite.refreshToken = entity.OauthRefreshToken{
		ID:                 1,
		OauthAccessTokenID: 1,
		Token:              "some token",
		ExpiresIn:          time.Time{},
		CreatedAt:          time.Time{},
		RevokedAT:          mysql.NullTime{},
	}
}

func (suite *ClientCredentialSuiteTest) TestValidateTokenGrantSuiteTest() {
	suite.Run("Positive cases", func() {
		suite.Subtest("When all parameters is valid, then it would return nil", func() {
			suite.refreshTokenJSON = suite.formatter.RefreshToken(suite.refreshToken)
			gomock.InOrder(
				suite.oauthApplicationRepo.EXPECT().OneByUIDandSecret(suite.ktx, *suite.accessTokenRequestJSON.ClientUID, *suite.accessTokenRequestJSON.ClientSecret, suite.sqldb).Return(suite.oauthApplication, nil),
				suite.sqldb.EXPECT().Transaction(suite.ktx, "authorization-grant-client-credential", gomock.Any()).DoAndReturn(func(ctx kontext.Context, transactionKey string, f func(tx db.TX) exception.Exception) exception.Exception {
					suite.oauthAccessTokenRepo.EXPECT().Create(suite.ktx, gomock.Any(), suite.sqldb).DoAndReturn(func(ktx kontext.Context, data entity.OauthAccessTokenInsertable, tx db.TX) (int, exception.Exception) {
						suite.Assert().Equal(*suite.accessTokenRequestJSON.Scope, data.Scopes)
						suite.Assert().Equal(suite.oauthApplication.ID, data.OauthApplicationID)
						return 1, nil
					})
					suite.oauthAccessTokenRepo.EXPECT().One(suite.ktx, 1, suite.sqldb).Return(suite.accessToken, nil)
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
			byteExpectedAccessToken, _ := json.Marshal(suite.formatter.AccessToken(suite.accessToken, "", &suite.refreshTokenJSON))
			suite.Assert().Nil(err)
			suite.Equal(string(byteExpectedAccessToken), string(byteAccessToken))
		})

		suite.Subtest("When all parameters is valid and refresh token config inactive, then it would return nil without access token", func() {
			suite.config.Config.RefreshToken.Active = false
			suite.authorization = usecase.NewAuthorization(suite.oauthApplicationRepo, suite.oauthAccessTokenRepo, suite.oauthAccessGrantRepo, suite.oauthRefreshTokenRepo, suite.formatter, suite.config, suite.sqldb, suite.apiError)

			suite.refreshTokenJSON = suite.formatter.RefreshToken(suite.refreshToken)
			gomock.InOrder(
				suite.oauthApplicationRepo.EXPECT().OneByUIDandSecret(suite.ktx, *suite.accessTokenRequestJSON.ClientUID, *suite.accessTokenRequestJSON.ClientSecret, suite.sqldb).Return(suite.oauthApplication, nil),
				suite.sqldb.EXPECT().Transaction(suite.ktx, "authorization-grant-client-credential", gomock.Any()).DoAndReturn(func(ctx kontext.Context, transactionKey string, f func(tx db.TX) exception.Exception) exception.Exception {
					suite.oauthAccessTokenRepo.EXPECT().Create(suite.ktx, gomock.Any(), suite.sqldb).DoAndReturn(func(ktx kontext.Context, data entity.OauthAccessTokenInsertable, tx db.TX) (int, exception.Exception) {
						suite.Assert().Equal(*suite.accessTokenRequestJSON.Scope, data.Scopes)
						suite.Assert().Equal(suite.oauthApplication.ID, data.OauthApplicationID)
						return 1, nil
					})
					suite.oauthAccessTokenRepo.EXPECT().One(suite.ktx, 1, suite.sqldb).Return(suite.accessToken, nil)
					return f(suite.sqldb)
				}),
			)

			accessTokenJSON, err := suite.authorization.GrantToken(suite.ktx, suite.accessTokenRequestJSON)
			byteAccessToken, _ := json.Marshal(accessTokenJSON)
			byteExpectedAccessToken, _ := json.Marshal(suite.formatter.AccessToken(suite.accessToken, "", nil))
			suite.Assert().Nil(err)
			suite.Equal(string(byteExpectedAccessToken), string(byteAccessToken))
		})
	})

	suite.Run("Negative cases", func() {
		suite.Subtest("When create access token failed, then it would return error", func() {
			suite.refreshTokenJSON = suite.formatter.RefreshToken(suite.refreshToken)
			gomock.InOrder(
				suite.oauthApplicationRepo.EXPECT().OneByUIDandSecret(suite.ktx, *suite.accessTokenRequestJSON.ClientUID, *suite.accessTokenRequestJSON.ClientSecret, suite.sqldb).Return(suite.oauthApplication, nil),
				suite.sqldb.EXPECT().Transaction(suite.ktx, "authorization-grant-client-credential", gomock.Any()).DoAndReturn(func(ctx kontext.Context, transactionKey string, f func(tx db.TX) exception.Exception) exception.Exception {
					suite.oauthAccessTokenRepo.EXPECT().Create(suite.ktx, gomock.Any(), suite.sqldb).DoAndReturn(func(ktx kontext.Context, data entity.OauthAccessTokenInsertable, tx db.TX) (int, exception.Exception) {
						suite.Assert().Equal(*suite.accessTokenRequestJSON.Scope, data.Scopes)
						suite.Assert().Equal(suite.oauthApplication.ID, data.OauthApplicationID)
						return 0, exception.Throw(errors.New("unexpected"))
					})
					return f(suite.sqldb)
				}),
			)

			_, err := suite.authorization.GrantToken(suite.ktx, suite.accessTokenRequestJSON)
			suite.Assert().NotNil(err)
			suite.Assert().Equal("JSONAPI Error:\n[Internal server error] Detail: Something is not right, help us fix this problem. Contribute to https://github.com/kodefluence/altair. Tracing code: '<nil>', Code: ERR0500\n", err.Error())
			suite.Assert().Equal(http.StatusInternalServerError, err.HTTPStatus())
		})

		suite.Subtest("When get newly created access token failed, then it would return error", func() {
			suite.refreshTokenJSON = suite.formatter.RefreshToken(suite.refreshToken)
			gomock.InOrder(
				suite.oauthApplicationRepo.EXPECT().OneByUIDandSecret(suite.ktx, *suite.accessTokenRequestJSON.ClientUID, *suite.accessTokenRequestJSON.ClientSecret, suite.sqldb).Return(suite.oauthApplication, nil),
				suite.sqldb.EXPECT().Transaction(suite.ktx, "authorization-grant-client-credential", gomock.Any()).DoAndReturn(func(ctx kontext.Context, transactionKey string, f func(tx db.TX) exception.Exception) exception.Exception {
					suite.oauthAccessTokenRepo.EXPECT().Create(suite.ktx, gomock.Any(), suite.sqldb).DoAndReturn(func(ktx kontext.Context, data entity.OauthAccessTokenInsertable, tx db.TX) (int, exception.Exception) {
						suite.Assert().Equal(*suite.accessTokenRequestJSON.Scope, data.Scopes)
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

		suite.Subtest("When refresh token creation failure, then it would return error", func() {
			suite.refreshTokenJSON = suite.formatter.RefreshToken(suite.refreshToken)
			gomock.InOrder(
				suite.oauthApplicationRepo.EXPECT().OneByUIDandSecret(suite.ktx, *suite.accessTokenRequestJSON.ClientUID, *suite.accessTokenRequestJSON.ClientSecret, suite.sqldb).Return(suite.oauthApplication, nil),
				suite.sqldb.EXPECT().Transaction(suite.ktx, "authorization-grant-client-credential", gomock.Any()).DoAndReturn(func(ctx kontext.Context, transactionKey string, f func(tx db.TX) exception.Exception) exception.Exception {
					suite.oauthAccessTokenRepo.EXPECT().Create(suite.ktx, gomock.Any(), suite.sqldb).DoAndReturn(func(ktx kontext.Context, data entity.OauthAccessTokenInsertable, tx db.TX) (int, exception.Exception) {
						suite.Assert().Equal(*suite.accessTokenRequestJSON.Scope, data.Scopes)
						suite.Assert().Equal(suite.oauthApplication.ID, data.OauthApplicationID)
						return 1, nil
					})
					suite.oauthAccessTokenRepo.EXPECT().One(suite.ktx, 1, suite.sqldb).Return(suite.accessToken, nil)
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

func (suite *ClientCredentialSuiteTest) Subtest(testcase string, subtest func()) {
	suite.SetupTest()
	suite.AuthorizationBaseSuiteTest.Subtest(testcase, subtest)
	suite.TearDownTest()
}
