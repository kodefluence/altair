package formatter

import (
	"context"

	"github.com/codefluence-x/altair/provider/plugin/oauth/entity"
	"github.com/codefluence-x/altair/util"
)

type oauthApplication struct{}

func OauthApplication() *oauthApplication {
	return &oauthApplication{}
}

func (oa *oauthApplication) ApplicationList(ctx context.Context, applications []entity.OauthApplication) []entity.OauthApplicationJSON {
	oauthApplicationJSON := make([]entity.OauthApplicationJSON, len(applications))

	for k, v := range applications {
		oauthApplicationJSON[k] = oa.Application(ctx, v)
	}

	return oauthApplicationJSON
}

func (oa oauthApplication) Application(ctx context.Context, application entity.OauthApplication) entity.OauthApplicationJSON {
	oauthApplicationJSON := entity.OauthApplicationJSON{
		ID:           &application.ID,
		OwnerType:    &application.OwnerType,
		ClientUID:    &application.ClientUID,
		ClientSecret: &application.ClientSecret,
		CreatedAt:    &application.CreatedAt,
		UpdatedAt:    &application.UpdatedAt,
	}

	if application.OwnerID.Valid {
		oauthApplicationJSON.OwnerID = util.IntToPointer(int(application.OwnerID.Int64))
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
