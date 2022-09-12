package usecase

import (
	"time"

	"github.com/kodefluence/altair/plugin/oauth/entity"
	"github.com/kodefluence/altair/util"
)

func (*Formatter) AccessToken(e entity.OauthAccessToken, redirectURI string, refreshTokenJSON *entity.OauthRefreshTokenJSON) entity.OauthAccessTokenJSON {
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
