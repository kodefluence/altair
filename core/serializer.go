package core

import (
	"context"

	"github.com/codefluence-x/altair/entity"
)

type OauthApplicationSerializer interface {
	ApplicationList(ctx context.Context, applications []entity.OauthApplication) []entity.OauthApplicationJSON
}
