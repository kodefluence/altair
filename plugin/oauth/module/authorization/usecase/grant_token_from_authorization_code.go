package usecase

import (
	"github.com/kodefluence/altair/plugin/oauth/entity"
	"github.com/kodefluence/monorepo/db"
	"github.com/kodefluence/monorepo/exception"
	"github.com/kodefluence/monorepo/jsonapi"
	"github.com/kodefluence/monorepo/kontext"
	"github.com/rs/zerolog"
)

func (a *Authorization) GrantTokenFromAuthorizationCode(ktx kontext.Context, accessTokenReq entity.AccessTokenRequestJSON, oauthApplication entity.OauthApplication) (entity.OauthAccessToken, string, jsonapi.Errors) {
	var finalOauthAccessToken entity.OauthAccessToken
	var finalRedirectURI string

	exc := a.sqldb.Transaction(ktx, "authorization-grant-token-from-refresh-token", func(tx db.TX) exception.Exception {
		oauthAccessGrant, err := a.oauthAccessGrantRepo.OneByCode(ktx, *accessTokenReq.Code, tx)
		if err != nil {
			if err.Type() == exception.NotFound {
				errorObject := jsonapi.BuildResponse(a.apiError.NotFoundError(ktx, "authorization_code")).Errors[0]
				return exception.Throw(err, exception.WithType(exception.NotFound), exception.WithDetail(errorObject.Detail), exception.WithTitle(errorObject.Title))
			}

			return exception.Throw(err, exception.WithType(exception.Unexpected), exception.WithTitle("Internal Server Error"), exception.WithDetail("authorization code cannot be found because there was an error"))
		}

		if exc := a.ValidateTokenAuthorizationCode(ktx, accessTokenReq, oauthAccessGrant); exc != nil {
			return exc
		}

		id, err := a.oauthAccessTokenRepo.Create(ktx, a.formatter.AccessTokenFromOauthAccessGrantInsertable(oauthAccessGrant, oauthApplication), tx)
		if err != nil {
			return exception.Throw(err, exception.WithType(exception.Unexpected), exception.WithTitle("Internal Server Error"), exception.WithDetail("error creating access token data"))
		}

		oauthAccessToken, err := a.oauthAccessTokenRepo.One(ktx, id, tx)
		if err != nil {
			return exception.Throw(err, exception.WithType(exception.Unexpected), exception.WithTitle("Internal Server Error"), exception.WithDetail("error selecting newly created access token"))
		}

		err = a.oauthAccessGrantRepo.Revoke(ktx, *accessTokenReq.Code, tx)
		if err != nil {
			return exception.Throw(err, exception.WithType(exception.Unexpected), exception.WithTitle("Internal Server Error"), exception.WithDetail("error revoking oauth access grant"))
		}

		finalOauthAccessToken = oauthAccessToken
		finalRedirectURI = oauthAccessGrant.RedirectURI.String

		return nil
	})
	if exc != nil {
		return entity.OauthAccessToken{}, "", a.exceptionMapping(ktx, exc, zerolog.Arr().Str("service").Str("authorization").Str("refresh_token"))
	}

	return finalOauthAccessToken, finalRedirectURI, nil
}
