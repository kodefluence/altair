package interfaces

//go:generate mockgen -destination ./../mock/mock_interfaces.go -package mock -source ./interfaces.go

import (
	"context"
	"database/sql"

	"github.com/codefluence-x/altair/core"
	"github.com/codefluence-x/altair/provider/plugin/oauth/entity"
)

type OauthApplicationModel interface {
	Name() string
	Paginate(ctx context.Context, offset, limit int) ([]entity.OauthApplication, error)
	One(ctx context.Context, ID int) (entity.OauthApplication, error)
	OneByUIDandSecret(ctx context.Context, clientUID, clientSecret string) (entity.OauthApplication, error)
	Count(ctx context.Context) (int, error)
	Create(ctx context.Context, data entity.OauthApplicationInsertable, txs ...*sql.Tx) (int, error)
}

type OauthAccessTokenModel interface {
	Name() string
	One(ctx context.Context, ID int) (entity.OauthAccessToken, error)
	OneByToken(ctx context.Context, token string) (entity.OauthAccessToken, error)
	Create(ctx context.Context, data entity.OauthAccessTokenInsertable, txs ...*sql.Tx) (int, error)
	Revoke(ctx context.Context, token string) error
}

type OauthAccessGrantModel interface {
	Name() string
	One(ctx context.Context, ID int) (entity.OauthAccessGrant, error)
	Create(ctx context.Context, data entity.OauthAccessGrantInsertable, txs ...*sql.Tx) (int, error)
	OneByCode(ctx context.Context, code string) (entity.OauthAccessGrant, error)
}

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
	RevokeToken(ctx context.Context, revokeAccessTokenReq entity.RevokeAccessTokenRequestJSON) *entity.Error
}

type OauthApplicationFormater interface {
	ApplicationList(ctx context.Context, applications []entity.OauthApplication) []entity.OauthApplicationJSON
	Application(ctx context.Context, application entity.OauthApplication) entity.OauthApplicationJSON
}

type OauthFormatter interface {
	AccessGrant(e entity.OauthAccessGrant) entity.OauthAccessGrantJSON
	AccessToken(e entity.OauthAccessToken, redirectURI string) entity.OauthAccessTokenJSON
}

type ModelFormater interface {
	AccessTokenFromAuthorizationRequest(r entity.AuthorizationRequestJSON, application entity.OauthApplication) entity.OauthAccessTokenInsertable
	AccessTokenFromOauthAccessGrant(oauthAccessGrant entity.OauthAccessGrant, application entity.OauthApplication) entity.OauthAccessTokenInsertable
	AccessGrantFromAuthorizationRequest(r entity.AuthorizationRequestJSON, application entity.OauthApplication) entity.OauthAccessGrantInsertable
	OauthApplication(r entity.OauthApplicationJSON) entity.OauthApplicationInsertable
}

type OauthValidator interface {
	ValidateApplication(ctx context.Context, data entity.OauthApplicationJSON) *entity.Error
	ValidateAuthorizationGrant(ctx context.Context, r entity.AuthorizationRequestJSON, application entity.OauthApplication) *entity.Error
	ValidateTokenGrant(ctx context.Context, r entity.AccessTokenRequestJSON) *entity.Error
}

type OauthDispatcher interface {
	Application() OauthApplicationDispatcher
	Authorization() AuthorizationDispatcher
}

type OauthApplicationDispatcher interface {
	List(applicationManager ApplicationManager) core.Controller
	One(applicationManager ApplicationManager) core.Controller
	Create(applicationManager ApplicationManager) core.Controller
}

type AuthorizationDispatcher interface {
	Grant(authorization Authorization) core.Controller
	Revoke(authorization Authorization) core.Controller
}
