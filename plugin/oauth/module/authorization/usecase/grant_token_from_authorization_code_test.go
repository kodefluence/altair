package usecase_test

import (
	"database/sql"
	"encoding/json"
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

type GrantTokenFromAuthorizationCodeTest struct {
	*AuthorizationBaseSuiteTest

	oauthApplication       entity.OauthApplication
	accessTokenRequestJSON entity.AccessTokenRequestJSON
	accessGrant            entity.OauthAccessGrant
	accessToken            entity.OauthAccessToken
	refreshToken           entity.OauthRefreshToken
	refreshTokenJSON       entity.OauthRefreshTokenJSON
}

func TestGrantTokenFromAuthorizationCode(t *testing.T) {
	suite.Run(t, &GrantTokenFromAuthorizationCodeTest{
		AuthorizationBaseSuiteTest: &AuthorizationBaseSuiteTest{},
	})
}

func (suite *GrantTokenFromAuthorizationCodeTest) SetupTest() {
	suite.oauthApplication = entity.OauthApplication{
		ID: 1,
		Scopes: sql.NullString{
			String: "public users",
			Valid:  true,
		},
		OwnerType: "confidential",
	}
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
		RevokedAT: mysql.NullTime{},
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
		Token:              "some token",
		ExpiresIn:          time.Time{},
		CreatedAt:          time.Time{},
		RevokedAT:          mysql.NullTime{},
	}
}

func (suite *GrantTokenFromAuthorizationCodeTest) TestValidateTokenGrantSuiteTest() {
	suite.Run("Positive cases", func() {
		suite.Subtest("When all parameters is valid, then it would return nil", func() {
			suite.refreshTokenJSON = suite.formatter.RefreshToken(suite.refreshToken)
			gomock.InOrder(
				suite.oauthApplicationRepo.EXPECT().OneByUIDandSecret(suite.ktx, *suite.accessTokenRequestJSON.ClientUID, *suite.accessTokenRequestJSON.ClientSecret, suite.sqldb).Return(suite.oauthApplication, nil),
				suite.sqldb.EXPECT().Transaction(suite.ktx, "authorization-grant-token-from-refresh-token", gomock.Any()).DoAndReturn(func(ctx kontext.Context, transactionKey string, f func(tx db.TX) exception.Exception) exception.Exception {
					suite.oauthAccessGrantRepo.EXPECT().OneByCode(suite.ktx, *suite.accessTokenRequestJSON.Code, suite.sqldb).Return(suite.accessGrant, nil)
					suite.oauthAccessTokenRepo.EXPECT().Create(suite.ktx, gomock.Any(), suite.sqldb).DoAndReturn(func(ktx kontext.Context, data entity.OauthAccessTokenInsertable, tx db.TX) (int, exception.Exception) {
						suite.Assert().Equal(suite.accessGrant.Scopes.String, data.Scopes)
						suite.Assert().Equal(suite.oauthApplication.ID, data.OauthApplicationID)
						return 1, nil
					})
					suite.oauthAccessTokenRepo.EXPECT().One(suite.ktx, 1, suite.sqldb).Return(suite.accessToken, nil)
					suite.oauthAccessGrantRepo.EXPECT().Revoke(suite.ktx, *suite.accessTokenRequestJSON.Code, suite.sqldb).Return(nil)
					suite.oauthRefreshTokenRepo.EXPECT().Create(suite.ktx, gomock.Any(), suite.sqldb).DoAndReturn(func(ktx kontext.Context, data entity.OauthRefreshTokenInsertable, tx db.TX) (int, exception.Exception) {
						suite.Assert().Equal(suite.accessToken.ID, data.OauthAccessTokenID)
						return 1, nil
					})
					suite.oauthRefreshTokenRepo.EXPECT().One(suite.ktx, suite.refreshToken.ID, suite.sqldb).Return(suite.refreshToken, nil)
					f(suite.sqldb)
					return nil
				}),
			)

			accessTokenJSON, err := suite.authorization.Token(suite.ktx, suite.accessTokenRequestJSON)
			byteAccessToken, _ := json.Marshal(accessTokenJSON)
			byteExpectedAccessToken, _ := json.Marshal(suite.formatter.AccessToken(suite.accessToken, suite.accessGrant.RedirectURI.String, &suite.refreshTokenJSON))
			suite.Assert().Nil(err)
			suite.Equal(string(byteExpectedAccessToken), string(byteAccessToken))
		})
	})

	suite.Run("Negative cases", func() {

	})
}

func (suite *GrantTokenFromAuthorizationCodeTest) Subtest(testcase string, subtest func()) {
	suite.SetupTest()
	suite.AuthorizationBaseSuiteTest.Subtest(testcase, subtest)
	suite.TearDownTest()
}
