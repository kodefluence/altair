package usecase

import (
	"github.com/kodefluence/altair/plugin/oauth/entity"
	"github.com/kodefluence/monorepo/jsonapi"
	"github.com/kodefluence/monorepo/kontext"
)

func (a *Authorization) Token(ktx kontext.Context, accessTokenReq entity.AccessTokenRequestJSON) (entity.OauthAccessTokenJSON, jsonapi.Errors) {
	_, jsonapiErr := a.FindAndValidateApplication(ktx, accessTokenReq.ClientUID, accessTokenReq.ClientSecret)
	if jsonapiErr != nil {
		return entity.OauthAccessTokenJSON{}, jsonapiErr
	}

	if jsonapiErr := a.ValidateTokenGrant(accessTokenReq); jsonapiErr != nil {
		return entity.OauthAccessTokenJSON{}, jsonapiErr
	}

	switch *accessTokenReq.GrantType {

	case "authorization_code":
		// Grant authorization code here
	case "refresh_token":
		if a.config.Config.RefreshToken.Active {
			// Grant refresh token here
		}
	case "client_credentials":
		// Grant client credentials here
	}

	return entity.OauthAccessTokenJSON{}, jsonapi.BuildResponse(
		a.apiError.ValidationError(`grant_type can't be empty`),
	).Errors
}
