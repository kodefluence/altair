package usecase

import (
	"github.com/kodefluence/altair/plugin/oauth/entity"
	"github.com/kodefluence/monorepo/exception"
	"github.com/kodefluence/monorepo/jsonapi"
	"github.com/kodefluence/monorepo/kontext"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// One retrieve oauth application
func (am *ApplicationManager) One(ktx kontext.Context, ID int) (entity.OauthApplicationJSON, jsonapi.Errors) {
	oauthApplication, err := am.oauthApplicationRepo.One(ktx, ID, am.sqldb)
	if err != nil {
		if err.Type() == exception.NotFound {

			return entity.OauthApplicationJSON{},
				jsonapi.BuildResponse(am.apiError.NotFoundError(ktx, "oauth_application")).Errors
		}

		log.Error().
			Err(err).
			Stack().
			Int("id", ID).
			Array("tags", zerolog.Arr().Str("service").Str("application_manager").Str("one").Str("model_one")).
			Msg("Error when fetching single oauth application")

		return entity.OauthApplicationJSON{},
			jsonapi.BuildResponse(am.apiError.InternalServerError(ktx)).Errors
	}

	formattedResult := am.formatter.Application(oauthApplication)
	return formattedResult, nil
}
