package formatter

import (
	"context"

	"github.com/codefluence-x/altair/core"
	"github.com/codefluence-x/altair/entity"
	"github.com/codefluence-x/altair/util"
)

type oauthApplication struct{}

func OauthApplication() core.OauthApplicationFormater {
	return oauthApplication{}
}

func (oa oauthApplication) ApplicationList(ctx context.Context, applications []entity.OauthApplication) []entity.OauthApplicationJSON {
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
		Description:  &application.Description,
		Scopes:       &application.Scopes,
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

	return oauthApplicationJSON
}
