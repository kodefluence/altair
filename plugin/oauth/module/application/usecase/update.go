package usecase

import (
	"github.com/kodefluence/altair/plugin/oauth/entity"
	"github.com/kodefluence/monorepo/jsonapi"
	"github.com/kodefluence/monorepo/kontext"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Update oauth application
func (am *ApplicationManager) Update(ktx kontext.Context, ID int, e entity.OauthApplicationUpdateJSON) (entity.OauthApplicationJSON, jsonapi.Errors) {
	err := am.oauthApplicationRepo.Update(ktx, ID, entity.OauthApplicationUpdateable{
		Description: e.Description,
		Scopes:      e.Scopes,
	}, am.sqldb)
	if err != nil {
		log.Error().
			Err(err).
			Stack().
			Interface("data", e).
			Array("tags", zerolog.Arr().Str("service").Str("application_manager").Str("update").Str("model_update")).
			Msg("Error when updating oauth application data")

		return entity.OauthApplicationJSON{},
			jsonapi.BuildResponse(am.apiError.InternalServerError(ktx)).Errors
	}

	return am.One(ktx, ID)
}
