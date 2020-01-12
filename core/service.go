package core

import (
	"context"

	"github.com/codefluence-x/altair/entity"
)

type ApplicationManager interface {
	List(ctx context.Context, offset, limit int) ([]entity.OauthApplicationJSON, int, *entity.Error)
	One(ctx context.Context, ID int) (entity.OauthApplicationJSON, *entity.Error)
	Create(ctx context.Context, e entity.OauthApplicationJSON) (entity.OauthApplicationJSON, *entity.Error)
}

type Authorization interface {
	Grant(ctx context.Context, authorizationReq entity.AuthorizationRequestJSON) (error, entity.OauthAccessGrantJSON)
	Token(ctx context.Context, accessTokenReq entity.AccessTokenRequestJSON) (error, entity.OauthAccessTokenJSON)
}
