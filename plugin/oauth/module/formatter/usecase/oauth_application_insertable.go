package usecase

import (
	"github.com/google/uuid"
	"github.com/kodefluence/altair/plugin/oauth/entity"
	"github.com/kodefluence/altair/util"
	"github.com/kodefluence/aurelia"
)

func (*Formatter) OauthApplicationInsertable(r entity.OauthApplicationJSON) entity.OauthApplicationInsertable {
	var oauthApplicationInsertable entity.OauthApplicationInsertable

	oauthApplicationInsertable.OwnerID = util.PointerToValue(r.OwnerID)
	oauthApplicationInsertable.OwnerType = *r.OwnerType
	oauthApplicationInsertable.Description = util.PointerToValue(r.Description)
	oauthApplicationInsertable.Scopes = util.PointerToValue(r.Scopes)
	oauthApplicationInsertable.ClientUID = util.SHA1()
	oauthApplicationInsertable.ClientSecret = aurelia.Hash(oauthApplicationInsertable.ClientUID, uuid.New().String())

	return oauthApplicationInsertable
}
