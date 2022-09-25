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
	"github.com/kodefluence/altair/util"
	"github.com/kodefluence/monorepo/db"
	"github.com/kodefluence/monorepo/exception"
	"github.com/kodefluence/monorepo/kontext"
	"github.com/stretchr/testify/suite"
)

type ImplicitGrantSuiteTest struct {
	*AuthorizationBaseSuiteTest

	authorizationRequestJSON entity.AuthorizationRequestJSON
	oauthApplication         entity.OauthApplication
	accessToken              entity.OauthAccessToken
}

func TestImplicitGrant(t *testing.T) {
	suite.Run(t, &ImplicitGrantSuiteTest{
		AuthorizationBaseSuiteTest: &AuthorizationBaseSuiteTest{},
	})
}

func (suite *ImplicitGrantSuiteTest) SetupTest() {
	suite.authorizationRequestJSON = entity.AuthorizationRequestJSON{
		ResponseType:    util.StringToPointer("token"),
		ResourceOwnerID: util.IntToPointer(1),
		RedirectURI:     util.StringToPointer("www.github.com"),
		ClientUID:       util.StringToPointer("client_uid"),
		ClientSecret:    util.StringToPointer("client_secret"),
		Scopes:          util.StringToPointer(""),
	}
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
}

func (suite *ImplicitGrantSuiteTest) TestImplicitGrant() {
	suite.Run("Positive cases", func() {
		suite.Subtest("When all parameters is valid, it would return oauth access token", func() {
			gomock.InOrder(
				suite.oauthApplicationRepo.EXPECT().OneByUIDandSecret(suite.ktx, *suite.authorizationRequestJSON.ClientUID, *suite.authorizationRequestJSON.ClientSecret, suite.sqldb).Return(suite.oauthApplication, nil),
				suite.sqldb.EXPECT().Transaction(suite.ktx, "authorization-implicit-grant", gomock.Any()).DoAndReturn(func(ktx kontext.Context, transactionKey string, f func(tx db.TX) exception.Exception) exception.Exception {
					suite.oauthAccessTokenRepo.EXPECT().Create(ktx, gomock.Any(), suite.sqldb).DoAndReturn(func(ktx kontext.Context, data entity.OauthAccessTokenInsertable, tx db.TX) (int, exception.Exception) {
						suite.Assert().Equal(suite.oauthApplication.ID, data.OauthApplicationID)
						return suite.accessToken.ID, nil
					})
					suite.oauthAccessTokenRepo.EXPECT().One(ktx, suite.accessToken.ID, suite.sqldb).Return(suite.accessToken, nil)

					suite.Assert().Nil(f(suite.sqldb))
					return nil
				}),
			)

			finalJson, err := suite.authorization.ImplicitGrant(suite.ktx, suite.authorizationRequestJSON)
			suite.Assert().Nil(err)
			suite.Assert().Equal(suite.formatter.AccessToken(suite.accessToken, *suite.authorizationRequestJSON.RedirectURI, nil), finalJson)
		})
	})

	suite.Run("Negative cases", func() {
		suite.Subtest("When client_uid is nil, then it would return error", func() {
			suite.authorizationRequestJSON.ClientUID = nil
			finalJson, err := suite.authorization.ImplicitGrant(suite.ktx, suite.authorizationRequestJSON)
			suite.Assert().Equal("JSONAPI Error:\n[Validation error] Detail: Validation error because of: client_uid cannot be empty, Code: ERR1442\n", err.Error())
			suite.Assert().Equal(http.StatusUnprocessableEntity, err.HTTPStatus())
			suite.Assert().Equal(entity.OauthAccessTokenJSON{}, finalJson)
		})

		suite.Subtest("When all parameter is valid, but scope is available in oauth application then it would return error", func() {
			suite.authorizationRequestJSON.Scopes = util.StringToPointer("public users admin")
			suite.oauthApplicationRepo.EXPECT().OneByUIDandSecret(suite.ktx, *suite.authorizationRequestJSON.ClientUID, *suite.authorizationRequestJSON.ClientSecret, suite.sqldb).Return(suite.oauthApplication, nil)
			finalJson, err := suite.authorization.ImplicitGrant(suite.ktx, suite.authorizationRequestJSON)
			suite.Assert().Equal("JSONAPI Error:\n[Forbidden error] Detail: Resource of `application` is forbidden to be accessed, because of: your requested scopes `([admin])` is not exists in application. Tracing code: `<nil>`, Code: ERR0403\n", err.Error())
			suite.Assert().Equal(http.StatusForbidden, err.HTTPStatus())
			suite.Assert().Equal(entity.OauthAccessTokenJSON{}, finalJson)
		})

		suite.Subtest("When all parameters is valid, but create token return an error then it would return error", func() {
			gomock.InOrder(
				suite.oauthApplicationRepo.EXPECT().OneByUIDandSecret(suite.ktx, *suite.authorizationRequestJSON.ClientUID, *suite.authorizationRequestJSON.ClientSecret, suite.sqldb).Return(suite.oauthApplication, nil),
				suite.sqldb.EXPECT().Transaction(suite.ktx, "authorization-implicit-grant", gomock.Any()).DoAndReturn(func(ktx kontext.Context, transactionKey string, f func(tx db.TX) exception.Exception) exception.Exception {
					exc := exception.Throw(errors.New("unexpected"))
					suite.oauthAccessTokenRepo.EXPECT().Create(ktx, gomock.Any(), suite.sqldb).DoAndReturn(func(ktx kontext.Context, data entity.OauthAccessTokenInsertable, tx db.TX) (int, exception.Exception) {
						suite.Assert().Equal(suite.oauthApplication.ID, data.OauthApplicationID)
						return 0, exc
					})

					suite.Assert().NotNil(f(suite.sqldb))
					return exc
				}),
			)

			finalJson, err := suite.authorization.ImplicitGrant(suite.ktx, suite.authorizationRequestJSON)
			suite.Assert().Equal("JSONAPI Error:\n[Internal server error] Detail: Something is not right, help us fix this problem. Contribute to https://github.com/kodefluence/altair. Tracing code: '<nil>', Code: ERR0500\n", err.Error())
			suite.Assert().Equal(http.StatusInternalServerError, err.HTTPStatus())
			suite.Assert().Equal(entity.OauthAccessTokenJSON{}, finalJson)
		})

		suite.Subtest("When all parameters is valid, but find one access token return error then it would return error", func() {
			gomock.InOrder(
				suite.oauthApplicationRepo.EXPECT().OneByUIDandSecret(suite.ktx, *suite.authorizationRequestJSON.ClientUID, *suite.authorizationRequestJSON.ClientSecret, suite.sqldb).Return(suite.oauthApplication, nil),
				suite.sqldb.EXPECT().Transaction(suite.ktx, "authorization-implicit-grant", gomock.Any()).DoAndReturn(func(ktx kontext.Context, transactionKey string, f func(tx db.TX) exception.Exception) exception.Exception {
					exc := exception.Throw(errors.New("unexpected"))
					suite.oauthAccessTokenRepo.EXPECT().Create(ktx, gomock.Any(), suite.sqldb).DoAndReturn(func(ktx kontext.Context, data entity.OauthAccessTokenInsertable, tx db.TX) (int, exception.Exception) {
						suite.Assert().Equal(suite.oauthApplication.ID, data.OauthApplicationID)
						return suite.accessToken.ID, nil
					})
					suite.oauthAccessTokenRepo.EXPECT().One(ktx, suite.accessToken.ID, suite.sqldb).Return(entity.OauthAccessToken{}, exc)

					suite.Assert().NotNil(f(suite.sqldb))
					return exc
				}),
			)

			finalJson, err := suite.authorization.ImplicitGrant(suite.ktx, suite.authorizationRequestJSON)
			suite.Assert().Equal("JSONAPI Error:\n[Internal server error] Detail: Something is not right, help us fix this problem. Contribute to https://github.com/kodefluence/altair. Tracing code: '<nil>', Code: ERR0500\n", err.Error())
			suite.Assert().Equal(http.StatusInternalServerError, err.HTTPStatus())
			suite.Assert().Equal(entity.OauthAccessTokenJSON{}, finalJson)
		})
	})
}

func (suite *ImplicitGrantSuiteTest) Subtest(testcase string, subtest func()) {
	suite.SetupTest()
	suite.AuthorizationBaseSuiteTest.Subtest(testcase, subtest)
	suite.TearDownTest()
}
