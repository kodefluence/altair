package usecase

import (
	"github.com/kodefluence/altair/plugin/oauth/entity"
	"github.com/kodefluence/monorepo/jsonapi"
	"github.com/kodefluence/monorepo/kontext"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// List of oauth applications
func (am *ApplicationManager) List(ktx kontext.Context, offset, limit int) ([]entity.OauthApplicationJSON, int, jsonapi.Errors) {
	oauthApplications, err := am.oauthApplicationRepo.Paginate(ktx, offset, limit, am.sqldb)
	if err != nil {
		log.Error().
			Err(err).
			Stack().
			Interface("request_id", ktx.GetWithoutCheck("request_id")).
			Int("offset", offset).
			Int("limit", limit).
			Array("tags", zerolog.Arr().Str("service").Str("application_manager").Str("list").Str("paginate")).
			Msg("Error paginating oauth applications")

		return []entity.OauthApplicationJSON(nil), 0, jsonapi.BuildResponse(
			am.apiError.InternalServerError(ktx),
		).Errors
	}

	total, err := am.oauthApplicationRepo.Count(ktx, am.sqldb)
	if err != nil {
		log.Error().
			Err(err).
			Stack().
			Interface("request_id", ktx.GetWithoutCheck("request_id")).
			Int("offset", offset).
			Int("limit", limit).
			Array("tags", zerolog.Arr().Str("service").Str("application_manager").Str("list").Str("count")).
			Msg("Error count total of oauth applications")

		return []entity.OauthApplicationJSON(nil), 0, jsonapi.BuildResponse(
			am.apiError.InternalServerError(ktx),
		).Errors
	}

	return am.formatter.ApplicationList(oauthApplications), total, nil
}
