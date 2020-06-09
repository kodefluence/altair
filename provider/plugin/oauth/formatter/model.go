package formatter

import (
	"strconv"
	"time"

	"github.com/codefluence-x/altair/entity"
	"github.com/codefluence-x/altair/provider/plugin/oauth"
	"github.com/codefluence-x/altair/util"
	"github.com/codefluence-x/aurelia"
	"github.com/google/uuid"
)

type model struct {
	tokenExpiresIn time.Duration
	codeExpiresIn  time.Duration
}

func Model(tokenExpiresIn time.Duration, codeExpiresIn time.Duration) oauth.ModelFormater {
	return &model{
		tokenExpiresIn: tokenExpiresIn,
		codeExpiresIn:  codeExpiresIn,
	}
}

func (m *model) OauthApplication(r entity.OauthApplicationJSON) entity.OauthApplicationInsertable {
	var oauthApplicationInsertable entity.OauthApplicationInsertable

	oauthApplicationInsertable.OwnerID = util.PointerToInt(r.OwnerID)
	oauthApplicationInsertable.OwnerType = *r.OwnerType
	oauthApplicationInsertable.Description = util.PointerToString(r.Description)
	oauthApplicationInsertable.Scopes = util.PointerToString(r.Scopes)
	oauthApplicationInsertable.ClientUID = util.SHA1()
	oauthApplicationInsertable.ClientSecret = aurelia.Hash(oauthApplicationInsertable.ClientUID, uuid.New().String())

	return oauthApplicationInsertable
}

func (m *model) AccessGrantFromAuthorizationRequest(r entity.AuthorizationRequestJSON, application entity.OauthApplication) entity.OauthAccessGrantInsertable {
	var accessGrantInsertable entity.OauthAccessGrantInsertable

	accessGrantInsertable.OauthApplicationID = application.ID
	accessGrantInsertable.ResourceOwnerID = *r.ResourceOwnerID
	accessGrantInsertable.Scopes = util.PointerToString(r.Scopes)
	accessGrantInsertable.Code = util.SHA1()
	accessGrantInsertable.RedirectURI = util.PointerToString(r.RedirectURI)
	accessGrantInsertable.ExpiresIn = time.Now().Add(m.codeExpiresIn)

	return accessGrantInsertable
}

func (m *model) AccessTokenFromAuthorizationRequest(r entity.AuthorizationRequestJSON, application entity.OauthApplication) entity.OauthAccessTokenInsertable {
	var accessTokenInsertable entity.OauthAccessTokenInsertable

	accessTokenInsertable.OauthApplicationID = application.ID
	accessTokenInsertable.ResourceOwnerID = *r.ResourceOwnerID
	accessTokenInsertable.Token = aurelia.Hash(application.ClientUID, application.ClientSecret+strconv.Itoa(*r.ResourceOwnerID))
	accessTokenInsertable.Scopes = util.PointerToString(r.Scopes)
	accessTokenInsertable.ExpiresIn = time.Now().Add(m.tokenExpiresIn)

	return accessTokenInsertable
}
