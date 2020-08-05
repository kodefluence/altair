package service

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/codefluence-x/altair/provider/plugin/oauth/entity"
	"github.com/codefluence-x/altair/provider/plugin/oauth/eobject"
	"github.com/codefluence-x/altair/provider/plugin/oauth/interfaces"
	"github.com/codefluence-x/journal"
)

type applicationManager struct {
	formatter             interfaces.OauthApplicationFormater
	modelFormatter        interfaces.ModelFormater
	oauthApplicationModel interfaces.OauthApplicationModel
	applicationValidator  interfaces.OauthValidator
}

func ApplicationManager(formatter interfaces.OauthApplicationFormater, modelFormatter interfaces.ModelFormater, oauthApplicationModel interfaces.OauthApplicationModel, applicationValidator interfaces.OauthValidator) *applicationManager {
	return &applicationManager{
		formatter:             formatter,
		modelFormatter:        modelFormatter,
		oauthApplicationModel: oauthApplicationModel,
		applicationValidator:  applicationValidator,
	}
}

func (am *applicationManager) List(ctx context.Context, offset, limit int) ([]entity.OauthApplicationJSON, int, *entity.Error) {
	oauthApplications, err := am.oauthApplicationModel.Paginate(ctx, offset, limit)
	if err != nil {
		journal.Error("Error paginating oauth applications", err).
			AddField("offset", offset).
			AddField("limit", limit).
			SetTags("service", "application_manager", "list", "paginate").
			SetTrackId(ctx.Value("track_id")).
			Log()

		return []entity.OauthApplicationJSON(nil), 0, &entity.Error{
			HttpStatus: http.StatusInternalServerError,
			Errors:     eobject.Wrap(eobject.InternalServerError(ctx)),
		}
	}

	total, err := am.oauthApplicationModel.Count(ctx)
	if err != nil {
		journal.Error("Error count total of oauth applications", err).
			AddField("offset", offset).
			AddField("limit", limit).
			SetTags("service", "application_manager", "list", "count").
			SetTrackId(ctx.Value("track_id")).
			Log()

		return []entity.OauthApplicationJSON(nil), 0, &entity.Error{
			HttpStatus: http.StatusInternalServerError,
			Errors:     eobject.Wrap(eobject.InternalServerError(ctx)),
		}
	}

	formattedResult := am.formatter.ApplicationList(ctx, oauthApplications)
	return formattedResult, total, nil
}

func (am *applicationManager) Create(ctx context.Context, e entity.OauthApplicationJSON) (entity.OauthApplicationJSON, *entity.Error) {
	if err := am.applicationValidator.ValidateApplication(ctx, e); err != nil {
		journal.Error("Got validation error from oauth application validator", err).
			AddField("data", e).
			SetTags("service", "application_manager", "create", "model_create").
			SetTrackId(ctx.Value("track_id")).
			Log()

		return entity.OauthApplicationJSON{}, err
	}

	id, err := am.oauthApplicationModel.Create(ctx, am.modelFormatter.OauthApplication(e))
	if err != nil {
		journal.Error("Error when creating oauth application data", err).
			AddField("data", e).
			SetTags("service", "application_manager", "create", "model_create").
			SetTrackId(ctx.Value("track_id")).
			Log()

		return entity.OauthApplicationJSON{}, &entity.Error{
			HttpStatus: http.StatusInternalServerError,
			Errors:     eobject.Wrap(eobject.InternalServerError(ctx)),
		}
	}

	return am.One(ctx, id)
}

func (am *applicationManager) Update(ctx context.Context, ID int, e entity.OauthApplicationUpdateJSON) (entity.OauthApplicationJSON, *entity.Error) {

	err := am.oauthApplicationModel.Update(ctx, ID, entity.OauthApplicationUpdateable{
		Description: e.Description,
		Scopes:      e.Scopes,
	})
	if err != nil {
		journal.Error("Error when updating oauth application data", err).
			AddField("data", e).
			SetTags("service", "application_manager", "update", "model_update").
			SetTrackId(ctx.Value("track_id")).
			Log()

		return entity.OauthApplicationJSON{}, &entity.Error{
			HttpStatus: http.StatusInternalServerError,
			Errors:     eobject.Wrap(eobject.InternalServerError(ctx)),
		}
	}

	return am.One(ctx, ID)
}

func (am *applicationManager) One(ctx context.Context, ID int) (entity.OauthApplicationJSON, *entity.Error) {
	oauthApplication, err := am.oauthApplicationModel.One(ctx, ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return entity.OauthApplicationJSON{}, &entity.Error{
				HttpStatus: http.StatusNotFound,
				Errors:     eobject.Wrap(eobject.NotFoundError(ctx, "oauth_application")),
			}
		}

		journal.Error("Error when fetching single oauth application", err).
			AddField("id", ID).
			SetTags("service", "application_manager", "one", "model_one").
			SetTrackId(ctx.Value("track_id")).
			Log()

		return entity.OauthApplicationJSON{}, &entity.Error{
			HttpStatus: http.StatusInternalServerError,
			Errors:     eobject.Wrap(eobject.InternalServerError(ctx)),
		}
	}

	formattedResult := am.formatter.Application(ctx, oauthApplication)
	return formattedResult, nil
}
