package service

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/codefluence-x/altair/core"
	"github.com/codefluence-x/altair/entity"
	"github.com/codefluence-x/altair/eobject"
	"github.com/codefluence-x/journal"
)

type authorization struct {
	oauthApplicationModel core.OauthApplicationModel
	modelFormatter        core.ModelFormater
	oauthValidator        core.OauthValidator
}

func Authorization(
	oauthApplicationModel core.OauthApplicationModel,
	modelFormatter core.ModelFormater,
	oauthValidator core.OauthValidator,
) core.Authorization {
	return &authorization{
		oauthApplicationModel: oauthApplicationModel,
		modelFormatter:        modelFormatter,
		oauthValidator:        oauthValidator,
	}
}

func (a *authorization) Grantor(ctx context.Context, authorizationReq entity.AuthorizationRequestJSON) (interface{}, *entity.Error) {
	if *authorizationReq.ResponseType == "token" {
		return a.GrantToken(ctx, authorizationReq)
	} else if *authorizationReq.ResponseType == "code" {
		return a.Grant(ctx, authorizationReq)
	}

	err := &entity.Error{
		HttpStatus: http.StatusUnprocessableEntity,
		Errors:     eobject.Wrap(eobject.ValidationError("response_type is invalid. Should be either `token` or `code`.")),
	}

	journal.Error("invalid response type sent by client", err).
		AddField("request", authorizationReq).
		SetTags("service", "authorization", "grantor").
		Log()

	return nil, err
}

// WIP
func (a *authorization) Grant(ctx context.Context, authorizationReq entity.AuthorizationRequestJSON) (entity.OauthAccessGrantJSON, *entity.Error) {
	var oauthAccessGrantJSON entity.OauthAccessGrantJSON

	_, err := a.findAndValidateApplication(ctx, authorizationReq)
	if err != nil {
		return oauthAccessGrantJSON, err
	}

	return oauthAccessGrantJSON, nil
}

// WIP
func (a *authorization) GrantToken(ctx context.Context, authorizationReq entity.AuthorizationRequestJSON) (entity.OauthAccessTokenJSON, *entity.Error) {
	var oauthAccessTokenJSON entity.OauthAccessTokenJSON

	_, err := a.findAndValidateApplication(ctx, authorizationReq)
	if err != nil {
		return oauthAccessTokenJSON, err
	}

	return oauthAccessTokenJSON, nil
}

func (a *authorization) findAndValidateApplication(ctx context.Context, authorizationReq entity.AuthorizationRequestJSON) (entity.OauthApplication, *entity.Error) {
	oauthApplication, err := a.oauthApplicationModel.OneByUIDandSecret(ctx, *authorizationReq.ClientUID, *authorizationReq.ClientSecret)
	if err != nil {

		journal.Error("application cannot be found because there was an error", err).
			AddField("request", authorizationReq).
			SetTags("service", "authorization", "find_secret").
			Log()

		if err == sql.ErrNoRows {
			return entity.OauthApplication{}, &entity.Error{
				HttpStatus: http.StatusNotFound,
				Errors:     eobject.Wrap(eobject.NotFoundError(ctx, "client_uid & client_secret")),
			}
		}

		return entity.OauthApplication{}, &entity.Error{
			HttpStatus: http.StatusInternalServerError,
			Errors:     eobject.Wrap(eobject.InternalServerError(ctx)),
		}
	}

	if err := a.oauthValidator.ValidateAuthorizationGrant(ctx, authorizationReq, oauthApplication); err != nil {
		return entity.OauthApplication{}, err
	}

	return oauthApplication, nil
}

// WIP
func (a *authorization) Token(ctx context.Context, accessTokenReq entity.AccessTokenRequestJSON) (entity.OauthAccessTokenJSON, *entity.Error) {
	var oauthAccessTokenJSON entity.OauthAccessTokenJSON

	return oauthAccessTokenJSON, nil
}
