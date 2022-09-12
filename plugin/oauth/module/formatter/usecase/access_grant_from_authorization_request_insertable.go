package usecase

import (
	"time"

	"github.com/kodefluence/altair/plugin/oauth/entity"
	"github.com/kodefluence/altair/util"
)

func (f *Formatter) AccessGrantFromAuthorizationRequestInsertable(r entity.AuthorizationRequestJSON, application entity.OauthApplication) entity.OauthAccessGrantInsertable {
	var accessGrantInsertable entity.OauthAccessGrantInsertable

	accessGrantInsertable.OauthApplicationID = application.ID
	accessGrantInsertable.ResourceOwnerID = *r.ResourceOwnerID
	accessGrantInsertable.Scopes = util.PointerToString(r.Scopes)
	accessGrantInsertable.Code = util.SHA1()
	accessGrantInsertable.RedirectURI = util.PointerToString(r.RedirectURI)
	accessGrantInsertable.ExpiresIn = time.Now().Add(f.codeExpiresIn)

	return accessGrantInsertable
}
