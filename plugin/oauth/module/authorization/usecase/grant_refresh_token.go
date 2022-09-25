package usecase

import (
	"github.com/kodefluence/altair/plugin/oauth/entity"
	"github.com/kodefluence/monorepo/db"
	"github.com/kodefluence/monorepo/jsonapi"
	"github.com/kodefluence/monorepo/kontext"
)

func (a *Authorization) GrantRefreshToken(ktx kontext.Context, oauthAccessToken entity.OauthAccessToken, oauthApplication entity.OauthApplication, tx db.TX) (entity.OauthRefreshToken, jsonapi.Errors) {
	refreshTokenID, err := a.oauthRefreshTokenRepo.Create(ktx, a.formatter.RefreshTokenInsertable(oauthApplication, oauthAccessToken), tx)
	if err != nil {
		return entity.OauthRefreshToken{}, jsonapi.BuildResponse(a.apiError.InternalServerError(ktx)).Errors
	}

	oauthRefreshToken, err := a.oauthRefreshTokenRepo.One(ktx, refreshTokenID, tx)
	if err != nil {
		return entity.OauthRefreshToken{}, jsonapi.BuildResponse(a.apiError.InternalServerError(ktx)).Errors
	}

	return oauthRefreshToken, nil
}
