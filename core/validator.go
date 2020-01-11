package core

import (
	"context"

	"github.com/codefluence-x/altair/entity"
)

type OauthApplicationValidator interface {
	ValidateCreate(ctx context.Context, data entity.OauthApplicationJSON) *entity.Error
}
