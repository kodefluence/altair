package validator

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/codefluence-x/altair/provider/plugin/oauth/entity"
	"github.com/codefluence-x/altair/provider/plugin/oauth/eobject"
	"github.com/codefluence-x/altair/util"
)

type Oauth struct {
	refreshTokenToggle bool
}

// NewOauth create Oauth struct for validation
func NewOauth(refreshTokenToggle bool) *Oauth {
	return &Oauth{
		refreshTokenToggle: refreshTokenToggle,
	}
}

// ValidateApplication will validate oauth application json
func (a *Oauth) ValidateApplication(ctx context.Context, data entity.OauthApplicationJSON) *entity.Error {
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

// ValidateAuthorizationGrant will validate authorization grant request
func (a *Oauth) ValidateAuthorizationGrant(ctx context.Context, r entity.AuthorizationRequestJSON, application entity.OauthApplication) *entity.Error {
	var entityErr = &entity.Error{}

	if r.ResponseType == nil {
		entityErr.Errors = append(entityErr.Errors, eobject.ValidationError(`response_type can't be empty`))
	}

	if r.ResourceOwnerID == nil {
		entityErr.Errors = append(entityErr.Errors, eobject.ValidationError(`resource_owner_id can't be empty`))
	}

	if r.RedirectURI == nil {
		entityErr.Errors = append(entityErr.Errors, eobject.ValidationError(`redirect_uri can't be empty`))
	}

	if len(entityErr.Errors) > 0 {
		entityErr.HttpStatus = http.StatusUnprocessableEntity
		return entityErr
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
		return &entity.Error{
			HttpStatus: http.StatusForbidden,
			Errors:     eobject.Wrap(eobject.ForbiddenError(ctx, "application", fmt.Sprintf("your requested scopes `(%v)` is not exists in application", invalidScope))),
		}
	}

	if *r.ResponseType == "token" && application.OwnerType != "confidential" {
		return &entity.Error{
			HttpStatus: http.StatusForbidden,
			Errors:     eobject.Wrap(eobject.ForbiddenError(ctx, "access_token", "your response type is not allowed in this application")),
		}
	}

	return nil
}

// ValidateTokenGrant will validate token grant
func (a *Oauth) ValidateTokenGrant(ctx context.Context, r entity.AccessTokenRequestJSON) *entity.Error {
	var entityErr = &entity.Error{}

	if r.GrantType == nil {
		entityErr.HttpStatus = http.StatusUnprocessableEntity
		entityErr.Errors = append(entityErr.Errors, eobject.ValidationError(`grant_type can't be empty`))
		return entityErr
	}

	switch *r.GrantType {
	case "authorization_code":
		if r.Code == nil {
			entityErr.Errors = append(entityErr.Errors, eobject.ValidationError(`code can't be empty`))
		}

		if r.RedirectURI == nil {
			entityErr.Errors = append(entityErr.Errors, eobject.ValidationError(`redirect_uri can't be empty`))
		}
	case "refresh_token":
		if a.refreshTokenToggle {
			if r.RefreshToken == nil {
				entityErr.Errors = append(entityErr.Errors, eobject.ValidationError(`refresh token can't be empty`))
			}
		}
	default:
		entityErr.Errors = append(entityErr.Errors, eobject.ValidationError(`grant_type is not a valid value`))
	}

	if len(entityErr.Errors) > 0 {
		entityErr.HttpStatus = http.StatusUnprocessableEntity
		return entityErr
	}

	return nil
}

// ValidateTokenAuthorizationCode will validate oauth access grant
func (a *Oauth) ValidateTokenAuthorizationCode(ctx context.Context, r entity.AccessTokenRequestJSON, data entity.OauthAccessGrant) *entity.Error {
	if data.RevokedAT.Valid {
		return &entity.Error{
			HttpStatus: http.StatusForbidden,
			Errors:     eobject.Wrap(eobject.ForbiddenError(ctx, "access_token", "authorization code already used")),
		}
	}

	if time.Now().After(data.ExpiresIn) {
		return &entity.Error{
			HttpStatus: http.StatusForbidden,
			Errors:     eobject.Wrap(eobject.ForbiddenError(ctx, "access_token", "authorization code already expired")),
		}
	}

	if data.RedirectURI.String != *r.RedirectURI {
		return &entity.Error{
			HttpStatus: http.StatusForbidden,
			Errors:     eobject.Wrap(eobject.ForbiddenError(ctx, "access_token", "redirect uri is different from one that generated before")),
		}
	}

	return nil
}

// ValidateTokenRefreshToken will validate refresh token
func (a *Oauth) ValidateTokenRefreshToken(ctx context.Context, data entity.OauthRefreshToken) *entity.Error {

	if data.RevokedAT.Valid {
		return &entity.Error{
			HttpStatus: http.StatusForbidden,
			Errors:     eobject.Wrap(eobject.ForbiddenError(ctx, "access_token", "refresh token already used")),
		}
	}

	return nil
}
