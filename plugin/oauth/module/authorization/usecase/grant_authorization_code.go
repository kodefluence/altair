package usecase

import (
	"github.com/kodefluence/altair/plugin/oauth/entity"
	"github.com/kodefluence/monorepo/jsonapi"
	"github.com/kodefluence/monorepo/kontext"
)

func (a *Authorization) GrantAuthorizationCode(ktx kontext.Context, authorizationReq entity.AuthorizationRequestJSON) (interface{}, jsonapi.Errors) {
	if authorizationReq.ResponseType == nil {
		return nil, jsonapi.BuildResponse(
			a.apiError.ValidationError("response_type cannot be empty"),
		).Errors
	}

	switch *authorizationReq.ResponseType {
	case "token":
		if a.config.Config.ImplicitGrant.Active == false {
			break
		}

		return a.ImplicitGrant(ktx, authorizationReq)
	case "code":
		return a.Grant(ktx, authorizationReq)
	}

	return nil, jsonapi.BuildResponse(
		a.apiError.ValidationError("response_type is invalid. Should be either `token` or `code`"),
	).Errors
}
