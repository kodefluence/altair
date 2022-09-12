package usecase

import (
	"strconv"
	"time"

	"github.com/kodefluence/altair/plugin/oauth/entity"
	"github.com/kodefluence/aurelia"
)

func (f *Formatter) AccessTokenFromOauthRefreshTokenInsertable(application entity.OauthApplication, accessToken entity.OauthAccessToken) entity.OauthAccessTokenInsertable {
	var accessTokenInsertable entity.OauthAccessTokenInsertable

	accessTokenInsertable.OauthApplicationID = application.ID
	accessTokenInsertable.ResourceOwnerID = accessToken.ResourceOwnerID
	accessTokenInsertable.Token = aurelia.Hash(application.ClientUID, application.ClientSecret+strconv.Itoa(accessToken.ResourceOwnerID))
	accessTokenInsertable.Scopes = accessToken.Scopes.String
	accessTokenInsertable.ExpiresIn = time.Now().Add(f.tokenExpiresIn)

	return accessTokenInsertable
}
