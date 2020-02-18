package formatter

import (
	"time"

	"github.com/codefluence-x/altair/core"
	"github.com/codefluence-x/altair/entity"
	"github.com/codefluence-x/altair/util"
)

type oauthFormatter struct{}

func Oauth() core.OauthFormatter {
	return &oauthFormatter{}
}

func (*oauthFormatter) AccessGrant(r entity.AuthorizationRequestJSON, e entity.OauthAccessGrant) entity.OauthAccessGrantJSON {
	var data entity.OauthAccessGrantJSON

	return data
}

func (*oauthFormatter) AccessToken(r entity.AuthorizationRequestJSON, e entity.OauthAccessToken) entity.OauthAccessTokenJSON {
	var data entity.OauthAccessTokenJSON

	data.ID = &e.ID
	data.OauthApplicationID = &e.OauthApplicationID
	data.ResourceOwnerID = &e.ResourceOwnerID
	data.Token = &e.Token
	data.Scopes = r.Scopes
	data.RedirectURI = r.RedirectURI
	data.CreatedAt = &e.CreatedAt

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
