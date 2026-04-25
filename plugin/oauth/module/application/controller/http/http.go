package http

import (
	"github.com/kodefluence/monorepo/jsonapi"
	"github.com/kodefluence/monorepo/kontext"

	"github.com/kodefluence/altair/plugin/oauth/entity"
)

//go:generate mockgen -destination ./mock/mock.go -package mock -source ./http.go

type ApplicationManager interface {
	Create(ktx kontext.Context, e entity.OauthApplicationJSON) (entity.OauthApplicationJSON, jsonapi.Errors)
	List(ktx kontext.Context, offset, limit int) ([]entity.OauthApplicationJSON, int, jsonapi.Errors)
	One(ktx kontext.Context, ID int) (entity.OauthApplicationJSON, jsonapi.Errors)
	Update(ktx kontext.Context, ID int, e entity.OauthApplicationUpdateJSON) (entity.OauthApplicationJSON, jsonapi.Errors)
}
