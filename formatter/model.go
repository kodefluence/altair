package formatter

import (
	"strconv"
	"time"

	"github.com/codefluence-x/altair/core"
	"github.com/codefluence-x/altair/entity"
	"github.com/codefluence-x/altair/util"
	"github.com/codefluence-x/aurelia"
)

type model struct {
	tokenExpiresIn time.Duration
	codeExpiresIn  time.Duration
}

func Model(tokenExpiresIn time.Duration, codeExpiresIn time.Duration) core.ModelFormater {
	return &model{
		tokenExpiresIn: tokenExpiresIn,
		codeExpiresIn:  codeExpiresIn,
	}
}

func (m *model) AccessGrantFromAuthorizationRequest(r entity.AuthorizationRequestJSON, application entity.OauthApplication) entity.OauthAccessGrantInsertable {
	var accessGrantInsertable entity.OauthAccessGrantInsertable

	accessGrantInsertable.OauthApplicationID = application.ID
	accessGrantInsertable.ResourceOwnerID = *r.ResourceOwnerID
	accessGrantInsertable.Scopes = *r.Scopes
	accessGrantInsertable.Code = util.SHA1()
	accessGrantInsertable.RedirectURI = *r.RedirectURI
	accessGrantInsertable.ExpiresIn = time.Now().Add(m.codeExpiresIn)

	return accessGrantInsertable
}

func (m *model) AccessTokenFromAuthorizationRequest(r entity.AuthorizationRequestJSON, application entity.OauthApplication) entity.OauthAccessTokenInsertable {
	var accessTokenInsertable entity.OauthAccessTokenInsertable

	accessTokenInsertable.OauthApplicationID = application.ID
	accessTokenInsertable.ResourceOwnerID = *r.ResourceOwnerID
	accessTokenInsertable.Token = aurelia.Hash(application.ClientUID, application.ClientSecret+strconv.Itoa(*r.ResourceOwnerID))
	accessTokenInsertable.Scopes = *r.Scopes
	accessTokenInsertable.ExpiresIn = time.Now().Add(m.tokenExpiresIn)

	return accessTokenInsertable
}
