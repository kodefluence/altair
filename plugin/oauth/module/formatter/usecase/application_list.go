package usecase

import (
	"github.com/kodefluence/altair/plugin/oauth/entity"
)

func (f *Formatter) ApplicationList(applications []entity.OauthApplication) []entity.OauthApplicationJSON {
	oauthApplicationJSON := make([]entity.OauthApplicationJSON, len(applications))

	for k, v := range applications {
		oauthApplicationJSON[k] = f.Application(v)
	}

	return oauthApplicationJSON
}
