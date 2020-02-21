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
	One(ctx context.Context, ID int) (entity.OauthApplication, error)
	OneByUIDandSecret(ctx context.Context, clientUID, clientSecret string) (entity.OauthApplication, error)
	Count(ctx context.Context) (int, error)
	Create(ctx context.Context, data entity.OauthApplicationJSON, txs ...*sql.Tx) (int, error)
}

type OauthAccessTokenModel interface {
	HasName
	One(ctx context.Context, ID int) (entity.OauthAccessToken, error)
	Create(ctx context.Context, data entity.OauthAccessTokenInsertable, txs ...*sql.Tx) (int, error)
}

type OauthAccessGrantModel interface {
	HasName
	One(ctx context.Context, ID int) (entity.OauthAccessGrant, error)
}
