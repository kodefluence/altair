package usecase_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/kodefluence/altair/plugin/oauth/entity"
	"github.com/kodefluence/altair/util"
	"github.com/kodefluence/monorepo/exception"
	"github.com/stretchr/testify/suite"
)

type RevokeTokenSuiteTest struct {
	*AuthorizationBaseSuiteTest

	revokeRequest entity.RevokeAccessTokenRequestJSON
}

func TestRevokeToken(t *testing.T) {
	suite.Run(t, &RevokeTokenSuiteTest{
		AuthorizationBaseSuiteTest: &AuthorizationBaseSuiteTest{},
	})
}

func (suite *RevokeTokenSuiteTest) SetupTest() {
	suite.revokeRequest = entity.RevokeAccessTokenRequestJSON{
		Token: util.ValueToPointer("some-token"),
	}
}

func (suite *RevokeTokenSuiteTest) Subtest(testcase string, subtest func()) {
	suite.SetupTest()
	suite.AuthorizationBaseSuiteTest.Subtest(testcase, subtest)
	suite.TearDownTest()
}

func (suite *RevokeTokenSuiteTest) TestRevokeToken() {
	suite.Run("Positive cases", func() {
		suite.Subtest("When all parameters is valid, then it would return nil", func() {
			suite.oauthAccessTokenRepo.EXPECT().Revoke(suite.ktx, *suite.revokeRequest.Token, suite.sqldb).Return(nil)
			err := suite.authorization.RevokeToken(suite.ktx, suite.revokeRequest)
			suite.Assert().Nil(err)
		})
	})

	suite.Run("Negative cases", func() {
		suite.Subtest("When token parameter is nil, then it would return error", func() {
			suite.revokeRequest.Token = nil
			err := suite.authorization.RevokeToken(suite.ktx, suite.revokeRequest)
			suite.Assert().Equal("JSONAPI Error:\n[Validation error] Detail: Validation error because of: token cannot be empty, Code: ERR1442\n", err.Error())
			suite.Assert().Equal(http.StatusUnprocessableEntity, err.HTTPStatus())
		})

		suite.Subtest("When all parameters is valid but revoke return not found, then it would return error", func() {
			suite.oauthAccessTokenRepo.EXPECT().Revoke(suite.ktx, *suite.revokeRequest.Token, suite.sqldb).Return(exception.Throw(errors.New("not found"), exception.WithType(exception.NotFound), exception.WithDetail("oauth access token is not found"), exception.WithTitle("Not Found")))
			err := suite.authorization.RevokeToken(suite.ktx, suite.revokeRequest)
			suite.Assert().Equal("JSONAPI Error:\n[Not Found] Detail: oauth access token is not found, Code: ERR0404\n", err.Error())
			suite.Assert().Equal(http.StatusNotFound, err.HTTPStatus())
		})

		suite.Subtest("When all parameters is valid but revoke return unexpected error, then it would return error", func() {
			suite.oauthAccessTokenRepo.EXPECT().Revoke(suite.ktx, *suite.revokeRequest.Token, suite.sqldb).Return(exception.Throw(errors.New("unexpected"), exception.WithType(exception.Unexpected)))
			err := suite.authorization.RevokeToken(suite.ktx, suite.revokeRequest)
			suite.Assert().Equal("JSONAPI Error:\n[Internal server error] Detail: Something is not right, help us fix this problem. Contribute to https://github.com/kodefluence/altair. Tracing code: '<nil>', Code: ERR0500\n", err.Error())
			suite.Assert().Equal(http.StatusInternalServerError, err.HTTPStatus())
		})
	})
}
