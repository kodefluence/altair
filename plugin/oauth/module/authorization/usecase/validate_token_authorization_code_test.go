package usecase_test

import (
	"database/sql"
	"testing"
	"time"

	"github.com/kodefluence/altair/plugin/oauth/entity"
	"github.com/kodefluence/altair/util"
	"github.com/stretchr/testify/suite"
)

type ValidateTokenAuthorizationCodeSuiteTest struct {
	*AuthorizationBaseSuiteTest

	accessTokenRequestJSON entity.AccessTokenRequestJSON
	accessGrant            entity.OauthAccessGrant
}

func TestValidateTokenAuthorizationCode(t *testing.T) {
	suite.Run(t, &ValidateTokenAuthorizationCodeSuiteTest{
		AuthorizationBaseSuiteTest: &AuthorizationBaseSuiteTest{},
	})
}

func (suite *ValidateTokenAuthorizationCodeSuiteTest) SetupTest() {
	suite.accessTokenRequestJSON = entity.AccessTokenRequestJSON{
		GrantType:    util.StringToPointer("authorization_code"),
		ClientUID:    util.StringToPointer("client_uid"),
		ClientSecret: util.StringToPointer("client_secret"),
		RefreshToken: util.StringToPointer("some-refresh-token"),
		Code:         util.StringToPointer("some-code"),
		RedirectURI:  util.StringToPointer("https://github.com/kodefluence/altair"),
	}
	suite.accessGrant = entity.OauthAccessGrant{
		ID:                 1,
		OauthApplicationID: 1,
		ResourceOwnerID:    0,
		Code:               "some-authorization-code",
		RedirectURI: sql.NullString{
			String: *suite.accessTokenRequestJSON.RedirectURI,
			Valid:  true,
		},
		Scopes:    sql.NullString{},
		ExpiresIn: time.Now().Add(time.Hour),
		CreatedAt: time.Now().Add(-24 * time.Hour),
		RevokedAT: sql.NullTime{},
	}
}

func (suite *ValidateTokenAuthorizationCodeSuiteTest) TestValidateTokenGrantSuiteTest() {
	suite.Run("Positive cases", func() {
		suite.Subtest("When all parameters is valid, then it would return nil", func() {
			exc := suite.authorization.ValidateTokenAuthorizationCode(suite.ktx, suite.accessTokenRequestJSON, suite.accessGrant)
			suite.Assert().Nil(exc)
		})
	})

	suite.Run("Negative cases", func() {
		suite.Subtest("When access grant is already revoked, then it would return jsonapi option", func() {
			suite.accessGrant.RevokedAT = sql.NullTime{
				Time:  time.Now().Add(-1 * time.Hour),
				Valid: true,
			}
			exc := suite.authorization.ValidateTokenAuthorizationCode(suite.ktx, suite.accessTokenRequestJSON, suite.accessGrant)
			suite.Assert().NotNil(exc)
			suite.Assert().Equal("forbidden", exc.Error())
			suite.Assert().Equal("Forbidden resource access", exc.Title())
			suite.Assert().Equal("authorization code already used", exc.Detail())
		})

		suite.Subtest("When access grant is already expired, then it would return jsonapi option", func() {
			suite.accessGrant.ExpiresIn = time.Now().Add(-1 * time.Hour)
			exc := suite.authorization.ValidateTokenAuthorizationCode(suite.ktx, suite.accessTokenRequestJSON, suite.accessGrant)
			suite.Assert().NotNil(exc)
			suite.Assert().Equal("forbidden", exc.Error())
			suite.Assert().Equal("Forbidden resource access", exc.Title())
			suite.Assert().Equal("authorization code already expired", exc.Detail())
		})

		suite.Subtest("When access grant redirect uri is different, then it would return jsonapi option", func() {
			suite.accessGrant.RedirectURI.String = ""
			exc := suite.authorization.ValidateTokenAuthorizationCode(suite.ktx, suite.accessTokenRequestJSON, suite.accessGrant)
			suite.Assert().NotNil(exc)
			suite.Assert().Equal("forbidden", exc.Error())
			suite.Assert().Equal("Forbidden resource access", exc.Title())
			suite.Assert().Equal("redirect uri is different from one that generated before", exc.Detail())
		})
	})
}

func (suite *ValidateTokenAuthorizationCodeSuiteTest) Subtest(testcase string, subtest func()) {
	suite.SetupTest()
	suite.AuthorizationBaseSuiteTest.Subtest(testcase, subtest)
	suite.TearDownTest()
}
