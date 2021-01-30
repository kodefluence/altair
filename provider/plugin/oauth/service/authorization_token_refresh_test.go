package service_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/codefluence-x/altair/provider/plugin/oauth/entity"
	"github.com/codefluence-x/altair/provider/plugin/oauth/formatter"
	"github.com/codefluence-x/altair/provider/plugin/oauth/mock"
	"github.com/codefluence-x/altair/provider/plugin/oauth/service"
	"github.com/codefluence-x/altair/util"
	"github.com/codefluence-x/aurelia"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestAuthorizationRefreshToken(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	t.Run("Token", func(t *testing.T) {
		t.Run("Given context and refresh token request", func(t *testing.T) {
			t.Run("When refresh token request valid and there is no error in database side", func(t *testing.T) {
				t.Run("Then it will return access token response", func(t *testing.T) {
					oauthApplicationModel := mock.NewMockOauthApplicationModel(mockCtrl)
					oauthAccessTokenModel := mock.NewMockOauthAccessTokenModel(mockCtrl)
					oauthAccessGrantModel := mock.NewMockOauthAccessGrantModel(mockCtrl)
					oauthRefreshTokenModel := mock.NewMockOauthRefreshTokenModel(mockCtrl)
					oauthValidator := mock.NewMockOauthValidator(mockCtrl)
					modelFormatter := formatter.NewModel(time.Hour*4, time.Hour*2, time.Hour*2)
					oauthFormatter := formatter.Oauth()

					ctx := context.Background()

					accessTokenRequest := entity.AccessTokenRequestJSON{
						ClientSecret: util.StringToPointer("client_secret"),
						ClientUID:    util.StringToPointer("client_uid"),
						RefreshToken: util.StringToPointer("abcdef_123456"),
						GrantType:    util.StringToPointer("refresh_token"),
						RedirectURI:  util.StringToPointer("http://localhost:8000/oauth_redirect"),
					}

					oauthApplication := entity.OauthApplication{
						ID: 1,
						OwnerID: sql.NullInt64{
							Int64: 1,
							Valid: true,
						},
						OwnerType: "confidential",
						Description: sql.NullString{
							String: "Application 01",
							Valid:  true,
						},
						Scopes: sql.NullString{
							String: "public users",
							Valid:  true,
						},
						ClientUID:    *accessTokenRequest.ClientUID,
						ClientSecret: *accessTokenRequest.ClientSecret,
						CreatedAt:    time.Now().Add(-time.Hour * 4),
						UpdatedAt:    time.Now(),
					}

					oldAccessToken := entity.OauthAccessToken{
						ID: 999,
					}

					oauthRefreshToken := entity.OauthRefreshToken{
						ID:                 1,
						OauthAccessTokenID: oldAccessToken.ID,
						Token:              *accessTokenRequest.RefreshToken,
					}

					oauthAccessToken := entity.OauthAccessToken{
						ID:                 1000,
						OauthApplicationID: oauthApplication.ID,
						ResourceOwnerID:    oldAccessToken.ResourceOwnerID,
						Token:              aurelia.Hash("x", "y"),
						Scopes: sql.NullString{
							String: oldAccessToken.Scopes.String,
							Valid:  true,
						},
						ExpiresIn: time.Now().Add(time.Hour * 4),
						CreatedAt: time.Now(),
					}

					oauthAccessTokenInsertable := modelFormatter.AccessTokenFromOauthRefreshToken(oauthApplication, oldAccessToken)
					oauthRefreshTokenInsertable := modelFormatter.RefreshToken(oauthApplication, oauthAccessToken)
					oauthAccessTokenJSON := oauthFormatter.AccessToken(oauthAccessToken, *accessTokenRequest.RedirectURI)

					gomock.InOrder(
						oauthApplicationModel.EXPECT().
							OneByUIDandSecret(ctx, *accessTokenRequest.ClientUID, *accessTokenRequest.ClientSecret).
							Return(oauthApplication, nil),
						oauthValidator.EXPECT().ValidateTokenGrant(ctx, accessTokenRequest).Return(nil),
						oauthRefreshTokenModel.EXPECT().OneByToken(ctx, *accessTokenRequest.RefreshToken).Return(oauthRefreshToken, nil),
						oauthValidator.EXPECT().ValidateTokenRefreshToken(ctx, oauthRefreshToken).Return(nil),
						oauthAccessTokenModel.EXPECT().One(ctx, oauthRefreshToken.OauthAccessTokenID).Return(oldAccessToken, nil),
						oauthAccessTokenModel.EXPECT().Create(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, insertable entity.OauthAccessTokenInsertable) (int, error) {
							assert.Equal(t, oauthAccessTokenInsertable.ResourceOwnerID, insertable.ResourceOwnerID)
							assert.Equal(t, oauthAccessTokenInsertable.OauthApplicationID, insertable.OauthApplicationID)
							return 1000, nil
						}),
						oauthAccessTokenModel.EXPECT().One(ctx, 1000).Return(oauthAccessToken, nil),
						oauthRefreshTokenModel.EXPECT().Revoke(ctx, oauthRefreshToken.Token).Return(nil),
						oauthRefreshTokenModel.EXPECT().Create(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, insertable entity.OauthRefreshTokenInsertable) (int, error) {
							assert.Equal(t, oauthRefreshTokenInsertable.OauthAccessTokenID, insertable.OauthAccessTokenID)
							return 2, nil
						}),
					)

					authorizationService := service.NewAuthorization(oauthApplicationModel, oauthAccessTokenModel, oauthAccessGrantModel, oauthRefreshTokenModel, modelFormatter, oauthValidator, oauthFormatter, true)
					oauthAccessTokenOutput, err := authorizationService.Token(ctx, accessTokenRequest)

					assert.Nil(t, err)
					assert.Equal(t, oauthAccessTokenJSON, oauthAccessTokenOutput)
				})
			})
		})
	})
}
