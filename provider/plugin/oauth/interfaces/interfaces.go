package interfaces

//go:generate mockgen -destination ./../mock/mock_interfaces.go -package mock -source ./interfaces.go

import (
	"context"
	"database/sql"

	"github.com/codefluence-x/altair/provider/plugin/oauth/entity"
)

// OauthApplicationModel handle all database connection to oauth_applications table
type OauthApplicationModel interface {
	Name() string
	Paginate(ctx context.Context, offset, limit int) ([]entity.OauthApplication, error)
	One(ctx context.Context, ID int) (entity.OauthApplication, error)
	OneByUIDandSecret(ctx context.Context, clientUID, clientSecret string) (entity.OauthApplication, error)
	Count(ctx context.Context) (int, error)
	Create(ctx context.Context, data entity.OauthApplicationInsertable, txs ...*sql.Tx) (int, error)
	Update(ctx context.Context, ID int, data entity.OauthApplicationUpdateable, txs ...*sql.Tx) error
}

// OauthAccessTokenModel handle all database connection to oauth_access_tokens table
type OauthAccessTokenModel interface {
	Name() string
	One(ctx context.Context, ID int) (entity.OauthAccessToken, error)
	OneByToken(ctx context.Context, token string) (entity.OauthAccessToken, error)
	Create(ctx context.Context, data entity.OauthAccessTokenInsertable, txs ...*sql.Tx) (int, error)
	Revoke(ctx context.Context, token string) error
}

// OauthAccessGrantModel handle all database connection to oauth_access_grants table
type OauthAccessGrantModel interface {
	Name() string
	One(ctx context.Context, ID int) (entity.OauthAccessGrant, error)
	Create(ctx context.Context, data entity.OauthAccessGrantInsertable, txs ...*sql.Tx) (int, error)
	OneByCode(ctx context.Context, code string) (entity.OauthAccessGrant, error)
	Revoke(ctx context.Context, code string, txs ...*sql.Tx) error
}

// OauthRefreshTokenModel handle all database connection to oauth_refresh_tokens table
type OauthRefreshTokenModel interface {
	Name() string
	One(ctx context.Context, ID int) (entity.OauthRefreshToken, error)
	Create(ctx context.Context, data entity.OauthRefreshTokenInsertable, txs ...*sql.Tx) (int, error)
	OneByToken(ctx context.Context, token string) (entity.OauthRefreshToken, error)
	Revoke(ctx context.Context, token string) error
}

// ApplicationManager manage all flow related oauth application CRUD
type ApplicationManager interface {
	List(ctx context.Context, offset, limit int) ([]entity.OauthApplicationJSON, int, *entity.Error)
	One(ctx context.Context, ID int) (entity.OauthApplicationJSON, *entity.Error)
	Create(ctx context.Context, e entity.OauthApplicationJSON) (entity.OauthApplicationJSON, *entity.Error)
	Update(ctx context.Context, ID int, e entity.OauthApplicationUpdateJSON) (entity.OauthApplicationJSON, *entity.Error)
}

// Authorization handle all service related to access token authorization
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

// ModelFormater format compiled entity into insertable
type ModelFormater interface {
	AccessTokenFromAuthorizationRequest(r entity.AuthorizationRequestJSON, application entity.OauthApplication) entity.OauthAccessTokenInsertable
	AccessTokenFromOauthAccessGrant(oauthAccessGrant entity.OauthAccessGrant, application entity.OauthApplication) entity.OauthAccessTokenInsertable
	AccessGrantFromAuthorizationRequest(r entity.AuthorizationRequestJSON, application entity.OauthApplication) entity.OauthAccessGrantInsertable
	OauthApplication(r entity.OauthApplicationJSON) entity.OauthApplicationInsertable
	AccessTokenFromOauthRefreshToken(application entity.OauthApplication, accessToken entity.OauthAccessToken) entity.OauthAccessTokenInsertable
	RefreshToken(application entity.OauthApplication, accessToken entity.OauthAccessToken) entity.OauthRefreshTokenInsertable
}

type OauthValidator interface {
	ValidateApplication(ctx context.Context, data entity.OauthApplicationJSON) *entity.Error
	ValidateAuthorizationGrant(ctx context.Context, r entity.AuthorizationRequestJSON, application entity.OauthApplication) *entity.Error
	ValidateTokenGrant(ctx context.Context, r entity.AccessTokenRequestJSON) *entity.Error
	ValidateTokenRefreshToken(ctx context.Context, data entity.OauthRefreshToken) *entity.Error
	ValidateTokenAuthorizationCode(ctx context.Context, r entity.AccessTokenRequestJSON, data entity.OauthAccessGrant) *entity.Error
}
