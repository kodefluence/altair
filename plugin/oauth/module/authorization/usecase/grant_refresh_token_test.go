package usecase_test

import (
	"database/sql"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/golang/mock/gomock"
	"github.com/kodefluence/altair/plugin/oauth/entity"
	"github.com/kodefluence/monorepo/db"
	"github.com/kodefluence/monorepo/exception"
	"github.com/kodefluence/monorepo/kontext"
	"github.com/stretchr/testify/suite"
)

type GrantRefreshTokenSuiteTest struct {
	*AuthorizationBaseSuiteTest

	oauthApplication entity.OauthApplication
	accessToken      entity.OauthAccessToken
	refreshToken     entity.OauthRefreshToken
}

func TestGrantRefreshToken(t *testing.T) {
	suite.Run(t, &GrantRefreshTokenSuiteTest{
		AuthorizationBaseSuiteTest: &AuthorizationBaseSuiteTest{},
	})
}

func (suite *GrantRefreshTokenSuiteTest) SetupTest() {
	suite.oauthApplication = entity.OauthApplication{
		ID: 1,
		Scopes: sql.NullString{
			String: "public users",
			Valid:  true,
		},
		OwnerType: "confidential",
	}
	suite.accessToken = entity.OauthAccessToken{
		ID:                 1,
		OauthApplicationID: 1,
		ResourceOwnerID:    0,
		Token:              "some random token",
		Scopes:             sql.NullString{},
		ExpiresIn:          time.Time{},
		CreatedAt:          time.Time{},
		RevokedAT:          mysql.NullTime{},
	}
	suite.refreshToken = entity.OauthRefreshToken{
		ID:                 1,
		OauthAccessTokenID: 1,
		Token:              "",
		ExpiresIn:          time.Time{},
		CreatedAt:          time.Time{},
		RevokedAT:          mysql.NullTime{},
	}
}

func (suite *GrantRefreshTokenSuiteTest) TestGrantRefreshTokenSuiteTest() {
	suite.Run("Positive cases", func() {
		suite.Subtest("When all parameters are valid, then it would return nil", func() {
			gomock.InOrder(
				suite.oauthRefreshTokenRepo.EXPECT().Create(suite.ktx, gomock.Any(), suite.sqldb).DoAndReturn(func(ktx kontext.Context, data entity.OauthRefreshTokenInsertable, tx db.TX) (int, exception.Exception) {
					suite.Assert().Equal(suite.accessToken.ID, data.OauthAccessTokenID)
					return 1, nil
				}),
				suite.oauthRefreshTokenRepo.EXPECT().One(suite.ktx, suite.refreshToken.ID, suite.sqldb).Return(suite.refreshToken, nil),
			)

			refreshToken, err := suite.authorization.GrantRefreshToken(suite.ktx, suite.accessToken, suite.oauthApplication, suite.sqldb)
			suite.Assert().Nil(err)
			suite.Assert().Equal(suite.refreshToken, refreshToken)
		})
	})

	suite.Run("Negative cases", func() {
		suite.Subtest("When create refresh token error, then it would return error", func() {
			gomock.InOrder(
				suite.oauthRefreshTokenRepo.EXPECT().Create(suite.ktx, gomock.Any(), suite.sqldb).DoAndReturn(func(ktx kontext.Context, data entity.OauthRefreshTokenInsertable, tx db.TX) (int, exception.Exception) {
					suite.Assert().Equal(suite.accessToken.ID, data.OauthAccessTokenID)
					return 1, exception.Throw(errors.New("unexpected"))
				}),
			)

			_, err := suite.authorization.GrantRefreshToken(suite.ktx, suite.accessToken, suite.oauthApplication, suite.sqldb)
			suite.Assert().NotNil(err)
			suite.Assert().Equal("JSONAPI Error:\n[Internal server error] Detail: Something is not right, help us fix this problem. Contribute to https://github.com/kodefluence/altair. Tracing code: '<nil>', Code: ERR0500\n", err.Error())
			suite.Assert().Equal(http.StatusInternalServerError, err.HTTPStatus())
		})

		suite.Subtest("When find token error, then it would return error", func() {
			gomock.InOrder(
				suite.oauthRefreshTokenRepo.EXPECT().Create(suite.ktx, gomock.Any(), suite.sqldb).DoAndReturn(func(ktx kontext.Context, data entity.OauthRefreshTokenInsertable, tx db.TX) (int, exception.Exception) {
					suite.Assert().Equal(suite.accessToken.ID, data.OauthAccessTokenID)
					return 1, nil
				}),
				suite.oauthRefreshTokenRepo.EXPECT().One(suite.ktx, suite.refreshToken.ID, suite.sqldb).Return(entity.OauthRefreshToken{}, exception.Throw(errors.New("unexpected"))),
			)

			_, err := suite.authorization.GrantRefreshToken(suite.ktx, suite.accessToken, suite.oauthApplication, suite.sqldb)
			suite.Assert().NotNil(err)
			suite.Assert().Equal("JSONAPI Error:\n[Internal server error] Detail: Something is not right, help us fix this problem. Contribute to https://github.com/kodefluence/altair. Tracing code: '<nil>', Code: ERR0500\n", err.Error())
			suite.Assert().Equal(http.StatusInternalServerError, err.HTTPStatus())
		})
	})
}

func (suite *GrantRefreshTokenSuiteTest) Subtest(testcase string, subtest func()) {
	suite.SetupTest()
	suite.AuthorizationBaseSuiteTest.Subtest(testcase, subtest)
	suite.TearDownTest()
}
