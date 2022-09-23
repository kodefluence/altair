package usecase

import (
	"github.com/kodefluence/altair/plugin/oauth/entity"
	"github.com/kodefluence/monorepo/jsonapi"
	"github.com/kodefluence/monorepo/kontext"
)

func (a *Authorization) GrantToken(ktx kontext.Context, accessTokenReq entity.AccessTokenRequestJSON) (entity.OauthAccessTokenJSON, jsonapi.Errors) {
	if jsonapiErr := a.ValidateTokenGrant(accessTokenReq); jsonapiErr != nil {
		return entity.OauthAccessTokenJSON{}, jsonapiErr
	}

	var oauthApplication entity.OauthApplication
	var jsonapierr jsonapi.Errors

	if *accessTokenReq.GrantType != "refresh_token" {
		oauthApplication, jsonapierr = a.FindAndValidateApplication(ktx, accessTokenReq.ClientUID, accessTokenReq.ClientSecret)
		if jsonapierr != nil {
			return entity.OauthAccessTokenJSON{}, jsonapierr
		}
	}

	switch *accessTokenReq.GrantType {
	case "authorization_code":
		oauthAccessToken, oauthRefreshToken, redirectURI, jsonapierr := a.GrantTokenFromAuthorizationCode(ktx, accessTokenReq, oauthApplication)
		if jsonapierr != nil {
			return entity.OauthAccessTokenJSON{}, jsonapierr
		}

		if oauthRefreshToken == nil {
			return a.formatter.AccessToken(oauthAccessToken, redirectURI, nil), nil
		}

		refreshTokenJSON := a.formatter.RefreshToken(*oauthRefreshToken)
		return a.formatter.AccessToken(oauthAccessToken, redirectURI, &refreshTokenJSON), nil
	case "refresh_token":
		if a.config.Config.RefreshToken.Active {
			oauthAccessToken, oauthRefreshToken, jsonapierr := a.GrantTokenFromRefreshToken(ktx, accessTokenReq)
			if jsonapierr != nil {
				return entity.OauthAccessTokenJSON{}, jsonapierr
			}

			refreshTokenJSON := a.formatter.RefreshToken(oauthRefreshToken)
			return a.formatter.AccessToken(oauthAccessToken, "", &refreshTokenJSON), nil
		}
	case "client_credentials":
		oauthAccessToken, oauthRefreshToken, jsonapierr := a.ClientCredential(ktx, accessTokenReq, oauthApplication)
		if jsonapierr != nil {
			return entity.OauthAccessTokenJSON{}, jsonapierr
		}

		if oauthRefreshToken == nil {
			return a.formatter.AccessToken(oauthAccessToken, "", nil), nil
		}

		refreshTokenJSON := a.formatter.RefreshToken(*oauthRefreshToken)
		return a.formatter.AccessToken(oauthAccessToken, "", &refreshTokenJSON), nil
	}

	// This code is actually unreachable since there are already validation put in place in ValidateTokenGrant lol
	// But I'll keep here just in case
	return entity.OauthAccessTokenJSON{}, jsonapi.BuildResponse(
		a.apiError.ValidationError(`grant_type can't be empty`),
	).Errors
}
