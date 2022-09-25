package http

import (
	"github.com/kodefluence/altair/plugin/oauth/entity"
	"github.com/kodefluence/monorepo/jsonapi"
	"github.com/kodefluence/monorepo/kontext"
)

//go:generate mockgen -destination ./mock/mock.go -package mock -source ./http.go
type Authorization interface {
	GrantAuthorizationCode(ktx kontext.Context, authorizationReq entity.AuthorizationRequestJSON) (interface{}, jsonapi.Errors)
	GrantToken(ktx kontext.Context, accessTokenReq entity.AccessTokenRequestJSON) (entity.OauthAccessTokenJSON, jsonapi.Errors)
	RevokeToken(ktx kontext.Context, revokeAccessTokenReq entity.RevokeAccessTokenRequestJSON) jsonapi.Errors
}
