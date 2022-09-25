package usecase

import (
	"github.com/kodefluence/altair/plugin/oauth/entity"
	"github.com/kodefluence/monorepo/jsonapi"
)

func (a *Authorization) ValidateTokenGrant(r entity.AccessTokenRequestJSON) jsonapi.Errors {
	var errorOptions []jsonapi.Option

	if r.GrantType == nil {
		return jsonapi.BuildResponse(a.apiError.ValidationError(`grant_type can't be empty`)).Errors
	}

	switch *r.GrantType {
	case "authorization_code":
		if r.Code == nil {
			errorOptions = append(errorOptions, a.apiError.ValidationError(`code is not valid value`))
		}

		if r.RedirectURI == nil {
			errorOptions = append(errorOptions, a.apiError.ValidationError(`redirect_uri is not valid value`))
		}
	case "refresh_token":
		if r.RefreshToken == nil {
			errorOptions = append(errorOptions, a.apiError.ValidationError(`refresh_token is not valid value`))
		}
	case "client_credentials":
		// No validations, since client_uid and client_secret validation already validated before
	default:
		errorOptions = append(errorOptions, a.apiError.ValidationError(`grant_type is not valid value`))
	}

	if len(errorOptions) > 0 {
		return jsonapi.BuildResponse(errorOptions...).Errors
	}

	return nil
}
