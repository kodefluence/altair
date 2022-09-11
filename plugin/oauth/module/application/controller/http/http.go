package http

import (
	"github.com/kodefluence/altair/plugin/oauth/entity"
	"github.com/kodefluence/monorepo/jsonapi"
	"github.com/kodefluence/monorepo/kontext"
)

//go:generate mockgen -destination ./mock/mock.go -package mock -source ./http.go

type ApplicationManager interface {
	Create(ktx kontext.Context, e entity.OauthApplicationJSON) (entity.OauthApplicationJSON, jsonapi.Errors)
	List(ktx kontext.Context, offset, limit int) ([]entity.OauthApplicationJSON, int, jsonapi.Errors)
	One(ktx kontext.Context, ID int) (entity.OauthApplicationJSON, jsonapi.Errors)
	Update(ktx kontext.Context, ID int, e entity.OauthApplicationUpdateJSON) (entity.OauthApplicationJSON, jsonapi.Errors)
}

type ApiError interface {
	// InternalServerError(ktx kontext.Context) jsonapi.Option
	BadRequestError(in string) jsonapi.Option
	// NotFoundError(ktx kontext.Context, entityType string) jsonapi.Option
	// UnauthorizedError() jsonapi.Option
	// ForbiddenError(ktx kontext.Context, entityType, reason string) jsonapi.Option
	// ValidationError(msg string) jsonapi.Option
}
