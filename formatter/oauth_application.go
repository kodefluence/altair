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

		if v.OwnerID.Valid {
			oauthApplicationJSON[k].OwnerID = util.IntToPointer(int(v.OwnerID.Int64))
		}

		if v.RevokedAt.Valid {
			oauthApplicationJSON[k].RevokedAt = &v.RevokedAt.Time
		}
	}

	return oauthApplicationJSON
}

func (oa oauthApplication) Application(ctx context.Context, application entity.OauthApplication) entity.OauthApplicationJSON {
	return entity.OauthApplicationJSON{
		ID:           &application.ID,
		Description:  &application.Description,
		Scopes:       &application.Scopes,
		ClientUID:    &application.ClientUID,
		ClientSecret: &application.ClientSecret,
		CreatedAt:    &application.CreatedAt,
		UpdatedAt:    &application.UpdatedAt,
	}
}
