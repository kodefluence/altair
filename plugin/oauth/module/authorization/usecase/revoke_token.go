package usecase

import (
	"github.com/kodefluence/altair/plugin/oauth/entity"
	"github.com/kodefluence/monorepo/jsonapi"
	"github.com/kodefluence/monorepo/kontext"
	"github.com/rs/zerolog"
)

// RevokeToken revoke given access token request
func (a *Authorization) RevokeToken(ktx kontext.Context, revokeAccessTokenReq entity.RevokeAccessTokenRequestJSON) jsonapi.Errors {

	if revokeAccessTokenReq.Token == nil {
		return jsonapi.BuildResponse(
			a.apiError.ValidationError("token cannot be empty"),
		).Errors
	}

	exc := a.oauthAccessTokenRepo.Revoke(ktx, *revokeAccessTokenReq.Token, a.sqldb)
	if exc != nil {
		return a.exceptionMapping(ktx, exc, zerolog.Arr().Str("service").Str("authorization").Str("revoke_token"))
	}

	return nil
}
