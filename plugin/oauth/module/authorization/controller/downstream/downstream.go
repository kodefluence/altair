package downstream

import (
	"github.com/kodefluence/altair/plugin/oauth/entity"
	"github.com/kodefluence/monorepo/db"
	"github.com/kodefluence/monorepo/exception"
	"github.com/kodefluence/monorepo/kontext"
)

//go:generate mockgen -destination ./mock/mock.go -package mock -source ./downstream.go
type OauthAccessTokenRepository interface {
	OneByToken(ktx kontext.Context, token string, tx db.TX) (entity.OauthAccessToken, exception.Exception)
}
