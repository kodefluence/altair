package usecase_test

import (
	"database/sql"
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

type GrantorSuiteTest struct {
	*AuthorizationBaseSuiteTest

	authorizationRequestJSON entity.AuthorizationRequestJSON
	oauthApplication         entity.OauthApplication
	accessToken              entity.OauthAccessToken
	accessGrant              entity.OauthAccessGrant
}

func TestGrantor(t *testing.T) {
	suite.Run(t, &GrantorSuiteTest{
		AuthorizationBaseSuiteTest: &AuthorizationBaseSuiteTest{},
	})
}

func (suite *GrantorSuiteTest) SetupTest() {
	suite.authorizationRequestJSON = entity.AuthorizationRequestJSON{
		ResponseType:    util.ValueToPointer("token"),
		ResourceOwnerID: util.ValueToPointer(1),
		RedirectURI:     util.ValueToPointer("www.github.com"),
		ClientUID:       util.ValueToPointer("client_uid"),
		ClientSecret:    util.ValueToPointer("client_secret"),
		Scopes:          util.ValueToPointer(""),
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
		RevokedAT:          sql.NullTime{},
	}
	suite.accessGrant = entity.OauthAccessGrant{
		ID:                 1,
		OauthApplicationID: 1,
		ResourceOwnerID:    0,
		Code:               "",
		RedirectURI:        sql.NullString{},
		Scopes:             sql.NullString{},
		ExpiresIn:          time.Time{},
		CreatedAt:          time.Time{},
		RevokedAT:          sql.NullTime{},
	}
}

func (suite *GrantorSuiteTest) Subtest(testcase string, subtest func()) {
	suite.SetupTest()
	suite.AuthorizationBaseSuiteTest.Subtest(testcase, subtest)
	suite.TearDownTest()
}

func (suite *GrantorSuiteTest) TestGrantor() {
	suite.Run("Positive cases", func() {
		suite.Subtest("When all parameters is valid and response type is token, it would return oauth access token", func() {
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

			finalJson, err := suite.authorization.GrantAuthorizationCode(suite.ktx, suite.authorizationRequestJSON)
			suite.Assert().Nil(err)
			suite.Assert().Equal(suite.formatter.AccessToken(suite.accessToken, *suite.authorizationRequestJSON.RedirectURI, nil), finalJson)
		})

		suite.Subtest("When all parameters is valid and response type is code, it would return oauth access grant", func() {
			suite.authorizationRequestJSON.ResponseType = util.ValueToPointer("code")
			gomock.InOrder(
				suite.oauthApplicationRepo.EXPECT().OneByUIDandSecret(suite.ktx, *suite.authorizationRequestJSON.ClientUID, *suite.authorizationRequestJSON.ClientSecret, suite.sqldb).Return(suite.oauthApplication, nil),
				suite.sqldb.EXPECT().Transaction(suite.ktx, "authorization-grant-authorization-code", gomock.Any()).DoAndReturn(func(ktx kontext.Context, transactionKey string, f func(tx db.TX) exception.Exception) exception.Exception {
					suite.oauthAccessGrantRepo.EXPECT().Create(ktx, gomock.Any(), suite.sqldb).DoAndReturn(func(ktx kontext.Context, data entity.OauthAccessGrantInsertable, tx db.TX) (int, exception.Exception) {
						suite.Assert().Equal(suite.oauthApplication.ID, data.OauthApplicationID)
						return suite.accessGrant.ID, nil
					})
					suite.oauthAccessGrantRepo.EXPECT().One(ktx, suite.accessGrant.ID, suite.sqldb).Return(suite.accessGrant, nil)

					suite.Assert().Nil(f(suite.sqldb))
					return nil
				}),
			)

			finalJson, err := suite.authorization.GrantAuthorizationCode(suite.ktx, suite.authorizationRequestJSON)
			suite.Assert().Nil(err)
			suite.Assert().Equal(suite.formatter.AccessGrant(suite.accessGrant), finalJson)
		})
	})

	suite.Run("Negative cases", func() {
		suite.Subtest("When response type is nil, then it would return error", func() {
			suite.authorizationRequestJSON.ResponseType = nil
			finalJson, err := suite.authorization.GrantAuthorizationCode(suite.ktx, suite.authorizationRequestJSON)
			suite.Assert().Equal("JSONAPI Error:\n[Validation error] Detail: Validation error because of: response_type cannot be empty, Code: ERR1442\n", err.Error())
			suite.Assert().Equal(http.StatusUnprocessableEntity, err.HTTPStatus())
			suite.Assert().Equal(nil, finalJson)
		})

		suite.Subtest("When response type is empty, then it would return error", func() {
			suite.authorizationRequestJSON.ResponseType = util.ValueToPointer("")
			finalJson, err := suite.authorization.GrantAuthorizationCode(suite.ktx, suite.authorizationRequestJSON)
			suite.Assert().Equal("JSONAPI Error:\n[Validation error] Detail: Validation error because of: response_type is invalid. Should be either `token` or `code`, Code: ERR1442\n", err.Error())
			suite.Assert().Equal(http.StatusUnprocessableEntity, err.HTTPStatus())
			suite.Assert().Equal(nil, finalJson)
		})

		suite.Subtest("When response type is invalid, then it would return error", func() {
			suite.authorizationRequestJSON.ResponseType = util.ValueToPointer("client_credentials")
			finalJson, err := suite.authorization.GrantAuthorizationCode(suite.ktx, suite.authorizationRequestJSON)
			suite.Assert().Equal("JSONAPI Error:\n[Validation error] Detail: Validation error because of: response_type is invalid. Should be either `token` or `code`, Code: ERR1442\n", err.Error())
			suite.Assert().Equal(http.StatusUnprocessableEntity, err.HTTPStatus())
			suite.Assert().Equal(nil, finalJson)
		})

		suite.Subtest("When response type is token but implicit grant feature is inactive, then it would return error", func() {
			suite.config.Config.ImplicitGrant.Active = false
			suite.authorization = usecase.NewAuthorization(suite.oauthApplicationRepo, suite.oauthAccessTokenRepo, suite.oauthAccessGrantRepo, suite.oauthRefreshTokenRepo, suite.formatter, suite.config, suite.sqldb, suite.apiError)
			finalJson, err := suite.authorization.GrantAuthorizationCode(suite.ktx, suite.authorizationRequestJSON)
			suite.Assert().Equal("JSONAPI Error:\n[Validation error] Detail: Validation error because of: response_type is invalid. Should be either `token` or `code`, Code: ERR1442\n", err.Error())
			suite.Assert().Equal(http.StatusUnprocessableEntity, err.HTTPStatus())
			suite.Assert().Equal(nil, finalJson)
		})
	})
}
