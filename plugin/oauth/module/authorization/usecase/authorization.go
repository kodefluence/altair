package usecase

import (
	"github.com/kodefluence/altair/plugin/oauth/entity"
	"github.com/kodefluence/monorepo/db"
	"github.com/kodefluence/monorepo/exception"
	"github.com/kodefluence/monorepo/kontext"
)

//go:generate mockgen -destination ./mock/mock.go -package mock -source ./authorization.go

type Formatter interface {
	AccessTokenFromAuthorizationRequestInsertable(r entity.AuthorizationRequestJSON, application entity.OauthApplication) entity.OauthAccessTokenInsertable
	AccessTokenFromOauthAccessGrantInsertable(oauthAccessGrant entity.OauthAccessGrant, application entity.OauthApplication) entity.OauthAccessTokenInsertable
	AccessGrantFromAuthorizationRequestInsertable(r entity.AuthorizationRequestJSON, application entity.OauthApplication) entity.OauthAccessGrantInsertable
	OauthApplicationInsertable(r entity.OauthApplicationJSON) entity.OauthApplicationInsertable
	AccessTokenFromOauthRefreshTokenInsertable(application entity.OauthApplication, accessToken entity.OauthAccessToken) entity.OauthAccessTokenInsertable
	RefreshTokenInsertable(application entity.OauthApplication, accessToken entity.OauthAccessToken) entity.OauthRefreshTokenInsertable
	AccessGrant(e entity.OauthAccessGrant) entity.OauthAccessGrantJSON
	AccessToken(e entity.OauthAccessToken, redirectURI string, refreshTokenJSON *entity.OauthRefreshTokenJSON) entity.OauthAccessTokenJSON
	RefreshToken(e entity.OauthRefreshToken) entity.OauthRefreshTokenJSON
}

type OauthAccessGrantRepository interface {
	One(ktx kontext.Context, ID int, tx db.TX) (entity.OauthAccessGrant, exception.Exception)
	OneByCode(ktx kontext.Context, code string, tx db.TX) (entity.OauthAccessGrant, exception.Exception)
	Create(ktx kontext.Context, data entity.OauthAccessGrantInsertable, tx db.TX) (int, exception.Exception)
	Revoke(ktx kontext.Context, code string, tx db.TX) exception.Exception
}

type OauthAccessTokenRepository interface {
	OneByToken(ktx kontext.Context, token string, tx db.TX) (entity.OauthAccessToken, exception.Exception)
	One(ktx kontext.Context, ID int, tx db.TX) (entity.OauthAccessToken, exception.Exception)
	Create(ktx kontext.Context, data entity.OauthAccessTokenInsertable, tx db.TX) (int, exception.Exception)
	Revoke(ktx kontext.Context, token string, tx db.TX) exception.Exception
}

type OauthApplicationRepository interface {
	Paginate(ktx kontext.Context, offset, limit int, tx db.TX) ([]entity.OauthApplication, exception.Exception)
	Count(ktx kontext.Context, tx db.TX) (int, exception.Exception)
	One(ktx kontext.Context, ID int, tx db.TX) (entity.OauthApplication, exception.Exception)
	OneByUIDandSecret(ktx kontext.Context, clientUID, clientSecret string, tx db.TX) (entity.OauthApplication, exception.Exception)
	Create(ktx kontext.Context, data entity.OauthApplicationInsertable, tx db.TX) (int, exception.Exception)
	Update(ktx kontext.Context, ID int, data entity.OauthApplicationUpdateable, tx db.TX) exception.Exception
}

type OauthRefreshTokenRepository interface {
	OneByToken(ktx kontext.Context, token string, tx db.TX) (entity.OauthRefreshToken, exception.Exception)
	One(ktx kontext.Context, ID int, tx db.TX) (entity.OauthRefreshToken, exception.Exception)
	Create(ktx kontext.Context, data entity.OauthRefreshTokenInsertable, tx db.TX) (int, exception.Exception)
	Revoke(ktx kontext.Context, token string, tx db.TX) exception.Exception
}

// Authorization struct handle all of things related to oauth2 authorization
type Authorization struct {
	oauthApplicationModel  OauthApplicationRepository
	oauthAccessTokenModel  OauthAccessTokenRepository
	oauthAccessGrantModel  OauthAccessGrantRepository
	oauthRefreshTokenModel OauthRefreshTokenRepository

	formatter Formatter

	refreshTokenToggle bool

	sqldb db.DB
}

func NewAuthorization(
	oauthApplicationModel OauthApplicationRepository,
	oauthAccessTokenModel OauthAccessTokenRepository,
	oauthAccessGrantModel OauthAccessGrantRepository,
	oauthRefreshTokenModel OauthRefreshTokenRepository,
	formatter Formatter,
	refreshTokenToggle bool,
	sqldb db.DB,
) *Authorization {
	return &Authorization{
		oauthApplicationModel:  oauthApplicationModel,
		oauthAccessTokenModel:  oauthAccessTokenModel,
		oauthAccessGrantModel:  oauthAccessGrantModel,
		oauthRefreshTokenModel: oauthRefreshTokenModel,
		formatter:              formatter,
		refreshTokenToggle:     refreshTokenToggle,
		sqldb:                  sqldb,
	}
}
