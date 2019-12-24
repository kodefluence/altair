package core

import (
	"context"

	"github.com/codefluence-x/altair/entity"
)

type ApplicationManager interface {
	List(ctx context.Context, offset, limit int) ([]entity.OauthApplicationJSON, int, *entity.Error)
}
