package usecase

import (
	"strconv"
	"time"

	"github.com/kodefluence/altair/plugin/oauth/entity"
	"github.com/kodefluence/altair/util"
	"github.com/kodefluence/aurelia"
)

func (f *Formatter) AccessTokenFromAuthorizationRequestInsertable(r entity.AuthorizationRequestJSON, application entity.OauthApplication) entity.OauthAccessTokenInsertable {
	var accessTokenInsertable entity.OauthAccessTokenInsertable

	accessTokenInsertable.OauthApplicationID = application.ID
	accessTokenInsertable.ResourceOwnerID = *r.ResourceOwnerID
	accessTokenInsertable.Token = aurelia.Hash(application.ClientUID, application.ClientSecret+strconv.Itoa(*r.ResourceOwnerID))
	accessTokenInsertable.Scopes = util.PointerToString(r.Scopes)
	accessTokenInsertable.ExpiresIn = time.Now().Add(f.tokenExpiresIn)

	return accessTokenInsertable
}
