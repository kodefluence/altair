package usecase

import (
	"github.com/kodefluence/altair/plugin/oauth/entity"
	"github.com/kodefluence/monorepo/jsonapi"
	"github.com/kodefluence/monorepo/kontext"
)

func (a *Authorization) Grantor(ktx kontext.Context, authorizationReq entity.AuthorizationRequestJSON) (interface{}, jsonapi.Errors) {
	if authorizationReq.ResponseType == nil {
		// 	return nil, &entity.Error{
		// 		HttpStatus: http.StatusUnprocessableEntity,
		// 		Errors:     eobject.Wrap(eobject.ValidationError("response_type cannot be empty")),
		// 	}
	}

	switch *authorizationReq.ResponseType {
	case "token":
		// Implicit grant token in here
	case "code":
		// Grant code in here
	}

	// err := &entity.Error{
	// 	HttpStatus: http.StatusUnprocessableEntity,
	// 	Errors:     eobject.Wrap(eobject.ValidationError("response_type is invalid. Should be either `token` or `code`.")),
	// }

	// log.Error().
	// 	Err(err).
	// 	Stack().
	// 	Interface("request_id", ctx.Value("request_id")).
	// 	Interface("request", authorizationReq).
	// 	Array("tags", zerolog.Arr().Str("service").Str("authorization").Str("grantor")).
	// 	Msg("invalid response type sent by client")
	// return nil, err

	return nil, nil
}
