package usecase_test

import (
	"time"

	"github.com/golang/mock/gomock"
	"github.com/kodefluence/altair/module"
	"github.com/kodefluence/altair/module/apierror"
	"github.com/kodefluence/altair/plugin/oauth/entity"
	"github.com/kodefluence/altair/plugin/oauth/module/authorization/usecase"
	"github.com/kodefluence/altair/plugin/oauth/module/authorization/usecase/mock"
	"github.com/kodefluence/altair/plugin/oauth/module/formatter"
	mockdb "github.com/kodefluence/monorepo/db/mock"
	"github.com/kodefluence/monorepo/kontext"
	"github.com/stretchr/testify/suite"
)

type AuthorizationBaseSuiteTest struct {
	mockCtrl *gomock.Controller

	ktx kontext.Context

	oauthApplicationRepo  *mock.MockOauthApplicationRepository
	oauthAccessTokenRepo  *mock.MockOauthAccessTokenRepository
	oauthAccessGrantRepo  *mock.MockOauthAccessGrantRepository
	oauthRefreshTokenRepo *mock.MockOauthRefreshTokenRepository
	formatter             usecase.Formatter
	config                entity.OauthPlugin
	apiError              module.ApiError
	authorization         usecase.Authorization
	sqldb                 *mockdb.MockDB

	suite.Suite
}

func (suite *AuthorizationBaseSuiteTest) SetupTest() {
	suite.mockCtrl = gomock.NewController(suite.T())

	suite.ktx = kontext.Fabricate()

	suite.config = entity.OauthPlugin{
		Config: entity.PluginConfig{
			Database:                    "main_database",
			AccessTokenTimeoutRaw:       "24h",
			AuthorizationCodeTimeoutRaw: "24h",
			RefreshToken: struct {
				Timeout string "yaml:\"timeout\""
				Active  bool   "yaml:\"active\""
			}{
				Timeout: "24h",
				Active:  true,
			},
			ImplicitGrant: struct {
				Active bool "yaml:\"active\""
			}{
				Active: true,
			},
		},
	}

	suite.oauthApplicationRepo = mock.NewMockOauthApplicationRepository(suite.mockCtrl)
	suite.oauthAccessTokenRepo = mock.NewMockOauthAccessTokenRepository(suite.mockCtrl)
	suite.oauthAccessGrantRepo = mock.NewMockOauthAccessGrantRepository(suite.mockCtrl)
	suite.oauthRefreshTokenRepo = mock.NewMockOauthRefreshTokenRepository(suite.mockCtrl)
	suite.formatter = formatter.Provide(24*time.Hour, 24*time.Hour, 24*time.Hour)
	suite.sqldb = mockdb.NewMockDB(suite.mockCtrl)
	suite.apiError = apierror.Provide()
	suite.authorization = *usecase.NewAuthorization(suite.oauthApplicationRepo, suite.oauthAccessTokenRepo, suite.oauthAccessGrantRepo, suite.oauthRefreshTokenRepo, suite.formatter, suite.config, suite.sqldb, suite.apiError)
}

func (suite *AuthorizationBaseSuiteTest) TearDownTest() {
	suite.mockCtrl.Finish()
}

func (suite *AuthorizationBaseSuiteTest) Subtest(testcase string, subtest func()) {
	suite.SetupTest()
	suite.Run("When all parameter is valid it would return nil", subtest)
	suite.TearDownTest()
}
