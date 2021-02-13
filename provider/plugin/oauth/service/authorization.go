package service

import (
	"context"
	"net/http"

	"github.com/codefluence-x/altair/provider/plugin/oauth/entity"
	"github.com/codefluence-x/altair/provider/plugin/oauth/eobject"
	"github.com/codefluence-x/altair/provider/plugin/oauth/interfaces"
	"github.com/codefluence-x/monorepo/db"
	"github.com/codefluence-x/monorepo/exception"
	"github.com/codefluence-x/monorepo/kontext"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Authorization struct handle all of things related to oauth2 authorization
type Authorization struct {
	oauthApplicationModel  interfaces.OauthApplicationModel
	oauthAccessTokenModel  interfaces.OauthAccessTokenModel
	oauthAccessGrantModel  interfaces.OauthAccessGrantModel
	oauthRefreshTokenModel interfaces.OauthRefreshTokenModel

	oauthValidator interfaces.OauthValidator

	modelFormatter interfaces.ModelFormater
	oauthFormatter interfaces.OauthFormatter

	refreshTokenToggle bool

	sqldb db.DB
}

// NewAuthorization create new service to handler authorize related flow
func NewAuthorization(
	oauthApplicationModel interfaces.OauthApplicationModel,
	oauthAccessTokenModel interfaces.OauthAccessTokenModel,
	oauthAccessGrantModel interfaces.OauthAccessGrantModel,
	oauthRefreshTokenModel interfaces.OauthRefreshTokenModel,
	modelFormatter interfaces.ModelFormater,
	oauthValidator interfaces.OauthValidator,
	oauthFormatter interfaces.OauthFormatter,
	refreshTokenToggle bool,
	sqldb db.DB,
) *Authorization {
	return &Authorization{
		oauthApplicationModel:  oauthApplicationModel,
		oauthAccessTokenModel:  oauthAccessTokenModel,
		oauthAccessGrantModel:  oauthAccessGrantModel,
		oauthRefreshTokenModel: oauthRefreshTokenModel,
		modelFormatter:         modelFormatter,
		oauthValidator:         oauthValidator,
		oauthFormatter:         oauthFormatter,
		refreshTokenToggle:     refreshTokenToggle,
		sqldb:                  sqldb,
	}
}

// Grantor provide granting logic for authorization request
func (a *Authorization) Grantor(ctx context.Context, authorizationReq entity.AuthorizationRequestJSON) (interface{}, *entity.Error) {
	if authorizationReq.ResponseType == nil {
		return nil, &entity.Error{
			HttpStatus: http.StatusUnprocessableEntity,
			Errors:     eobject.Wrap(eobject.ValidationError("response_type cannot be empty")),
		}
	}

	switch *authorizationReq.ResponseType {
	case "token":
		return a.GrantToken(ctx, authorizationReq)
	case "code":
		return a.Grant(ctx, authorizationReq)
	}

	err := &entity.Error{
		HttpStatus: http.StatusUnprocessableEntity,
		Errors:     eobject.Wrap(eobject.ValidationError("response_type is invalid. Should be either `token` or `code`.")),
	}

	log.Error().
		Err(err).
		Stack().
		Interface("request_id", ctx.Value("request_id")).
		Interface("request", authorizationReq).
		Array("tags", zerolog.Arr().Str("service").Str("authorization").Str("grantor")).
		Msg("invalid response type sent by client")
	return nil, err
}

// Grant authorization an access code
func (a *Authorization) Grant(ctx context.Context, authorizationReq entity.AuthorizationRequestJSON) (entity.OauthAccessGrantJSON, *entity.Error) {
	var oauthAccessGrantJSON entity.OauthAccessGrantJSON

	oauthApplication, entityErr := a.findAndValidateApplication(ctx, authorizationReq.ClientUID, authorizationReq.ClientSecret)
	if entityErr != nil {
		return oauthAccessGrantJSON, entityErr
	}

	if err := a.oauthValidator.ValidateAuthorizationGrant(ctx, authorizationReq, oauthApplication); err != nil {
		return oauthAccessGrantJSON, err
	}

	exc := a.sqldb.Transaction(kontext.Fabricate(kontext.WithDefaultContext(ctx)), "authorization-grant-authorization-code", func(tx db.TX) exception.Exception {
		id, err := a.oauthAccessGrantModel.Create(kontext.Fabricate(kontext.WithDefaultContext(ctx)), a.modelFormatter.AccessGrantFromAuthorizationRequest(authorizationReq, oauthApplication), tx)
		if err != nil {
			return exception.Throw(err, exception.WithDetail("error creating authorization code"))
		}

		oauthAccessGrant, err := a.oauthAccessGrantModel.One(kontext.Fabricate(kontext.WithDefaultContext(ctx)), id, tx)
		if err != nil {
			return exception.Throw(err, exception.WithDetail("error selecting newly created authorization code"))
		}

		oauthAccessGrantJSON = a.oauthFormatter.AccessGrant(oauthAccessGrant)
		return nil
	})
	if exc != nil {
		return oauthAccessGrantJSON, a.exceptionMapping(ctx, exc, zerolog.Arr().Str("service").Str("authorization").Str("grant"))
	}

	return oauthAccessGrantJSON, nil
}

// GrantToken will grant an access token
func (a *Authorization) GrantToken(ctx context.Context, authorizationReq entity.AuthorizationRequestJSON) (entity.OauthAccessTokenJSON, *entity.Error) {
	var finalOauthTokenJSON entity.OauthAccessTokenJSON
	oauthApplication, entityErr := a.findAndValidateApplication(ctx, authorizationReq.ClientUID, authorizationReq.ClientSecret)
	if entityErr != nil {
		return entity.OauthAccessTokenJSON{}, entityErr
	}

	if err := a.oauthValidator.ValidateAuthorizationGrant(ctx, authorizationReq, oauthApplication); err != nil {
		return entity.OauthAccessTokenJSON{}, err
	}

	exc := a.sqldb.Transaction(kontext.Fabricate(kontext.WithDefaultContext(ctx)), "authorization-grant-token", func(tx db.TX) exception.Exception {
		id, err := a.oauthAccessTokenModel.Create(kontext.Fabricate(kontext.WithDefaultContext(ctx)), a.modelFormatter.AccessTokenFromAuthorizationRequest(authorizationReq, oauthApplication), tx)
		if err != nil {
			return exception.Throw(err)
		}

		oauthAccessToken, err := a.oauthAccessTokenModel.One(kontext.Fabricate(kontext.WithDefaultContext(ctx)), id, tx)
		if err != nil {
			return exception.Throw(err)
		}

		var oauthRefreshTokenJSON *entity.OauthRefreshTokenJSON

		if a.refreshTokenToggle {
			refreshTokenID, err := a.oauthRefreshTokenModel.Create(kontext.Fabricate(kontext.WithDefaultContext(ctx)), a.modelFormatter.RefreshToken(oauthApplication, oauthAccessToken), tx)
			if err != nil {
				return exception.Throw(err)
			}

			oauthRefreshToken, err := a.oauthRefreshTokenModel.One(kontext.Fabricate(kontext.WithDefaultContext(ctx)), refreshTokenID, tx)
			if err != nil {
				return exception.Throw(err)
			}

			newOauthRefreshTokenJSON := a.oauthFormatter.RefreshToken(oauthRefreshToken)
			oauthRefreshTokenJSON = &newOauthRefreshTokenJSON
		}

		finalOauthTokenJSON = a.oauthFormatter.AccessToken(oauthAccessToken, *authorizationReq.RedirectURI, oauthRefreshTokenJSON)

		return nil
	})
	if exc != nil {
		return entity.OauthAccessTokenJSON{}, a.exceptionMapping(ctx, exc, zerolog.Arr().Str("service").Str("authorization").Str("grant_token"))
	}

	return finalOauthTokenJSON, nil
}

func (a *Authorization) findAndValidateApplication(ctx context.Context, clientUID, clientSecret *string) (entity.OauthApplication, *entity.Error) {
	if clientUID == nil {
		return entity.OauthApplication{}, &entity.Error{
			HttpStatus: http.StatusUnprocessableEntity,
			Errors:     eobject.Wrap(eobject.ValidationError("client_uid cannot be empty")),
		}
	}

	if clientSecret == nil {
		return entity.OauthApplication{}, &entity.Error{
			HttpStatus: http.StatusUnprocessableEntity,
			Errors:     eobject.Wrap(eobject.ValidationError("client_secret cannot be empty")),
		}
	}

	oauthApplication, err := a.oauthApplicationModel.OneByUIDandSecret(kontext.Fabricate(kontext.WithDefaultContext(ctx)), *clientUID, *clientSecret, a.sqldb)
	if err != nil {
		log.Error().
			Err(err).
			Stack().
			Interface("request_id", ctx.Value("request_id")).
			Str("client_uid", *clientUID).
			Array("tags", zerolog.Arr().Str("service").Str("authorization").Str("find_secret")).
			Msg("application cannot be found because there was an error")
		if err.Type() == exception.NotFound {
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

	return oauthApplication, nil
}

// Token will grant a token from authorization code
func (a *Authorization) Token(ctx context.Context, accessTokenReq entity.AccessTokenRequestJSON) (entity.OauthAccessTokenJSON, *entity.Error) {
	oauthApplication, entityErr := a.findAndValidateApplication(ctx, accessTokenReq.ClientUID, accessTokenReq.ClientSecret)
	if entityErr != nil {
		return entity.OauthAccessTokenJSON{}, entityErr
	}

	if entityErr := a.oauthValidator.ValidateTokenGrant(ctx, accessTokenReq); entityErr != nil {
		return entity.OauthAccessTokenJSON{}, entityErr
	}

	switch *accessTokenReq.GrantType {

	case "authorization_code":
		oauthAccessToken, redirectURI, entityErr := a.grantTokenFromAuthorizationCode(ctx, accessTokenReq, oauthApplication)
		if entityErr != nil {
			return entity.OauthAccessTokenJSON{}, entityErr
		}

		return a.oauthFormatter.AccessToken(oauthAccessToken, redirectURI, nil), nil
	case "refresh_token":
		if a.refreshTokenToggle {
			oauthAccessToken, oauthRefreshToken, entityErr := a.grantTokenFromRefreshToken(ctx, accessTokenReq, oauthApplication)
			if entityErr != nil {
				return entity.OauthAccessTokenJSON{}, entityErr
			}

			oauthRefreshTokenJSON := a.oauthFormatter.RefreshToken(oauthRefreshToken)
			return a.oauthFormatter.AccessToken(oauthAccessToken, *accessTokenReq.RedirectURI, &oauthRefreshTokenJSON), nil
		}
	}

	return entity.OauthAccessTokenJSON{}, &entity.Error{
		HttpStatus: http.StatusUnprocessableEntity,
		Errors:     eobject.Wrap(eobject.ValidationError(`grant_type can't be empty`)),
	}
}

func (a *Authorization) grantTokenFromRefreshToken(ctx context.Context, accessTokenReq entity.AccessTokenRequestJSON, oauthApplication entity.OauthApplication) (entity.OauthAccessToken, entity.OauthRefreshToken, *entity.Error) {
	var finalOauthAccessToken entity.OauthAccessToken
	var finalOauthRefreshToken entity.OauthRefreshToken

	exc := a.sqldb.Transaction(kontext.Fabricate(kontext.WithDefaultContext(ctx)), "authorization-grant-token-from-refresh-token", func(tx db.TX) exception.Exception {
		oldOauthRefreshToken, err := a.oauthRefreshTokenModel.OneByToken(kontext.Fabricate(kontext.WithDefaultContext(ctx)), *accessTokenReq.RefreshToken, tx)
		if err != nil {
			if err.Type() == exception.NotFound {
				errorObject := eobject.NotFoundError(ctx, "refresh_token")
				return exception.Throw(err, exception.WithType(exception.NotFound), exception.WithDetail(errorObject.Message), exception.WithTitle(errorObject.Code))
			}

			return exception.Throw(err, exception.WithDetail("refresh token cannot be found because there was an error"))
		}

		if err := a.oauthValidator.ValidateTokenRefreshToken(ctx, oldOauthRefreshToken); err != nil {
			return err
		}

		oldAccessToken, err := a.oauthAccessTokenModel.One(kontext.Fabricate(kontext.WithDefaultContext(ctx)), oldOauthRefreshToken.OauthAccessTokenID, tx)
		if err != nil {
			if err.Type() == exception.NotFound {
				return exception.Throw(err, exception.WithType(exception.Unauthorized), exception.WithDetail("access token cannot be found because there was an error"))
			}

			return exception.Throw(err, exception.WithDetail("access token cannot be found because there was an error"))
		}

		oauthAccessTokenID, err := a.oauthAccessTokenModel.Create(kontext.Fabricate(kontext.WithDefaultContext(ctx)), a.modelFormatter.AccessTokenFromOauthRefreshToken(oauthApplication, oldAccessToken), tx)
		if err != nil {
			return exception.Throw(err, exception.WithDetail("error creating access token"))
		}

		oauthAccessToken, err := a.oauthAccessTokenModel.One(kontext.Fabricate(kontext.WithDefaultContext(ctx)), oauthAccessTokenID, tx)
		if err != nil {
			return exception.Throw(err, exception.WithDetail("error when selecting newly created access token"))
		}

		err = a.oauthRefreshTokenModel.Revoke(kontext.Fabricate(kontext.WithDefaultContext(ctx)), *accessTokenReq.RefreshToken, tx)
		if err != nil {
			return exception.Throw(err, exception.WithDetail("error revoke refresh token"))
		}

		oauthRefreshTokenID, err := a.oauthRefreshTokenModel.Create(kontext.Fabricate(kontext.WithDefaultContext(ctx)), a.modelFormatter.RefreshToken(oauthApplication, oauthAccessToken), tx)
		if err != nil {
			return exception.Throw(err, exception.WithDetail("error creating refresh token"))
		}

		oauthRefreshToken, err := a.oauthRefreshTokenModel.One(kontext.Fabricate(kontext.WithDefaultContext(ctx)), oauthRefreshTokenID, tx)
		if err != nil {
			return exception.Throw(err, exception.WithDetail("error when selecting newly created refresh token"))
		}

		finalOauthAccessToken = oauthAccessToken
		finalOauthRefreshToken = oauthRefreshToken

		return nil
	})

	if exc != nil {
		return entity.OauthAccessToken{}, entity.OauthRefreshToken{}, a.exceptionMapping(ctx, exc, zerolog.Arr().Str("service").Str("authorization").Str("refresh_token"))
	}

	return finalOauthAccessToken, finalOauthRefreshToken, nil
}

func (a *Authorization) grantTokenFromAuthorizationCode(ctx context.Context, accessTokenReq entity.AccessTokenRequestJSON, oauthApplication entity.OauthApplication) (entity.OauthAccessToken, string, *entity.Error) {
	oauthAccessGrant, err := a.oauthAccessGrantModel.OneByCode(kontext.Fabricate(kontext.WithDefaultContext(ctx)), *accessTokenReq.Code, a.sqldb)
	if err != nil {
		log.Error().
			Err(err).
			Stack().
			Interface("request_id", ctx.Value("request_id")).
			Array("tags", zerolog.Arr().Str("service").Str("authorization").Str("authorization_code").Str("one_by_code")).
			Msg("authorization code cannot be found because there was an error")
		if err.Type() == exception.NotFound {
			return entity.OauthAccessToken{}, "", &entity.Error{
				HttpStatus: http.StatusNotFound,
				Errors:     eobject.Wrap(eobject.NotFoundError(ctx, "authorization_code")),
			}
		}

		return entity.OauthAccessToken{}, oauthAccessGrant.RedirectURI.String, &entity.Error{
			HttpStatus: http.StatusInternalServerError,
			Errors:     eobject.Wrap(eobject.InternalServerError(ctx)),
		}
	}

	if err := a.oauthValidator.ValidateTokenAuthorizationCode(ctx, accessTokenReq, oauthAccessGrant); err != nil {
		return entity.OauthAccessToken{}, oauthAccessGrant.RedirectURI.String, err
	}

	id, err := a.oauthAccessTokenModel.Create(kontext.Fabricate(kontext.WithDefaultContext(ctx)), a.modelFormatter.AccessTokenFromOauthAccessGrant(oauthAccessGrant, oauthApplication), a.sqldb)
	if err != nil {

		log.Error().
			Err(err).
			Stack().
			Interface("request_id", ctx.Value("request_id")).
			Array("tags", zerolog.Arr().Str("service").Str("authorization").Str("grant_token")).
			Msg("Error creating access token after creating the data")

		return entity.OauthAccessToken{}, oauthAccessGrant.RedirectURI.String, &entity.Error{
			HttpStatus: http.StatusInternalServerError,
			Errors:     eobject.Wrap(eobject.InternalServerError(ctx)),
		}
	}

	oauthAccessToken, err := a.oauthAccessTokenModel.One(kontext.Fabricate(kontext.WithDefaultContext(ctx)), id, a.sqldb)
	if err != nil {

		log.Error().
			Err(err).
			Stack().
			Interface("request_id", ctx.Value("request_id")).
			Array("tags", zerolog.Arr().Str("service").Str("authorization").Str("token")).
			Msg("Error selecting one access token")

		return entity.OauthAccessToken{}, oauthAccessGrant.RedirectURI.String, &entity.Error{
			HttpStatus: http.StatusInternalServerError,
			Errors:     eobject.Wrap(eobject.InternalServerError(ctx)),
		}
	}

	err = a.oauthAccessGrantModel.Revoke(kontext.Fabricate(kontext.WithDefaultContext(ctx)), *accessTokenReq.Code, a.sqldb)
	if err != nil {
		// TODO: Error is intended to be suppressed until database transaction is implemented. After database transaction is implemented, then it will be rollbacked if there is error in revoke oauth access grants process
		log.Error().
			Err(err).
			Stack().
			Interface("request_id", ctx.Value("request_id")).
			Array("tags", zerolog.Arr().Str("service").Str("authorization").Str("token")).
			Msg("Error revoke oauth access grant")
	}

	return oauthAccessToken, oauthAccessGrant.RedirectURI.String, nil
}

// RevokeToken revoke given access token request
func (a *Authorization) RevokeToken(ctx context.Context, revokeAccessTokenReq entity.RevokeAccessTokenRequestJSON) *entity.Error {

	if revokeAccessTokenReq.Token == nil {
		return &entity.Error{
			HttpStatus: http.StatusUnprocessableEntity,
			Errors:     eobject.Wrap(eobject.ValidationError("token is empty")),
		}
	}

	err := a.oauthAccessTokenModel.Revoke(kontext.Fabricate(kontext.WithDefaultContext(ctx)), *revokeAccessTokenReq.Token, a.sqldb)
	if err != nil && err.Type() == exception.NotFound {
		return &entity.Error{
			HttpStatus: http.StatusNotFound,
			Errors:     eobject.Wrap(eobject.NotFoundError(ctx, "token")),
		}
	} else if err != nil {
		return &entity.Error{
			HttpStatus: http.StatusInternalServerError,
			Errors:     eobject.Wrap(eobject.InternalServerError(ctx)),
		}
	}

	return nil
}

func (a *Authorization) exceptionMapping(ctx context.Context, exc exception.Exception, tags *zerolog.Array) *entity.Error {
	log.Error().
		Err(exc).
		Stack().
		Interface("request_id", ctx.Value("request_id")).
		Array("tags", tags).
		Msg(exc.Detail())

	switch exc.Type() {
	case exception.NotFound:
		return &entity.Error{
			HttpStatus: http.StatusNotFound,
			Errors: eobject.Wrap(entity.ErrorObject{
				Code:    exc.Title(),
				Message: exc.Detail(),
			}),
		}

	case exception.Unauthorized:
		return &entity.Error{
			HttpStatus: http.StatusUnauthorized,
			Errors:     eobject.Wrap(eobject.UnauthorizedError()),
		}

	case exception.Forbidden:
		return &entity.Error{
			HttpStatus: http.StatusForbidden,
			Errors: eobject.Wrap(entity.ErrorObject{
				Code:    exc.Title(),
				Message: exc.Detail(),
			}),
		}

	default:
		return &entity.Error{
			HttpStatus: http.StatusInternalServerError,
			Errors:     eobject.Wrap(eobject.InternalServerError(ctx)),
		}
	}
}
