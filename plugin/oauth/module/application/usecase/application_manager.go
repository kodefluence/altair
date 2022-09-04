package usecase

import (
	"context"
	"net/http"

	"github.com/kodefluence/altair/provider/plugin/oauth/entity"
	"github.com/kodefluence/altair/provider/plugin/oauth/eobject"
	"github.com/kodefluence/altair/provider/plugin/oauth/interfaces"
	"github.com/kodefluence/monorepo/db"
	"github.com/kodefluence/monorepo/exception"
	"github.com/kodefluence/monorepo/jsonapi"
	"github.com/kodefluence/monorepo/kontext"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// ApplicationManager manage all oauth_applications CRUD
type ApplicationManager struct {
	formatter             interfaces.OauthApplicationFormater
	modelFormatter        interfaces.ModelFormater
	oauthApplicationModel interfaces.OauthApplicationModel
	applicationValidator  interfaces.OauthValidator
	sqldb                 db.DB
}

// NewApplicationManager manage all oauth application data business logic
func NewApplicationManager(formatter interfaces.OauthApplicationFormater, modelFormatter interfaces.ModelFormater, oauthApplicationModel interfaces.OauthApplicationModel, applicationValidator interfaces.OauthValidator, sqldb db.DB) *ApplicationManager {
	return &ApplicationManager{
		formatter:             formatter,
		modelFormatter:        modelFormatter,
		oauthApplicationModel: oauthApplicationModel,
		applicationValidator:  applicationValidator,
		sqldb:                 sqldb,
	}
}

// List of oauth applications
func (am *ApplicationManager) List(ctx context.Context, offset, limit int) ([]entity.OauthApplicationJSON, int, jsonapi.Errors) {
	oauthApplications, err := am.oauthApplicationModel.Paginate(kontext.Fabricate(kontext.WithDefaultContext(ctx)), offset, limit, am.sqldb)
	if err != nil {
		log.Error().
			Err(err).
			Stack().
			Interface("request_id", ctx.Value("request_id")).
			Int("offset", offset).
			Int("limit", limit).
			Array("tags", zerolog.Arr().Str("service").Str("application_manager").Str("list").Str("paginate")).
			Msg("Error paginating oauth applications")

		// return []entity.OauthApplicationJSON(nil), 0, &entity.Error{
		// 	HttpStatus: http.StatusInternalServerError,
		// 	Errors:     eobject.Wrap(eobject.InternalServerError(ctx)),
		// }
		return []entity.OauthApplicationJSON(nil), 0, jsonapi.BuildResponse().Errors
	}

	total, err := am.oauthApplicationModel.Count(kontext.Fabricate(kontext.WithDefaultContext(ctx)), am.sqldb)
	if err != nil {
		log.Error().
			Err(err).
			Stack().
			Interface("request_id", ctx.Value("request_id")).
			Int("offset", offset).
			Int("limit", limit).
			Array("tags", zerolog.Arr().Str("service").Str("application_manager").Str("list").Str("count")).
			Msg("Error count total of oauth applications")
		// return []entity.OauthApplicationJSON(nil), 0, &entity.Error{
		// 	HttpStatus: http.StatusInternalServerError,
		// 	Errors:     eobject.Wrap(eobject.InternalServerError(ctx)),
		// }
		return []entity.OauthApplicationJSON(nil), 0, jsonapi.BuildResponse().Errors
	}

	formattedResult := am.formatter.ApplicationList(ctx, oauthApplications)
	return formattedResult, total, nil
}

// Create oauth application
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

	id, err := am.oauthApplicationModel.Create(kontext.Fabricate(kontext.WithDefaultContext(ctx)), am.modelFormatter.OauthApplication(e), am.sqldb)
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

// Update oauth application
func (am *ApplicationManager) Update(ctx context.Context, ID int, e entity.OauthApplicationUpdateJSON) (entity.OauthApplicationJSON, *entity.Error) {

	err := am.oauthApplicationModel.Update(kontext.Fabricate(kontext.WithDefaultContext(ctx)), ID, entity.OauthApplicationUpdateable{
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
		return entity.OauthApplicationJSON{}, &entity.Error{
			HttpStatus: http.StatusInternalServerError,
			Errors:     eobject.Wrap(eobject.InternalServerError(ctx)),
		}
	}

	return am.One(ctx, ID)
}

// One retrieve oauth application
func (am *ApplicationManager) One(ctx context.Context, ID int) (entity.OauthApplicationJSON, *entity.Error) {
	oauthApplication, err := am.oauthApplicationModel.One(kontext.Fabricate(kontext.WithDefaultContext(ctx)), ID, am.sqldb)
	if err != nil {
		if err.Type() == exception.NotFound {
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
