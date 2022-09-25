package usecase

import (
	"strconv"
	"time"

	"github.com/kodefluence/altair/plugin/oauth/entity"
	"github.com/kodefluence/aurelia"
)

func (f *Formatter) AccessTokenClientCredentialInsertable(application entity.OauthApplication, scope *string) entity.OauthAccessTokenInsertable {
	var accessTokenInsertable entity.OauthAccessTokenInsertable

	accessTokenInsertable.OauthApplicationID = application.ID
	accessTokenInsertable.ResourceOwnerID = 0
	accessTokenInsertable.Token = aurelia.Hash(application.ClientUID, application.ClientSecret+strconv.Itoa(0))

	if scope != nil {
		accessTokenInsertable.Scopes = *scope
	}
	accessTokenInsertable.ExpiresIn = time.Now().Add(f.tokenExpiresIn)

	return accessTokenInsertable
}
