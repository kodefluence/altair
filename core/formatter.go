package core

import (
	"context"

	"github.com/codefluence-x/altair/entity"
)

type OauthApplicationFormater interface {
	ApplicationList(ctx context.Context, applications []entity.OauthApplication) []entity.OauthApplicationJSON
}
