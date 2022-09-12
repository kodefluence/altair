package usecase

import (
	"strconv"
	"time"

	"github.com/kodefluence/altair/plugin/oauth/entity"
	"github.com/kodefluence/aurelia"
)

func (f *Formatter) RefreshTokenInsertable(application entity.OauthApplication, accessToken entity.OauthAccessToken) entity.OauthRefreshTokenInsertable {
	var refreshTokenInsertable entity.OauthRefreshTokenInsertable

	refreshTokenInsertable.Token = aurelia.Hash(application.ClientUID, application.ClientSecret+strconv.Itoa(accessToken.ResourceOwnerID))
	refreshTokenInsertable.OauthAccessTokenID = accessToken.ID
	refreshTokenInsertable.ExpiresIn = time.Now().Add(f.refreshTokenExpiresIn)

	return refreshTokenInsertable
}
