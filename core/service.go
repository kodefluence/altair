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
	Grantor(ctx context.Context, authorizationReq entity.AuthorizationRequestJSON) (interface{}, *entity.Error)
	Grant(ctx context.Context, authorizationReq entity.AuthorizationRequestJSON) (entity.OauthAccessGrantJSON, *entity.Error)
	GrantToken(ctx context.Context, authorizationReq entity.AuthorizationRequestJSON) (entity.OauthAccessTokenJSON, *entity.Error)
	Token(ctx context.Context, accessTokenReq entity.AccessTokenRequestJSON) (entity.OauthAccessTokenJSON, *entity.Error)
}
