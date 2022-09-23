package usecase

import (
	"github.com/kodefluence/altair/plugin/oauth/entity"
	"github.com/kodefluence/monorepo/jsonapi"
	"github.com/kodefluence/monorepo/kontext"
)

func (a *Authorization) Token(ktx kontext.Context, accessTokenReq entity.AccessTokenRequestJSON) (entity.OauthAccessTokenJSON, jsonapi.Errors) {
	oauthApplication, jsonapiErr := a.FindAndValidateApplication(ktx, accessTokenReq.ClientUID, accessTokenReq.ClientSecret)
	if jsonapiErr != nil {
		return entity.OauthAccessTokenJSON{}, jsonapiErr
	}

	if jsonapiErr := a.ValidateTokenGrant(accessTokenReq); jsonapiErr != nil {
		return entity.OauthAccessTokenJSON{}, jsonapiErr
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
			// Grant refresh token here
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

	return entity.OauthAccessTokenJSON{}, jsonapi.BuildResponse(
		a.apiError.ValidationError(`grant_type can't be empty`),
	).Errors
}
