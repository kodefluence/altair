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
	Application(application entity.OauthApplication) entity.OauthApplicationJSON
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
