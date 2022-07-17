package formatter

import (
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/kodefluence/altair/provider/plugin/oauth/entity"
	"github.com/kodefluence/altair/util"
	"github.com/kodefluence/aurelia"
)

// Model implement converter from other struct into model insertable struct
type Model struct {
	tokenExpiresIn        time.Duration
	codeExpiresIn         time.Duration
	refreshTokenExpiresIn time.Duration
}

// NewModel create new model object
func NewModel(tokenExpiresIn time.Duration, codeExpiresIn time.Duration, refreshTokenExpiresIn time.Duration) *Model {
	return &Model{
		tokenExpiresIn:        tokenExpiresIn,
		codeExpiresIn:         codeExpiresIn,
		refreshTokenExpiresIn: refreshTokenExpiresIn,
	}
}

// OauthApplication format into OauthApplicationInsertable
func (m *Model) OauthApplication(r entity.OauthApplicationJSON) entity.OauthApplicationInsertable {
	var oauthApplicationInsertable entity.OauthApplicationInsertable

	oauthApplicationInsertable.OwnerID = util.PointerToInt(r.OwnerID)
	oauthApplicationInsertable.OwnerType = *r.OwnerType
	oauthApplicationInsertable.Description = util.PointerToString(r.Description)
	oauthApplicationInsertable.Scopes = util.PointerToString(r.Scopes)
	oauthApplicationInsertable.ClientUID = util.SHA1()
	oauthApplicationInsertable.ClientSecret = aurelia.Hash(oauthApplicationInsertable.ClientUID, uuid.New().String())

	return oauthApplicationInsertable
}

// AccessGrantFromAuthorizationRequest _
func (m *Model) AccessGrantFromAuthorizationRequest(r entity.AuthorizationRequestJSON, application entity.OauthApplication) entity.OauthAccessGrantInsertable {
	var accessGrantInsertable entity.OauthAccessGrantInsertable

	accessGrantInsertable.OauthApplicationID = application.ID
	accessGrantInsertable.ResourceOwnerID = *r.ResourceOwnerID
	accessGrantInsertable.Scopes = util.PointerToString(r.Scopes)
	accessGrantInsertable.Code = util.SHA1()
	accessGrantInsertable.RedirectURI = util.PointerToString(r.RedirectURI)
	accessGrantInsertable.ExpiresIn = time.Now().Add(m.codeExpiresIn)

	return accessGrantInsertable
}

// AccessTokenFromAuthorizationRequest _
func (m *Model) AccessTokenFromAuthorizationRequest(r entity.AuthorizationRequestJSON, application entity.OauthApplication) entity.OauthAccessTokenInsertable {
	var accessTokenInsertable entity.OauthAccessTokenInsertable

	accessTokenInsertable.OauthApplicationID = application.ID
	accessTokenInsertable.ResourceOwnerID = *r.ResourceOwnerID
	accessTokenInsertable.Token = aurelia.Hash(application.ClientUID, application.ClientSecret+strconv.Itoa(*r.ResourceOwnerID))
	accessTokenInsertable.Scopes = util.PointerToString(r.Scopes)
	accessTokenInsertable.ExpiresIn = time.Now().Add(m.tokenExpiresIn)

	return accessTokenInsertable
}

// AccessTokenFromOauthAccessGrant _
func (m *Model) AccessTokenFromOauthAccessGrant(oauthAccessGrant entity.OauthAccessGrant, application entity.OauthApplication) entity.OauthAccessTokenInsertable {
	var accessTokenInsertable entity.OauthAccessTokenInsertable

	accessTokenInsertable.OauthApplicationID = application.ID
	accessTokenInsertable.ResourceOwnerID = oauthAccessGrant.ResourceOwnerID
	accessTokenInsertable.Token = aurelia.Hash(application.ClientUID, application.ClientSecret+strconv.Itoa(oauthAccessGrant.ResourceOwnerID))
	accessTokenInsertable.Scopes = oauthAccessGrant.Scopes.String
	accessTokenInsertable.ExpiresIn = time.Now().Add(m.tokenExpiresIn)

	return accessTokenInsertable
}

// AccessTokenFromOauthRefreshToken _
func (m *Model) AccessTokenFromOauthRefreshToken(application entity.OauthApplication, accessToken entity.OauthAccessToken) entity.OauthAccessTokenInsertable {
	var accessTokenInsertable entity.OauthAccessTokenInsertable

	accessTokenInsertable.OauthApplicationID = application.ID
	accessTokenInsertable.ResourceOwnerID = accessToken.ResourceOwnerID
	accessTokenInsertable.Token = aurelia.Hash(application.ClientUID, application.ClientSecret+strconv.Itoa(accessToken.ResourceOwnerID))
	accessTokenInsertable.Scopes = accessToken.Scopes.String
	accessTokenInsertable.ExpiresIn = time.Now().Add(m.tokenExpiresIn)

	return accessTokenInsertable
}

// RefreshToken _
func (m *Model) RefreshToken(application entity.OauthApplication, accessToken entity.OauthAccessToken) entity.OauthRefreshTokenInsertable {
	var refreshTokenInsertable entity.OauthRefreshTokenInsertable

	refreshTokenInsertable.Token = aurelia.Hash(application.ClientUID, application.ClientSecret+strconv.Itoa(accessToken.ResourceOwnerID))
	refreshTokenInsertable.OauthAccessTokenID = accessToken.ID
	refreshTokenInsertable.ExpiresIn = time.Now().Add(m.refreshTokenExpiresIn)

	return refreshTokenInsertable
}
