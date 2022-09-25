package usecase

import (
	"github.com/kodefluence/altair/plugin/oauth/entity"
	"github.com/kodefluence/monorepo/db"
	"github.com/kodefluence/monorepo/exception"
	"github.com/kodefluence/monorepo/jsonapi"
	"github.com/kodefluence/monorepo/kontext"
	"github.com/rs/zerolog"
)

// Grant authorization an access code
func (a *Authorization) Grant(ktx kontext.Context, authorizationReq entity.AuthorizationRequestJSON) (entity.OauthAccessGrantJSON, jsonapi.Errors) {
	var oauthAccessGrantJSON entity.OauthAccessGrantJSON

	oauthApplication, jsonapiErr := a.FindAndValidateApplication(ktx, authorizationReq.ClientUID, authorizationReq.ClientSecret)
	if jsonapiErr != nil {
		return oauthAccessGrantJSON, jsonapiErr
	}

	if err := a.ValidateAuthorizationGrant(ktx, authorizationReq, oauthApplication); err != nil {
		return oauthAccessGrantJSON, err
	}

	exc := a.sqldb.Transaction(ktx, "authorization-grant-authorization-code", func(tx db.TX) exception.Exception {
		id, err := a.oauthAccessGrantRepo.Create(ktx, a.formatter.AccessGrantFromAuthorizationRequestInsertable(authorizationReq, oauthApplication), tx)
		if err != nil {
			return exception.Throw(err, exception.WithDetail("error creating authorization code"))
		}

		oauthAccessGrant, err := a.oauthAccessGrantRepo.One(ktx, id, tx)
		if err != nil {
			return exception.Throw(err, exception.WithDetail("error selecting newly created authorization code"))
		}

		oauthAccessGrantJSON = a.formatter.AccessGrant(oauthAccessGrant)
		return nil
	})
	if exc != nil {
		return oauthAccessGrantJSON, a.exceptionMapping(ktx, exc, zerolog.Arr().Str("service").Str("authorization").Str("grant"))
	}

	return oauthAccessGrantJSON, nil
}
