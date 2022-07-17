package formatter

import (
	"time"

	"github.com/kodefluence/altair/provider/plugin/oauth/entity"
	"github.com/kodefluence/altair/util"
)

type OauthFormatter struct{}

func Oauth() *OauthFormatter {
	return &OauthFormatter{}
}

func (*OauthFormatter) AccessGrant(e entity.OauthAccessGrant) entity.OauthAccessGrantJSON {
	var data entity.OauthAccessGrantJSON

	data.ID = &e.ID
	data.OauthApplicationID = &e.OauthApplicationID
	data.ResourceOwnerID = &e.ResourceOwnerID
	data.Code = &e.Code
	data.RedirectURI = &e.RedirectURI.String
	data.Scopes = &e.Scopes.String

	if time.Now().Before(e.ExpiresIn) {
		data.ExpiresIn = util.IntToPointer(int(e.ExpiresIn.Sub(time.Now()).Seconds()))
	} else {
		data.ExpiresIn = util.IntToPointer(0)
	}

	data.CreatedAt = &e.CreatedAt

	if e.RevokedAT.Valid {
		data.RevokedAT = &e.RevokedAT.Time
	}

	return data
}

func (*OauthFormatter) AccessToken(e entity.OauthAccessToken, redirectURI string, refreshTokenJSON *entity.OauthRefreshTokenJSON) entity.OauthAccessTokenJSON {
	var data entity.OauthAccessTokenJSON

	data.ID = &e.ID
	data.OauthApplicationID = &e.OauthApplicationID
	data.ResourceOwnerID = &e.ResourceOwnerID
	data.Token = &e.Token
	data.Scopes = &e.Scopes.String
	data.RedirectURI = &redirectURI
	data.CreatedAt = &e.CreatedAt

	if time.Now().Before(e.ExpiresIn) {
		data.ExpiresIn = util.IntToPointer(int(e.ExpiresIn.Sub(time.Now()).Seconds()))
	} else {
		data.ExpiresIn = util.IntToPointer(0)
	}

	if e.RevokedAT.Valid {
		data.RevokedAT = &e.RevokedAT.Time
	}

	if refreshTokenJSON != nil {
		data.RefreshToken = refreshTokenJSON
	}

	return data
}

func (*OauthFormatter) RefreshToken(e entity.OauthRefreshToken) entity.OauthRefreshTokenJSON {
	var data entity.OauthRefreshTokenJSON

	data.CreatedAt = &e.CreatedAt
	data.Token = &e.Token

	if time.Now().Before(e.ExpiresIn) {
		data.ExpiresIn = util.IntToPointer(int(e.ExpiresIn.Sub(time.Now()).Seconds()))
	} else {
		data.ExpiresIn = util.IntToPointer(0)
	}

	if e.RevokedAT.Valid {
		data.RevokedAT = &e.RevokedAT.Time
	}

	return data
}
