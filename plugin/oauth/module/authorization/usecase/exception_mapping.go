package usecase

import (
	"net/http"

	"github.com/kodefluence/monorepo/exception"
	"github.com/kodefluence/monorepo/jsonapi"
	"github.com/kodefluence/monorepo/kontext"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func (a *Authorization) exceptionMapping(ktx kontext.Context, exc exception.Exception, tags *zerolog.Array) jsonapi.Errors {
	log.Error().
		Err(exc).
		Stack().
		Interface("request_id", ktx.GetWithoutCheck("request_id")).
		Array("tags", tags).
		Msg(exc.Detail())

	switch exc.Type() {
	case exception.NotFound:
		return jsonapi.BuildResponse(
			jsonapi.WithException(
				"ERR0404",
				http.StatusNotFound,
				exc,
			),
		).Errors
	case exception.Forbidden:
		return jsonapi.BuildResponse(
			jsonapi.WithException(
				"ERR0403",
				http.StatusForbidden,
				exc,
			),
		).Errors
	default:
		return jsonapi.BuildResponse(a.apiError.InternalServerError(ktx)).Errors
	}
}
