package entity

import (
	"database/sql"
	"time"

	"github.com/go-sql-driver/mysql"
)

type OauthApplication struct {
	ID           int
	OwnerID      sql.NullInt64
	Description  string
	Scopes       string
	ClientUID    string
	ClientSecret string
	RevokedAt    mysql.NullTime
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
