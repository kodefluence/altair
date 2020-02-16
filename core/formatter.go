package core

import (
	"context"

	"github.com/codefluence-x/altair/entity"
)

type OauthApplicationFormater interface {
	ApplicationList(ctx context.Context, applications []entity.OauthApplication) []entity.OauthApplicationJSON
	Application(ctx context.Context, application entity.OauthApplication) entity.OauthApplicationJSON
}

type OauthFormatter interface {
	AccessGrant(r entity.AuthorizationRequestJSON, e entity.OauthAccessGrant) entity.OauthAccessGrantJSON
	AccessToken(r entity.AuthorizationRequestJSON, e entity.OauthAccessToken) entity.OauthAccessTokenJSON
}
