package core

import (
	"context"

	"github.com/codefluence-x/altair/entity"
)

type HasName interface {
	Name() string
}

type OauthApplicationModel interface {
	HasName
	Paginate(ctx context.Context, offset, limit int) ([]entity.OauthApplication, error)
}
