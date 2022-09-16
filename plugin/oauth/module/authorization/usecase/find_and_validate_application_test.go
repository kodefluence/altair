package usecase_test

import (
	"database/sql"
	"errors"
	"net/http"
	"testing"

	"github.com/kodefluence/altair/plugin/oauth/entity"
	"github.com/kodefluence/altair/util"
	"github.com/kodefluence/monorepo/exception"
	"github.com/stretchr/testify/suite"
)

type FindAndValidateApplicationSuiteTest struct {
	*AuthorizationBaseSuiteTest
	oauthApplication entity.OauthApplication
	clientUID        *string
	clientSecret     *string
}

func TestFindAndValidateApplication(t *testing.T) {
	suite.Run(t, &FindAndValidateApplicationSuiteTest{
		AuthorizationBaseSuiteTest: &AuthorizationBaseSuiteTest{},
	})
}

func (suite *FindAndValidateApplicationSuiteTest) SetupTest() {
	suite.oauthApplication = entity.OauthApplication{
		ID: 1,
		Scopes: sql.NullString{
			String: "public users",
			Valid:  true,
		},
	}

	suite.clientUID = util.StringToPointer("client_uid")
	suite.clientSecret = util.StringToPointer("client_secret")
}

func (suite *FindAndValidateApplicationSuiteTest) TestFindAndValidateApplication() {
	suite.Run("Positive cases", func() {
		suite.Subtest("When all parameter is valid and oauth application found in database, then it would return nil", func() {
			suite.oauthApplicationRepo.EXPECT().OneByUIDandSecret(suite.ktx, *suite.clientUID, *suite.clientSecret, suite.sqldb).Return(suite.oauthApplication, nil)
			oauthApplication, err := suite.authorization.FindAndValidateApplication(suite.ktx, suite.clientUID, suite.clientSecret)
			suite.Assert().Nil(err)
			suite.Assert().Equal(suite.oauthApplication, oauthApplication)
		})
	})

	suite.Run("Negative cases", func() {
		suite.Subtest("When client_uid is nil, then it would return error", func() {
			suite.clientUID = nil
			oauthApplication, err := suite.authorization.FindAndValidateApplication(suite.ktx, suite.clientUID, suite.clientSecret)
			suite.Assert().Equal("JSONAPI Error:\n[Validation error] Detail: Validation error because of: client_uid cannot be empty, Code: ERR1442\n", err.Error())
			suite.Assert().Equal(http.StatusUnprocessableEntity, err.HTTPStatus())
			suite.Assert().Equal(entity.OauthApplication{}, oauthApplication)
		})

		suite.Subtest("When client_secret is nil, then it would return error", func() {
			suite.clientSecret = nil
			oauthApplication, err := suite.authorization.FindAndValidateApplication(suite.ktx, suite.clientUID, suite.clientSecret)
			suite.Assert().Equal("JSONAPI Error:\n[Validation error] Detail: Validation error because of: client_secret cannot be empty, Code: ERR1442\n", err.Error())
			suite.Assert().Equal(http.StatusUnprocessableEntity, err.HTTPStatus())
			suite.Assert().Equal(entity.OauthApplication{}, oauthApplication)
		})

		suite.Subtest("When all parameter is valid but oauth application repo return unexpected error, then it would return error", func() {
			suite.oauthApplicationRepo.EXPECT().OneByUIDandSecret(suite.ktx, *suite.clientUID, *suite.clientSecret, suite.sqldb).Return(entity.OauthApplication{}, exception.Throw(
				errors.New("unexpected"),
				exception.WithType(exception.Unexpected),
			))
			oauthApplication, err := suite.authorization.FindAndValidateApplication(suite.ktx, suite.clientUID, suite.clientSecret)
			suite.Assert().Equal("JSONAPI Error:\n[Internal server error] Detail: Something is not right, help us fix this problem. Contribute to https://github.com/kodefluence/altair. Tracing code: '<nil>', Code: ERR0500\n", err.Error())
			suite.Assert().Equal(http.StatusInternalServerError, err.HTTPStatus())
			suite.Assert().Equal(entity.OauthApplication{}, oauthApplication)
		})

		suite.Subtest("When all parameter is valid but oauth application repo return notfound, then it would return error", func() {
			suite.oauthApplicationRepo.EXPECT().OneByUIDandSecret(suite.ktx, *suite.clientUID, *suite.clientSecret, suite.sqldb).Return(entity.OauthApplication{}, exception.Throw(
				errors.New("not found"),
				exception.WithType(exception.NotFound),
			))
			oauthApplication, err := suite.authorization.FindAndValidateApplication(suite.ktx, suite.clientUID, suite.clientSecret)
			suite.Assert().Equal("JSONAPI Error:\n[Not found error] Detail: Resource of `client_uid & client_secret` is not found. Tracing code: `<nil>`, Code: ERR0404\n", err.Error())
			suite.Assert().Equal(http.StatusNotFound, err.HTTPStatus())
			suite.Assert().Equal(entity.OauthApplication{}, oauthApplication)
		})
	})
}

func (suite *FindAndValidateApplicationSuiteTest) Subtest(testcase string, subtest func()) {
	suite.SetupTest()
	suite.AuthorizationBaseSuiteTest.Subtest(testcase, subtest)
	suite.TearDownTest()
}
