package validator

import (
	"context"
	"net/http"

	"github.com/codefluence-x/altair/core"
	"github.com/codefluence-x/altair/entity"
	"github.com/codefluence-x/altair/eobject"
)

type application struct {
}

func Application() core.OauthApplicationValidator {
	return application{}
}

func (a application) ValidateCreate(ctx context.Context, data entity.OauthApplicationJSON) *entity.Error {
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
