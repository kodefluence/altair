package usecase

import (
	"fmt"
	"strings"

	"github.com/kodefluence/altair/plugin/oauth/entity"
	"github.com/kodefluence/altair/util"
	"github.com/kodefluence/monorepo/jsonapi"
	"github.com/kodefluence/monorepo/kontext"
)

func (a *Authorization) ValidateAuthorizationGrant(ktx kontext.Context, r entity.AuthorizationRequestJSON, application entity.OauthApplication) jsonapi.Errors {
	var errorOptions []jsonapi.Option

	if r.ResponseType == nil {
		errorOptions = append(errorOptions, a.apiError.ValidationError("response_type can't be empty"))
	}

	if r.ResourceOwnerID == nil {
		errorOptions = append(errorOptions, a.apiError.ValidationError("resource_owner_id can't be empty"))
	}

	if r.RedirectURI == nil {
		errorOptions = append(errorOptions, a.apiError.ValidationError("redirect_uri can't be empty"))
	}

	if len(errorOptions) > 0 {
		return jsonapi.BuildResponse(errorOptions...).Errors
	}

	if r.Scopes == nil {
		r.Scopes = util.StringToPointer("")
	}

	requestScopes := strings.Fields(*r.Scopes)
	applicationScopes := strings.Fields(application.Scopes.String)

	var invalidScope []string

	for _, rs := range requestScopes {

		scopeNotExists := true

		for _, as := range applicationScopes {
			if rs == as {
				scopeNotExists = false
				break
			}
		}

		if scopeNotExists {
			invalidScope = append(invalidScope, rs)
		}
	}

	if len(invalidScope) > 0 {
		return jsonapi.BuildResponse(
			a.apiError.ForbiddenError(ktx, "application", fmt.Sprintf("your requested scopes `(%v)` is not exists in application", invalidScope)),
		).Errors
	}

	if *r.ResponseType == "token" && application.OwnerType != "confidential" {
		return jsonapi.BuildResponse(
			a.apiError.ForbiddenError(ktx, "access_token", "your response type is not allowed in this application"),
		).Errors
	}

	return nil
}
