package usecase

import (
	"time"

	"github.com/kodefluence/altair/plugin/oauth/entity"
	"github.com/kodefluence/altair/util"
)

func (*Formatter) AccessGrant(e entity.OauthAccessGrant) entity.OauthAccessGrantJSON {
	var data entity.OauthAccessGrantJSON

	data.ID = &e.ID
	data.OauthApplicationID = &e.OauthApplicationID
	data.ResourceOwnerID = &e.ResourceOwnerID
	data.Code = &e.Code
	data.RedirectURI = &e.RedirectURI.String
	data.Scopes = &e.Scopes.String

	if time.Now().Before(e.ExpiresIn) {
		data.ExpiresIn = util.IntToPointer(int(time.Until(e.ExpiresIn).Seconds()))
	} else {
		data.ExpiresIn = util.IntToPointer(0)
	}

	data.CreatedAt = &e.CreatedAt

	if e.RevokedAT.Valid {
		data.RevokedAT = &e.RevokedAT.Time
	}

	return data
}
