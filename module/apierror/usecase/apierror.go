package usecase

import (
	"fmt"
	"net/http"

	"github.com/kodefluence/monorepo/exception"
	"github.com/kodefluence/monorepo/jsonapi"
	"github.com/kodefluence/monorepo/kontext"
)

type ApiError struct{}

func NewApiError() *ApiError {
	return &ApiError{}
}

func (*ApiError) InternalServerError(ktx kontext.Context) jsonapi.Option {
	err := fmt.Errorf("Something is not right, help us fix this problem. Contribute to https://github.com/kodefluence/altair. Tracing code: '%v'", ktx.GetWithoutCheck("request_id"))
	return jsonapi.WithException(
		"ERR0500",
		http.StatusInternalServerError,
		exception.Throw(
			err,
			exception.WithTitle("Internal server error"),
			exception.WithDetail(err.Error()),
			exception.WithType(exception.Unexpected),
		),
	)
}

func (*ApiError) BadRequestError(in string) jsonapi.Option {
	err := fmt.Errorf("You've send malformed request in your `%s`", in)
	return jsonapi.WithException(
		"ERR0400",
		http.StatusBadRequest,
		exception.Throw(
			err,
			exception.WithTitle("Bad request error"),
			exception.WithDetail(err.Error()),
			exception.WithType(exception.BadInput),
		),
	)
}

func (*ApiError) NotFoundError(ktx kontext.Context, entityType string) jsonapi.Option {
	err := fmt.Errorf("Resource of `%s` is not found. Tracing code: `%v`", entityType, ktx.GetWithoutCheck("request_id"))
	return jsonapi.WithException(
		"ERR0404",
		http.StatusNotFound,
		exception.Throw(
			err,
			exception.WithTitle("Not found error"),
			exception.WithDetail(err.Error()),
			exception.WithType(exception.NotFound),
		),
	)
}

func (*ApiError) UnauthorizedError() jsonapi.Option {
	err := fmt.Errorf("You are unauthorized")
	return jsonapi.WithException(
		"ERR0401",
		http.StatusUnauthorized,
		exception.Throw(
			err,
			exception.WithTitle("Unauthorized error"),
			exception.WithDetail(err.Error()),
			exception.WithType(exception.Unauthorized),
		),
	)
}

func (*ApiError) ForbiddenError(ktx kontext.Context, entityType, reason string) jsonapi.Option {
	err := fmt.Errorf("Resource of `%s` is forbidden to be accessed, because of: %s. Tracing code: `%v`", entityType, reason, ktx.GetWithoutCheck("request_id"))
	return jsonapi.WithException(
		"ERR0403",
		http.StatusForbidden,
		exception.Throw(
			err,
			exception.WithTitle("Forbidden error"),
			exception.WithDetail(err.Error()),
			exception.WithType(exception.Forbidden),
		),
	)
}

func (*ApiError) ValidationError(msg string) jsonapi.Option {
	err := fmt.Errorf(fmt.Sprintf("Validation error because of: %s", msg))
	return jsonapi.WithException(
		"ERR1442",
		http.StatusUnprocessableEntity,
		exception.Throw(
			err,
			exception.WithTitle("Validation error"),
			exception.WithDetail(err.Error()),
			exception.WithType(exception.BadInput),
		),
	)
}
