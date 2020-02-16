package formatter

import (
	"strconv"
	"time"

	"github.com/codefluence-x/altair/core"
	"github.com/codefluence-x/altair/entity"
	"github.com/codefluence-x/aurelia"
)

type model struct {
	expiresIn time.Duration
}

func Model(expiresIn time.Duration) core.ModelFormater {
	return &model{
		expiresIn: expiresIn,
	}
}

func (m *model) AccessTokenFromAuthorizationRequest(r entity.AuthorizationRequestJSON, application entity.OauthApplication) entity.OauthAccessTokenInsertable {
	var accessTokenInsertable entity.OauthAccessTokenInsertable

	accessTokenInsertable.OauthApplicationID = application.ID
	accessTokenInsertable.ResourceOwnerID = *r.ResourceOwnerID
	accessTokenInsertable.Token = aurelia.Hash(application.ClientUID, application.ClientSecret+strconv.Itoa(*r.ResourceOwnerID))
	accessTokenInsertable.Scopes = *r.Scopes
	accessTokenInsertable.ExpiresIn = time.Now().Add(m.expiresIn)

	return accessTokenInsertable
}
