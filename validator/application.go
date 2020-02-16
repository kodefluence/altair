package validator

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/codefluence-x/altair/core"
	"github.com/codefluence-x/altair/entity"
	"github.com/codefluence-x/altair/eobject"
)

type application struct {
}

func Oauth() core.OauthValidator {
	return &application{}
}

func (a *application) ValidateApplication(ctx context.Context, data entity.OauthApplicationJSON) *entity.Error {
	var entityError = &entity.Error{}

	if data.OwnerType == nil {
		entityError.Errors = append(entityError.Errors, eobject.ValidationError("object `owner_type` is nil or not exists"))
	} else {
		if *data.OwnerType != "confidential" && *data.OwnerType != "public" {
			entityError.Errors = append(entityError.Errors, eobject.ValidationError("object `owner_type` must be either of `confidential` or `public`"))
		}
	}

	if len(entityError.Errors) > 0 {
		entityError.HttpStatus = http.StatusUnprocessableEntity
		return entityError
	}

	return nil
}

func (a *application) ValidateAuthorizationGrant(ctx context.Context, r entity.AuthorizationRequestJSON, application entity.OauthApplication) *entity.Error {
	requestScopes := strings.Fields(*r.Scopes)
	applicationScopes := strings.Fields(application.Scopes)

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
		return &entity.Error{
			HttpStatus: http.StatusForbidden,
			Errors:     eobject.Wrap(eobject.ForbiddenError(ctx, "application", fmt.Sprintf("your requested scopes `(%v)` is not exists in application", invalidScope))),
		}
	}
	return nil
}
