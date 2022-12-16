package usecase_test

import (
	"database/sql"
	"errors"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/kodefluence/altair/plugin/oauth/entity"
	"github.com/kodefluence/altair/plugin/oauth/module/authorization/usecase"
	"github.com/kodefluence/altair/util"
	"github.com/kodefluence/monorepo/exception"
	"github.com/stretchr/testify/suite"
)

type GrantTokenSuiteTest struct {
	*AuthorizationBaseSuiteTest

	oauthApplication       entity.OauthApplication
	accessTokenRequestJSON entity.AccessTokenRequestJSON
}

func TestGrantToken(t *testing.T) {
	suite.Run(t, &GrantTokenSuiteTest{
		AuthorizationBaseSuiteTest: &AuthorizationBaseSuiteTest{},
	})
}

func (suite *GrantTokenSuiteTest) SetupTest() {
	suite.oauthApplication = entity.OauthApplication{
		ID: 1,
		Scopes: sql.NullString{
			String: "public users",
			Valid:  true,
		},
		OwnerType: "confidential",
	}
	suite.accessTokenRequestJSON = entity.AccessTokenRequestJSON{
		GrantType:    util.ValueToPointer("client_credentials"),
		ClientUID:    util.ValueToPointer("client_uid"),
		ClientSecret: util.ValueToPointer("client_secret"),
		Scope:        util.ValueToPointer("public"),
	}
}

func (suite *GrantTokenSuiteTest) Subtest(testcase string, subtest func()) {
	suite.SetupTest()
	suite.AuthorizationBaseSuiteTest.Subtest(testcase, subtest)
	suite.TearDownTest()
}

func (suite *GrantTokenSuiteTest) TestGrantTokenSuiteTest() {
	suite.Run("Negative cases", func() {
		suite.Subtest("When find application return an error, then it would return error", func() {
			gomock.InOrder(
				suite.oauthApplicationRepo.EXPECT().OneByUIDandSecret(suite.ktx, *suite.accessTokenRequestJSON.ClientUID, *suite.accessTokenRequestJSON.ClientSecret, suite.sqldb).Return(entity.OauthApplication{}, exception.Throw(errors.New("unexpected"))),
			)

			_, err := suite.authorization.GrantToken(suite.ktx, suite.accessTokenRequestJSON)
			suite.Assert().NotNil(err)
			suite.Assert().Equal("JSONAPI Error:\n[Internal server error] Detail: Something is not right, help us fix this problem. Contribute to https://github.com/kodefluence/altair. Tracing code: '<nil>', Code: ERR0500\n", err.Error())
			suite.Assert().Equal(http.StatusInternalServerError, err.HTTPStatus())
		})

		suite.Subtest("When grant type is empty, then it would return error", func() {
			suite.accessTokenRequestJSON.GrantType = util.ValueToPointer("")

			_, err := suite.authorization.GrantToken(suite.ktx, suite.accessTokenRequestJSON)
			suite.Assert().NotNil(err)
			suite.Assert().Equal("JSONAPI Error:\n[Validation error] Detail: Validation error because of: grant_type is not valid value, Code: ERR1442\n", err.Error())
			suite.Assert().Equal(http.StatusUnprocessableEntity, err.HTTPStatus())
		})

		suite.Subtest("When grant type is refresh token but refresh token config is inactive, then it would return error", func() {
			suite.accessTokenRequestJSON.GrantType = util.ValueToPointer("refresh_token")
			suite.config.Config.RefreshToken.Active = false
			suite.authorization = usecase.NewAuthorization(suite.oauthApplicationRepo, suite.oauthAccessTokenRepo, suite.oauthAccessGrantRepo, suite.oauthRefreshTokenRepo, suite.formatter, suite.config, suite.sqldb, suite.apiError)

			_, err := suite.authorization.GrantToken(suite.ktx, suite.accessTokenRequestJSON)
			suite.Assert().NotNil(err)
			suite.Assert().Equal("JSONAPI Error:\n[Validation error] Detail: Validation error because of: refresh_token is not valid value, Code: ERR1442\n", err.Error())
			suite.Assert().Equal(http.StatusUnprocessableEntity, err.HTTPStatus())
		})
	})
}
