package usecase

import (
	"github.com/kodefluence/altair/plugin/oauth/entity"
	"github.com/kodefluence/monorepo/jsonapi"
	"github.com/kodefluence/monorepo/kontext"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Create oauth application
func (am *ApplicationManager) Create(ktx kontext.Context, e entity.OauthApplicationJSON) (entity.OauthApplicationJSON, jsonapi.Errors) {
	if err := am.ValidateApplication(e); err != nil {
		log.Error().
			Err(err).
			Stack().
			Interface("data", e).
			Array("tags", zerolog.Arr().Str("service").Str("application_manager").Str("create").Str("validate_application")).
			Msg("Got validation error from oauth application validator")
		return entity.OauthApplicationJSON{}, err
	}

	id, err := am.oauthApplicationRepo.Create(ktx, am.formatter.OauthApplicationInsertable(e), am.sqldb)
	if err != nil {
		log.Error().
			Err(err).
			Stack().
			Interface("data", e).
			Array("tags", zerolog.Arr().Str("service").Str("application_manager").Str("create").Str("model_create")).
			Msg("Error when creating oauth application data")

		return entity.OauthApplicationJSON{},
			jsonapi.BuildResponse(am.apiError.InternalServerError(ktx)).Errors
	}

	return am.One(ktx, id)
}
