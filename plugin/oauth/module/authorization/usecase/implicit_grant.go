package usecase

import (
	"github.com/kodefluence/altair/plugin/oauth/entity"
	"github.com/kodefluence/monorepo/db"
	"github.com/kodefluence/monorepo/exception"
	"github.com/kodefluence/monorepo/jsonapi"
	"github.com/kodefluence/monorepo/kontext"
	"github.com/rs/zerolog"
)

func (a *Authorization) ImplicitGrant(ktx kontext.Context, authorizationReq entity.AuthorizationRequestJSON) (entity.OauthAccessTokenJSON, jsonapi.Errors) {
	var finalOauthTokenJSON entity.OauthAccessTokenJSON

	oauthApplication, jsonError := a.FindAndValidateApplication(ktx, authorizationReq.ClientUID, authorizationReq.ClientSecret)
	if jsonError != nil {
		return entity.OauthAccessTokenJSON{}, jsonError
	}

	if err := a.ValidateAuthorizationGrant(ktx, authorizationReq, oauthApplication); err != nil {
		return entity.OauthAccessTokenJSON{}, err
	}

	exc := a.sqldb.Transaction(ktx, "authorization-grant-token", func(tx db.TX) exception.Exception {
		id, err := a.oauthAccessTokenRepo.Create(ktx, a.formatter.AccessTokenFromAuthorizationRequestInsertable(authorizationReq, oauthApplication), tx)
		if err != nil {
			return exception.Throw(err, exception.WithDetail("error creating new oauth access token"), exception.WithType(exception.Unexpected), exception.WithTitle("access token creation error"))
		}

		oauthAccessToken, err := a.oauthAccessTokenRepo.One(ktx, id, tx)
		if err != nil {
			return exception.Throw(err, exception.WithDetail("error selecting newly created access token"), exception.WithType(exception.Unexpected), exception.WithTitle("access token creation error"))
		}

		var oauthRefreshTokenJSON *entity.OauthRefreshTokenJSON

		if a.config.Config.RefreshToken.Active {
			refreshTokenID, err := a.oauthRefreshTokenRepo.Create(ktx, a.formatter.RefreshTokenInsertable(oauthApplication, oauthAccessToken), tx)
			if err != nil {
				return exception.Throw(err, exception.WithDetail("error creating new oauth refresh token"), exception.WithType(exception.Unexpected), exception.WithTitle("access token creation error"))
			}

			oauthRefreshToken, err := a.oauthRefreshTokenRepo.One(ktx, refreshTokenID, tx)
			if err != nil {
				return exception.Throw(err, exception.WithDetail("error selecting newly created refresh token"), exception.WithType(exception.Unexpected), exception.WithTitle("access token creation error"))
			}

			newOauthRefreshTokenJSON := a.formatter.RefreshToken(oauthRefreshToken)
			oauthRefreshTokenJSON = &newOauthRefreshTokenJSON
		}

		finalOauthTokenJSON = a.formatter.AccessToken(oauthAccessToken, *authorizationReq.RedirectURI, oauthRefreshTokenJSON)

		return nil
	})
	if exc != nil {
		return entity.OauthAccessTokenJSON{}, a.exceptionMapping(ktx, exc, zerolog.Arr().Str("service").Str("authorization").Str("grant_token"))
	}

	return finalOauthTokenJSON, nil
}
