package usecase

import (
	"time"

	"github.com/kodefluence/altair/plugin/oauth/entity"
	"github.com/kodefluence/monorepo/db"
	"github.com/kodefluence/monorepo/exception"
	"github.com/kodefluence/monorepo/jsonapi"
	"github.com/kodefluence/monorepo/kontext"
	"github.com/rs/zerolog"
)

func (a *Authorization) GrantTokenFromRefreshToken(ktx kontext.Context, accessTokenReq entity.AccessTokenRequestJSON) (entity.OauthAccessToken, entity.OauthRefreshToken, jsonapi.Errors) {
	var finalOauthAccessToken entity.OauthAccessToken
	var finalOauthRefreshToken entity.OauthRefreshToken

	exc := a.sqldb.Transaction(ktx, "authorization-grant-token-from-refresh-token", func(tx db.TX) exception.Exception {
		oldOauthRefreshToken, err := a.oauthRefreshTokenRepo.OneByToken(ktx, *accessTokenReq.RefreshToken, tx)
		if err != nil {
			if err.Type() == exception.NotFound {
				errorObject := jsonapi.BuildResponse(a.apiError.NotFoundError(ktx, "refresh_token")).Errors[0]
				return exception.Throw(err, exception.WithType(exception.NotFound), exception.WithDetail(errorObject.Detail), exception.WithTitle(errorObject.Title))
			}

			return exception.Throw(err, exception.WithType(exception.Unexpected), exception.WithTitle("Internal Server Error"), exception.WithDetail("refresh token cannot be found because there was an error"))
		}

		if oldOauthRefreshToken.RevokedAT.Valid || time.Now().After(oldOauthRefreshToken.ExpiresIn) {
			errorObject := jsonapi.BuildResponse(a.apiError.ForbiddenError(ktx, "access_token", "refresh token already used")).Errors[0]
			return exception.Throw(errorObject, exception.WithType(exception.Forbidden), exception.WithDetail(errorObject.Detail), exception.WithTitle(errorObject.Title))
		}

		oldAccessToken, err := a.oauthAccessTokenRepo.One(ktx, oldOauthRefreshToken.OauthAccessTokenID, tx)
		if err != nil {
			return exception.Throw(err, exception.WithType(exception.Unexpected), exception.WithTitle("Internal Server Error"), exception.WithDetail("access token cannot be found because there was an error"))
		}

		oauthApplication, err := a.oauthApplicationRepo.One(ktx, oldAccessToken.OauthApplicationID, tx)
		if err != nil {
			return exception.Throw(err, exception.WithType(exception.Unexpected), exception.WithTitle("Internal Server Error"), exception.WithDetail("error find oauth applications"))
		}

		oauthAccessTokenID, err := a.oauthAccessTokenRepo.Create(ktx, a.formatter.AccessTokenFromOauthRefreshTokenInsertable(oauthApplication, oldAccessToken), tx)
		if err != nil {
			return exception.Throw(err, exception.WithType(exception.Unexpected), exception.WithTitle("Internal Server Error"), exception.WithDetail("error creating access token"))
		}

		oauthAccessToken, err := a.oauthAccessTokenRepo.One(ktx, oauthAccessTokenID, tx)
		if err != nil {
			return exception.Throw(err, exception.WithType(exception.Unexpected), exception.WithTitle("Internal Server Error"), exception.WithDetail("error when selecting newly created access token"))
		}

		err = a.oauthRefreshTokenRepo.Revoke(ktx, *accessTokenReq.RefreshToken, tx)
		if err != nil {
			return exception.Throw(err, exception.WithType(exception.Unexpected), exception.WithTitle("Internal Server Error"), exception.WithDetail("error revoke refresh token"))
		}

		oauthRefreshTokenID, err := a.oauthRefreshTokenRepo.Create(ktx, a.formatter.RefreshTokenInsertable(oauthApplication, oauthAccessToken), tx)
		if err != nil {
			return exception.Throw(err, exception.WithType(exception.Unexpected), exception.WithTitle("Internal Server Error"), exception.WithDetail("error creating refresh token"))
		}

		oauthRefreshToken, err := a.oauthRefreshTokenRepo.One(ktx, oauthRefreshTokenID, tx)
		if err != nil {
			return exception.Throw(err, exception.WithType(exception.Unexpected), exception.WithTitle("Internal Server Error"), exception.WithDetail("error when selecting newly created refresh token"))
		}

		finalOauthAccessToken = oauthAccessToken
		finalOauthRefreshToken = oauthRefreshToken

		return nil
	})
	if exc != nil {
		return entity.OauthAccessToken{}, entity.OauthRefreshToken{}, a.exceptionMapping(ktx, exc, zerolog.Arr().Str("service").Str("authorization").Str("refresh_token"))
	}

	return finalOauthAccessToken, finalOauthRefreshToken, nil
}
