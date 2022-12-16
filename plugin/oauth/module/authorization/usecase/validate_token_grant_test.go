package usecase_test

import (
	"testing"

	"github.com/kodefluence/altair/plugin/oauth/entity"
	"github.com/kodefluence/altair/util"
	"github.com/stretchr/testify/suite"
)

type ValidateTokenGrantSuiteTest struct {
	*AuthorizationBaseSuiteTest

	accessTokenRequestJSON entity.AccessTokenRequestJSON
}

func TestValidateTokenGrant(t *testing.T) {
	suite.Run(t, &ValidateTokenGrantSuiteTest{
		AuthorizationBaseSuiteTest: &AuthorizationBaseSuiteTest{},
	})
}

func (suite *ValidateTokenGrantSuiteTest) SetupTest() {
	suite.accessTokenRequestJSON = entity.AccessTokenRequestJSON{
		GrantType:    util.ValueToPointer("authorization_code"),
		ClientUID:    util.ValueToPointer("client_uid"),
		ClientSecret: util.ValueToPointer("client_secret"),
		RefreshToken: util.ValueToPointer("some-refresh-token"),
		Code:         util.ValueToPointer("some-code"),
		RedirectURI:  util.ValueToPointer("https://github.com/kodefluence/altair"),
	}
}

func (suite *ValidateTokenGrantSuiteTest) TestValidateTokenGrantSuiteTest() {
	suite.Run("Positive cases", func() {
		suite.Subtest("When request param is valid and grant_type is authorization_code, then it would return nil", func() {
			err := suite.authorization.ValidateTokenGrant(suite.accessTokenRequestJSON)
			suite.Nil(err)
		})

		suite.Subtest("When request param is valid and grant_type is client_credentials, then it would return nil", func() {
			suite.accessTokenRequestJSON.GrantType = util.ValueToPointer("client_credentials")
			err := suite.authorization.ValidateTokenGrant(suite.accessTokenRequestJSON)
			suite.Nil(err)
		})
	})

	suite.Run("Negative cases", func() {
		suite.Subtest("When grant type is nil and grant_type is authorization_code, then it would return error", func() {
			suite.accessTokenRequestJSON.GrantType = nil
			err := suite.authorization.ValidateTokenGrant(suite.accessTokenRequestJSON)
			suite.Equal("JSONAPI Error:\n[Validation error] Detail: Validation error because of: grant_type can't be empty, Code: ERR1442\n", err.Error())
		})

		suite.Subtest("When grant type is empty, then it would return error", func() {
			suite.accessTokenRequestJSON.GrantType = util.ValueToPointer("")
			err := suite.authorization.ValidateTokenGrant(suite.accessTokenRequestJSON)
			suite.Equal("JSONAPI Error:\n[Validation error] Detail: Validation error because of: grant_type is not valid value, Code: ERR1442\n", err.Error())
		})

		suite.Subtest("When grant type is invalid value, then it would return error", func() {
			suite.accessTokenRequestJSON.GrantType = util.ValueToPointer("invalid value")
			err := suite.authorization.ValidateTokenGrant(suite.accessTokenRequestJSON)
			suite.Equal("JSONAPI Error:\n[Validation error] Detail: Validation error because of: grant_type is not valid value, Code: ERR1442\n", err.Error())
		})

		suite.Subtest("When grant_type is authorization_code and code is nil, then it would return error", func() {
			suite.accessTokenRequestJSON.Code = nil
			err := suite.authorization.ValidateTokenGrant(suite.accessTokenRequestJSON)
			suite.Equal("JSONAPI Error:\n[Validation error] Detail: Validation error because of: code is not valid value, Code: ERR1442\n", err.Error())
		})

		suite.Subtest("When grant_type is authorization_code and redirect_uri is nil, then it would return error", func() {
			suite.accessTokenRequestJSON.RedirectURI = nil
			err := suite.authorization.ValidateTokenGrant(suite.accessTokenRequestJSON)
			suite.Equal("JSONAPI Error:\n[Validation error] Detail: Validation error because of: redirect_uri is not valid value, Code: ERR1442\n", err.Error())
		})

		suite.Subtest("When grant_type is authorization_code with nil code and redirect_uri, then it would return error", func() {
			suite.accessTokenRequestJSON.RedirectURI = nil
			suite.accessTokenRequestJSON.Code = nil
			err := suite.authorization.ValidateTokenGrant(suite.accessTokenRequestJSON)
			suite.Equal("JSONAPI Error:\n[Validation error] Detail: Validation error because of: code is not valid value, Code: ERR1442\n[Validation error] Detail: Validation error because of: redirect_uri is not valid value, Code: ERR1442\n", err.Error())
		})

		suite.Subtest("When grant_type is refresh_token with nil refresh_token, then it would return error", func() {
			suite.accessTokenRequestJSON.GrantType = util.ValueToPointer("refresh_token")
			suite.accessTokenRequestJSON.RefreshToken = nil
			err := suite.authorization.ValidateTokenGrant(suite.accessTokenRequestJSON)
			suite.Equal("JSONAPI Error:\n[Validation error] Detail: Validation error because of: refresh_token is not valid value, Code: ERR1442\n", err.Error())
		})

	})
}

func (suite *ValidateTokenGrantSuiteTest) Subtest(testcase string, subtest func()) {
	suite.SetupTest()
	suite.AuthorizationBaseSuiteTest.Subtest(testcase, subtest)
	suite.TearDownTest()
}
