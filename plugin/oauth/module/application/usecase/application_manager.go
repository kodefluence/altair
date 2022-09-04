package usecase

import (
	"github.com/kodefluence/altair/plugin/oauth/entity"
	"github.com/kodefluence/monorepo/db"
	"github.com/kodefluence/monorepo/exception"
	"github.com/kodefluence/monorepo/jsonapi"
	"github.com/kodefluence/monorepo/kontext"
)

//go:generate mockgen -destination ./mock/mock.go -package mock -source ./application_manager.go

type ApiError interface {
	InternalServerError(ktx kontext.Context) jsonapi.Option
	BadRequestError(in string) jsonapi.Option
	NotFoundError(ktx kontext.Context, entityType string) jsonapi.Option
	UnauthorizedError() jsonapi.Option
	ForbiddenError(ktx kontext.Context, entityType, reason string) jsonapi.Option
	ValidationError(msg string) jsonapi.Option
}

type OauthApplicationRepository interface {
	Paginate(ktx kontext.Context, offset, limit int, tx db.TX) ([]entity.OauthApplication, exception.Exception)
	Count(ktx kontext.Context, tx db.TX) (int, exception.Exception)
	One(ktx kontext.Context, ID int, tx db.TX) (entity.OauthApplication, exception.Exception)
	OneByUIDandSecret(ktx kontext.Context, clientUID, clientSecret string, tx db.TX) (entity.OauthApplication, exception.Exception)
	Create(ktx kontext.Context, data entity.OauthApplicationInsertable, tx db.TX) (int, exception.Exception)
	Update(ktx kontext.Context, ID int, data entity.OauthApplicationUpdateable, tx db.TX) exception.Exception
}

type Formatter interface {
	ApplicationList(applications []entity.OauthApplication) []entity.OauthApplicationJSON
}

// ApplicationManager manage all oauth_applications CRUD
type ApplicationManager struct {
	sqldb                db.DB
	oauthApplicationRepo OauthApplicationRepository
	apiError             ApiError
	formatter            Formatter
}

// NewApplicationManager manage all oauth application data business logic
func NewApplicationManager(sqldb db.DB, oauthApplicationRepo OauthApplicationRepository, apiError ApiError, formatter Formatter) *ApplicationManager {
	return &ApplicationManager{
		sqldb:                sqldb,
		oauthApplicationRepo: oauthApplicationRepo,
		apiError:             apiError,
		formatter:            formatter,
		// applicationValidator:  applicationValidator,
	}
}

// // Create oauth application
// func (am *ApplicationManager) Create(ctx context.Context, e entity.OauthApplicationJSON) (entity.OauthApplicationJSON, *entity.Error) {
// 	if err := am.applicationValidator.ValidateApplication(ctx, e); err != nil {
// 		log.Error().
// 			Err(err).
// 			Stack().
// 			Interface("data", e).
// 			Array("tags", zerolog.Arr().Str("service").Str("application_manager").Str("create").Str("validate_application")).
// 			Msg("Got validation error from oauth application validator")
// 		return entity.OauthApplicationJSON{}, err
// 	}

// 	id, err := am.oauthApplicationModel.Create(kontext.Fabricate(kontext.WithDefaultContext(ctx)), am.modelFormatter.OauthApplication(e), am.sqldb)
// 	if err != nil {
// 		log.Error().
// 			Err(err).
// 			Stack().
// 			Interface("data", e).
// 			Array("tags", zerolog.Arr().Str("service").Str("application_manager").Str("create").Str("model_create")).
// 			Msg("Error when creating oauth application data")

// 		return entity.OauthApplicationJSON{}, &entity.Error{
// 			HttpStatus: http.StatusInternalServerError,
// 			Errors:     eobject.Wrap(eobject.InternalServerError(ctx)),
// 		}
// 	}

// 	return am.One(ctx, id)
// }

// // Update oauth application
// func (am *ApplicationManager) Update(ctx context.Context, ID int, e entity.OauthApplicationUpdateJSON) (entity.OauthApplicationJSON, *entity.Error) {

// 	err := am.oauthApplicationModel.Update(kontext.Fabricate(kontext.WithDefaultContext(ctx)), ID, entity.OauthApplicationUpdateable{
// 		Description: e.Description,
// 		Scopes:      e.Scopes,
// 	}, am.sqldb)
// 	if err != nil {
// 		log.Error().
// 			Err(err).
// 			Stack().
// 			Interface("data", e).
// 			Array("tags", zerolog.Arr().Str("service").Str("application_manager").Str("update").Str("model_update")).
// 			Msg("Error when updating oauth application data")
// 		return entity.OauthApplicationJSON{}, &entity.Error{
// 			HttpStatus: http.StatusInternalServerError,
// 			Errors:     eobject.Wrap(eobject.InternalServerError(ctx)),
// 		}
// 	}

// 	return am.One(ctx, ID)
// }

// // One retrieve oauth application
// func (am *ApplicationManager) One(ctx context.Context, ID int) (entity.OauthApplicationJSON, *entity.Error) {
// 	oauthApplication, err := am.oauthApplicationModel.One(kontext.Fabricate(kontext.WithDefaultContext(ctx)), ID, am.sqldb)
// 	if err != nil {
// 		if err.Type() == exception.NotFound {
// 			return entity.OauthApplicationJSON{}, &entity.Error{
// 				HttpStatus: http.StatusNotFound,
// 				Errors:     eobject.Wrap(eobject.NotFoundError(ctx, "oauth_application")),
// 			}
// 		}

// 		log.Error().
// 			Err(err).
// 			Stack().
// 			Int("id", ID).
// 			Array("tags", zerolog.Arr().Str("service").Str("application_manager").Str("one").Str("model_one")).
// 			Msg("Error when fetching single oauth application")

// 		return entity.OauthApplicationJSON{}, &entity.Error{
// 			HttpStatus: http.StatusInternalServerError,
// 			Errors:     eobject.Wrap(eobject.InternalServerError(ctx)),
// 		}
// 	}

// 	formattedResult := am.formatter.Application(ctx, oauthApplication)
// 	return formattedResult, nil
// }
