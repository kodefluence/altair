package service

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/codefluence-x/altair/provider/plugin/oauth/entity"
	"github.com/codefluence-x/altair/provider/plugin/oauth/eobject"
	"github.com/codefluence-x/altair/provider/plugin/oauth/interfaces"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type ApplicationManager struct {
	formatter             interfaces.OauthApplicationFormater
	modelFormatter        interfaces.ModelFormater
	oauthApplicationModel interfaces.OauthApplicationModel
	applicationValidator  interfaces.OauthValidator
}

// NewApplicationManager manage all oauth application data business logic
func NewApplicationManager(formatter interfaces.OauthApplicationFormater, modelFormatter interfaces.ModelFormater, oauthApplicationModel interfaces.OauthApplicationModel, applicationValidator interfaces.OauthValidator) *ApplicationManager {
	return &ApplicationManager{
		formatter:             formatter,
		modelFormatter:        modelFormatter,
		oauthApplicationModel: oauthApplicationModel,
		applicationValidator:  applicationValidator,
	}
}

func (am *ApplicationManager) List(ctx context.Context, offset, limit int) ([]entity.OauthApplicationJSON, int, *entity.Error) {
	oauthApplications, err := am.oauthApplicationModel.Paginate(ctx, offset, limit)
	if err != nil {
		log.Error().
			Err(err).
			Stack().
			Interface("request_id", ctx.Value("request_id")).
			Int("offset", offset).
			Int("limit", limit).
			Array("tags", zerolog.Arr().Str("service").Str("application_manager").Str("list").Str("paginate")).
			Msg("Error paginating oauth applications")

		return []entity.OauthApplicationJSON(nil), 0, &entity.Error{
			HttpStatus: http.StatusInternalServerError,
			Errors:     eobject.Wrap(eobject.InternalServerError(ctx)),
		}
	}

	total, err := am.oauthApplicationModel.Count(ctx)
	if err != nil {
		log.Error().
			Err(err).
			Stack().
			Interface("request_id", ctx.Value("request_id")).
			Int("offset", offset).
			Int("limit", limit).
			Array("tags", zerolog.Arr().Str("service").Str("application_manager").Str("list").Str("count")).
			Msg("Error count total of oauth applications")
		return []entity.OauthApplicationJSON(nil), 0, &entity.Error{
			HttpStatus: http.StatusInternalServerError,
			Errors:     eobject.Wrap(eobject.InternalServerError(ctx)),
		}
	}

	formattedResult := am.formatter.ApplicationList(ctx, oauthApplications)
	return formattedResult, total, nil
}

func (am *ApplicationManager) Create(ctx context.Context, e entity.OauthApplicationJSON) (entity.OauthApplicationJSON, *entity.Error) {
	if err := am.applicationValidator.ValidateApplication(ctx, e); err != nil {
		log.Error().
			Err(err).
			Stack().
			Interface("data", e).
			Array("tags", zerolog.Arr().Str("service").Str("application_manager").Str("create").Str("validate_application")).
			Msg("Got validation error from oauth application validator")
		return entity.OauthApplicationJSON{}, err
	}

	id, err := am.oauthApplicationModel.Create(ctx, am.modelFormatter.OauthApplication(e))
	if err != nil {
		log.Error().
			Err(err).
			Stack().
			Interface("data", e).
			Array("tags", zerolog.Arr().Str("service").Str("application_manager").Str("create").Str("model_create")).
			Msg("Error when creating oauth application data")

		return entity.OauthApplicationJSON{}, &entity.Error{
			HttpStatus: http.StatusInternalServerError,
			Errors:     eobject.Wrap(eobject.InternalServerError(ctx)),
		}
	}

	return am.One(ctx, id)
}

func (am *ApplicationManager) Update(ctx context.Context, ID int, e entity.OauthApplicationUpdateJSON) (entity.OauthApplicationJSON, *entity.Error) {

	err := am.oauthApplicationModel.Update(ctx, ID, entity.OauthApplicationUpdateable{
		Description: e.Description,
		Scopes:      e.Scopes,
	})
	if err != nil {
		log.Error().
			Err(err).
			Stack().
			Interface("data", e).
			Array("tags", zerolog.Arr().Str("service").Str("application_manager").Str("update").Str("model_update")).
			Msg("Error when updating oauth application data")
		return entity.OauthApplicationJSON{}, &entity.Error{
			HttpStatus: http.StatusInternalServerError,
			Errors:     eobject.Wrap(eobject.InternalServerError(ctx)),
		}
	}

	return am.One(ctx, ID)
}

func (am *ApplicationManager) One(ctx context.Context, ID int) (entity.OauthApplicationJSON, *entity.Error) {
	oauthApplication, err := am.oauthApplicationModel.One(ctx, ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return entity.OauthApplicationJSON{}, &entity.Error{
				HttpStatus: http.StatusNotFound,
				Errors:     eobject.Wrap(eobject.NotFoundError(ctx, "oauth_application")),
			}
		}

		log.Error().
			Err(err).
			Stack().
			Int("id", ID).
			Array("tags", zerolog.Arr().Str("service").Str("application_manager").Str("one").Str("model_one")).
			Msg("Error when fetching single oauth application")

		return entity.OauthApplicationJSON{}, &entity.Error{
			HttpStatus: http.StatusInternalServerError,
			Errors:     eobject.Wrap(eobject.InternalServerError(ctx)),
		}
	}

	formattedResult := am.formatter.Application(ctx, oauthApplication)
	return formattedResult, nil
}
