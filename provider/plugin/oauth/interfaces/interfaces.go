package interfaces

//go:generate mockgen -destination ./../mock/mock_interfaces.go -package mock -source ./interfaces.go

import (
	"context"

	"github.com/codefluence-x/altair/provider/plugin/oauth/entity"
	"github.com/codefluence-x/monorepo/db"
	"github.com/codefluence-x/monorepo/exception"
	"github.com/codefluence-x/monorepo/kontext"
)

// OauthApplicationModel handle all database connection to oauth_applications table
type OauthApplicationModel interface {
	Paginate(ktx kontext.Context, offset, limit int, tx db.TX) ([]entity.OauthApplication, exception.Exception)
	Count(ktx kontext.Context, tx db.TX) (int, exception.Exception)
	One(ktx kontext.Context, ID int, tx db.TX) (entity.OauthApplication, exception.Exception)
	OneByUIDandSecret(ktx kontext.Context, clientUID, clientSecret string, tx db.TX) (entity.OauthApplication, exception.Exception)
	Create(ktx kontext.Context, data entity.OauthApplicationInsertable, tx db.TX) (int, exception.Exception)
	Update(ktx kontext.Context, ID int, data entity.OauthApplicationUpdateable, tx db.TX) exception.Exception
}

// OauthAccessTokenModel handle all database connection to oauth_access_tokens table
type OauthAccessTokenModel interface {
	OneByToken(ktx kontext.Context, token string, tx db.TX) (entity.OauthAccessToken, exception.Exception)
	One(ktx kontext.Context, ID int, tx db.TX) (entity.OauthAccessToken, exception.Exception)
	Create(ktx kontext.Context, data entity.OauthAccessTokenInsertable, tx db.TX) (int, exception.Exception)
	Revoke(ktx kontext.Context, token string, tx db.TX) exception.Exception
}

// OauthAccessGrantModel handle all database connection to oauth_access_grants table
type OauthAccessGrantModel interface {
	One(ktx kontext.Context, ID int, tx db.TX) (entity.OauthAccessGrant, exception.Exception)
	OneByCode(ktx kontext.Context, code string, tx db.TX) (entity.OauthAccessGrant, exception.Exception)
	Create(ktx kontext.Context, data entity.OauthAccessGrantInsertable, tx db.TX) (int, exception.Exception)
	Revoke(ktx kontext.Context, code string, tx db.TX) exception.Exception
}

// OauthRefreshTokenModel handle all database connection to oauth_refresh_tokens table
type OauthRefreshTokenModel interface {
	OneByToken(ktx kontext.Context, token string, tx db.TX) (entity.OauthRefreshToken, exception.Exception)
	One(ktx kontext.Context, ID int, tx db.TX) (entity.OauthRefreshToken, exception.Exception)
	Create(ktx kontext.Context, data entity.OauthRefreshTokenInsertable, tx db.TX) (int, exception.Exception)
	Revoke(ktx kontext.Context, token string, tx db.TX) exception.Exception
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
	AccessToken(e entity.OauthAccessToken, redirectURI string, refreshTokenJSON *entity.OauthRefreshTokenJSON) entity.OauthAccessTokenJSON
	RefreshToken(e entity.OauthRefreshToken) entity.OauthRefreshTokenJSON
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
	ValidateTokenRefreshToken(ctx context.Context, data entity.OauthRefreshToken) exception.Exception
	ValidateTokenAuthorizationCode(ctx context.Context, r entity.AccessTokenRequestJSON, data entity.OauthAccessGrant) *entity.Error
}
