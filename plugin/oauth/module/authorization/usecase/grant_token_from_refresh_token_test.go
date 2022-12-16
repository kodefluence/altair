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
	"github.com/kodefluence/altair/util"
	"github.com/kodefluence/monorepo/db"
	"github.com/kodefluence/monorepo/exception"
	"github.com/kodefluence/monorepo/kontext"
	"github.com/stretchr/testify/suite"
)

type GrantTokenFromRefreshTokenSuiteTest struct {
	*AuthorizationBaseSuiteTest

	oauthApplication       entity.OauthApplication
	accessTokenRequestJSON entity.AccessTokenRequestJSON
	oldaccessToken         entity.OauthAccessToken
	accessToken            entity.OauthAccessToken
	oldrefreshToken        entity.OauthRefreshToken
	refreshToken           entity.OauthRefreshToken
	refreshTokenJSON       entity.OauthRefreshTokenJSON
}

func TestGrantTokenFromRefreshToken(t *testing.T) {
	suite.Run(t, &GrantTokenFromRefreshTokenSuiteTest{
		AuthorizationBaseSuiteTest: &AuthorizationBaseSuiteTest{},
	})
}

func (suite *GrantTokenFromRefreshTokenSuiteTest) SetupTest() {
	suite.oauthApplication = entity.OauthApplication{
		ID: 1,
		Scopes: sql.NullString{
			String: "public users",
			Valid:  true,
		},
		OwnerType: "confidential",
	}
	suite.accessTokenRequestJSON = entity.AccessTokenRequestJSON{
		GrantType:    util.StringToPointer("refresh_token"),
		RefreshToken: util.StringToPointer("some refresh token"),
		Scope:        util.StringToPointer("public"),
	}
	suite.oldaccessToken = entity.OauthAccessToken{
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
		RevokedAT: sql.NullTime{},
	}
	suite.accessToken = entity.OauthAccessToken{
		ID:                 2,
		OauthApplicationID: 1,
		ResourceOwnerID:    0,
		Token:              "some random token",
		Scopes: sql.NullString{
			String: "public",
			Valid:  false,
		},
		ExpiresIn: time.Time{},
		CreatedAt: time.Time{},
		RevokedAT: sql.NullTime{},
	}

	suite.oldrefreshToken = entity.OauthRefreshToken{
		ID:                 1,
		OauthAccessTokenID: 1,
		Token:              "some token",
		ExpiresIn:          time.Now().Add(time.Hour),
		CreatedAt:          time.Time{},
		RevokedAT:          sql.NullTime{},
	}

	suite.refreshToken = entity.OauthRefreshToken{
		ID:                 2,
		OauthAccessTokenID: 1,
		Token:              "some token",
		ExpiresIn:          time.Time{},
		CreatedAt:          time.Time{},
		RevokedAT:          sql.NullTime{},
	}
}

func (suite *GrantTokenFromRefreshTokenSuiteTest) TestValidateTokenGrantSuiteTest() {
	suite.Run("Positive cases", func() {
		suite.Subtest("When all parameters is valid, then it would return nil", func() {
			suite.refreshTokenJSON = suite.formatter.RefreshToken(suite.refreshToken)
			gomock.InOrder(
				suite.sqldb.EXPECT().Transaction(suite.ktx, "authorization-grant-token-from-refresh-token", gomock.Any()).DoAndReturn(func(ctx kontext.Context, transactionKey string, f func(tx db.TX) exception.Exception) exception.Exception {
					suite.oauthRefreshTokenRepo.EXPECT().OneByToken(suite.ktx, *suite.accessTokenRequestJSON.RefreshToken, suite.sqldb).Return(suite.oldrefreshToken, nil)
					suite.oauthAccessTokenRepo.EXPECT().One(suite.ktx, 1, suite.sqldb).Return(suite.oldaccessToken, nil)
					suite.oauthApplicationRepo.EXPECT().One(suite.ktx, suite.oldaccessToken.OauthApplicationID, suite.sqldb).Return(suite.oauthApplication, nil)
					suite.oauthAccessTokenRepo.EXPECT().Create(suite.ktx, gomock.Any(), suite.sqldb).DoAndReturn(func(ktx kontext.Context, data entity.OauthAccessTokenInsertable, tx db.TX) (int, exception.Exception) {
						suite.Assert().Equal(suite.oldaccessToken.Scopes.String, data.Scopes)
						suite.Assert().Equal(suite.oauthApplication.ID, data.OauthApplicationID)
						return 2, nil
					})
					suite.oauthAccessTokenRepo.EXPECT().One(suite.ktx, 2, suite.sqldb).Return(suite.accessToken, nil)
					suite.oauthRefreshTokenRepo.EXPECT().Revoke(suite.ktx, *suite.accessTokenRequestJSON.RefreshToken, suite.sqldb).Return(nil)
					suite.oauthRefreshTokenRepo.EXPECT().Create(suite.ktx, gomock.Any(), suite.sqldb).DoAndReturn(func(ktx kontext.Context, data entity.OauthRefreshTokenInsertable, tx db.TX) (int, exception.Exception) {
						suite.Assert().Equal(suite.accessToken.ID, data.OauthAccessTokenID)
						return 2, nil
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
	})

	suite.Run("Negative cases", func() {
		suite.Subtest("When select old access token failed not found, then it would return error", func() {
			suite.refreshTokenJSON = suite.formatter.RefreshToken(suite.refreshToken)
			gomock.InOrder(
				suite.sqldb.EXPECT().Transaction(suite.ktx, "authorization-grant-token-from-refresh-token", gomock.Any()).DoAndReturn(func(ctx kontext.Context, transactionKey string, f func(tx db.TX) exception.Exception) exception.Exception {
					suite.oauthRefreshTokenRepo.EXPECT().OneByToken(suite.ktx, *suite.accessTokenRequestJSON.RefreshToken, suite.sqldb).Return(entity.OauthRefreshToken{}, exception.Throw(errors.New("not found"), exception.WithType(exception.NotFound)))
					return f(suite.sqldb)
				}),
			)

			_, err := suite.authorization.GrantToken(suite.ktx, suite.accessTokenRequestJSON)
			suite.Assert().NotNil(err)
			suite.Assert().Equal("JSONAPI Error:\n[Not found error] Detail: Resource of `refresh_token` is not found. Tracing code: `<nil>`, Code: ERR0404\n", err.Error())
			suite.Assert().Equal(http.StatusNotFound, err.HTTPStatus())
		})

		suite.Subtest("When select old access token failed with unexpected error, then it would return error", func() {
			suite.refreshTokenJSON = suite.formatter.RefreshToken(suite.refreshToken)
			gomock.InOrder(
				suite.sqldb.EXPECT().Transaction(suite.ktx, "authorization-grant-token-from-refresh-token", gomock.Any()).DoAndReturn(func(ctx kontext.Context, transactionKey string, f func(tx db.TX) exception.Exception) exception.Exception {
					suite.oauthRefreshTokenRepo.EXPECT().OneByToken(suite.ktx, *suite.accessTokenRequestJSON.RefreshToken, suite.sqldb).Return(entity.OauthRefreshToken{}, exception.Throw(errors.New("unexpected")))
					return f(suite.sqldb)
				}),
			)

			_, err := suite.authorization.GrantToken(suite.ktx, suite.accessTokenRequestJSON)
			suite.Assert().NotNil(err)
			suite.Assert().Equal("JSONAPI Error:\n[Internal server error] Detail: Something is not right, help us fix this problem. Contribute to https://github.com/kodefluence/altair. Tracing code: '<nil>', Code: ERR0500\n", err.Error())
			suite.Assert().Equal(http.StatusInternalServerError, err.HTTPStatus())
		})

		suite.Subtest("When refresh token already expired, then it would return error", func() {
			suite.refreshTokenJSON = suite.formatter.RefreshToken(suite.refreshToken)
			suite.oldrefreshToken.ExpiresIn = time.Now().Add(-1 * time.Hour)

			gomock.InOrder(
				suite.sqldb.EXPECT().Transaction(suite.ktx, "authorization-grant-token-from-refresh-token", gomock.Any()).DoAndReturn(func(ctx kontext.Context, transactionKey string, f func(tx db.TX) exception.Exception) exception.Exception {
					suite.oauthRefreshTokenRepo.EXPECT().OneByToken(suite.ktx, *suite.accessTokenRequestJSON.RefreshToken, suite.sqldb).Return(suite.oldrefreshToken, nil)
					return f(suite.sqldb)
				}),
			)

			_, err := suite.authorization.GrantToken(suite.ktx, suite.accessTokenRequestJSON)
			suite.Assert().NotNil(err)
			suite.Assert().Equal("JSONAPI Error:\n[Forbidden error] Detail: Resource of `access_token` is forbidden to be accessed, because of: refresh token already used. Tracing code: `<nil>`, Code: ERR0403\n", err.Error())
			suite.Assert().Equal(http.StatusForbidden, err.HTTPStatus())
		})

		suite.Subtest("When refresh token already revoked, then it would return error", func() {
			suite.refreshTokenJSON = suite.formatter.RefreshToken(suite.refreshToken)
			suite.oldrefreshToken.RevokedAT.Time = time.Now().Add(-1 * time.Hour)
			suite.oldrefreshToken.RevokedAT.Valid = true

			gomock.InOrder(
				suite.sqldb.EXPECT().Transaction(suite.ktx, "authorization-grant-token-from-refresh-token", gomock.Any()).DoAndReturn(func(ctx kontext.Context, transactionKey string, f func(tx db.TX) exception.Exception) exception.Exception {
					suite.oauthRefreshTokenRepo.EXPECT().OneByToken(suite.ktx, *suite.accessTokenRequestJSON.RefreshToken, suite.sqldb).Return(suite.oldrefreshToken, nil)
					return f(suite.sqldb)
				}),
			)

			_, err := suite.authorization.GrantToken(suite.ktx, suite.accessTokenRequestJSON)
			suite.Assert().NotNil(err)
			suite.Assert().Equal("JSONAPI Error:\n[Forbidden error] Detail: Resource of `access_token` is forbidden to be accessed, because of: refresh token already used. Tracing code: `<nil>`, Code: ERR0403\n", err.Error())
			suite.Assert().Equal(http.StatusForbidden, err.HTTPStatus())
		})

		suite.Subtest("When select old access token failed, then it would return error", func() {
			suite.refreshTokenJSON = suite.formatter.RefreshToken(suite.refreshToken)
			gomock.InOrder(
				suite.sqldb.EXPECT().Transaction(suite.ktx, "authorization-grant-token-from-refresh-token", gomock.Any()).DoAndReturn(func(ctx kontext.Context, transactionKey string, f func(tx db.TX) exception.Exception) exception.Exception {
					suite.oauthRefreshTokenRepo.EXPECT().OneByToken(suite.ktx, *suite.accessTokenRequestJSON.RefreshToken, suite.sqldb).Return(suite.oldrefreshToken, nil)
					suite.oauthAccessTokenRepo.EXPECT().One(suite.ktx, 1, suite.sqldb).Return(entity.OauthAccessToken{}, exception.Throw(errors.New("unexpected")))
					return f(suite.sqldb)
				}),
			)

			_, err := suite.authorization.GrantToken(suite.ktx, suite.accessTokenRequestJSON)
			suite.Assert().NotNil(err)
			suite.Assert().Equal("JSONAPI Error:\n[Internal server error] Detail: Something is not right, help us fix this problem. Contribute to https://github.com/kodefluence/altair. Tracing code: '<nil>', Code: ERR0500\n", err.Error())
			suite.Assert().Equal(http.StatusInternalServerError, err.HTTPStatus())
		})

		suite.Subtest("When select oauth application failed, then it would return error", func() {
			suite.refreshTokenJSON = suite.formatter.RefreshToken(suite.refreshToken)
			gomock.InOrder(
				suite.sqldb.EXPECT().Transaction(suite.ktx, "authorization-grant-token-from-refresh-token", gomock.Any()).DoAndReturn(func(ctx kontext.Context, transactionKey string, f func(tx db.TX) exception.Exception) exception.Exception {
					suite.oauthRefreshTokenRepo.EXPECT().OneByToken(suite.ktx, *suite.accessTokenRequestJSON.RefreshToken, suite.sqldb).Return(suite.oldrefreshToken, nil)
					suite.oauthAccessTokenRepo.EXPECT().One(suite.ktx, 1, suite.sqldb).Return(suite.oldaccessToken, nil)
					suite.oauthApplicationRepo.EXPECT().One(suite.ktx, suite.oldaccessToken.OauthApplicationID, suite.sqldb).Return(entity.OauthApplication{}, exception.Throw(errors.New("unexpected")))
					return f(suite.sqldb)
				}),
			)

			_, err := suite.authorization.GrantToken(suite.ktx, suite.accessTokenRequestJSON)
			suite.Assert().NotNil(err)
			suite.Assert().Equal("JSONAPI Error:\n[Internal server error] Detail: Something is not right, help us fix this problem. Contribute to https://github.com/kodefluence/altair. Tracing code: '<nil>', Code: ERR0500\n", err.Error())
			suite.Assert().Equal(http.StatusInternalServerError, err.HTTPStatus())
		})

		suite.Subtest("When create new access token is failed, then it would return error", func() {
			suite.refreshTokenJSON = suite.formatter.RefreshToken(suite.refreshToken)
			gomock.InOrder(
				suite.sqldb.EXPECT().Transaction(suite.ktx, "authorization-grant-token-from-refresh-token", gomock.Any()).DoAndReturn(func(ctx kontext.Context, transactionKey string, f func(tx db.TX) exception.Exception) exception.Exception {
					suite.oauthRefreshTokenRepo.EXPECT().OneByToken(suite.ktx, *suite.accessTokenRequestJSON.RefreshToken, suite.sqldb).Return(suite.oldrefreshToken, nil)
					suite.oauthAccessTokenRepo.EXPECT().One(suite.ktx, 1, suite.sqldb).Return(suite.oldaccessToken, nil)
					suite.oauthApplicationRepo.EXPECT().One(suite.ktx, suite.oldaccessToken.OauthApplicationID, suite.sqldb).Return(suite.oauthApplication, nil)
					suite.oauthAccessTokenRepo.EXPECT().Create(suite.ktx, gomock.Any(), suite.sqldb).DoAndReturn(func(ktx kontext.Context, data entity.OauthAccessTokenInsertable, tx db.TX) (int, exception.Exception) {
						suite.Assert().Equal(suite.oldaccessToken.Scopes.String, data.Scopes)
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

		suite.Subtest("When get newly created access token is failed, then it would return error", func() {
			suite.refreshTokenJSON = suite.formatter.RefreshToken(suite.refreshToken)
			gomock.InOrder(
				suite.sqldb.EXPECT().Transaction(suite.ktx, "authorization-grant-token-from-refresh-token", gomock.Any()).DoAndReturn(func(ctx kontext.Context, transactionKey string, f func(tx db.TX) exception.Exception) exception.Exception {
					suite.oauthRefreshTokenRepo.EXPECT().OneByToken(suite.ktx, *suite.accessTokenRequestJSON.RefreshToken, suite.sqldb).Return(suite.oldrefreshToken, nil)
					suite.oauthAccessTokenRepo.EXPECT().One(suite.ktx, 1, suite.sqldb).Return(suite.oldaccessToken, nil)
					suite.oauthApplicationRepo.EXPECT().One(suite.ktx, suite.oldaccessToken.OauthApplicationID, suite.sqldb).Return(suite.oauthApplication, nil)
					suite.oauthAccessTokenRepo.EXPECT().Create(suite.ktx, gomock.Any(), suite.sqldb).DoAndReturn(func(ktx kontext.Context, data entity.OauthAccessTokenInsertable, tx db.TX) (int, exception.Exception) {
						suite.Assert().Equal(suite.oldaccessToken.Scopes.String, data.Scopes)
						suite.Assert().Equal(suite.oauthApplication.ID, data.OauthApplicationID)
						return 2, nil
					})
					suite.oauthAccessTokenRepo.EXPECT().One(suite.ktx, 2, suite.sqldb).Return(entity.OauthAccessToken{}, exception.Throw(errors.New("unexpected")))
					return f(suite.sqldb)
				}),
			)

			_, err := suite.authorization.GrantToken(suite.ktx, suite.accessTokenRequestJSON)
			suite.Assert().NotNil(err)
			suite.Assert().Equal("JSONAPI Error:\n[Internal server error] Detail: Something is not right, help us fix this problem. Contribute to https://github.com/kodefluence/altair. Tracing code: '<nil>', Code: ERR0500\n", err.Error())
			suite.Assert().Equal(http.StatusInternalServerError, err.HTTPStatus())
		})

		suite.Subtest("When revoke refresh token failed, then it would return error", func() {
			suite.refreshTokenJSON = suite.formatter.RefreshToken(suite.refreshToken)
			gomock.InOrder(
				suite.sqldb.EXPECT().Transaction(suite.ktx, "authorization-grant-token-from-refresh-token", gomock.Any()).DoAndReturn(func(ctx kontext.Context, transactionKey string, f func(tx db.TX) exception.Exception) exception.Exception {
					suite.oauthRefreshTokenRepo.EXPECT().OneByToken(suite.ktx, *suite.accessTokenRequestJSON.RefreshToken, suite.sqldb).Return(suite.oldrefreshToken, nil)
					suite.oauthAccessTokenRepo.EXPECT().One(suite.ktx, 1, suite.sqldb).Return(suite.oldaccessToken, nil)
					suite.oauthApplicationRepo.EXPECT().One(suite.ktx, suite.oldaccessToken.OauthApplicationID, suite.sqldb).Return(suite.oauthApplication, nil)
					suite.oauthAccessTokenRepo.EXPECT().Create(suite.ktx, gomock.Any(), suite.sqldb).DoAndReturn(func(ktx kontext.Context, data entity.OauthAccessTokenInsertable, tx db.TX) (int, exception.Exception) {
						suite.Assert().Equal(suite.oldaccessToken.Scopes.String, data.Scopes)
						suite.Assert().Equal(suite.oauthApplication.ID, data.OauthApplicationID)
						return 2, nil
					})
					suite.oauthAccessTokenRepo.EXPECT().One(suite.ktx, 2, suite.sqldb).Return(suite.accessToken, nil)
					suite.oauthRefreshTokenRepo.EXPECT().Revoke(suite.ktx, *suite.accessTokenRequestJSON.RefreshToken, suite.sqldb).Return(exception.Throw(errors.New("unexpected")))
					return f(suite.sqldb)
				}),
			)

			_, err := suite.authorization.GrantToken(suite.ktx, suite.accessTokenRequestJSON)
			suite.Assert().NotNil(err)
			suite.Assert().Equal("JSONAPI Error:\n[Internal server error] Detail: Something is not right, help us fix this problem. Contribute to https://github.com/kodefluence/altair. Tracing code: '<nil>', Code: ERR0500\n", err.Error())
			suite.Assert().Equal(http.StatusInternalServerError, err.HTTPStatus())
		})

		suite.Subtest("When create new refresh token is failed, then it would return error", func() {
			suite.refreshTokenJSON = suite.formatter.RefreshToken(suite.refreshToken)
			gomock.InOrder(
				suite.sqldb.EXPECT().Transaction(suite.ktx, "authorization-grant-token-from-refresh-token", gomock.Any()).DoAndReturn(func(ctx kontext.Context, transactionKey string, f func(tx db.TX) exception.Exception) exception.Exception {
					suite.oauthRefreshTokenRepo.EXPECT().OneByToken(suite.ktx, *suite.accessTokenRequestJSON.RefreshToken, suite.sqldb).Return(suite.oldrefreshToken, nil)
					suite.oauthAccessTokenRepo.EXPECT().One(suite.ktx, 1, suite.sqldb).Return(suite.oldaccessToken, nil)
					suite.oauthApplicationRepo.EXPECT().One(suite.ktx, suite.oldaccessToken.OauthApplicationID, suite.sqldb).Return(suite.oauthApplication, nil)
					suite.oauthAccessTokenRepo.EXPECT().Create(suite.ktx, gomock.Any(), suite.sqldb).DoAndReturn(func(ktx kontext.Context, data entity.OauthAccessTokenInsertable, tx db.TX) (int, exception.Exception) {
						suite.Assert().Equal(suite.oldaccessToken.Scopes.String, data.Scopes)
						suite.Assert().Equal(suite.oauthApplication.ID, data.OauthApplicationID)
						return 2, nil
					})
					suite.oauthAccessTokenRepo.EXPECT().One(suite.ktx, 2, suite.sqldb).Return(suite.accessToken, nil)
					suite.oauthRefreshTokenRepo.EXPECT().Revoke(suite.ktx, *suite.accessTokenRequestJSON.RefreshToken, suite.sqldb).Return(nil)
					suite.oauthRefreshTokenRepo.EXPECT().Create(suite.ktx, gomock.Any(), suite.sqldb).DoAndReturn(func(ktx kontext.Context, data entity.OauthRefreshTokenInsertable, tx db.TX) (int, exception.Exception) {
						suite.Assert().Equal(suite.accessToken.ID, data.OauthAccessTokenID)
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
				suite.sqldb.EXPECT().Transaction(suite.ktx, "authorization-grant-token-from-refresh-token", gomock.Any()).DoAndReturn(func(ctx kontext.Context, transactionKey string, f func(tx db.TX) exception.Exception) exception.Exception {
					suite.oauthRefreshTokenRepo.EXPECT().OneByToken(suite.ktx, *suite.accessTokenRequestJSON.RefreshToken, suite.sqldb).Return(suite.oldrefreshToken, nil)
					suite.oauthAccessTokenRepo.EXPECT().One(suite.ktx, 1, suite.sqldb).Return(suite.oldaccessToken, nil)
					suite.oauthApplicationRepo.EXPECT().One(suite.ktx, suite.oldaccessToken.OauthApplicationID, suite.sqldb).Return(suite.oauthApplication, nil)
					suite.oauthAccessTokenRepo.EXPECT().Create(suite.ktx, gomock.Any(), suite.sqldb).DoAndReturn(func(ktx kontext.Context, data entity.OauthAccessTokenInsertable, tx db.TX) (int, exception.Exception) {
						suite.Assert().Equal(suite.oldaccessToken.Scopes.String, data.Scopes)
						suite.Assert().Equal(suite.oauthApplication.ID, data.OauthApplicationID)
						return 2, nil
					})
					suite.oauthAccessTokenRepo.EXPECT().One(suite.ktx, 2, suite.sqldb).Return(suite.accessToken, nil)
					suite.oauthRefreshTokenRepo.EXPECT().Revoke(suite.ktx, *suite.accessTokenRequestJSON.RefreshToken, suite.sqldb).Return(nil)
					suite.oauthRefreshTokenRepo.EXPECT().Create(suite.ktx, gomock.Any(), suite.sqldb).DoAndReturn(func(ktx kontext.Context, data entity.OauthRefreshTokenInsertable, tx db.TX) (int, exception.Exception) {
						suite.Assert().Equal(suite.accessToken.ID, data.OauthAccessTokenID)
						return 2, nil
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

func (suite *GrantTokenFromRefreshTokenSuiteTest) Subtest(testcase string, subtest func()) {
	suite.SetupTest()
	suite.AuthorizationBaseSuiteTest.Subtest(testcase, subtest)
	suite.TearDownTest()
}
