package service_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/codefluence-x/altair/entity"
	"github.com/codefluence-x/altair/formatter"
	"github.com/codefluence-x/altair/mock"
	"github.com/codefluence-x/altair/service"
	"github.com/codefluence-x/altair/util"
	"github.com/codefluence-x/aurelia"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAuthorization(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	t.Run("Grantor", func(t *testing.T) {
		t.Run("Given context and authorization request with a response type of token", func(t *testing.T) {
			t.Run("Return entity.OauthAccessTokenJSON and nil", func(t *testing.T) {
				oauthApplicationModel := mock.NewMockOauthApplicationModel(mockCtrl)
				oauthAccessTokenModel := mock.NewMockOauthAccessTokenModel(mockCtrl)
				oauthValidator := mock.NewMockOauthValidator(mockCtrl)
				modelFormatter := formatter.Model(time.Hour * 4)
				modelFormatterMock := mock.NewMockModelFormater(mockCtrl)
				oauthFormatter := formatter.Oauth()

				ctx := context.WithValue(context.Background(), "track_id", uuid.New().String())

				authorizationRequest := entity.AuthorizationRequestJSON{
					ResponseType:    util.StringToPointer("token"),
					ResourceOwnerID: util.IntToPointer(1),
					ClientUID:       util.StringToPointer(aurelia.Hash("x", "y")),
					ClientSecret:    util.StringToPointer(aurelia.Hash("z", "a")),
					RedirectURI:     util.StringToPointer("http://github.com"),
					Scopes:          util.StringToPointer("public users"),
				}

				oauthApplication := entity.OauthApplication{
					ID: 1,
					OwnerID: sql.NullInt64{
						Int64: 1,
						Valid: true,
					},
					OwnerType:    "confidential",
					Description:  "Blablabla",
					Scopes:       "public users stores",
					ClientUID:    *authorizationRequest.ClientUID,
					ClientSecret: *authorizationRequest.ClientSecret,
					CreatedAt:    time.Now().Add(-time.Hour * 4),
					UpdatedAt:    time.Now(),
				}

				oauthAccessToken := entity.OauthAccessToken{
					ID:                 1,
					OauthApplicationID: oauthApplication.ID,
					ResourceOwnerID:    *authorizationRequest.ResourceOwnerID,
					Token:              aurelia.Hash("x", "y"),
					Scopes:             *authorizationRequest.Scopes,
					ExpiresIn:          time.Now().Add(time.Hour * 4),
					CreatedAt:          time.Now(),
				}

				oauthAccessTokenInsertable := modelFormatter.AccessTokenFromAuthorizationRequest(authorizationRequest, oauthApplication)

				gomock.InOrder(
					oauthApplicationModel.EXPECT().
						OneByUIDandSecret(ctx, *authorizationRequest.ClientUID, *authorizationRequest.ClientSecret).
						Return(oauthApplication, nil),
					oauthValidator.EXPECT().ValidateAuthorizationGrant(ctx, authorizationRequest, oauthApplication).Return(nil),
					modelFormatterMock.EXPECT().AccessTokenFromAuthorizationRequest(authorizationRequest, oauthApplication).Return(oauthAccessTokenInsertable),
					oauthAccessTokenModel.EXPECT().Create(ctx, oauthAccessTokenInsertable).Return(oauthAccessToken.ID, nil),
					oauthAccessTokenModel.EXPECT().One(ctx, oauthAccessToken.ID).Return(oauthAccessToken, nil),
				)

				authorizationService := service.Authorization(oauthApplicationModel, oauthAccessTokenModel, modelFormatterMock, oauthValidator, oauthFormatter)
				_, err := authorizationService.Grantor(ctx, authorizationRequest)
				assert.Nil(t, err)
			})
		})

		t.Run("Given context and authorization request with a response type of code", func(t *testing.T) {

		})
	})
}
