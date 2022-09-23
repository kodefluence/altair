package usecase

import (
	"github.com/kodefluence/altair/plugin/oauth/entity"
	"github.com/kodefluence/monorepo/db"
	"github.com/kodefluence/monorepo/exception"
	"github.com/kodefluence/monorepo/jsonapi"
	"github.com/kodefluence/monorepo/kontext"
	"github.com/rs/zerolog"
)

func (a *Authorization) ClientCredential(ktx kontext.Context, accessTokenReq entity.AccessTokenRequestJSON, oauthApplication entity.OauthApplication) (entity.OauthAccessToken, *entity.OauthRefreshToken, jsonapi.Errors) {
	var finalOauthAccessToken entity.OauthAccessToken
	var finalRefreshToken *entity.OauthRefreshToken

	exc := a.sqldb.Transaction(ktx, "authorization-grant-client-credential", func(tx db.TX) exception.Exception {
		id, err := a.oauthAccessTokenRepo.Create(ktx, a.formatter.AccessTokenClientCredentialInsertable(oauthApplication, accessTokenReq.Scope), tx)
		if err != nil {
			return exception.Throw(err, exception.WithDetail("error creating new oauth access token"), exception.WithType(exception.Unexpected), exception.WithTitle("access token creation error"))
		}

		oauthAccessToken, err := a.oauthAccessTokenRepo.One(ktx, id, tx)
		if err != nil {
			return exception.Throw(err, exception.WithDetail("error selecting newly created access token"), exception.WithType(exception.Unexpected), exception.WithTitle("access token creation error"))
		}

		if a.config.Config.RefreshToken.Active {
			if refreshToken, err := a.GrantRefreshToken(ktx, oauthAccessToken, oauthApplication, tx); err != nil {
				return exception.Throw(err, exception.WithType(exception.Unexpected), exception.WithTitle("Internal Server Error"), exception.WithDetail("error creating refresh token data"))
			} else {
				finalRefreshToken = &refreshToken
			}
		}

		finalOauthAccessToken = oauthAccessToken

		return nil
	})

	if exc != nil {
		return entity.OauthAccessToken{}, nil, a.exceptionMapping(ktx, exc, zerolog.Arr().Str("service").Str("authorization").Str("refresh_token"))
	}
	return finalOauthAccessToken, finalRefreshToken, nil
}
