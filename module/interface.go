package module

import (
	"github.com/kodefluence/monorepo/jsonapi"
	"github.com/kodefluence/monorepo/kontext"
)

type ApiError interface {
	InternalServerError(ktx kontext.Context) jsonapi.Option
	BadRequestError(in string) jsonapi.Option
	NotFoundError(ktx kontext.Context, entityType string) jsonapi.Option
	UnauthorizedError() jsonapi.Option
	ForbiddenError(ktx kontext.Context, entityType, reason string) jsonapi.Option
	ValidationError(msg string) jsonapi.Option
}
