package usecase

import (
	"github.com/kodefluence/altair/plugin/oauth/entity"
	"github.com/kodefluence/monorepo/exception"
	"github.com/kodefluence/monorepo/jsonapi"
	"github.com/kodefluence/monorepo/kontext"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func (a *Authorization) FindAndValidateApplication(ktx kontext.Context, clientUID, clientSecret *string) (entity.OauthApplication, jsonapi.Errors) {
	if clientUID == nil {
		return entity.OauthApplication{}, jsonapi.BuildResponse(
			a.apiError.ValidationError("client_uid cannot be empty"),
		).Errors
	}

	if clientSecret == nil {
		return entity.OauthApplication{}, jsonapi.BuildResponse(
			a.apiError.ValidationError("client_secret cannot be empty"),
		).Errors
	}

	oauthApplication, err := a.oauthApplicationRepo.OneByUIDandSecret(ktx, *clientUID, *clientSecret, a.sqldb)
	if err != nil {
		log.Error().
			Err(err).
			Stack().
			Interface("request_id", ktx.GetWithoutCheck("request_id")).
			Str("client_uid", *clientUID).
			Array("tags", zerolog.Arr().Str("service").Str("authorization").Str("find_secret")).
			Msg("application cannot be found because there was an error")

		if err.Type() == exception.NotFound {
			return entity.OauthApplication{},
				jsonapi.BuildResponse(a.apiError.NotFoundError(ktx, "client_uid & client_secret")).Errors
		}

		return entity.OauthApplication{},
			jsonapi.BuildResponse(a.apiError.InternalServerError(ktx)).Errors
	}

	return oauthApplication, nil
}
