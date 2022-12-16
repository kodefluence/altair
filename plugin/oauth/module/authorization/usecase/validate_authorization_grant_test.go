package usecase_test

import (
	"database/sql"
	"net/http"
	"testing"

	"github.com/kodefluence/altair/plugin/oauth/entity"
	"github.com/kodefluence/altair/util"
	"github.com/stretchr/testify/suite"
)

type ValidateAuthorizationGrantSuiteTest struct {
	*AuthorizationBaseSuiteTest

	authorizationRequestJSON entity.AuthorizationRequestJSON
	oauthApplication         entity.OauthApplication
}

func TestValidateAuthorizationGrant(t *testing.T) {
	suite.Run(t, &ValidateAuthorizationGrantSuiteTest{
		AuthorizationBaseSuiteTest: &AuthorizationBaseSuiteTest{},
	})
}

func (suite *ValidateAuthorizationGrantSuiteTest) SetupTest() {
	suite.authorizationRequestJSON = entity.AuthorizationRequestJSON{
		ResponseType:    util.ValueToPointer("code"),
		ResourceOwnerID: util.ValueToPointer(1),
		RedirectURI:     util.ValueToPointer("www.github.com"),
		Scopes:          util.ValueToPointer(""),
	}
	suite.oauthApplication = entity.OauthApplication{
		Scopes: sql.NullString{
			String: "public users",
			Valid:  true,
		},
	}
}

func (suite *ValidateAuthorizationGrantSuiteTest) TestValidateAuthorizationGrant() {
	suite.Run("Positive cases", func() {
		suite.Subtest("When all parameter is valid, but scope is empty then it would return nil", func() {
			err := suite.authorization.ValidateAuthorizationGrant(suite.ktx, suite.authorizationRequestJSON, suite.oauthApplication)
			suite.Assert().Nil(err)
		})

		suite.Subtest("When all parameter is valid, but scope is nil then it would return nil", func() {
			suite.authorizationRequestJSON.Scopes = nil
			err := suite.authorization.ValidateAuthorizationGrant(suite.ktx, suite.authorizationRequestJSON, suite.oauthApplication)
			suite.Assert().Nil(err)
		})

		suite.Subtest("When all parameter is valid, and scope is available in oauth application then it would return nil", func() {
			suite.authorizationRequestJSON.Scopes = util.ValueToPointer("public users")
			err := suite.authorization.ValidateAuthorizationGrant(suite.ktx, suite.authorizationRequestJSON, suite.oauthApplication)
			suite.Assert().Nil(err)
		})
	})

	suite.Run("Negative cases", func() {
		suite.Subtest("When all parameter is valid, but scope is available in oauth application then it would return error", func() {
			suite.authorizationRequestJSON.Scopes = util.ValueToPointer("public users admin")
			err := suite.authorization.ValidateAuthorizationGrant(suite.ktx, suite.authorizationRequestJSON, suite.oauthApplication)
			suite.Assert().Equal("JSONAPI Error:\n[Forbidden error] Detail: Resource of `application` is forbidden to be accessed, because of: your requested scopes `([admin])` is not exists in application. Tracing code: `<nil>`, Code: ERR0403\n", err.Error())
			suite.Assert().Equal(http.StatusForbidden, err.HTTPStatus())
		})

		suite.Subtest("When all parameter is invalid, then it would return error", func() {
			suite.authorizationRequestJSON = entity.AuthorizationRequestJSON{
				Scopes: util.ValueToPointer(""),
			}
			err := suite.authorization.ValidateAuthorizationGrant(suite.ktx, suite.authorizationRequestJSON, suite.oauthApplication)
			suite.Assert().Equal("JSONAPI Error:\n[Validation error] Detail: Validation error because of: response_type can't be empty, Code: ERR1442\n[Validation error] Detail: Validation error because of: resource_owner_id can't be empty, Code: ERR1442\n[Validation error] Detail: Validation error because of: redirect_uri can't be empty, Code: ERR1442\n", err.Error())
			suite.Assert().Equal(http.StatusUnprocessableEntity, err.HTTPStatus())
		})
		suite.Subtest("When response type is token but application is not confidential, then it would return error", func() {
			suite.authorizationRequestJSON.ResponseType = util.ValueToPointer("token")
			suite.oauthApplication.OwnerType = "public"
			err := suite.authorization.ValidateAuthorizationGrant(suite.ktx, suite.authorizationRequestJSON, suite.oauthApplication)
			suite.Assert().Equal("JSONAPI Error:\n[Forbidden error] Detail: Resource of `access_token` is forbidden to be accessed, because of: your response type is not allowed in this application. Tracing code: `<nil>`, Code: ERR0403\n", err.Error())
			suite.Assert().Equal(http.StatusForbidden, err.HTTPStatus())
		})
	})
}

func (suite *ValidateAuthorizationGrantSuiteTest) Subtest(testcase string, subtest func()) {
	suite.SetupTest()
	suite.AuthorizationBaseSuiteTest.Subtest(testcase, subtest)
	suite.TearDownTest()
}
