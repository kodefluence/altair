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
		oauthApplicationJSON[k] = entity.OauthApplicationJSON{
			ID:           &v.ID,
			Description:  &v.Description,
			Scopes:       &v.Scopes,
			ClientUID:    &v.ClientUID,
			ClientSecret: &v.ClientSecret,
			CreatedAt:    &v.CreatedAt,
			UpdatedAt:    &v.UpdatedAt,
		}

		if v.OwnerID.Valid {
			oauthApplicationJSON[k].OwnerID = util.IntToPointer(int(v.OwnerID.Int64))
		}

		if v.RevokedAt.Valid {
			oauthApplicationJSON[k].RevokedAt = &v.RevokedAt.Time
		}
	}

	return oauthApplicationJSON
}
