package usecase

import (
	"errors"
	"time"

	"github.com/kodefluence/altair/plugin/oauth/entity"
	"github.com/kodefluence/monorepo/exception"
	"github.com/kodefluence/monorepo/kontext"
)

func (a *Authorization) ValidateTokenAuthorizationCode(ktx kontext.Context, r entity.AccessTokenRequestJSON, data entity.OauthAccessGrant) exception.Exception {
	if data.RevokedAT.Valid {
		return exception.Throw(errors.New("forbidden"), exception.WithType(exception.Forbidden), exception.WithDetail("authorization code already used"), exception.WithTitle("Forbidden resource access"))
	}

	if time.Now().After(data.ExpiresIn) {
		return exception.Throw(errors.New("forbidden"), exception.WithType(exception.Forbidden), exception.WithDetail("authorization code already expired"), exception.WithTitle("Forbidden resource access"))
	}

	if data.RedirectURI.String != *r.RedirectURI {
		return exception.Throw(errors.New("forbidden"), exception.WithType(exception.Forbidden), exception.WithDetail("redirect uri is different from one that generated before"), exception.WithTitle("Forbidden resource access"))
	}

	return nil
}
