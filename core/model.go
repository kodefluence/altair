package core

import (
	"context"
	"database/sql"

	"github.com/codefluence-x/altair/entity"
)

type HasName interface {
	Name() string
}

type DBExecutable interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
}

type OauthApplicationModel interface {
	HasName
	Paginate(ctx context.Context, offset, limit int) ([]entity.OauthApplication, error)
	Count(ctx context.Context) (int, error)
	Create(ctx context.Context, data *entity.OauthApplicationJSON, txs ...*sql.Tx) (int, error)
}
