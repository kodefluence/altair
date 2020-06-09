package model

import "database/sql"

type DBExecutable interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
}
