package command

import (
	"github.com/kodefluence/altair/plugin/oauth/entity"
	"github.com/kodefluence/monorepo/jsonapi"
	"github.com/kodefluence/monorepo/kontext"
)

//go:generate mockgen -destination ./mock/mock.go -package mock -source ./command.go
type ApplicationManager interface {
	Create(ktx kontext.Context, e entity.OauthApplicationJSON) (entity.OauthApplicationJSON, jsonapi.Errors)
}
