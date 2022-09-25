package usecase

import (
	"strconv"
	"time"

	"github.com/kodefluence/altair/plugin/oauth/entity"
	"github.com/kodefluence/aurelia"
)

func (f *Formatter) AccessTokenFromOauthAccessGrantInsertable(oauthAccessGrant entity.OauthAccessGrant, application entity.OauthApplication) entity.OauthAccessTokenInsertable {
	var accessTokenInsertable entity.OauthAccessTokenInsertable

	accessTokenInsertable.OauthApplicationID = application.ID
	accessTokenInsertable.ResourceOwnerID = oauthAccessGrant.ResourceOwnerID
	accessTokenInsertable.Token = aurelia.Hash(application.ClientUID, application.ClientSecret+strconv.Itoa(oauthAccessGrant.ResourceOwnerID))
	accessTokenInsertable.Scopes = oauthAccessGrant.Scopes.String
	accessTokenInsertable.ExpiresIn = time.Now().Add(f.tokenExpiresIn)

	return accessTokenInsertable
}
