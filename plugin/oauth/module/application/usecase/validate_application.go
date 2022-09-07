package usecase

import (
	"github.com/kodefluence/altair/plugin/oauth/entity"
	"github.com/kodefluence/monorepo/jsonapi"
)

func (am *ApplicationManager) ValidateApplication(data entity.OauthApplicationJSON) jsonapi.Errors {
	var errorOptions []jsonapi.Option

	if data.OwnerType == nil {
		errorOptions = append(errorOptions, am.apiError.ValidationError("object `owner_type` is nil or not exists"))
	} else {
		if *data.OwnerType != "confidential" && *data.OwnerType != "public" {
			errorOptions = append(errorOptions, am.apiError.ValidationError("object `owner_type` must be either of `confidential` or `public`"))
		}
	}

	if len(errorOptions) > 0 {
		return jsonapi.BuildResponse(errorOptions...).Errors
	}

	return nil
}
