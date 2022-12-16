package usecase

import (
	"time"

	"github.com/kodefluence/altair/plugin/oauth/entity"
	"github.com/kodefluence/altair/util"
)

func (*Formatter) RefreshToken(e entity.OauthRefreshToken) entity.OauthRefreshTokenJSON {
	var data entity.OauthRefreshTokenJSON

	data.CreatedAt = &e.CreatedAt
	data.Token = &e.Token

	if time.Now().Before(e.ExpiresIn) {
		data.ExpiresIn = util.ValueToPointer(int(time.Until(e.ExpiresIn).Seconds()))
	} else {
		data.ExpiresIn = util.ValueToPointer(0)
	}

	if e.RevokedAT.Valid {
		data.RevokedAT = &e.RevokedAT.Time
	}

	return data
}
