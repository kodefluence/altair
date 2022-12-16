package usecase

import (
	"github.com/kodefluence/altair/plugin/oauth/entity"
	"github.com/kodefluence/altair/util"
)

func (f *Formatter) Application(application entity.OauthApplication) entity.OauthApplicationJSON {
	oauthApplicationJSON := entity.OauthApplicationJSON{
		ID:           &application.ID,
		OwnerType:    &application.OwnerType,
		ClientUID:    &application.ClientUID,
		ClientSecret: &application.ClientSecret,
		CreatedAt:    &application.CreatedAt,
		UpdatedAt:    &application.UpdatedAt,
	}

	if application.OwnerID.Valid {
		oauthApplicationJSON.OwnerID = util.ValueToPointer(int(application.OwnerID.Int64))
	}

	if application.RevokedAt.Valid {
		oauthApplicationJSON.RevokedAt = &application.RevokedAt.Time
	}

	if application.Description.Valid {
		oauthApplicationJSON.Description = &application.Description.String
	}

	if application.Scopes.Valid {
		oauthApplicationJSON.Scopes = &application.Scopes.String
	}

	return oauthApplicationJSON
}
