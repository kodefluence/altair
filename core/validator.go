package core

import (
	"context"

	"github.com/codefluence-x/altair/entity"
)

type OauthValidator interface {
	ValidateApplication(ctx context.Context, data entity.OauthApplicationJSON) *entity.Error
}
